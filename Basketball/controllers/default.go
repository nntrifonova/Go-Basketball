package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
)

type BaseController struct {
	beego.Controller
	isLogin bool
}

func (c *BaseController) Resp(status int, data interface{}, err error) {
	if err == nil {
		c.Data["json"] = data

	} else {
		c.Data["json"] = err.Error()
	}
	c.Ctx.Output.SetStatus(status)
	c.ServeJSON()
	c.StopRun()
}

// GetString returns the input value by key string or the default value while it's present and input is blank
func (c *BaseController) GetString(key string, def ...string) string {
	if v := c.Ctx.Input.Query(key); v != "" {
		return v
	}
	if len(def) > 0 {
		return def[0]
	}
	return ""
}

//func (this *BaseController) Prepare() {
//	userLogin := this.GetSession("userLogin")
//	fmt.Println("userLogin", userLogin)
//	if userLogin == nil {
//		this.isLogin = false
//	} else {
//		this.isLogin = true
//	}
//	this.Data["isLogin"] = this.isLogin
//}
