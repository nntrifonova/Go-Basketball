package db

import (
	"Basketball/conf"
	"fmt"
	"github.com/beego/beego/v2/adapter/orm"
	"github.com/beego/beego/v2/core/logs"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var dbServer = conf.GetEnvConst("DB_SERVER")
var dbPort = conf.GetEnvConst("DB_PORT")
var dbName = conf.GetEnvConst("DB_NAME")
var dbUser = conf.GetEnvConst("DB_USER")
var dbUserPass = conf.GetEnvConst("DB_USER_PASS")

func InitSql() {
	if err := orm.RegisterDriver("mysql", orm.DRMySQL); err != nil {
		logs.Error(err)
	}

	path := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
		dbUser,
		dbUserPass,
		dbServer,
		dbPort,
		dbName,
	)

	if err := orm.RegisterDataBase("default", "mysql", path); err != nil {
		logs.Error(err)
	}

	orm.DefaultTimeLoc = time.UTC
}
