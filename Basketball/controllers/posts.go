package controllers

import (
	"Basketball/models"
	"encoding/json"
	beego "github.com/beego/beego/v2/server/web"
	"strconv"
)

// Operations about article
type PostController struct {
	beego.Controller
}

// @Title Create
// @Description create article
// @Param	body		body 	models.Post	true		"The post content"
// @Success 200 {string} models.Post.Id
// @Failure 403 body is empty
// @router / [post]
func (o *PostController) Post() {
	var ob models.Post
	json.Unmarshal(o.Ctx.Input.RequestBody, &ob)
	if ob.Title == "" || ob.Content == "" {
		o.Ctx.Output.SetStatus(403)
		return
	}
	postid := models.AddOne(ob)
	s := strconv.Itoa(int(postid))
	o.Data["json"] = map[string]string{"Id": s}
	o.ServeJSON()
}

// @Title Get
// @Description find article by articleid
// @Param	postId		path 	string	true		"the articleid you want to get"
// @Success 200 {post} models.Post
// @Failure 403 :postId is empty
// @router /:postId [get]
func (o *PostController) Get() {
	postId := o.Ctx.Input.Param(":postId")
	if postId != "" {
		ob, err := models.GetOne(postId)
		if err != nil {
			o.Data["json"] = err.Error()
		} else {
			o.Data["json"] = ob
		}
	}
	o.ServeJSON()
}

// @Title GetAll
// @Description get all articles
// @Success 200 {post} models.Post
// @Failure 403 : postId is empty
// @router / [get]
func (o *PostController) GetAll() {
	obs := models.GetAll()
	o.Data["json"] = obs
	o.ServeJSON()
}

// @Title Update
// @Description update the article
// @Param	postId		path 	string	true		"The articleid you want to update"
// @Param	body		body 	models.Post	true		"The body"
// @Success 200 {article} models.Post
// @Failure 403 :postId is empty
// @router /:postId [put]
func (o *PostController) Put() {
	postId := o.Ctx.Input.Param(":postId")
	var ob models.Post
	json.Unmarshal(o.Ctx.Input.RequestBody, &ob)

	err := models.Update(postId, ob.Title, ob.Content)
	if err != nil {
		o.Data["json"] = err.Error()
	} else {
		o.Data["json"] = "update success!"
	}
	o.ServeJSON()
}

// @Title Delete
// @Description delete the article
// @Param	postId		path 	string	true		"The articleId you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 postId is empty
// @router /:postId [delete]
func (o *PostController) Delete() {
	postId := o.Ctx.Input.Param(":articleId")
	models.Delete(postId)
	o.Data["json"] = "delete success!"
	o.ServeJSON()
}
