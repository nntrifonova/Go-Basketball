package routers

import (
	"basketball/controllers"
	"github.com/astaxie/beego"
	"github.com/beego/admin"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	admin.Run()

}
