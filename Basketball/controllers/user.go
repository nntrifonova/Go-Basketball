package controllers

import (
	"Basketball/conf"
	"Basketball/models"
	"Basketball/services/mailgun"
	"Basketball/utiles"
	"encoding/json"
	"errors"
	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
)

// Create a struct to read the email from the request body
type UserEmailData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type uEmail struct {
	Email string `json:"email"`
}

// Create a struct to read password from the request body
type UserPassword struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type ValidationCode struct {
	ValidationCode string `json:"code"`
}

type CurrentUser struct {
	ID             int64  `orm:"column(id);pk;auto" json:"id"`
	Name           string `orm:"size(128)" json:"name"`
	Email          string `orm:"size(128);unique" json:"email"`
	Phone          string `orm:"size(128);unique" json:"phone"`
	PinCode        string `orm:"size(128)" json:"pin_code"`
	EmailConfirmed bool   `orm:"size(128)" json:"email_confirmed"`
}

// UsersController operations for Users
type UsersController struct {
	BaseController
}

// URLMapping ...
func (c *UsersController) URLMapping() {
	c.Mapping("Put", c.PutEmail)
	c.Mapping("Put", c.PutPassword)
	c.Mapping("Get", c.CheckEmail)
	c.Mapping("Get", c.GetCurrent)
	c.Mapping("Get", c.ValidateEmail)
}

func (c *UsersController) PutEmail() {
	var err error
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	s := string(c.Ctx.Input.RequestBody)
	var user *models.User
	var userEmail UserEmailData

	if user, err = models.GetUsersById(id); err != nil {
		log.Error(err)
		c.Resp(http.StatusInternalServerError, nil, err)
	}

	if err = json.Unmarshal([]byte(s), &userEmail); err != nil {
		log.Error(err)
		c.Resp(http.StatusBadRequest, nil, err)
	}

	// Compare the stored hashed password, with the hashed version of the password that was received
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userEmail.Password)); err != nil {
		// If the two passwords don't match, return a 401 status
		err := errors.New("wrong password, please enter the correct password")
		c.Resp(http.StatusUnauthorized, nil, err)
	}
	var canChanged, _ = utiles.CanRegisteredOrChanged(userEmail.Email)

	var uEmail uEmail

	if canChanged && user.Email != userEmail.Email {
		user.Email = userEmail.Email

		if err = models.UpdateUsersById(user); err != nil {
			log.Error(err)
			c.Resp(http.StatusInternalServerError, nil, err)
		}
		uEmail.Email = user.Email
		c.Resp(http.StatusOK, uEmail, nil)

	} else if user.Email == userEmail.Email {
		uEmail.Email = user.Email
		c.Resp(http.StatusOK, uEmail, nil)

	} else {
		err := errors.New("such email already exists")
		c.Resp(http.StatusConflict, nil, err)
	}
}

func (c *UsersController) PutPassword() {
	var err error
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	s := string(c.Ctx.Input.RequestBody)
	var user *models.User
	var password UserPassword

	if user, err = models.GetUsersById(id); err != nil {
		log.Error(err)
		c.Resp(http.StatusInternalServerError, nil, err)
	}

	if err = json.Unmarshal([]byte(s), &password); err != nil {
		log.Error(err)
		c.Resp(http.StatusBadRequest, nil, err)
	}

	// Compare the stored hashed password, with the hashed version of the password that was received
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password.OldPassword)); err != nil {
		// If the two passwords don't match, return a 401 status
		err := errors.New("wrong password, please enter the correct password")
		c.Resp(http.StatusUnauthorized, nil, err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password.NewPassword), 8)

	if err != nil {
		log.Error(err)
		c.Resp(http.StatusBadRequest, nil, err)
	}
	user.Password = string(hashedPassword)

	if err = models.UpdateUsersById(user); err != nil {
		log.Error(err)
		c.Resp(http.StatusInternalServerError, nil, err)
	}
	c.Resp(http.StatusOK, "password updated", nil)
}

func (c *UsersController) CheckEmail() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	var user *models.User
	var err error

	if user, err = models.GetUsersById(id); err != nil {
		log.Error(err)
		c.Resp(http.StatusInternalServerError, nil, err)
	}

	if user.EmailConfirmed == true {
		c.Resp(http.StatusOK, true, nil)
	}
	c.Resp(http.StatusBadRequest, false, nil)
}

func (c *UsersController) GetCurrent() {
	var err error
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	var user *models.User
	var currentUser CurrentUser

	if user, err = models.GetUsersById(id); err != nil {
		log.Error(err)
		c.Resp(http.StatusInternalServerError, nil, err)
	}
	currentUser.ID = user.Id
	currentUser.Name = user.Name
	currentUser.Email = user.Email
	currentUser.EmailConfirmed = user.EmailConfirmed
	c.Resp(http.StatusOK, currentUser, nil)
}

func (c *UsersController) ValidateEmail() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	var user *models.User
	var err error
	var emailConfirmationCode string

	if user, err = models.GetUsersById(id); err != nil {
		log.Error(err)
		c.Resp(http.StatusInternalServerError, nil, err)
	}
	emailConfirmationCode = utiles.GetEmailConfirmationCode(user)
	url := conf.GetEnvConst("APP_URL") + "/active/" + emailConfirmationCode

	// send Email to forward user email
	_, err = mailgun.SendMail(
		conf.GetEnvConst("NOTIFICATION_EMAIL"),
		user.Email,
		"Email validation code",
		url,
	)

	if err != nil {
		log.Error(err)
		c.Resp(http.StatusInternalServerError, nil, err)
	}
	c.Resp(http.StatusOK, "Email validation url is sent", nil)
}
