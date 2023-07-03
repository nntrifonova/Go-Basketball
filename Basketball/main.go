package main

import (
	_ "Basketball/routers"
	"fmt"
	"github.com/beego/beego/v2/client/orm"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	// Register `mysql` driver
	_ = orm.RegisterDriver("mysql", orm.DRMySQL)
	// Register `default` database
	_ = orm.RegisterDataBase("default", "mysql", "root:@tcp(127.0.0.1:3306)/beego_jwt_test?charset=utf8")
	// Run migrations to create tables
	force := true // Drop old table and create new
	err := orm.RunSyncdb("default", force, beego.BConfig.RunMode == "dev")
	if err != nil {
		fmt.Printf("An Error Occurred: %v", err)
	}
}

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}
