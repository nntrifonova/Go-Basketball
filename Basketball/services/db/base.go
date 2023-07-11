package db

import (
	"Basketball/conf"
	"github.com/beego/beego/v2/adapter/orm"
	"github.com/beego/beego/v2/core/logs"
	"time"
)

var dbServer = conf.GetEnvConst("DB_SERVER")
var dbPort = conf.GetEnvConst("DB_PORT")
var dbName = conf.GetEnvConst("DB_NAME")
var dbUser = conf.GetEnvConst("DB_USER")
var dbUserPass = conf.GetEnvConst("DB_USER_PASS")

func InitSql() {
	if err := orm.RegisterDriver("postgres", orm.DRPostgres); err != nil {
		logs.Error(err)
	}
	dbparams := "user=" + dbUser +
		" password=" + dbUserPass +
		" host=" + dbServer +
		" port=" + dbPort +
		" dbname=" + dbName +
		" sslmode=disable"

	if err := orm.RegisterDataBase("default", "postgres", dbparams); err != nil {
		logs.Error(err)
	}

	orm.DefaultTimeLoc = time.UTC
}
