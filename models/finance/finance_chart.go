package finance

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type FinData struct{
	Id int `orm:"pk;auto"`
	FinData string `orm:"description(财务月份);size(32)"`
	SalesVolume float64 `orm:"description(本月销售额);digits(10);decimals(2)"`
	StudentIncress int `orm:"description(学员增加数量)"`
	Django int `orm:"description(Django课程数量)"`
	VueDjango int `orm:"description(VueDjango课程数量)"`
	Celery int `orm:"description(celery课程数量)"`
	CreateDate time.Time `orm:"type(datetime);auto_now"`
}

func (t * FinData)TableName() string {
	return "sys_fin_data"
}


func init() {
	orm.RegisterModel(new(FinData))
}