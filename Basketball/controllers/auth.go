package controllers

import (
	"Basketball/conf"
	"Basketball/models"
	"Basketball/services"
	"Basketball/services/mailgun"
	"Basketball/utiles"
	"encoding/hex"
	"encoding/json"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/gommon/log"
)

type AuthController struct {
	BaseController
}

// URLMapping ...
func (c *AuthController) URLMapping() {
	c.Mapping("Post", c.Register)
	c.Mapping("Post", c.Login)
	c.Mapping("Post", c.CheckAccessToken)
	c.Mapping("Get", c.ConfirmEmail)
	c.Mapping("Get", c.Logout)

}

// Create a struct to read the email or phone and password from the request body
type RegisterCredentials struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	PasswordRep string `json:"password-rep"`
}

// Create a struct to read the email or phone and password from the request body
type LoginCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Create a struct to read the email or phone and password from the request body
type EmailValidationCredentials struct {
	Email               string `json:"email"`
	EmailValidationCode string `json:"code"`
}

type AuthorizedResponse struct {
	Message string       `json:"message"`
	User    *models.User `json:"user"`
	Token   string       `json:"token"`
}

type AccessToken struct {
	AccessToken string `json:"access_token"`
	UserID      int64  `json:"user_id"`
}

func (c *AuthController) Register() {
	var err error
	var user models.User
	var credentials RegisterCredentials
	s := string(c.Ctx.Input.RequestBody)
	var emailConfirmationCode string

	if err = json.Unmarshal([]byte(s), &credentials); err != nil {
		log.Error(err)
		c.Response(http.StatusBadRequest, nil, err)
	}
	// user credentials validation
	var canRegisteredEmail, _ = utiles.CanRegisteredOrChanged(credentials.Email)

	if !utiles.ValidateEmail(credentials.Email) {
		var err = errors.New("email address is invalid")
		c.Response(http.StatusInternalServerError, nil, err)
	}
	user.Role = "user"

	if canRegisteredEmail {
		user.Email = strings.ToLower(credentials.Email)
		user.Name = credentials.Name
		user.EmailConfirmed = false
		var accessToken string
		var userID int64

		if userID, err = models.AddUsers(&user); err != nil {
			log.Error(err)
			c.Response(http.StatusInternalServerError, nil, err)
		}

		if accessToken, err = CreateAccessToken(int(user.Id), user.Role); err != nil {
			log.Error(err)
			c.Response(http.StatusInternalServerError, nil, err)
		}
		user.AccessToken = accessToken

		if err = models.UpdateUsersById(&user); err != nil {
			log.Error(err)
			c.Response(http.StatusInternalServerError, nil, err)
		}
		emailConfirmationCode = utiles.GetEmailConfirmationCode(&user, nil)
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
			c.Response(http.StatusInternalServerError, nil, err)
		}
		var token AccessToken
		token.AccessToken = accessToken
		token.UserID = userID
		c.Response(http.StatusCreated, token, nil)

	} else {

	}
}

// Create the SignIn handler
func (c *AuthController) Login() {
	var credentials LoginCredentials
	var userByEmail *models.User
	var user *models.User
	var err error
	// Get the JSON body and decode into credentials
	s := string(c.Ctx.Input.RequestBody)

	if err = json.Unmarshal([]byte(s), &credentials); err != nil {
		log.Error(err)
		c.Response(http.StatusBadRequest, nil, err)
	}

	// Get the existing entry present in the database for the given email
	if userByEmail, err = models.GetUsersByEmail(credentials.Email); err != nil {
		log.Error(err)
		log.Info("no email provided")
	}

	if userByEmail != nil {

		log.Info("Logging by email")
		log.Info("email: ", credentials.Email)
		user = userByEmail
	}

	if user == nil {
		err := errors.New("no user found, please check your login data")
		c.Response(http.StatusBadRequest, nil, err)

	} else {
		var accessToken string

		if accessToken, err = CreateAccessToken(int(user.Id), user.Role); err != nil {
			log.Error(err)
			c.Response(http.StatusInternalServerError, nil, err)
		}

		// We create another instance of `Credentials` to store the credentials we get from the database
		storedCredentials := &LoginCredentials{}

		// Compare the stored hashed password, with the hashed version of the password that was received
		if err = bcrypt.CompareHashAndPassword([]byte(storedCredentials.Password), []byte(credentials.Password)); err != nil {
			// If the two passwords don't match, return a 401 status
			err := errors.New("wrong password, please enter the correct password")
			c.Response(http.StatusUnauthorized, nil, err)
		}
		// If we reach this point, that means the users password was correct, and that they are authorized
		// The default 200 status is sent
		var token AccessToken
		token.AccessToken = accessToken
		token.UserID = user.Id
		user.AccessToken = accessToken
		user.RecentLogin = time.Now()

		if err = models.UpdateUsersById(user); err != nil {
			log.Error(err)
			c.Response(http.StatusInternalServerError, nil, err)
		}
		c.Response(http.StatusOK, token, nil)
	}
}

func (c *AuthController) CheckAccessToken() {
	var err error
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	user := models.User{Id: id}
	s := string(c.Ctx.Input.RequestBody)
	var registeredUser *models.User

	if registeredUser, err = models.GetUsersById(id); err != nil {
		log.Error(err)
		c.Response(http.StatusInternalServerError, nil, err)
	}

	if err = json.Unmarshal([]byte(s), &user); err != nil {
		log.Error(err)
		c.Response(http.StatusBadRequest, nil, err)
	}
	t2 := time.Now()
	t1 := user.RecentLogin
	diff := t2.Sub(t1)

	if registeredUser.AccessToken == user.AccessToken && diff.Hours() < 12 {
		c.Response(http.StatusOK, true, nil)
	}
	c.Response(http.StatusBadRequest, false, nil)
}

func (c *AuthController) ConfirmEmail() {
	code := c.GetString(":code")
	var user *models.User
	hexStr := code[utiles.CodeLength:]
	b, err := hex.DecodeString(hexStr)

	if err != nil {
		log.Error(err)
	}
	s := strings.Split(string(b), "|")
	ID, _ := strconv.ParseInt(s[0], 10, 64)

	if user, err = models.GetUsersById(ID); err != nil {
		log.Error(err)
	}
	data := strconv.Itoa(int(user.Id)) + user.Email
	prefix := code[:utiles.CodeLength]

	if utiles.VerifyEmailConfirmationCode(data, prefix) {
		user.EmailConfirmed = true

		if err = models.UpdateUsersById(user); err != nil {
			log.Error(err)
			c.Response(http.StatusInternalServerError, nil, err)
		}
		c.Response(http.StatusOK, "Email confirmed", nil)

	} else {
		err := errors.New("email validation code is wrong")
		c.Response(http.StatusBadRequest, nil, err)
	}
}

func CreateAccessToken(userID int, role string) (string, error) {
	var token *services.TokenDetails
	var err error

	// Get Token
	if token, err = services.CreateToken(int64(userID), role); err != nil {
		log.Error(err)
		return "", err
	}
	return token.AccessToken, nil
}

func (c *AuthController) Logout() {

}