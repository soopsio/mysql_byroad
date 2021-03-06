package main

import (
	"fmt"
	"mysql_byroad/model"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	rpcserver   *RPCServer
	rpcclient   *RPCClient
	sendClient  *SendClient
	taskManager *TaskConsumerManager
)

func initGlobal() {
	var err error
	rpcserver = NewRPCServer("tcp", fmt.Sprintf("%s:%d", Conf.RPCServerConf.Host, Conf.RPCServerConf.Port), Conf.RPCServerConf.Desc)
	rpcserver.startRpcServer()
	rpcclient = NewRPCClient("tcp", fmt.Sprintf("%s:%d", Conf.MonitorConf.Host, Conf.MonitorConf.RpcPort), "")
	_, err = rpcclient.RegisterClient(rpcserver.getSchema(), rpcserver.desc)
	if err != nil {
		log.Error("register rpc client error: ", err.Error())
	}
	sendClient = NewSendClient()
}

func main() {
	InitLog()
	InitAlert(&Conf.AlertConfig)
	log.Debugf("Conf: %+v", Conf)
	initGlobal()
	dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=true",
		Conf.MysqlConf.Username, Conf.MysqlConf.Password, Conf.MysqlConf.Host, Conf.MysqlConf.Port,
		Conf.MysqlConf.DBName)
	confdb, err := sqlx.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	model.Init(confdb)
	tasks, err := rpcclient.GetAllTasks("")
	if err != nil {
		log.Error(err.Error())
	}

	taskManager = NewTaskConsumerManager()
	taskManager.Init(tasks)
	HandleSignal()
}

// HandleSignal fetch signal from chan then do exit or reload.
func HandleSignal() {
	// Block until a signal is received.
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGSTOP)
	for {
		s := <-c
		log.Infof("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			_, err := rpcclient.DeregisterClient(rpcserver.getSchema(), rpcserver.desc)
			if err != nil {
				log.Error("rpc deregister error: ", err.Error())
			}
			taskManager.StopAllTask()
			time.Sleep(1 * time.Second)
			return
		case syscall.SIGHUP:
			// TODO reload
			//return
		default:
			return
		}
	}
}
