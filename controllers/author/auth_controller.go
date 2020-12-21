package author

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"math"
	models "shizhan/models/author"
	"shizhan/util"
)

type Authcontroller struct{
	beego.Controller
}

func (a * Authcontroller) List(){

	o := orm.NewOrm()
	auths := []models.Auth{}
	//分页算法
	PagePerNum := 5
	CurrentPage,_ := a.GetInt("page")
	OffsetNum := PagePerNum*(CurrentPage - 1)

	Count,_ := o.QueryTable("sys_auth").Filter("is_delete__exact",0).Count()
	o.QueryTable("sys_auth").Limit(PagePerNum).Offset(OffsetNum).Filter("is_delete__exact",0).All(&auths)

	TotalPage := int(math.Ceil(float64(Count)/float64(PagePerNum)))

	HeadPage := 1
	//前一页与后一页
	PrePage := 1
	if PrePage <= 1 {
		PrePage = 1
	}else{
		PrePage = CurrentPage - 1
	}
	//尾页
	LastPage := TotalPage
	NextPage := TotalPage
	if TotalPage <= CurrentPage{
		NextPage = TotalPage
	}else{
		NextPage = CurrentPage + 1
	}

	//分页算法
	page_map := util.Paginator(CurrentPage,PagePerNum,Count)

	a.Data["total_page"] = TotalPage
	a.Data["current_page"] = CurrentPage
	a.Data["pre_page"] = PrePage
	a.Data["next_page"] = NextPage
	a.Data["head_page"] = HeadPage
	a.Data["last_page"] = LastPage
	a.Data["pagemap"] = page_map
	a.Data["auths"] = auths
	a.Data["count"] = Count
	a.TplName = "author/auth_list.html"
}

func (a * Authcontroller) ToAdd(){

	o := orm.NewOrm()
	auths := []models.Auth{}
	o.QueryTable("sys_auth").Filter("is_delete__exact",0).All(&auths)
	a.Data["auths"] = auths
	a.TplName = "author/auth_add.html"
}
func (a * Authcontroller) DoAdd(){

	authName := a.GetString("auth_name")
	authUrl := a.GetString("auth_url")
	authDesc := a.GetString("auth_desc")
	isActive,_ := a.GetInt("is_active")
	authWeight,_ := a.GetInt("auth_weight")
	parentId,_ := a.GetInt("parent_id")

	//将数据插入数据库中
	o := orm.NewOrm()
	auth := models.Auth{AuthName: authName,UrlFor: authUrl,Pid: parentId,Desc: authDesc,IsActive: isActive,Weight: authWeight}
	_,err := o.Insert(&auth)
	message := map[string]interface{}{}
	if err != nil {
		message ["code"] = 10010
		message ["msg"] = "用户名重复"
		a.Data["json"] = message
		a.ServeJSON()
		return
	}else{
		message ["code"] = 200
		message["msg"] = "添加成功"
		a.Data["json"] = message

	}
	a.ServeJSON()
}

func (a * Authcontroller) IsActive(){

	isActive,_ := a.GetInt("is_active")
	id,_ := a.GetInt("id")

	o := orm.NewOrm()
	qs :=o.QueryTable("sys_auth").Filter("id__exact",id)
	message := map[string]interface{}{}
	if isActive == 1{
		qs.Update(orm.Params{
			"is_active":0,
		})
		message["msg"] = "停用成功"
		a.Data["json"] = message
	}else if isActive == 0{
		qs.Update(orm.Params{
			"is_active":1,
		})
		message["msg"] = "启用成功"
		a.Data["json"] = message
	}
	a.ServeJSON()

}

func (a * Authcontroller) Delete(){

	id,_ := a.GetInt("id")
	o := orm.NewOrm()
	qs := o.QueryTable("sys_auth").Filter("id",id)
	qs.Update(orm.Params{
		"is_delete":1,
	})
	message := map[string]interface{}{}
	message["msg"] = "删除成功"
	a.Data["json"] = message

	a.ServeJSON()
}
