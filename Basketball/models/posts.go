package models

import (
	"errors"
	"strconv"
	"time"
)

var (
	Posts map[string]*Post
)

type Post struct {
	PostId      string    `json:"id"`
	Title       string    `orm:"size(256)" json:"title"`
	PublishTime time.Time `orm:"column(published_at);auto_now_add;type(timestamp with time zone);null" json:"published_at"`
	Content     string    `orm:"size(1024)" json:"content"`
	Image       string    `orm:"size(128)" json:"image"`
}

func AddOne(post Post) (ArticleId string) {
	post.PostId = "astaxie" + strconv.FormatInt(time.Now().UnixNano(), 10)
	post.PublishTime = time.Now().Local()
	Posts[post.PostId] = &post
	return post.PostId
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
	return errors.New("PostId Not Exist")
}

func Delete(PostId string) {
	delete(Posts, PostId)
}
