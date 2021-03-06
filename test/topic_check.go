package main

import (
	"fmt"
	"time"

	"flag"

	"github.com/samuel/go-zookeeper/zk"
)

var path string
var zkaddr string

func init() {
	flag.Parse()
	args := flag.Args()
	if len(args) > 1 {
		path = args[1]
		zkaddr = args[0]
	} else {
		path = "/"
	}
}
func getAllTopics() ([]string, error) {
	zkAddrs := []string{zkaddr}
	conn, _, err := zk.Connect(zkAddrs, time.Second*10)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	children, _, err := conn.Children(path)
	if err != nil {
		return nil, err
	}
	fmt.Printf("get all topics: %+v\n", children)
	for _, child := range children {
		fmt.Println("get ", path+"/"+child)
		if exists, _, _ := conn.Exists(path + "/" + child); exists {
			data, _, _ := conn.Get(path + "/" + child)
			fmt.Printf("%s:%s\n", path+"/"+child, string(data))
		}
	}
	return children, nil
}

func main() {
	fmt.Printf("check path %s\n", path)
	_, err := getAllTopics()
	if err != nil {
		fmt.Println("error: ", err.Error())
	}
}
