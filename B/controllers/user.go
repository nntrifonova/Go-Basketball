package controllers

import (
	"B/models"
	"B/utils"
	"fmt"
	"github.com/beego/beego/v2/adapter/validation"
	"strings"
	"time"
)

var (
	ErrPhoneIsRegis     = ErrResponse{422001, "11"}
	ErrNicknameIsRegis  = ErrResponse{422002, "22"}
	ErrNicknameOrPasswd = ErrResponse{422003, "33"}
)

// Operations about Users
type UserController struct {
	BaseController
}
type LoginToken struct {
	User  models.User `json:"user"`
	Token string      `json:"token"`
}

func (this *UserController) Registered() {
	phone := this.GetString("phone")
	nickname := this.GetString("nickname")
	password := this.GetString("password")

	valid := validation.Validation{}
	//表单验证
	valid.Required(phone, "phone").Message("1")
	valid.Required(nickname, "nickname").Message("2")
	valid.Required(password, "password").Message("3")
	valid.Mobile(phone, "phone").Message("4")
	valid.MinSize(nickname, 2, "nickname").Message("min 2")
	valid.MaxSize(nickname, 40, "nickname").Message("max 40")
	valid.Length(password, 32, "password").Message("5")

	if valid.HasErrors() {

		for _, err := range valid.Errors {
			this.Ctx.ResponseWriter.WriteHeader(403)
			this.Data["json"] = ErrResponse{403001, map[string]string{err.Key: err.Message}}
			this.ServeJSON()
			return
		}
	}

	user := models.User{
		Phone:    phone,
		Nickname: nickname,
		Password: password,
	}
	this.Data["json"] = Response{0, "success.", models.CreateUser(user)}
	this.ServeJSON()

}

// @Title logout
// @Description Logs out current logged in user session
// @Success 200 {string} logout success
// @router /logout [get]
func (u *UserController) Logout() {
	u.Data["json"] = "logout success"
	u.ServeJSON()
}

// @Title 登录
// @Description 账号登录
// @Success 200 {object}
// @Failure 404 no enough input
// @Failure 401 No Admin
// @router /login [post]
func (this *UserController) Login() {
	nickname := this.GetString("nickname")
	password := this.GetString("password")

	user, ok := models.CheckUserAuth(nickname, password)
	if !ok {
		this.Data["json"] = ErrNicknameOrPasswd
		this.ServeJSON()
		return
	}

	et := utils.EasyToken{
		Username: user.Nickname,
		Uid:      user.Id,
		Expires:  time.Now().Unix() + 3600,
	}

	token, err := et.GetToken()
	if token == "" || err != nil {
		this.Data["json"] = ErrResponse{-0, err}
	} else {
		this.Data["json"] = Response{0, "success.", LoginToken{user, token}}
	}

	this.ServeJSON()
}

// @Title 认证测试
// @Description 测试错误码
// @Success 200 {object}
// @Failure 401 unauthorized
// @router /auth [get]
func (this *UserController) Auth() {
	et := utils.EasyToken{}
	authtoken := strings.TrimSpace(this.Ctx.Request.Header.Get("Authorization"))
	valido, err := et.ValidateToken(authtoken)
	if !valido {
		this.Ctx.ResponseWriter.WriteHeader(401)
		this.Data["json"] = ErrResponse{-1, fmt.Sprintf("%s", err)}
		this.ServeJSON()
		return
	}

	this.Data["json"] = Response{0, "success.", "is login"}
	this.ServeJSON()
}
