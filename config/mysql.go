package config

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gohouse/gorose/v2"
)

//var rds *redis.Client
//var keySsc = "z:ssc" // zadd [keySsc] [score] Data
var engin *gorose.Engin

func mysqlInit(conf *gorose.Config) {
	if conf == nil {
		return
	}
	var err error
	engin, err = gorose.Open(conf)
	if err != nil {
		panic(err.Error())
	}
}
func DB() gorose.IOrm {
	return engin.NewOrm()
}
