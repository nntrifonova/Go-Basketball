// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"Basketball/controllers"
	"github.com/beego/beego/v2/server/web/context"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	beego.Get("/", func(ctx *context.Context) {
		_ = ctx.Output.Body([]byte("This is a Beego + JWT API - Creator: Mehran Abghari (mehran.ab80@gmail.com)"))
	})
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/object",
			beego.NSInclude(
				&controllers.ObjectController{},
			),
		),
		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
