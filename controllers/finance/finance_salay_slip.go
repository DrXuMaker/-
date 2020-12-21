package finance

import (
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"math"
	models "shizhan/models/user"
	"shizhan/util"
	"strconv"
	"time"
)

type FinController struct {
	beego.Controller
}

func (f * FinController) List(){


	o := orm.NewOrm()
	salarys := []models.SalarySlip{}
	//分页算法
	PagePerNum := 7
	Headpage := 1
	CurrentPage, err := f.GetInt("page")


	var Count int64 = 0
	OffsetNum := PagePerNum*(CurrentPage - 1)
	month := f.GetString("month")
	if month == ""{
		month = time.Now().Format("2006-01")
		Count, _ = o.QueryTable("sys_salary_slip").Filter("pay_date", month).Count()
		o.QueryTable("sys_salary_slip").Filter("pay_date",month).All(&salarys)
	}else{
		Count, _ = o.QueryTable("sys_salary_slip").Filter("pay_date", month).Count()
		o.QueryTable("sys_salary_slip").Limit(PagePerNum).Offset(OffsetNum).All(&salarys)
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

	f.Data["NextPage"] = NextPage
	f.Data["PrePage"] = PrePage
	f.Data["TotalPage"] = TotalPage
	f.Data["CurrentPage"] = CurrentPage
	f.Data["Count"] = Count
	f.Data["pagemap"] = pageMap
	f.Data["headpage"] = Headpage
	f.Data["salarys"] = salarys
	f.Data["month"] = month
	f.TplName = "finance/finance_list.html"
}

func (f * FinController) ToAddExcel(){

	f.TplName = "finance/excel_import.html"
}
func (f * FinController) DoAddExcel(){

	w,h,err := f.GetFile("upload_file")
	defer func() {
		w.Close()
	}()
	fileName := h.Filename
	TimeUnixInt := time.Now().Unix()
	timeUnixStr := strconv.FormatInt(TimeUnixInt,10)
	filePath := "upload/salary_upload_file/" + timeUnixStr + "-" + fileName
	f.SaveToFile("upload_file", filePath)
	errDateArr := []string{}
	//json返回message
	message := map[string]interface{}{}
	if err != nil {
		message["code"] = 10001
		message["msg"] = "文件上传失败"
		f.Data["json"] = message
	}else{
		message["code"] = 200
		message["msg"] = "文件上传成功"
		f.Data["json"] = message
	}

	//读取并插入数据库
	o := orm.NewOrm()
	salayArr := []models.SalarySlip{}
	file,_ := excelize.OpenFile(filePath)
	rows := file.GetRows("Sheet1")
	i := 0
	for _, row := range rows{
		if i == 0{
			i ++
			continue
		}
		cardId := row[2]
		basePay, _ := strconv.ParseFloat(row[3], 64)
		workingDay, _ := strconv.ParseFloat(row[4], 64)
		daysOff, _ := strconv.ParseFloat(row[5], 64)
		daysOffNo, _ := strconv.ParseFloat(row[6], 64)
		reward, _ := strconv.ParseFloat(row[7], 64)
		rentSubsidy, _ := strconv.ParseFloat(row[8], 64)
		transSubsidy, _ := strconv.ParseFloat(row[9], 64)
		socialSecurity, _ := strconv.ParseFloat(row[10], 64)
		housProvidentFund, _ := strconv.ParseFloat(row[11], 64)
		personalPncomeTax, _ := strconv.ParseFloat(row[12], 64)
		fine, _ := strconv.ParseFloat(row[13], 64)
		netSalary, _ := strconv.ParseFloat(row[14], 64)
		payDate := row[15]
		salarySlip := models.SalarySlip{
			CardId:            cardId,
			BasePay:           basePay,
			WorkingDay:        workingDay,
			DaysOff:           daysOff,
			DaysOffNo:         daysOffNo,
			Reward:            reward,
			RentSubsidy:       rentSubsidy,
			TransSubsidy:      transSubsidy,
			SocialSecurity:    socialSecurity,
			HousProvidentFund: housProvidentFund,
			PersonalPncomeTax: personalPncomeTax,
			Fine:              fine,
			NetSalary:         netSalary,
			PayDate:           payDate,
		}
		//重复导入相同月份的数据，先删除已有的工资月份，再导入
		qs := o.QueryTable("sys_salary_slip")
		isExist := qs.Filter("pay_date",payDate).Exist()
		if isExist {
			qs.Filter("pay_date",payDate).Delete()
		}
		_, err3 := o.Insert(&salarySlip)
		//精确到导入失败的数据信息提示
		if err3 != nil{
			errDateArr = append(errDateArr,cardId)
		}else{
			salayArr = append(salayArr,salarySlip)
		}
		i++
	}

	message_map := map[string]interface{}{}
	if len(errDateArr) <= 0{
		message_map["code"] = 200
		message_map["msg"] = "文件上传成功"
		f.Data["json"] = message
	}else{
		o.InsertMulti(100,salayArr)
		message_map["code"] = 10002
		message_map["msg"] = "导入失败"
		message_map["err_Data"] = errDateArr
		f.Data["json"] = message
	}
	f.ServeJSON()

}