package finance

import (
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"math"
	models "shizhan/models/finance"
	"shizhan/util"
	"strconv"
	"time"
)

type EctController struct{
	beego.Controller
}

func (e * EctController) List(){


	o := orm.NewOrm()
	finDataset := []models.FinData{}
	//分页算法
	PagePerNum := 7
	Headpage := 1
	CurrentPage, err := e.GetInt("page")


	var Count int64 = 0
	OffsetNum := PagePerNum*(CurrentPage - 1)
	month := e.GetString("month")
	if month == ""{
		month = time.Now().Format("2006-01")
		Count, _ = o.QueryTable("sys_fin_data").Filter("fin_data", month).Count()
		o.QueryTable("sys_fin_data").Limit(PagePerNum).Offset(OffsetNum).All(&finDataset)
	}else{
		Count, _ = o.QueryTable("sys_fin_data").Filter("fin_data", month).Count()
		o.QueryTable("sys_fin_data").Limit(PagePerNum).Offset(OffsetNum).All(&finDataset)
	}


	if err != nil {
		CurrentPage = 1
	}

	PrePage := 1
	if CurrentPage > 1{
		PrePage = CurrentPage - 1
	}else {
		PrePage  = CurrentPage
	}

	TotalPage := int(math.Ceil(float64(Count) / float64(PagePerNum)))

	NextPage := 1
	if TotalPage > CurrentPage {
		NextPage = CurrentPage + 1
	}else if TotalPage < CurrentPage {
		NextPage = CurrentPage
	}

	pageMap := util.Paginator(CurrentPage,PagePerNum,Count)

	e.Data["NextPage"] = NextPage
	e.Data["PrePage"] = PrePage
	e.Data["TotalPage"] = TotalPage
	e.Data["CurrentPage"] = CurrentPage
	e.Data["Count"] = Count
	e.Data["pagemap"] = pageMap
	e.Data["headpage"] = Headpage
	e.Data["finDataset"] = finDataset
	e.Data["month"] = month
	e.TplName = "finance/finance_echart_list.html"
}

func (e * EctController) ToAddExcel(){

	e.TplName = "finance/excel_chart_import.html"
}
func (e * EctController) DoAddExcel(){

	w,h,err := e.GetFile("upload_file")
	defer func() {
		w.Close()
	}()
	fileName := h.Filename
	TimeUnixInt := time.Now().Unix()
	timeUnixStr := strconv.FormatInt(TimeUnixInt,10)
	filePath := "upload/salary_upload_file/" + timeUnixStr + "-" + fileName
	e.SaveToFile("upload_file", filePath)
	errDateArr := []string{}
	//json返回message
	message := map[string]interface{}{}
	if err != nil {
		message["code"] = 10001
		message["msg"] = "文件上传失败"
		e.Data["json"] = message
	}else{
		message["code"] = 200
		message["msg"] = "文件上传成功"
		e.Data["json"] = message
	}

	//读取并插入数据库
	o := orm.NewOrm()
	file,_ := excelize.OpenFile(filePath)
	rows := file.GetRows("Sheet1")
	i := 0
	for _, row := range rows{
		if i == 0{
			i ++
			continue
		}
		FinData := row[0]
		SalesVolume, _ := strconv.ParseFloat(row[1],64)
		StudentIncress := util.StrToInt(row[2])
		Django := util.StrToInt(row[3])
		VueDjango := util.StrToInt(row[4])
		Celery := util.StrToInt(row[5])
		findata := models.FinData{
			FinData:        FinData,
			SalesVolume:    SalesVolume,
			StudentIncress: StudentIncress,
			Django:         Django,
			VueDjango:      VueDjango,
			Celery:         Celery,
		}
		//重复导入相同月份的数据，先删除已有的工资月份，再导入
		qs := o.QueryTable("sys_fin_data")
		isExist := qs.Filter("fin_data",FinData).Exist()
		if isExist {
			qs.Filter("fin_data",FinData).Delete()
		}
		_,err := o.Insert(&findata)
		logs.Error(err)
		i++
	}

	message_map := map[string]interface{}{}
	if len(errDateArr) <= 0{
		message_map["code"] = 200
		message_map["msg"] = "文件上传成功"
		e.Data["json"] = message
	}else{
		message_map["code"] = 10002
		message_map["msg"] = "导入失败"
		message_map["err_Data"] = errDateArr
		e.Data["json"] = message
	}
	e.ServeJSON()

}