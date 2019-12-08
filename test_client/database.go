package main

import (
	"context"
	"encoding/json"

	"github.com/OhYee/blotter/database/proto"
	"github.com/micro/go-micro"
)

func makeQuery(sqlStr string) (m []map[string]interface{}, err error) {
	service := micro.NewService(micro.Name("site.client"))
	service.Init()
	database := proto.NewDatabaseService("database", service.Client())

	m = make([]map[string]interface{}, 0)
	rsp, err := database.Query(context.TODO(), &proto.QueryRequest{SqlString: sqlStr})
	if err != nil {
		return m, err
	}

	err = json.Unmarshal(rsp.Result, &m)

	return m, err
}
