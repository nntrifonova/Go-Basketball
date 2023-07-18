package utiles

import (
	"Basketball/conf"
	"Basketball/models"
	_ "Basketball/models"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"regexp"
	"strconv"
)

var SecretKey = conf.GetEnvConst("ACCESS_SECRET")

const CodeLength = 40

// CanRegistered checks if e-mail is available.
func CanRegisteredOrChanged(email string) (bool, error) {
	cond := orm.NewCondition()
	cond = cond.Or("Email", email)
	var maps []orm.Params
	o := orm.NewOrm()

	n, err := o.QueryTable("users").SetCond(cond).Values(&maps, "Email")

	if err != nil {
		return false, err
	}
	emailCheck := true

	if n > 0 {
		for _, m := range maps {
			if emailCheck && orm.ToStr(m["Email"]) == email {
				emailCheck = false
			}
		}
	}
	return emailCheck, nil
}

func ValidateEmail(email string) bool {
	Re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,6}$`)
	return Re.MatchString(email)
}

// get a code for email confirmation
func GetEmailConfirmationCode(user *models.User) string {
	data := strconv.Itoa(int(user.Id)) + user.Email
	hexData := strconv.Itoa(int(user.Id)) + "|" + conf.GetEnvConst("ACCESS_SECRET")
	code := CreateEmailConfirmationCode(data)
	// add tail hex user id and secret key
	code += hex.EncodeToString([]byte(hexData))
	return code
}

// create code, it's format: 40 sha1 encoded string
func CreateEmailConfirmationCode(data string) string {
	// create sha1 encode string
	sh := sha1.New()
	sh.Write([]byte(data + SecretKey))
	encoded := hex.EncodeToString(sh.Sum(nil))
	code := fmt.Sprintf("%s", encoded)
	return code
}

// verify code
func VerifyEmailConfirmationCode(data string, code string) bool {

	if len(code) <= 18 {
		return false
	}
	// right active code
	retCode := CreateEmailConfirmationCode(data)

	if retCode == code {
		return true
	}
	return false
}
