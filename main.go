package main

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "shizhan/routers"
	"github.com/astaxie/beego"
	_"github.com/go-sql-driver/mysql"
	_"shizhan/models/author"
	_"shizhan/models/user"
	_"shizhan/models/finance"
	_"shizhan/models/news"
	"shizhan/util"
)

func init() {

	username := beego.AppConfig.String("username")
	psw := beego.AppConfig.String("psw")
	host := beego.AppConfig.String("host")
	port := beego.AppConfig.String("port")
	db := beego.AppConfig.String("db")
	data_source := username + ":" + psw + "@tcp(" + host + ":" + port + ")/" + db + "?charset=utf8"
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default","mysql",data_source)
}



func main() {
	//orm数据库命令行迁移
	orm.RunCommand()
	//未登录请求拦截
	beego.InsertFilter("/main/*",beego.BeforeRouter,util.LoginFilter)

	//日志
	logs.SetLogger(logs.AdapterMultiFile,`{"filename": "logs/site.log","separate":["error","info"]}`)

	beego.SetStaticPath("/upload","upload")

	beego.Run()
}

