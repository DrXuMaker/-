package user

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	models2 "shizhan/models/author"
	models "shizhan/models/user"
	"time"
)

type SalaryController struct{
	beego.Controller
}


func (u * SalaryController)Get() {

	month := u.GetString("month")
	if month == ""{
		month = time.Now().Format("2006-01")
	}
	id := u.GetSession("id")
	user := models2.User{}
	o := orm.NewOrm()
	o.QueryTable("sys_user").Filter("id",id).One(&user)
	card_id := user.CardId
	user_salary := models.SalarySlip{}
	o.QueryTable("sys_salary_slip").Filter("card_id",card_id).Filter("pay_date",month).One(&user_salary)
	u.Data["user_salary"] = user_salary
	u.TplName = "user/salary_slip.html"
}

func (u * SalaryController)ShowList() {

	month := u.GetString("month")
	if month == ""{
		month = time.Now().Format("2006-01")
	}
	id := u.GetSession("id")
	user := models2.User{}
	o := orm.NewOrm()
	o.QueryTable("sys_user").Filter("id",id).One(&user)
	card_id := user.CardId
	user_salary := models.SalarySlip{}
	o.QueryTable("sys_salary_slip").Filter("card_id",card_id).Filter("pay_date",month).One(&user_salary)
	u.Data["user_salary"] = user_salary
	u.TplName = "user/salary_list.html"
}

