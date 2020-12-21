package author

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"math"
	models "shizhan/models/author"
	"shizhan/util"
	"strings"
)

type Rolecontroller struct{
	beego.Controller
}

func (r * Rolecontroller)List(){

	o := orm.NewOrm()
	roles := []models.Role{}
	//分页算法
	PagePerNum := 5
	CurrentPage,_ := r.GetInt("page")
	OffsetNum := PagePerNum*(CurrentPage - 1)

	Count,_ := o.QueryTable("sys_role").Filter("is_delete__exact",0).Count()
	o.QueryTable("sys_role").Limit(PagePerNum).Offset(OffsetNum).Filter("is_delete__exact",0).All(&roles)

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

	r.Data["total_page"] = TotalPage
	r.Data["current_page"] = CurrentPage
	r.Data["pre_page"] = PrePage
	r.Data["next_page"] = NextPage
	r.Data["head_page"] = HeadPage
	r.Data["last_page"] = LastPage
	r.Data["pagemap"] = page_map
	r.Data["roles"] = roles
	r.Data["count"] = Count
	r.TplName = "author/role_list.html"
}

func (r * Rolecontroller) ToAdd() {

	o := orm.NewOrm()
	roles := []models.Role{}
	o.QueryTable("sys_role").Filter("is_delete",0).All(&roles)
	r.Data["roles"] = roles
	r.TplName = "author/role_add.html"
}
func (r * Rolecontroller) DoAdd(){

	role_name := r.GetString("role_name")
	role_desc := r.GetString("role_desc")
	is_active,_ := r.GetInt("is_active")
	//将数据插入数据库中
	o := orm.NewOrm()
	role := models.Role{RoleName:role_name,Desc:role_desc,IsActive: is_active}
	_,err := o.Insert(&role)
	message := map[string]interface{}{}
	if err != nil {
		message ["code"] = 10010
		message ["msg"] = "角色名重复"
		r.Data["json"] = message
		r.ServeJSON()
		return
	}else{
		message ["code"] = 200
		message["msg"] = "添加成功"
		r.Data["json"] = message

	}
	r.ServeJSON()
}

func (r * Rolecontroller) IsActive(){

	is_active,_ := r.GetInt("is_active")
	id,_ := r.GetInt("id")

	o := orm.NewOrm()
	qs :=o.QueryTable("sys_role").Filter("id__exact",id)
	message := map[string]interface{}{}
	if is_active == 1{
		qs.Update(orm.Params{
			"is_active":0,
		})
		message["msg"] = "停用成功"
		r.Data["json"] = message
	}else if is_active == 0{
		qs.Update(orm.Params{
			"is_active":1,
		})
		message["msg"] = "启用成功"
		r.Data["json"] = message
	}
	r.ServeJSON()

}

func (r * Rolecontroller) Delete(){

	id,_ := r.GetInt("id")
	o := orm.NewOrm()
	qs := o.QueryTable("sys_role").Filter("id",id)
	qs.Update(orm.Params{
		"is_delete":1,
	})
	message := map[string]interface{}{}
	message["msg"] = "删除成功"
	r.Data["json"] = message

	r.ServeJSON()
}
//角色--用户配置
func (r * Rolecontroller) ToRoleUser() {

	id,_ := r.GetInt("role_id")
	role := models.Role{}
	o := orm.NewOrm()
	o.QueryTable("sys_role").Filter("id",id).One(&role)

	//已绑定的用户
	o.LoadRelated(&role,"User")

	//未绑定的用户
	users := []models.User{}
	if len(role.User) > 0{
		o.QueryTable("sys_user").Filter("is_active",1).Filter("is_delete",0).Exclude("id__in",role.User).All(&users)
	}else{
		o.QueryTable("sys_user").Filter("is_active",1).Filter("is_delete",0).All(&users)
	}
	r.Data["role"] = role
	r.Data["users"] = users
	r.TplName = "author/role_user_add.html"
}
func (r * Rolecontroller) DoRoleUser(){


	o := orm.NewOrm()
	user_ids := r.GetString("user_ids")
	role_id,_ := r.GetInt("role_id")
	user_ids_arr := strings.Split(user_ids,",")
	role := models.Role{Id: role_id}

	//查询已绑定数据
	m2m := o.QueryM2M(&role,"User")
	m2m.Clear()


	for _,user_id := range user_ids_arr{
		user := models.User{Id: util.StrToInt(user_id)}
		m2m := o.QueryM2M(&role,"User")
		m2m.Add(user)
	}
	message := map[string]interface{}{}
	message["code"] = 200
	message["msg"] = "添加成功"
	r.Data["json"] = message
	r.ServeJSON()
}

//角色--权限配置
var AuthRoleId int
func (r * Rolecontroller) ToRoleAuth() {

	id,_ := r.GetInt("role_id")
	AuthRoleId = id
	role := models.Role{}
	o := orm.NewOrm()
	o.QueryTable("sys_role").Filter("id",id).One(&role)
	r.Data["role"] = role
	r.TplName = "author/role_auth_add.html"
}
func (r * Rolecontroller) DoRoleAuth(){



	auth_ids := r.GetString("auth_ids")
	role_id,_ := r.GetInt("role_id")
	//new_auth_ids := auth_ids[1:len(auth_ids)-1]
	auth_ids_arr := strings.Split(auth_ids,",")

	//清除数据
	o := orm.NewOrm()
	role := models.Role{Id:role_id}
	m2m := o.QueryM2M(&role,"Auth")
	m2m.Clear()

	//添加数据
	for _,auth_id := range auth_ids_arr{
		auth_data := util.StrToInt(auth_id)
		auth := models.Auth{Id: auth_data}
		m2m := o.QueryM2M(&role,"Auth")
		m2m.Add(auth)
	}
	message := map[string]interface{}{"code":200,"msg":"操作成功"}
	r.Data["json"] = message
	r.ServeJSON()
}

func (r * Rolecontroller) GetAuth() {

	id := AuthRoleId
	o := orm.NewOrm()
	auths := []models.Auth{}
	qs := o.QueryTable("sys_auth")
	//已绑定的权限
	role := models.Role{Id:id}
	o.LoadRelated(&role,"Auth")

	var authIdHas []int
	for _,authData := range role.Auth {
		authIdHas = append(authIdHas,authData.Id)
	}

	//所有权限
	qs.Filter("is_delete",0).All(&auths)

	auth_arr_map := []map[string]interface{}{}

	for _, auth_data := range auths {
		id := auth_data.Id
		pId := auth_data.Pid
		name := auth_data.AuthName
		if (pId == 0){
			auth_map := map[string]interface{}{"id":id,"pId":pId,"name":name,"open":false}
			auth_arr_map = append(auth_arr_map,auth_map)
		}else{
			auth_map := map[string]interface{}{"id":id,"pId":pId,"name":name}
			auth_arr_map = append(auth_arr_map,auth_map)
		}
	}
	authMaps := map[string]interface{}{}
	authMaps["auth_arr_map"] = auth_arr_map
	authMaps["auth_ids_has"] = authIdHas
	r.Data["json"] = authMaps
	r.ServeJSON()

}
