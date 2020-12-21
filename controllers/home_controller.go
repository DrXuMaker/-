package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"math"
	"shizhan/controllers/author"
	models "shizhan/models/author"
	"shizhan/util"
	"time"
)

type HomeController struct {
	beego.Controller
}

func (h * HomeController)Get(){
	//后端首页

	//查询用户权限
	userId := h.GetSession("id")
	o := orm.NewOrm()
	authIdArr := []int{}
	//interface 转int
	curUser := models.User{}
	user := models.User{Id: userId.(int)}
	o.LoadRelated(&user,"Role")
	o.QueryTable("sys_user").Filter("id", userId).One(&curUser)
	//查询role的权限
	for _, roleData := range user.Role{
		roleId := roleData.Id
		role := models.Role{Id: roleId}
		o.LoadRelated(&role,"Auth")
		for _, authData := range role.Auth{
			authId := authData.Id
			authIdArr = append(authIdArr, authId)
		}
	}

	//递归查询
	auths := []models.Auth{}
	qs := o.QueryTable("sys_auth")
	_,err := qs.Filter("pid__exact",0).Filter("id__in", authIdArr).OrderBy("-weight").All(&auths)
	if err != nil {
		return
	}

	trees := []author.Tree{}

	//一级菜单
	for _, authData := range auths{

		pid := authData.Id //根据pid获取所有子节点
		treeData := author.Tree{Id: authData.Id,AuthName: authData.AuthName,UrlFor: authData.UrlFor,Weight: authData.Weight,Children: []*author.Tree{}}
		GetChildNode(pid,&treeData)
		trees = append(trees, treeData)
	}

	// 消息通知,发送消息，使用定时任务优化
	qs1 := o.QueryTable("sys_cars_apply")
	carsApply := []models.CarsApply{}
	qs1.Filter("user_id", userId.(int)).Filter("return_status",0).Filter("notify_tag",0).All(&carsApply)

	curTime,_ := time.Parse("2006-01-02",time.Now().Format("2006-01-02"))

	for _,apply := range carsApply {
		returnDate := apply.ReturnDate
		ret := curTime.Sub(returnDate)
		content := fmt.Sprintf("%s用户，你借的车辆归还时间为%v,已经预期，请尽快归还!!",user.Username, returnDate.Format("2006-01-02"))
		if ret > 0 {  // 已经逾期
			messageNotify := models.MessageNotify{
				Flag:1,
				Title:"车辆归还逾期",
				Content:content,
				User:&user,
				ReadTag:0,

			}
			o.Insert(&messageNotify)
		}

		apply.NotifyTag = 1

		o.Update(&apply)

	}

	// 展示消息,使用websocket优化
	qs2 := o.QueryTable("sys_message_notify")
	notifyCount,_ := qs2.Filter("read_tag",0).Count()
	h.Data["notify_count"] = notifyCount
	h.Data["trees"] = trees
	h.Data["user"] = curUser
	h.TplName = "index.html"
}

func (h * HomeController)Welcome(){

	h.TplName = "welcome.html"
}

//递归拿到子节点
func GetChildNode(pid int, treenode *author.Tree){

	o :=orm.NewOrm()
	qs :=o.QueryTable("sys_auth")
	auths := []models.Auth{}
	_,err := qs.Filter("pid",pid).All(&auths)
	if err != nil {
		return
	}
	lenth := len(auths)
	//查询三级菜单及以上

	for i:=0; i<lenth; i++{
		pid := auths[i].Id
		treeData := author.Tree{Id: auths[i].Id,AuthName: auths[i].AuthName,UrlFor: auths[i].UrlFor,Weight: auths[i].Weight,Children: []*author.Tree{}}
		treenode.Children = append(treenode.Children,&treeData)
		GetChildNode(pid,&treeData)
	}

	return

}


func (h *HomeController) NotifyList()  {
	o := orm.NewOrm()

	qs := o.QueryTable("sys_message_notify")

	nofities := []models.MessageNotify{}
	// 每页显示的条数
	pagePerNum := 8
	// 当前页
	currentPage,err := h.GetInt("page")

	offsetNum := pagePerNum * (currentPage - 1)


	kw := h.GetString("kw")
	var count int64 = 0

	ret := fmt.Sprintf("当前页;%d,查询条件：%s",currentPage,kw)
	logs.Info(ret)
	if kw != ""{   // 有查询条件的
		// 总数
		count,_ = qs.Filter("title__contains",kw).Count()
		qs.Filter("title__contains",kw).Limit(pagePerNum).Offset(offsetNum).All(&nofities)
	}else {
		count,_ = qs.Count()
		qs.Limit(pagePerNum).Offset(offsetNum).All(&nofities)

	}
	if err != nil {   // 说明没有获取到当前页
		currentPage = 1
	}

	// 总页数
	countPage := int(math.Ceil(float64(count) / float64(pagePerNum)))

	prePage := 1
	if currentPage == 1{
		prePage = currentPage
	}else if currentPage > 1{
		prePage = currentPage -1
	}

	nextPage := 1
	if currentPage < countPage {
		nextPage = currentPage + 1
	}else if currentPage >= countPage {
		nextPage = currentPage
	}


	pageMap := util.Paginator(currentPage,pagePerNum,count)

	h.Data["nofities"] = nofities
	h.Data["prePage"] =prePage
	h.Data["nextPage"] = nextPage
	h.Data["currentPage"] = currentPage
	h.Data["countPage"] = countPage
	h.Data["count"] = count
	h.Data["page_map"] = pageMap
	h.Data["kw"] = kw

	h.TplName = "notify_list.html"

}

func (h *HomeController) ReadNotify()  {
	id,_ := h.GetInt("id")
	o := orm.NewOrm()
	qs := o.QueryTable("sys_message_notify")
	qs.Filter("id",id).Update(orm.Params{
		"read_tag":1,
	})
	h.Redirect(beego.URLFor("HomeController.NotifyList"),302)


}