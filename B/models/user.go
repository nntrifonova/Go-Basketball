package models

import (
	"github.com/beego/beego/v2/adapter/orm"
	"time"
)

var (
	UserList map[string]*User
)

type User struct {
	Id       int64     `json:"id" orm:"column(id);pk;auto;unique"`
	Phone    string    `json:"phone" orm:"column(phone);unique;size(11)"`
	Nickname string    `json:"nickname" orm:"column(nickname);unique;size(40);"`
	Password string    `json:"-" orm:"column(password);size(40)"`
	Created  time.Time `json:"create_at" orm:"column(create_at);auto_now_add;type(datetime)"`
	Updated  time.Time `json:"-" orm:"column(update_at);auto_now;type(datetime)"`
	Profile  Profile
}

type Profile struct {
	Gender string `json:"-" orm:"column(gender);size(40)"`
	Age    int    `json:"-" orm:"column(age);size(40)"`
	Email  string `json:"-" orm:"column(email);size(40)"`
}

func (u *User) TableName() string {
	return TableName("user")
}

func init() {
	orm.RegisterModel(new(User))
}

func Users() orm.QuerySeter {
	return orm.NewOrm().QueryTable(new(User))
}

func CreateUser(user User) User {
	o := orm.NewOrm()
	o.Insert(&user)
	return user
}

func CheckUserAuth(nickname string, password string) (User, bool) {
	o := orm.NewOrm()
	user := User{
		Nickname: nickname,
		Password: password,
	}
	err := o.Read(&user, "Nickname", "Password")
	if err != nil {
		return user, false
	}
	return user, true
}

// User database CRUD methods include Insert, Read, Update and Delete
func (usr *User) Insert() error {
	if _, err := orm.NewOrm().Insert(usr); err != nil {
		return err
	}
	return nil
}

func (usr *User) Read(fields ...string) error {
	if err := orm.NewOrm().Read(usr, fields...); err != nil {
		return err
	}
	return nil
}

func (usr *User) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(usr, fields...); err != nil {
		return err
	}
	return nil
}

func (usr *User) Delete() error {
	if _, err := orm.NewOrm().Delete(usr); err != nil {
		return err
	}
	return nil
}

func Login(username, password string) bool {
	for _, u := range UserList {
		if u.Nickname == username && u.Password == password {
			return true
		}
	}
	return false
}

func DeleteUser(uid string) {
	delete(UserList, uid)
}
