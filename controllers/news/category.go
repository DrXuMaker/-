package news

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"math"
	models"shizhan/models/news"
	"shizhan/util"
)

type CarController struct{
	beego.Controller
}

func (n * CarController)List(){

	o := orm.NewOrm()
	categories := []models.Category{}
	o.QueryTable("sys_category").All(&categories)
	//每页显示的条数
	PagePerNum := 7
	Headpage := 1
	CurrentPage, err := n.GetInt("page")


	var Count int64 = 0
	OffsetNum := PagePerNum*(CurrentPage - 1)

	kw :=n.GetString("kw")
	if kw != ""{
		Count,_ = o.QueryTable("sys_category").Filter("is_delete__exact",0).Filter("name__contains",kw).Count()
		o.QueryTable("sys_category").Filter("is_delete__exact",0).Filter("name__contains",kw).Limit(PagePerNum).Offset(OffsetNum).All(&categories)
	}else{
		Count,_ = o.QueryTable("sys_category").Filter("is_delete",0).Count()
		o.QueryTable("sys_category").Filter("is_delete__exact",0).Limit(PagePerNum).Offset(OffsetNum).All(&categories)
	}
	if err != nil {
		CurrentPage = 1
	}
	/*
		分页逻辑
		prepage		offset		limit	   				limitnumber *(currentPage -1)
	*/

	//分页传统写法
	PrePage := 1
	if CurrentPage > 1{
		PrePage = CurrentPage - 1
	}else {
		PrePage  = CurrentPage
	}

	//总页数
	TotalPage := int(math.Ceil(float64(Count) / float64(PagePerNum)))

	NextPage := 1
	if TotalPage > CurrentPage {
		NextPage = CurrentPage + 1
	}else if TotalPage < CurrentPage {
		NextPage = CurrentPage
	}

	//分页算法
	page_map := util.Paginator(CurrentPage,PagePerNum,Count)
	ret :=fmt.Sprintf("分页信息，当前页：%d,总页数：%d，查询条件：%s",CurrentPage,TotalPage,kw)
	logs.Info(ret)

	n.Data["categories"] = categories
	n.Data["NextPage"] = NextPage
	n.Data["PrePage"] = PrePage
	n.Data["TotalPage"] = TotalPage
	n.Data["CurrentPage"] = CurrentPage
	n.Data["Count"] = Count
	n.Data["pagemap"] = page_map
	n.Data["headpage"] = Headpage
	n.Data["kw"] = kw
	n.TplName = "news/cartegory_list.html"
}

func (n * CarController) ToAddCar() {

	id,_ := n.GetInt("id")
	o := orm.NewOrm()
	category := models.Category{}
	o.QueryTable("sys_category").Filter("Id",id).One(&category)
	n.Data["category"] = category
	n.TplName = "news/cartegory_add.html"
}
func (n * CarController) DoAddCar() {

	name := n.GetString("name")
	description := n.GetString("description")
	is_active := util.StrToInt(n.GetString("is_active"))
	//添加到数据库中
	o := orm.NewOrm()
	category := models.Category{
		Name: name,
		Desc: description,
		IsActive: is_active,
	}
	_,err := o.Insert(&category)
	message := map[string]interface{}{}
	if err != nil{
		message["code"] = 10001
		message["msg"] = "添加失败"
		n.Data["json"] = message
	}else{
		message["code"] = 200
		message["msg"] = "添加成功"
		n.Data["json"] = message
	}
	n.ServeJSON()
}

func (n * CarController) Delete(){

	Id := util.StrToInt(n.GetString("id"))
	o := orm.NewOrm()
	qs := o.QueryTable("sys_category").Filter("id",Id)
	_,err :=qs.Update(orm.Params{
		"is_delete": 1,
	})
	message := map[string]interface{}{}
	if err != nil {
		message["code"] = 10001
		message["msg"] = "删除失败"
		n.Data["json"] = message
	}else{
		message["code"] = 200
		message["msg"] = "删除成功"
		n.Data["json"] = message
	}
	n.ServeJSON()
}

