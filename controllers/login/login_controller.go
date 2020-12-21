package login

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	models "shizhan/models/author"
	_ "shizhan/models/user"
	"shizhan/util"
	"strconv"
)

//global

var Identify string
var ConNum int = 0

type LoginController struct {
	beego.Controller
}

func (t *LoginController)Get() {

	id, bs64, err:= util.GetCaptcha()
	if err != nil {
		ret := fmt.Sprintf("登录的get请求，获取验证码错误，错误信息：%v",err)
		logs.Error(ret)
		return
	}

	t.Data["captcha"] = util.Captcha{Id: id, BS64: bs64}
	t.TplName = "login/login.html"
}

func (t *LoginController)ChangeCaptcha() {

	id, bs64, err := util.GetCaptcha()
	message := map[string]interface{}{}
	if err != nil {
		ret := fmt.Sprintf("登录的get请求，生成验证码失败，验证信息：id:%s，错误信息：%v",id,err)
		logs.Error(ret)
		message["msg"] = "生成验证码失败"
		message["code"] = 404
		t.Data["json"] = message
	}else{
		t.Data["json"] = util.Captcha{Id: id, BS64: bs64, Code: 200}
	}
	t.ServeJSON()
}

func (t *LoginController)Post(){
	username := t.GetString("username")
	password := util.Str2m5(t.GetString("password"))
	captcha := t.GetString("captcha")
	verifyId := t.GetString("verify_id")


	//验证码校验
	isOk := util.VerifyCaptcha(verifyId,captcha)

	//用户名密码校验
	o := orm.NewOrm()
	user := models.User{}
	o.QueryTable("sys_user").Filter("user_name__exact",username).Filter("password__exact",password).One(&user)
	isExist := o.QueryTable("sys_user").Filter("user_name__exact",username).Filter("password__exact",password).Exist()

	//Ajax回显

	message := map[string]interface{}{}
	if !isExist {
		ret := fmt.Sprintf("登录的post请求，用户名或者密码错误，登录信息：username:%s,password:%s",username,password)
		logs.Info(ret)
		message ["code"] = 1080
		message ["msg"] = "用户名或密码错误"
		t.Data["json"] = message
	}else if !isOk {
		ret := fmt.Sprintf("登录的post请求，登录信息：验证码信息：%t", isOk)
		logs.Info(ret)
		message ["code"] = 1080
		message ["msg"] = "验证码错误"
		t.Data["json"] = message
	}else if user.IsActive == 0{
		ret := fmt.Sprintf("登录的post请求，此用户已停用：用户名：%s，状态：%t",username,user.IsActive)
		logs.Info(ret)
		message ["code"] = 1080
		message ["msg"] = "此用户已停用，请与管理员联系"
		t.Data["json"] = message
	}else{

		ret := fmt.Sprintf("登录的post请求，登录成功：用户名：%s ",username)
		logs.Info(ret)
		//用户单一登录
		//设置session
		t.SetSession("id",user.Id)
		Identify = strconv.Itoa(user.Id)
		IsExist := util.RedisIsExist(Identify,"Connection")
		if IsExist != 1{
			util.Save2Redis(Identify,"Connection",strconv.Itoa(ConNum+1))
		}else{
			util.IncRedisValue(Identify,"Connection",1)
		}
		//检查session
		checkNum := util.ReadRedis(Identify,"Connection")
		if checkNum == "1" {
			message ["code"] = 200
			message ["msg"] = "登录成功"
			t.Data["json"] = message
		}else {
			message ["code"] = 100010
			message ["msg"] = "该用户已登录"
			t.Data["json"] = message
			t.DelSession("id")
		}
	}
	t.ServeJSON()
}

func (t *LoginController)Logout() {
	util.Save2Redis(Identify,"Connection",strconv.Itoa(ConNum))
	t.DelSession("id")
	t.Redirect(beego.URLFor("LoginController.Get"),302)
}