package user

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	models "shizhan/models/author"
	"shizhan/util"
	"strconv"
)

type CentreController struct{
	beego.Controller
}

func (u * CentreController) Get() {

	id := u.GetSession("id").(int)
	o := orm.NewOrm()
	user := models.User{}
	qs := o.QueryTable("sys_user")
	qs.Filter("id",id).One(&user)
	u.Data["user"] =user
	u.TplName = "user/my_centre.html"
}

func (u * CentreController) Post() {

	userid,_ := u.GetInt("user_id")
	username := u.GetString("username")
	password := u.GetString("userpassword")
	userpassword := util.Str2m5(password)
	userage,_ := u.GetInt("userage")
	usergendar,_ := u.GetInt("usergendar")
	userphone,_ := strconv.ParseInt(u.GetString("userphone"),10,64)
	useraddr := u.GetString("useraddr")
	is_active,_ := u.GetInt("is_active")
	//查询修改
	o := orm.NewOrm()
	qs :=o.QueryTable("sys_user").Filter("id",userid)
	_,err := qs.Update(orm.Params{
		"Username": username,
		"Password":userpassword,
		"Age":userage,
		"Gender": usergendar,
		"Phone": userphone,
		"Addr":useraddr,
		"IsActive": is_active,
	})
	message := map[string]interface{}{}
	if err != nil {
		message ["code"] = 10010
		message ["msg"] = "用户名重复"
		u.Data["json"] = message
		u.ServeJSON()
		return
	}else{
		message ["code"] = 200
		message ["msg"] = "修改成功"
		u.Data["json"] = message
	}
	u.ServeJSON()

}