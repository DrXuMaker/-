package news

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type Category struct{
	Id int `orm:"pk;auto"`
	Name string `orm:"size(64);description(分类名称)"`
	Desc string  `orm:"size(255);description(描述)"`
	IsActive int `orm:"description(1启用,0停用);default(1)"`
	IsDelete int `orm:"description(1删除,0未删除);default(0)"`
	CreateTime time.Time `orm:"description(创建时间);type(datetime);auto_now"`
	News []*News `orm:"reverse(many)"`
}

type News struct{
	Id int `orm:"pk;auto"`
	Title string `orm:"size(64);description(新闻标题)"`
	Content string  `orm:"size(255);description(新闻内容);type(text)"`
	IsActive int `orm:"description(1启用,0停用);default(1)"`
	IsDelete int `orm:"description(1删除,0未删除);default(0)"`
	CreateTime time.Time `orm:"description(创建时间);type(datetime);auto_now"`
	Category *Category `orm:"rel(fk)"`
}

func (c * Category)TableName() string{
	return "sys_category"
}

func (n * News)TableName() string{
	return "sys_news"
}

func init() {
	orm.RegisterModel(new(News),new(Category))
}