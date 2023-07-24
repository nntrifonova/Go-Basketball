package controllers

import (
	"encoding/json"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/pagination"
	"github.com/labstack/gommon/log"
	"net/http"
	"strconv"

	"Basketball/models"
	"fmt"
)

// 添加blog
type PostsController struct {
	BaseController
}

// URLMapping ...
func (p *PostsController) URLMapping() {
	p.Mapping("Posts", p.CreatePost)
	p.Mapping("Posts", p.PostEdit)
	p.Mapping("Get", p.GetPosts)
	p.Mapping("Get", p.GetEdit)
	p.Mapping("Get", p.GetShow)
	p.Mapping("Get", p.GetList)
}

type PostCredentials struct {
	title   string `json:"title"`
	content string `json:"content"`
	author  int64  `json:"author"`
}

func (p *PostsController) GetPosts() {
	//if !p.isLogin {
	//	p.Redirect("/auth/login", 302)
	//	return
	//}

	var post models.Post
	post.Status = 1
	p.Data["post"] = post

}

func (p *PostsController) CreatePost() {
	var err error
	//if !p.isLogin {
	//	p.Data["json"] = map[string]interface{}{"code": 0, "message": "请先登录"}
	//	p.ServeJSON()
	//	return
	//}

	var credentials PostCredentials
	s := string(p.Ctx.Input.RequestBody)

	if err = json.Unmarshal([]byte(s), &credentials); err != nil {
		fmt.Print("errorr1")
		log.Error(err)
		p.Resp(http.StatusBadRequest, nil, err)

		if "" == credentials.title {
			p.Data["json"] = map[string]interface{}{"code": 0, "message": "no title"}
			p.ServeJSON()
			return
		}
		if "" == credentials.content {
			p.Data["json"] = map[string]interface{}{"code": 0, "message": "no content"}
			p.ServeJSON()
			return
		}

		var post models.Post
		post.Title = credentials.title
		post.AuthorId = credentials.author

		id, err := models.AddArticle(post)
		if err == nil {
			p.Data["json"] = map[string]interface{}{"code": 1, "message": "add post", "id": id}
		} else {
			p.Data["json"] = map[string]interface{}{"code": 0, "message": "can't add post"}
		}
		p.Ctx.Output.SetStatus(http.StatusCreated)
		p.ServeJSON()

	}
}

func (p *PostsController) GetEdit() {
	//if !p.isLogin {
	//	p.Data["json"] = map[string]interface{}{"code": 0, "message": "ooo"}
	//	p.ServeJSON()
	//	return
	//}
	idstr := p.Ctx.Input.Param(":id")
	id_, err := strconv.Atoi(idstr)
	var id = int64(id_)
	art, err := models.GetArticle(id)
	if err != nil {
		p.Redirect("/404.html", 302)
	}
	//p.Data["json"] = map[string]interface{}{"code": 0, "message": err}
	//p.ServeJson()
	p.Data["art"] = art

}

func (p *PostsController) PostEdit() {
	id_, err := p.GetInt("id")
	title := p.GetString("title")
	content := p.GetString("content")
	fmt.Println("content", content)
	author, e := p.GetInt("author")
	status, _ := p.GetInt("status")

	if "" == title {
		p.Data["json"] = map[string]interface{}{"code": 0, "message": "请填写标题"}
		p.ServeJSON()
		return
	}
	if "" == content {
		p.Data["json"] = map[string]interface{}{"code": 0, "message": "请填写内容"}
		p.ServeJSON()
		return
	}
	var id = int64(id_)
	_, errAttr := models.GetArticle(id)
	if errAttr != nil {
		p.Data["json"] = map[string]interface{}{"code": 0, "message": "博客不存在"}
		p.ServeJSON()
		return
	}

	var post models.Post
	post.Title = title
	if e == nil {
		post.AuthorId = int64(author)
	} else {
		p.Data["json"] = map[string]interface{}{"code": 0, "message": "author error"}
	}
	post.Status = status

	err = models.UpdateArticle(id, post)
	if err == nil {
		p.Data["json"] = map[string]interface{}{"code": 1, "message": "博客修改成功", "id": id}
	} else {
		p.Data["json"] = map[string]interface{}{"code": 0, "message": "博客修改出错"}
	}
	p.ServeJSON()
}

func (p *PostsController) GetList() {
	page, err1 := p.GetInt("p")
	title := p.GetString("title")
	keywords := p.GetString("keywords")
	status := p.GetString("status")
	if err1 != nil {
		page = 1
	}

	offset, err2 := beego.AppConfig.Int("pageoffset")
	if err2 != nil {
		offset = 9
	}

	condArr := make(map[string]string)
	condArr["title"] = title
	condArr["keywords"] = keywords
	if !p.isLogin {
		condArr["status"] = "1"
	} else {
		condArr["status"] = status
	}
	countArticle := models.CountArticle(condArr)

	paginator := pagination.SetPaginator(p.Ctx, offset, countArticle)
	_, _, art := models.ListArticle(condArr, page, offset)

	p.Data["paginator"] = paginator
	p.Data["art"] = art

}

func (p *PostsController) GetShow() {
	idstr := p.Ctx.Input.Param(":id")
	id_, err := strconv.Atoi(idstr)
	var id = int64(id_)
	post, err := models.GetArticle(id)
	if err != nil {
		p.Redirect("/404.html", 302)
	}

	p.Data["post"] = post

	// 评论分页
	//page, err1 := p.GetInt("p")
	//if err1 != nil {
	//	page = 1
	//}
	//offset, err2 := beego.AppConfig.Int("pageoffset")
	//if err2 != nil {
	//	offset = 9
	//}
	//condCom := make(map[string]string)
	//condCom["article"] = idstr
	//if !p.isLogin {
	//	condCom["status"] = "1"
	//}

}
