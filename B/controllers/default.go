package controllers

type DefaultController struct {
	BaseController
}

// @Title default
// @Description API --
// @Success 200 {object}
// @router / [any]
func (o *DefaultController) GetAll() {
	o.Data["json"] = Response{0, "success.", "API 1.0"}
	o.ServeJSON()
}
