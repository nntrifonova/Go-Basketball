package utils

import (
	"B/models"

	"github.com/beego/beego/v2/adapter/orm"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
)

func init() {
	models.Init()
	corsHandler := func(ctx *context.Context) {
		ctx.Output.Header("Access-Control-Allow-Origin", ctx.Input.Domain())
		ctx.Output.Header("Access-Control-Allow-Methods", "*")
	}
	beego.InsertFilter("*", beego.BeforeRouter, corsHandler)
}

func Main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
		orm.Debug = true
	}
	beego.Run()
}

//social auth
