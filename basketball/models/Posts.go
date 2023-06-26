package models

import "github.com/GoAdminGroup/go-admin/template/types/display"

type Post struct {
	postId   int
	title    string
	content  string
	image    string
	datePost display.Date
	userId   int
}
