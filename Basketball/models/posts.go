package models

import (
	"github.com/beego/beego/v2/client/orm"
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
	Status   int   `orm:"size(100)" json:"status"`
	AuthorId int64 `json:"author_id"`
}

func (p *Post) TableName() string {
	// db table name
	return "posts"
}

func init() {
	orm.RegisterModel(new(Post))
}

func GetArticle(id int64) (Post, error) {
	o := orm.NewOrm()
	art := Post{Id: id}
	err := o.Read(&art)
	return art, err
}

func AddArticle(updPost Post) (int64, error) {
	o := orm.NewOrm()
	post := new(Post)

	post.Title = updPost.Title
	post.Content = updPost.Content
	post.AuthorId = updPost.AuthorId
	post.CreatedAt = time.Now()
	//	post.Viewnum = 1
	post.Status = updPost.Status

	id, err := o.Insert(post)
	return id, err
}

func UpdateArticle(id int64, updPost Post) error {
	o := orm.NewOrm()
	post := Post{Id: id}
	post.Title = updPost.Title
	post.Content = updPost.Content
	post.AuthorId = updPost.AuthorId
	post.CreatedAt = time.Now()
	post.Status = updPost.Status
	_, err := o.Update(&post)
	return err
}

func ListArticle(condArr map[string]string, page int, offset int) (num int64, err error, post []Post) {
	o := orm.NewOrm()
	qs := o.QueryTable("posts")
	cond := orm.NewCondition()
	if condArr["title"] != "" {
		cond = cond.And("title__icontains", condArr["title"])
	}
	if condArr["status"] != "" {
		cond = cond.And("status", condArr["status"])
	}
	qs = qs.SetCond(cond)
	if page < 1 {
		page = 1
	}
	if offset < 1 {
		offset = 10
	}
	start := (page - 1) * offset
	var articles []Post
	num, err1 := qs.Limit(offset, start).All(&articles)
	return num, err1, articles
}

func CountArticle(condArr map[string]string) int64 {
	o := orm.NewOrm()
	qs := o.QueryTable("article")
	cond := orm.NewCondition()
	if condArr["title"] != "" {
		cond = cond.And("title__icontains", condArr["title"])
	}
	if condArr["status"] != "" {
		cond = cond.And("status", condArr["status"])
	}
	num, _ := qs.SetCond(cond).Count()
	return num
}
