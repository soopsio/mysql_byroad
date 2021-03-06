package schema

import (
	"fmt"
	"log"
	"mysql_byroad/model"
	"sort"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type MysqlConfig struct {
	Name     string
	Host     string
	Port     uint16
	Username string
	Password string
	Include  []string
	Exclude  []string
	Interval time.Duration
}

type Inspector struct {
	config    *MysqlConfig
	db        *sqlx.DB
	columnMap *ColumnMap
	sync.RWMutex
}

func NewInspector(config *MysqlConfig) (*Inspector, error) {
	inspector := Inspector{
		config: config,
	}
	db, err := inspector.connect()
	if err != nil {
		return &inspector, err
	}
	inspector.db = db
	return &inspector, nil
}

func (this *Inspector) connect() (*sqlx.DB, error) {
	config := this.config
	dsn := fmt.Sprintf("%s:%s@(%s:%d)/information_schema", config.Username, config.Password, config.Host, config.Port)
	db, err := sqlx.Open("mysql", dsn)
	return db, err
}

/*
开始定时刷新字段信息
*/
func (this *Inspector) LookupLoop() {
	clist, err := this.getColumns()
	if err != nil {
		fmt.Printf("get column error: %s", err.Error())
	}
	this.buildColumnMap(clist)
	ticker := time.NewTicker(this.config.Interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				clist, err := this.getColumnsMultiTimes()
				if err != nil {
					log.Printf("[ERROR] get columns error: %s", err.Error())
					continue
				}
				this.buildColumnMap(clist)
			}
		}
	}()
}

func (this *Inspector) BuildColumnMap() error {
	clist, err := this.getColumns()
	if err != nil {
		return err
	}
	this.buildColumnMap(clist)
	return nil
}

func (this *Inspector) buildColumnMap(columnList ColumnList) {
	cm := BuildColumnMap(columnList)
	this.Lock()
	this.columnMap = cm
	this.Unlock()
}

func (this *Inspector) GetColumnMap() *ColumnMap {
	this.RLock()
	defer this.RUnlock()
	return this.columnMap
}

func (this *Inspector) getColumns() (ColumnList, error) {
	sqlStr := "SELECT TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, DATA_TYPE, COLUMN_TYPE FROM columns "
	nodisplay := this.getNoDisplaySchema()
	display := this.getDisplaySchema()
	if nodisplay != "" && display != "" {
		sqlStr += "WHERE " + nodisplay + " AND " + display
	} else if nodisplay != "" || display != "" {
		sqlStr += "WHERE " + nodisplay + display
	} else {
	}
	var columnList = make([]*Column, 0, 10)
	err := this.db.Select(&columnList, sqlStr)
	return columnList, err
}

/*
每次查询一个数据库的字段信息，防止一次查询时间过长
*/
func (this *Inspector) getColumnsMultiTimes() (ColumnList, error) {
	schemas, err := this.GetSchemas()
	if err != nil {
		return nil, err
	}
	sqlStr := "SELECT TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, DATA_TYPE, COLUMN_TYPE FROM columns WHERE TABLE_SCHEMA=? "
	var columnList = make([]*Column, 0, 100)
	for _, schema := range schemas {
		var columns = make([]*Column, 0, 10)
		err := this.db.Select(&columns, sqlStr, schema)
		if err != nil {
			return nil, err
		}
		columnList = append(columnList, columns...)
	}
	return columnList, err
}

func (this *Inspector) getNoDisplaySchema() string {
	var data string
	for _, schema := range this.config.Exclude {
		data = data + "'" + schema + "'" + ","
	}
	if data != "" {
		data = strings.TrimRight(data, ",")
		return "TABLE_SCHEMA NOT IN (" + data + ")"
	}
	return ""
}

func (this *Inspector) getDisplaySchema() string {
	var data string
	for _, schema := range this.config.Include {
		data = data + "'" + schema + "'" + ","
	}
	if data != "" {
		data = strings.TrimRight(data, ",")
		return "TABLE_SCHEMA IN (" + data + ")"
	}
	return ""
}

func (this *Inspector) Close() error {
	return this.db.Close()
}

func (this *Inspector) GetSchemas() (schemas []string, err error) {
	sqlStr := "SELECT DISTINCT TABLE_SCHEMA FROM columns "
	nodisplay := this.getNoDisplaySchema()
	if nodisplay != "" {
		sqlStr += "WHERE TABLE_SCHEMA NOT IN (" + nodisplay + ") "
	}
	sqlStr += "ORDER BY TABLE_SCHEMA"
	err = this.db.Select(&schemas, sqlStr)
	return
}

func (this *Inspector) GetTables(schema string) (tables []string, err error) {
	sqlStr := "SELECT DISTINCT TABLE_NAME FROM columns WHERE TABLE_SCHEMA=? ORDER BY TABLE_NAME"
	err = this.db.Select(&tables, sqlStr, schema)
	return
}

func (this *Inspector) GetColumns(schema, table string) (columns []string, err error) {
	sqlStr := "SELECT DISTINCT COLUMN_NAME FROM columns WHERE TABLE_SCHEMA=? AND TABLE_NAME=? ORDER BY COLUMN_NAME"
	err = this.db.Select(&columns, sqlStr, schema, table)
	return
}

func getOrderedColumnsList(columnMap *ColumnMap) model.OrderedSchemas {
	colslist := make(model.OrderedSchemas, 0, 10)
	for schema, tables := range columnMap.columns {
		os := new(model.OrderedSchema)
		os.Schema = schema
		os.OrderedTables = make([]*model.OrderedTable, 0, 10)
		for table, columns := range tables {
			ot := new(model.OrderedTable)
			ot.Table = table
			ot.Columns = make([]string, 0, 10)
			for _, column := range columns {
				ot.Columns = append(ot.Columns, column.Name)
			}
			os.OrderedTables = append(os.OrderedTables, ot)
		}
		colslist = append(colslist, os)
	}
	sort.Sort(colslist)
	for _, tab := range colslist {
		sort.Sort(model.OrderedTables(tab.OrderedTables))
	}
	return colslist
}

func (this *Inspector) GetOrderedColumns() model.OrderedSchemas {
	return getOrderedColumnsList(this.GetColumnMap())
}
