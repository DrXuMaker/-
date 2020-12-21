package news

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"math"
	models "shizhan/models/news"
	"shizhan/util"
	"strconv"
	"strings"
	"time"
)

type NewController struct {
	beego.Controller
}

func (n * NewController) List() {


	o := orm.NewOrm()
	news := []models.News{}
	o.QueryTable("sys_news").All(&news)
	//每页显示的条数
	PagePerNum := 7
	Headpage := 1
	CurrentPage, err := n.GetInt("page")


	var Count int64 = 0
	OffsetNum := PagePerNum*(CurrentPage - 1)

	kw :=n.GetString("kw")
	if kw != ""{
		Count,_ = o.QueryTable("sys_news").Filter("is_delete__exact",0).Filter("title__contains",kw).Count()
		o.QueryTable("sys_news").Filter("is_delete__exact",0).Filter("title__contains",kw).Limit(PagePerNum).Offset(OffsetNum).RelatedSel().All(&news)
	}else{
		Count,_ = o.QueryTable("sys_news").Filter("is_delete",0).Count()
		o.QueryTable("sys_news").Filter("is_delete__exact",0).Limit(PagePerNum).Offset(OffsetNum).RelatedSel().All(&news)
	}
	if err != nil {
		CurrentPage = 1
	}

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

	n.Data["news"] = news
	n.Data["NextPage"] = NextPage
	n.Data["PrePage"] = PrePage
	n.Data["TotalPage"] = TotalPage
	n.Data["CurrentPage"] = CurrentPage
	n.Data["Count"] = Count
	n.Data["pagemap"] = page_map
	n.Data["headpage"] = Headpage
	n.Data["kw"] = kw


	n.TplName = "news/news_list.html"
}

func (n * NewController) ToAddNews() {

	o := orm.NewOrm()
	categories := []models.Category{}
	o.QueryTable("sys_category").All(&categories)
	n.Data["categories"] = categories
	n.TplName = "news/news_add.html"
}
func (n * NewController) DoAddNews() {

	title := n.GetString("title")
	newsContent := n.GetString("content")
	isActive := util.StrToInt(n.GetString("is_active"))
	category_id := util.StrToInt(n.GetString("category_id"))
	//添加到数据库中
	o := orm.NewOrm()
	category := models.Category{Id: category_id}
	news := models.News{
		Title:    title,
		Content:  newsContent,
		IsActive: isActive,
		Category: &category,
	}
	_,err := o.Insert(&news)
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

func (n * NewController) UpLoadImg(){

	f,h,err := n.GetFile("file")

	message_map := map[string]interface{}{}

	defer func() {
		f.Close()
	}()

	fileName := h.Filename

	timeUnixInt := time.Now().Unix()
	timeUnitStr := strconv.FormatInt(timeUnixInt,10)

	filePath := "upload/news_img/"+ timeUnitStr + "-" + fileName

	imgLink := "http://localhost:8080/" + filePath

	if err != nil {
		message_map["code"] = 10001
		message_map["msg"] = "文件上传失败"
		message_map["link"] = imgLink

	}

	n.SaveToFile("file", filePath)

	message_map["code"] = 200
	message_map["msg"] = "文件上传成功"
	message_map["link"] = imgLink

	n.Data["json"] = message_map
	n.ServeJSON()
}

func (n *NewController) ToEdit()  {

	nId,_ := n.GetInt("id")
	o := orm.NewOrm()
	qs := o.QueryTable("sys_news")

	newsData := models.News{}
	qs.Filter("id", nId).RelatedSel().One(&newsData)
	categories := []models.Category{}

	o.QueryTable("sys_category").Exclude("id", newsData.Category.Id).All(&categories)

	n.Data["news_data"] = newsData
	n.Data["categories"] = categories
	n.TplName = "news/news_edit.html"

}
func (n *NewController) DoEdit()  {

	newsId := util.StrToInt(n.GetString("news_id"))
	content := n.GetString("content")
	title := n.GetString("title")
	categoryId,_ := n.GetInt("category_id")
	isActive,_ := n.GetInt("is_active")

	o := orm.NewOrm()
	qs := o.QueryTable("sys_news")
	_,err := qs.Filter("id", newsId).Update(orm.Params{
		"title":       title,
		"content":     content,
		"category_id": categoryId,
		"is_active":   isActive,
	})
	messageMap := map[string]interface{}{}
	if err != nil {
		messageMap["code"] = 10001
		messageMap["msg"] = "更新失败"
	}
	messageMap["code"] = 200
	messageMap["msg"] = "更新成功"

	n.Data["json"] = messageMap
	n.ServeJSON()
}

func (n * NewController) IsActive(){

	is_active,_ := n.GetInt("is_active")
	id:= util.StrToInt(n.GetString("id"))
	o := orm.NewOrm()
	qs :=o.QueryTable("sys_news").Filter("id__exact",id)
	message := map[string]interface{}{}
	if is_active == 1{
		qs.Update(orm.Params{
			"is_active":0,
		})
		message["msg"] = "停用成功"
		n.Data["json"] = message
	}else if is_active == 0{
		qs.Update(orm.Params{
			"is_active":1,
		})
		message["msg"] = "启用成功"
		n.Data["json"] = message
	}
	n.ServeJSON()

}

func (n * NewController) Delete(){

	Id := util.StrToInt(n.GetString("id"))
	o := orm.NewOrm()
	qs := o.QueryTable("sys_news").Filter("id",Id)
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
func (n * NewController) MuliDel(){

	ids := n.GetString("ids")
	new_ids := ids[1:len(ids)-1]
	id_arr := strings.Split(new_ids,",")

	o := orm.NewOrm()
	qs := o.QueryTable("sys_news")

	for _,v := range id_arr {
		id_int := util.StrToInt(v)
		qs.Filter("id",id_int).Update(orm.Params{
			"is_delete":1,
		})
	}
	message := map[string]interface{}{}
	message["code"] = 200
	message["msg"] = "批量删除成功"
	n.Data["json"] = message
	n.ServeJSON()

}

func (n * NewController) ShowNews(){


	o := orm.NewOrm()
	news := []models.News{}
	o.QueryTable("sys_news").All(&news)
	//每页显示的条数
	PagePerNum := 7
	Headpage := 1
	CurrentPage, err := n.GetInt("page")


	var Count int64 = 0
	OffsetNum := PagePerNum*(CurrentPage - 1)

	kw :=n.GetString("kw")
	if kw != ""{
		Count,_ = o.QueryTable("sys_news").Filter("is_delete__exact",0).Filter("title__contains",kw).Count()
		o.QueryTable("sys_news").Filter("is_delete__exact",0).Filter("title__contains",kw).Limit(PagePerNum).Offset(OffsetNum).RelatedSel().All(&news)
	}else{
		Count,_ = o.QueryTable("sys_news").Filter("is_delete",0).Count()
		o.QueryTable("sys_news").Filter("is_delete__exact",0).Limit(PagePerNum).Offset(OffsetNum).RelatedSel().All(&news)
	}
	if err != nil {
		CurrentPage = 1
	}

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

	n.Data["news"] = news
	n.Data["NextPage"] = NextPage
	n.Data["PrePage"] = PrePage
	n.Data["TotalPage"] = TotalPage
	n.Data["CurrentPage"] = CurrentPage
	n.Data["Count"] = Count
	n.Data["pagemap"] = page_map
	n.Data["headpage"] = Headpage
	n.Data["kw"] = kw
	n.TplName = "news/news_show.html"
}
func (n * NewController) DoShowNews(){

	id,_ := n.GetInt("id")
	new := models.News{}
	o := orm.NewOrm()
	o.QueryTable("sys_news").Filter("id",id).One(&new)
	n.Data["data"] = new
	n.TplName = "news/news_doShow.html"
}
