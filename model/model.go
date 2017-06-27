package model

import (
	"github.com/go-xorm/xorm"
	_ "github.com/go-sql-driver/mysql"
)

type Model struct {
	engine int
}

var engine *xorm.Engine

func getEngine()(*xorm.Engine){
	var err error
	if engine == nil{
		engine, err = xorm.NewEngine("mysql", "root@/test_gin")
		if err != nil{
			panic(err)
		}
	}
	return engine
}
