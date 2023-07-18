package models

import (
	"errors"
	"github.com/beego/beego/v2/adapter/orm"
	"strconv"
	"time"
)

var (
	Posts map[string]*Post
)

type Post struct {
	Id        int64     `json:"post_id"`
	Title     string    `orm:"size(256)" json:"title"`
	Content   string    `orm:"size(1024)" json:"content"`
	CreatedAt time.Time `orm:"column(created_at);auto_now_add;type(timestamp with time zone);null" json:"created_at"`
	//Image       string    `orm:"size(128)" json:"image"`
}

func (p *Post) TableName() string {
	// db table name
	return "posts"
}

func init() {
	orm.RegisterModel(new(Post))
}

func AddOne(post Post) (PostId int64) {
	post.Id = time.Now().UnixNano()
	post.CreatedAt = time.Now().Local()
	s := strconv.Itoa(int(post.Id))
	Posts[s] = &post
	return post.Id
}

func GetOne(PostId string) (post *Post, err error) {
	if v, ok := Posts[PostId]; ok {
		return v, nil
	}
	return nil, errors.New("ArticleId Not Exist")
}

func GetAll() []Post {
	a := make([]Post, 0, len(Posts))
	for _, v := range Posts {
		a = append(a, *v)
	}
	return a
}

func Update(PostId string, title, content string) (err error) {
	if v, ok := Posts[PostId]; ok {
		v.Title = title
		v.Content = content
		return nil
	}
	return errors.New("Id Not Exist")
}

func Delete(PostId string) {
	delete(Posts, PostId)
}
