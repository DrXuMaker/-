package echarts

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type EtcController struct{
	beego.Controller
}

func (e * EtcController)FList(){

	e.TplName = "echarts/echarts_finance.html"
}

func (e *EtcController) GetChart()  {

	var FiDate orm.ParamsList
	var studentIncress orm.ParamsList
	o := orm.NewOrm()
	o.Raw("select fin_data from sys_fin_data").ValuesFlat(&FiDate)
	o.Raw("select student_incress from sys_fin_data").ValuesFlat(&studentIncress)

	mapData := map[string]interface{}{}

	mapData["FinDate"]  = FiDate
	mapData["student"] = studentIncress

	e.Data["json"] = mapData
	e.ServeJSON()


}