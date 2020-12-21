package routers

import (
	"github.com/astaxie/beego"
	"shizhan/controllers"
	"shizhan/controllers/author"
	"shizhan/controllers/cars"
	"shizhan/controllers/echarts"
	"shizhan/controllers/finance"
	"shizhan/controllers/login"
	"shizhan/controllers/news"
	"shizhan/controllers/user"
)

func init() {
	//不需要拦截的
    beego.Router("/", &login.LoginController{})
    beego.Router("/changecaptcha", &login.LoginController{},"get:ChangeCaptcha")
    //必须登录才可访问

    //登出
	beego.Router("/main/log_out", &login.LoginController{},"get:Logout")

    //后台首页
    beego.Router("/main/home",&controllers.HomeController{})
	beego.Router("/main/index/notify",&controllers.HomeController{},"get:NotifyList")
	beego.Router("/main/index/read_notify",&controllers.HomeController{},"get:ReadNotify")
	beego.Router("/main/welcome",&controllers.HomeController{},"get:Welcome")

    //user模块
    beego.Router("/main/user/list",&user.Usercontroller{},"get:List")
    beego.Router("/main/user/to_add",&user.Usercontroller{},"get:ToAdd")
    beego.Router("/main/user/do_add", &user.Usercontroller{},"post:DoAdd")
	beego.Router("/main/user/is_active", &user.Usercontroller{},"post:IsActive")
	beego.Router("/main/user/delete", &user.Usercontroller{},"get:Delete")
	beego.Router("/main/user/showpsd", &user.Usercontroller{},"get:ToPsd")
	beego.Router("/main/user/reset_psw", &user.Usercontroller{},"post:RetPsd")
	beego.Router("/main/user/to_user", &user.Usercontroller{},"get:ToUser")
	beego.Router("/main/user/do_user", &user.Usercontroller{},"post:DoUser")
	beego.Router("/main/user/muli_del", &user.Usercontroller{},"post:MuliDel")

    //Auth模块
	beego.Router("/main/auth/list", &author.Authcontroller{},"get:List")
	beego.Router("/main/auth/to_add", &author.Authcontroller{},"get:ToAdd")
	beego.Router("/main/auth/do_add", &author.Authcontroller{},"post:DoAdd")
	beego.Router("/main/auth/is_active", &author.Authcontroller{},"post:IsActive")
	beego.Router("/main/auth/delete", &author.Authcontroller{},"get:Delete")

    //Role模块
	beego.Router("/main/role/list", &author.Rolecontroller{},"get:List")
	beego.Router("/main/role/to_add", &author.Rolecontroller{},"get:ToAdd")
	beego.Router("/main/role/to_add", &author.Rolecontroller{},"post:DoAdd")
	beego.Router("/main/role/is_active", &author.Rolecontroller{},"post:IsActive")
	beego.Router("/main/role/delete", &author.Rolecontroller{},"get:Delete")
	//角色--用户配置
	beego.Router("/main/role/to_role_user", &author.Rolecontroller{},"get:ToRoleUser")
	beego.Router("/main/role/do_role_user", &author.Rolecontroller{},"post:DoRoleUser")
	//角色--权限配置
	beego.Router("/main/role/to_role_auth", &author.Rolecontroller{},"get:ToRoleAuth")
	beego.Router("/main/role/do_role_auth", &author.Rolecontroller{},"post:DoRoleAuth")
	beego.Router("/main/role/get_auth", &author.Rolecontroller{},"get:GetAuth")
    //个人中心
    beego.Router("/main/role/my_centre",&user.CentreController{})
	beego.Router("/main/role/my_centre_post",&user.CentreController{},"post:Post")
    //工资条
	beego.Router("/main/role/salary_slip",&user.SalaryController{})
	beego.Router("/main/role/salary_list",&user.SalaryController{},"get:ShowList")

    //财务中心
	beego.Router("/main/role/finance_list",&finance.FinController{},"get:List")
	beego.Router("/main/role/to_add_file",&finance.FinController{},"get:ToAddExcel")
	beego.Router("/main/role/do_add_file",&finance.FinController{},"post:DoAddExcel")
	beego.Router("/main/role/finance_chart_list",&finance.EctController{},"get:List")
	beego.Router("/main/role/finance_chart_to_add",&finance.EctController{},"get:ToAddExcel")
	beego.Router("/main/role/finance_chart_do_add",&finance.EctController{},"post:DoAddExcel")

    //内容管理
	beego.Router("/main/role/cart_list",&news.CarController{},"get:List")
	beego.Router("/main/role/cart_list_to_add",&news.CarController{},"get:ToAddCar")
	beego.Router("/main/role/cart_list_do_add",&news.CarController{},"post:DoAddCar")
	beego.Router("/main/role/cart_delete",&news.CarController{},"get:Delete")
	beego.Router("/main/role/news_list",&news.NewController{},"get:List")
	beego.Router("/main/role/news_list_to_add",&news.NewController{},"get:ToAddNews")
	beego.Router("/main/role/news_list_do_add",&news.NewController{},"post:DoAddNews")
	beego.Router("/main/role/news_list_do_add_img",&news.NewController{},"post:UpLoadImg")
	beego.Router("/main/role/news_list_to_edit",&news.NewController{},"get:ToEdit")
	beego.Router("/main/role/news_list_do_edit",&news.NewController{},"post:DoEdit")
	beego.Router("/main/role/news_is_active", &news.NewController{},"post:IsActive")
	beego.Router("/main/role/news_delete",&news.NewController{},"get:Delete")
	beego.Router("/main/user/news_multi_del", &news.NewController{},"post:MuliDel")
	beego.Router("/main/news_show", &news.NewController{},"get:ShowNews")
	beego.Router("/main/news_DoShow", &news.NewController{},"get:DoShowNews")

	// 车辆管理模块
	beego.Router("/main/cars/car_brand_list",&cars.CarBrandController{})
	beego.Router("/main/cars/to_car_brand_add",&cars.CarBrandController{},"get:ToAdd")
	beego.Router("/main/cars/do_car_brand_add",&cars.CarBrandController{},"post:DoAdd")

	beego.Router("/main/cars/cars_list",&cars.CarsController{})
	beego.Router("/main/cars/to_cars_add",&cars.CarsController{},"get:ToAdd")
	beego.Router("/main/cars/do_cars_add",&cars.CarsController{},"post:DoAdd")

	beego.Router("/main/cars/cars_apply_list",&cars.CarsApplyController{})
	beego.Router("/main/cars/to_cars_apply",&cars.CarsApplyController{},"get:ToApply")
	beego.Router("/main/cars/do_cars_apply",&cars.CarsApplyController{},"post:DoApply")
	beego.Router("/main/cars/my_apply",&cars.CarsApplyController{},"get:MyApply")
	beego.Router("/main/cars/audit_apply",&cars.CarsApplyController{},"get:AuditApply")
	beego.Router("/main/cars/to_audit_apply",&cars.CarsApplyController{},"get:ToAuditApply")
	beego.Router("/main/cars/do_audit_apply",&cars.CarsApplyController{},"post:DoAuditApply")
	beego.Router("/main/cars/do_return",&cars.CarsApplyController{},"get:DoReturn")

    //报表模块
    beego.Router("/main/echarts/list",&echarts.EtcController{},"get:FList")
	beego.Router("/main/echarts/get_fin_data",&echarts.EtcController{},"get:GetChart")




}
