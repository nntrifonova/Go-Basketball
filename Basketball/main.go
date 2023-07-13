package main

import (
	_ "Basketball/routers"
	servicesDb "Basketball/services/db"
	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	servicesDb.InitSql()
	beego.Run()
}
