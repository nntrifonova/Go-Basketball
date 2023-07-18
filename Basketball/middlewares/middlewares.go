package middlewares

import (
	"Basketball/conf"
	"Basketball/models"
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web/context"

	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strconv"
	"strings"
)

type APIResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func Jwt(ctx *context.Context) {
	ctx.Output.Header("Content-Type", "application/json")
	uri := ctx.Input.URI()
	s := strings.Split(uri, "/")

	if uri == "/api" {
		return
	}

	if ctx.Input.Header("Authorization") == "" {
		ctx.Output.SetStatus(http.StatusForbidden)
		resBody, err := json.Marshal(APIResponse{http.StatusForbidden, "notAllowed"})

		if err = ctx.Output.Body(resBody); err != nil {
			logs.Error(err)
		}
	}

	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	tokenString := ctx.Input.Header("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(conf.GetEnvConst("ACCESS_SECRET")), nil
	})

	if err != nil {
		ctx.Output.SetStatus(http.StatusForbidden)
		var responseBody APIResponse = APIResponse{http.StatusForbidden, err.Error()}
		resBytes, err := json.Marshal(responseBody)

		if err = ctx.Output.Body(resBytes); err != nil {
			logs.Error(err)
		}
	}
	var userId float64

	//s[2] is user id in url
	if userId, err = strconv.ParseFloat(s[2], 32); err != nil {
		logs.Error(err)
	}

	//TokenValidation and user from url should be the same as in token
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid && claims != nil && userId == token.Claims.(jwt.MapClaims)["user_id"] {
		return

	} else {
		ctx.Output.SetStatus(http.StatusForbidden)
		resBody, err := json.Marshal(APIResponse{http.StatusForbidden, ctx.Input.Header("Authorization")})
		ctx.Output.Body(resBody)

		if err != nil {
			logs.Error(err)
		}
	}
}

func CheckEmailIsValid(ctx *context.Context) {
	ctx.Output.Header("Content-Type", "application/json")
	uri := ctx.Input.URI()
	s := strings.Split(uri, "/")
	var err error

	if uri == "/api" || s[3] == "me" || s[3] == "validate_email" {
		return
	}
	var user *models.User
	var userID int

	//s[2] is user id in url
	if userID, err = strconv.Atoi(s[2]); err != nil {
		logs.Error(err)
	}

	if user, err = models.GetUsersById(int64(userID)); err != nil {
		logs.Error(err)
	}

	if user.EmailConfirmed == true {
		return
	} else {
		ctx.Output.SetStatus(http.StatusForbidden)
		resBody, err := json.Marshal(APIResponse{http.StatusForbidden, "email isn't confirmed"})

		if err = ctx.Output.Body(resBody); err != nil {
			logs.Error(err)
		}
	}
}
