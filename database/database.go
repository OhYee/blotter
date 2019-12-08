package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/OhYee/gosql/utils/connect"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/OhYee/blotter/database/proto"
	"github.com/micro/go-micro"
)

//go:generate protoc --proto_path=$GOPATH/src:. --micro_out=. --go_out=. ./proto/proto.proto

// Database server object
type Database struct {
	dsn connect.DataSourceName
}

// Query database using SqlString
func (d *Database) Query(ctx context.Context, req *proto.QueryRequest, rsp *proto.QueryResponse) (err error) {
	conn, err := connect.NewConnection(&d.dsn)
	if err != nil {
		return
	}
	defer conn.Close()

	res, err := conn.Query(req.SqlString)
	if err != nil {
		return
	}
	rsp.Result, err = json.Marshal(res)
	return
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	service := micro.NewService(
		micro.Name("database"),
	)
	service.Init()

	data, err := ioutil.ReadFile("database.json")
	checkErr(err)

	m := make(map[string]string)
	err = json.Unmarshal(data, &m)
	checkErr(err)

	port, err := strconv.ParseInt(m["port"], 10, 32)
	checkErr(err)

	proto.RegisterDatabaseHandler(service.Server(), &Database{connect.DataSourceName{
		Username: m["username"],
		Password: m["password"],
		Address:  m["address"],
		Port:     int(port),
	}})
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
