package author

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type Auth struct{
	Id int `orm:"pk;auto"`
	AuthName string `orm:"description(权限名称)"`
	UrlFor string `orm:"description(url反转)"`
	Pid int	`orm:"description(父节点id)"`
	Desc string  `orm:"description(描述)"`
	CreateTime time.Time `orm:"type(data_time);auto_now;description(创建时间)"`
	IsActive int `orm:"columns(is_active);description(是否启用)"`
	IsDelete int `orm:"columns(is_deleted);description(是否删除)"`
	Weight int `orm:"description(权重 数值越大排序越前)"`
	Role []*Role `orm:"reverse(many)"`

}

type Role struct {
	Id int `orm:"pk;auto"`
	RoleName string `orm:"description(角色名称)"`
	Desc string `orm:"description(描述)"`
	IsActive int `orm:"columns(is_active);description(是否启用)"`
	IsDelete int `orm:"columns(is_deleted);description(是否删除)"`
	CreateTime time.Time `orm:"type(data_time);auto_now;description(创建时间)"`
	Auth []*Auth `orm:"rel(m2m)"`
	User []*User `orm:"rel(m2m)"`

}

type User struct {
	Id int `orm:"pk;auto"`
	CardId string `orm:"size(64);column(card_id);description(员工工号)"`
	Username string `orm:"unique;column(user_name);size(64);description(用户)"`
	Password string	`orm:"size(32);description(密码)"`
	Age int	`orm:"null;description(年龄)"`
	Gender int `orm:"null;description(1男2女3未知)"`
	Phone int64	`orm:"null;description(手机号)"`
	Addr string	`orm:"null;size(255);description(地址)"`
	IsActive int `orm:"description(1启用,0停用);default(1)"`
	IsDelete int `orm:"description(1删除,0未删除);default(0)"`
	CreateTime time.Time  `orm:"auto_now;type(datetime);description(创建时间);null"`
	Role []*Role `orm:"reverse(many)"`
	CarsApply []*CarsApply `orm:"reverse(many)"`
	MessageNotify []*MessageNotify `orm:"reverse(many)"`
}

// 车辆品牌
type CarBrand struct {
	Id int `orm:"pk;auto"`
	Name string `orm:"description(品牌名称);size(64)"`
	Desc string `orm:"description(品牌描述);size(255)"`
	Cars []*Cars `orm:"reverse(many)"`
	IsActive int `orm:"default(1);description(启用:1,停用:0)"`
	IsDelete int `orm:"default(0);description(删除:1,未删除:0)"`
	CreateTime time.Time `orm:"description(创建时间);type(datetime);auto_now"`
}

// 车辆列表
type Cars struct {
	Id int `orm:"pk;auto"`
	Name string `orm:"description(车辆名称);size(64)"`
	CarBrand *CarBrand `orm:"rel(fk);description(车辆品牌外键)"`
	CarsApply []*CarsApply `orm:"reverse(many)"`
	Status int `orm:"description(1:可借,2:不可借);default(1)"`
	IsActive int `orm:"default(1);description(启用:1,停用:0)"`
	IsDelete int `orm:"default(0);description(删除:1,未删除:0)"`
	CreateTime time.Time `orm:"description(创建时间);type(datetime);auto_now"`
}


// 车辆申请、审核
type CarsApply struct {
	Id int `orm:"pk;auto"`
	User *User `orm:"rel(fk)"`
	Cars *Cars `orm:"rel(fk)"`
	Reason string `orm:"description(申请理由);size(255)"`
	Destination string `orm:"description(目的地);size(64)"`
	ReturnDate time.Time `orm:"type(date);auto_now;description(归还日期)"`
	ReturnStatus int `orm:"description(1:已归还，0：未归还);default(0)"`
	AuditStatus int `orm:"description(1:同意，2:未同意，3:未审批);default(3)"`
	AuditOption string `orm:"description(审批意见);size(255)"`
	IsActive int `orm:"default(1);description(启用:1,停用:0)"`
	IsDelete int `orm:"default(0);description(删除:1,未删除:0)"`
	CreateTime time.Time `orm:"description(创建时间);type(datetime);auto_now"`
	NotifyTag int `orm:"description(1:已发送通知，0：未发送通知);default(0)"`

}

// 消息通知
type MessageNotify struct {
	Id int `orm:"pk;auto"`
	Flag int `orm:"description(1:车辆逾期，2:所有通知);default(1)"`
	Title string `orm:"size(64);description(消息标题)"`
	Content string `orm:"type(text);description(消息内容)"`
	User *User `orm:"rel(fk);description(用户外键)"`
	ReadTag int `orm:"description(1:已读，0:未读)"`

}



func (u * User)TableName() string {
	return "sys_user"
}


func (r * Role)TableName() string{
	return "sys_role"
}

func (r * Auth)TableName() string{
	return "sys_auth"
}


func (c * CarBrand)TableName() string {
	return "sys_cars_brand"
}


func (c * Cars)TableName() string{
	return "sys_cars"
}

func (c * CarsApply)TableName() string{
	return "sys_cars_apply"
}

func (m * MessageNotify)TableName() string{
	return "sys_message_notify"
}

func init() {
	orm.RegisterModel(
		new(Auth),
		new(Role),
		new(User),
		new(CarBrand),
		new(Cars),
		new(CarsApply),
		new(MessageNotify),
		)
}