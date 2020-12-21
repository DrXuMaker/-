package user

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"math"
	models "shizhan/models/author"
	"shizhan/util"
	"strconv"
	"strings"
)

type Usercontroller struct {
	beego.Controller
}

func (t * Usercontroller) List(){
	o := orm.NewOrm()
	users := []models.User{}
	//每页显示的条数
	PagePerNum := 7
	Headpage := 1
	CurrentPage, err := t.GetInt("page")


	var Count int64 = 0
	OffsetNum := PagePerNum*(CurrentPage - 1)

	kw :=t.GetString("kw")
	if kw != ""{
		Count,_ = o.QueryTable("sys_user").Filter("is_delete__exact",0).Filter("username__contains",kw).Count()
		o.QueryTable("sys_user").Filter("is_delete__exact",0).Filter("username__contains",kw).Limit(PagePerNum).Offset(OffsetNum).All(&users)
	}else{
		Count,_ = o.QueryTable("sys_user").Filter("is_delete",0).Count()
		o.QueryTable("sys_user").Filter("is_delete__exact",0).Limit(PagePerNum).Offset(OffsetNum).All(&users)
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

	t.Data["users"] = users
	t.Data["NextPage"] = NextPage
	t.Data["PrePage"] = PrePage
	t.Data["TotalPage"] = TotalPage
	t.Data["CurrentPage"] = CurrentPage
	t.Data["Count"] = Count
	t.Data["pagemap"] = page_map
	t.Data["headpage"] = Headpage
	t.Data["kw"] = kw
	t.TplName = "user/user_list.html"
}

func (t * Usercontroller) ToAdd(){

	t.TplName = "user/user_add.html"
}
func (t * Usercontroller) DoAdd(){

	username := t.GetString("username")
	password := t.GetString("userpassword")
	userpassword := util.Str2m5(password)
	userage,_ := t.GetInt("userage")
	usergendar,_ := t.GetInt("usergendar")
	userphone,_ := strconv.ParseInt(t.GetString("userphone"),10,64)
	useraddr := t.GetString("useraddr")
	is_active,_ := t.GetInt("is_active")


	o := orm.NewOrm()
	userone := models.User{Username: username,Password:userpassword,Age:userage,Gender: usergendar,Phone: userphone,Addr:useraddr,IsActive: is_active}
	_,err :=o.Insert(&userone)
	

	message := map[string]interface{}{}
	if err != nil {
		ret :=fmt.Sprintf("用户添加失败，错误信息：用户名重复,用户名：%s",username)
		logs.Error(ret)
		message ["code"] = 10010
		message ["msg"] = "用户名重复"
		t.Data["json"] = message
		t.ServeJSON()
		return
	}else{
		ret :=fmt.Sprintf("用户添加成功，用户名：%s",username)
		logs.Info(ret)
		message ["code"] = 200
		message ["msg"] = "注册成功"
		t.Data["json"] = message
	}

	t.ServeJSON()


}

func (t * Usercontroller) IsActive(){

	is_active,_ := t.GetInt("is_active")
	id,_ := t.GetInt("id")

	o := orm.NewOrm()
	qs :=o.QueryTable("sys_user").Filter("id__exact",id)
	message := map[string]interface{}{}
	if is_active == 1{
		qs.Update(orm.Params{
			"is_active":0,
		})
		message["msg"] = "停用成功"
		t.Data["json"] = message
	}else if is_active == 0{
		qs.Update(orm.Params{
			"is_active":1,
		})
		message["msg"] = "启用成功"
		t.Data["json"] = message
	}
	t.ServeJSON()

}

func (t * Usercontroller) Delete(){
	
	id,_ := t.GetInt("id")
	o := orm.NewOrm()
	qs := o.QueryTable("sys_user").Filter("id",id)
	qs.Update(orm.Params{
		"is_delete":1,
	})
	message := map[string]interface{}{}
	message["msg"] = "删除成功"
	t.Data["json"] = message

	t.ServeJSON()
}

func(t * Usercontroller) ToPsd() {
	id := t.GetString("id")
	t.Data["id"] = id
	t.TplName = "user/user_password.html"
}
func (t * Usercontroller) RetPsd(){

	new_psw := t.GetString("userpassword")
	id,_ := t.GetInt("id")
	new_psword := util.Str2m5(new_psw)

	o := orm.NewOrm()
	qs := o.QueryTable("sys_user").Filter("id",id)
	qs.Update(orm.Params{
		"password":new_psword,
	})

	message := map[string]interface{}{}
	message["msg"] = "密码已重置"
	t.Data["json"] = message
	t.ServeJSON()


}

func(t * Usercontroller) ToUser() {


	id := t.GetString("id")
	usr := models.User{}
	o := orm.NewOrm()
	o.QueryTable("sys_user").Filter("id",id).One(&usr)
	t.Data["user"] = usr
	t.TplName = "user/user_toAdd.html"

}
func (t * Usercontroller) DoUser(){

	user_id,_ := t.GetInt("user_id")
	username := t.GetString("username")
	password := t.GetString("userpassword")
	userpassword := util.Str2m5(password)
	userage,_ := t.GetInt("userage")
	usergendar,_ := t.GetInt("usergendar")
	userphone,_ := strconv.ParseInt(t.GetString("userphone"),10,64)
	useraddr := t.GetString("useraddr")
	is_active,_ := t.GetInt("is_active")

	o := orm.NewOrm()
	qs := o.QueryTable("sys_user").Filter("id",user_id)

	message := map[string]interface{}{}
	_,err := qs.Update(orm.Params{
		"user_name":username,
		"password":userpassword,
		"age":userage,
		"gender":usergendar,
		"phone":userphone,
		"addr":useraddr,
		"is_active":is_active,
	})
	if err != nil {
		ret := fmt.Sprintf("修改用户名失败，错误信息：%v,用户名重复：%s",err,username)
		logs.Error(ret)
		message ["code"] = 10010
		message ["msg"] = "修改失败用户名重复"
		t.Data["json"] = message
		t.ServeJSON()
		return
	}else{
		ret := fmt.Sprintf("修改用户名成功，用户名：%s",username)
		logs.Info(ret)
		message ["code"] = 200
		message ["msg"] = "修改成功"
		t.Data["json"] = message
	}

	t.ServeJSON()

}

func (t * Usercontroller) MuliDel(){
	
	ids := t.GetString("ids")
	new_ids := ids[1:len(ids)-1]
	id_arr := strings.Split(new_ids,",")

	o := orm.NewOrm()
	qs := o.QueryTable("sys_user")

	for _,v := range id_arr {
		id_int := util.StrToInt(v)
		qs.Filter("id",id_int).Update(orm.Params{
			"is_delete":1,
		})
	}
	message := map[string]interface{}{}
	message["code"] = 200
	message["msg"] = "批量删除成功"
	t.Data["json"] = message
	t.ServeJSON()

}