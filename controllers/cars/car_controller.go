package cars

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"math"
	models "shizhan/models/author"
	"shizhan/util"
)

type CarsController struct {
	beego.Controller
}

func (c *CarsController) Get()  {
	o := orm.NewOrm()

	qs := o.QueryTable("sys_cars")

	carsData := []models.Cars{}

	// 每页显示的条数
	pagePerNum := 8
	// 当前页
	currentPage,err := c.GetInt("page")
	if err != nil {   // 说明没有获取到当前页
		currentPage = 1
	}

	offsetNum := pagePerNum * (currentPage - 1)


	kw := c.GetString("kw")
	var count int64 = 0

	ret := fmt.Sprintf("当前页;%d,查询条件：%s",currentPage,kw)
	logs.Info(ret)
	if kw != ""{   // 有查询条件的
		// 总数
		count,_ = qs.Filter("is_delete",0).Filter("name__contains",kw).Count()
		qs.Filter("is_delete",0).Filter("name__contains",kw).Limit(pagePerNum).Offset(offsetNum).RelatedSel().All(&carsData)
	}else {
		count,_ = qs.Filter("is_delete",0).Count()
		qs.Filter("is_delete",0).Limit(pagePerNum).Offset(offsetNum).RelatedSel().All(&carsData)

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
	c.Data["cars_data"] = carsData
	c.Data["prePage"] =prePage
	c.Data["nextPage"] = nextPage
	c.Data["currentPage"] = currentPage
	c.Data["countPage"] = countPage
	c.Data["count"] = count
	c.Data["page_map"] = pageMap
	c.Data["kw"] = kw
	c.TplName = "cars/cars_list.html"

}

func (c *CarsController) ToAdd()  {
	o := orm.NewOrm()
	qs := o.QueryTable("sys_cars_brand")
	carsBrand := []models.CarBrand{}
	qs.Filter("is_delete",0).All(&carsBrand)
	c.Data["cars_brand"] = carsBrand
	c.TplName = "cars/cars_add.html"

}

func (c *CarsController) DoAdd() {
	carsBrandId, _ := c.GetInt("cars_brand_id")
	name := c.GetString("name")
	isActive, _ := c.GetInt("is_active")
	status := 1
	o := orm.NewOrm()

	carsBrand := models.CarBrand{Id: carsBrandId}
	carsData := models.Cars{
		Name:     name,
		CarBrand: &carsBrand,
		IsActive: isActive,
		Status:   status,
	}
	_, err := o.Insert(&carsData)

	messageMap := map[string]interface{}{}
	if err != nil {
		messageMap["code"] = 10001
		messageMap["msg"] = "添加失败"
	}

	messageMap["code"] = 200
	messageMap["msg"] = "添加成功"
	c.Data["json"] = messageMap
	c.ServeJSON()

}