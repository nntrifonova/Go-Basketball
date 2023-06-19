package main

import (
	_ "basketball/routers"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	// Register `mysql` driver
	_ = orm.RegisterDriver("mysql", orm.DRMySQL)
	// Register `default` database
	_ = orm.RegisterDataBase("default", "mysql", "admin:pass@mysql-m/my_db_name?charset=utf8")
	// Run migrations to create tables
	force := true // Drop old table and create new
	err := orm.RunSyncdb("default", force, beego.BConfig.RunMode == "dev")
	if err != nil {
		fmt.Printf("An Error Occurred: %v", err)
	}
}

func main() {
	beego.Run()
	//var BeeLogger = logs.GetBeeLogger()
}
