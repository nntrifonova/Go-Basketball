package models

import (
	"fmt"
	"github.com/beego/beego/v2/adapter/orm"
	"time"
)

type User struct {
	Id             int64     `json:"id"`
	Email          string    `json:"email" orm:"unique;index;size(191)"`
	AccessToken    string    `orm:"size(128)" json:"access_token"`
	Name           string    `orm:"size(128)" json:"name" `
	Role           string    `orm:"size(128)" json:"role"`
	Password       string    `orm:"size(128)" json:"password"`
	Created        time.Time `json:"created_on" orm:"auto_now_add;type(datetime)"`
	Updated        time.Time `json:"updated_on" orm:"auto_now;type(datetime)"`
	EmailConfirmed bool      `orm:"size(128)" json:"email_confirmed"`
	RecentLogin    time.Time `orm:"column(recent_login);type(timestamp with time zone);null" json:"recent_login"`
}

func init() {
	orm.RegisterModel(new(User))
}

// AddUsers insert a new User into database and returns
// last inserted Id on success.
func AddUsers(m *User) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetUsersById retrieves User by Id. Returns error if
// Id doesn't exist
func GetUsersById(id int64) (v *User, err error) {
	o := orm.NewOrm()
	v = &User{Id: id}

	if err = o.QueryTable(new(User)).Filter("Id", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetUsersByEmail retrieves Customer by Email. Returns error if
// Id doesn't exist
func GetUsersByEmail(email string) (v *User, err error) {
	o := orm.NewOrm()
	v = &User{Email: email}

	if err = o.QueryTable(new(User)).Filter("Email", email).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetUsersByEmail retrieves Customer by Email. Returns error if
// Id doesn't exist

// UpdateUsers updates Users by Id and returns error if
// the record to be updated doesn't exist
func UpdateUsersById(m *User) (err error) {
	o := orm.NewOrm()
	v := User{Id: m.Id}

	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64

		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteUser deletes User by Id and returns error if
// the record to be deleted doesn't exist
func DeleteUsers(id int64) (err error) {
	o := orm.NewOrm()
	v := User{Id: id}

	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64

		if num, err = o.Delete(&User{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
