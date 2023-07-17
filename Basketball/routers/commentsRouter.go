package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context/param"
)

func init() {

	beego.GlobalControllerRouter["Basketball/controllers:PostController"] = append(beego.GlobalControllerRouter["Basketball/controllers:PostController"],
		beego.ControllerComments{
			Method:           "Post",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["Basketball/controllers:PostController"] = append(beego.GlobalControllerRouter["Basketball/controllers:PostController"],
		beego.ControllerComments{
			Method:           "GetAll",
			Router:           `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["Basketball/controllers:PostController"] = append(beego.GlobalControllerRouter["Basketball/controllers:PostController"],
		beego.ControllerComments{
			Method:           "Get",
			Router:           `/:postId`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["Basketball/controllers:PostController"] = append(beego.GlobalControllerRouter["Basketball/controllers:PostController"],
		beego.ControllerComments{
			Method:           "Put",
			Router:           `/:postId`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["Basketball/controllers:PostController"] = append(beego.GlobalControllerRouter["Basketball/controllers:PostController"],
		beego.ControllerComments{
			Method:           "Delete",
			Router:           `/:postId`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
