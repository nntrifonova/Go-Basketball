package models

import (
	"fmt"
	"github.com/beego/beego/v2/adapter/config"
	"github.com/beego/beego/v2/adapter/orm"
	_ "github.com/go-sql-driver/mysql"
)

func Init() {
	appConf, err := config.NewConfig("ini", "conf/app.conf")
	if err != nil {
		panic(err)
	}

	conn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s",
		appConf.String("database::db_user"),
		appConf.String("database::db_passwd"),
		appConf.String("database::db_host"),
		appConf.String("database::db_port"),
		appConf.String("database::db_name"),
		appConf.String("database::db_charset"))
	orm.RegisterDataBase("default", "mysql", conn)

	name := "default"
	force := false
	verbose := true
	err = orm.RunSyncdb(name, force, verbose)
	if err != nil {
		fmt.Println(err)
	}
	orm.RunCommand()
}

func TableName(str string) string {
	appConf, err := config.NewConfig("ini", "conf/app.conf")
	if err != nil {
		panic(err)
	}
	return appConf.String("database::db_prefix") + str
}
