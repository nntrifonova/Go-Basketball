package controllers

import (
	"Basketball/models"
	services "Basketball/service"
	"encoding/json"
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
)

type UserController struct {
	beego.Controller
}

// Define a struct to return when user is authorized
type AuthorizedResponse struct {
	Message string       `json:"message"`
	User    *models.User `json:"user"`
	Token   string       `json:"token"`
}

// Define a struct to return when there is an error
type ErrorResponse struct {
	Message string `json:"message"`
}

func (cont *UserController) GetRegUser() {
	cont.TplName = "index.tpl"
}

// @Title Register User
// @Description Register a new User in system
// @Param	user	body	{InputUser}	true	"User initial data"
// @router /register [post]
func (cont *UserController) RegisterUser() {
	// Parse input data
	var iu models.InputUser
	_ = json.Unmarshal(cont.Ctx.Input.RequestBody, &iu)

	fmt.Println("Input User: ", iu)

	// Create User
	id, err := models.CreateNew(iu.Email, iu.Password, iu.Name)
	if err != nil {
		// Return response
		errResponse := ErrorResponse{
			Message: err.Error(),
		}
		cont.Data["json"] = errResponse
		cont.ServeJSON()
		cont.StopRun()
	}

	// Get User from database
	user, err := models.FindById(id)
	if err != nil {
		errResponse := ErrorResponse{
			Message: err.Error(),
		}
		cont.Data["json"] = errResponse
		cont.ServeJSON()
		cont.StopRun()
	}

	// Get Token
	token, err := services.MakeToken()
	if err != nil {
		errResponse := ErrorResponse{
			Message: err.Error(),
		}
		cont.Data["json"] = errResponse
		cont.ServeJSON()
		cont.StopRun()
	}

	// Return result
	successRes := AuthorizedResponse{
		Message: "User created successfully",
		User:    user,
		Token:   token,
	}
	cont.Data["json"] = successRes
	cont.ServeJSON()
}

// @Title Login User
// @Description Log in an existing User with credentials
// @Param	credentials	body	{BasicCredentials}	true	"User credentials"
// @router /login [post]
func (cont *UserController) LoginUser() {
	// Parse input data
	var credentials models.BasicCredentials
	_ = json.Unmarshal(cont.Ctx.Input.RequestBody, &credentials)

	// Try to login
	user, err := models.Login(credentials.Email, credentials.Password)
	if err != nil {
		errResponse := ErrorResponse{
			Message: err.Error(),
		}
		cont.Data["json"] = errResponse
		cont.ServeJSON()
		cont.StopRun()
	}

	// Get Token
	token, err := services.MakeToken()
	if err != nil {
		errResponse := ErrorResponse{
			Message: err.Error(),
		}
		cont.Data["json"] = errResponse
		cont.ServeJSON()
		cont.StopRun()
	}

	// Return result
	successRes := AuthorizedResponse{
		Message: "User logged in successfully",
		User:    user,
		Token:   token,
	}
	cont.Data["json"] = successRes
	cont.ServeJSON()
}
