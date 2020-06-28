package models

import (
	"bfimpl/models/forms"
	"bfimpl/services/util"

	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Name     string `gorm:"size:30;not null;comment:'姓名'" json:"name"`
	Email    string `gorm:"unique_index;size:50;not null;comment:'邮箱'" json:"email"`
	PassWd   string `gorm:"size:70;not null;comment:'密码'" json:"-"`
	Wx       string `gorm:"size:50;not null;comment:'企业微信'" json:"wx"`
	Phone    string `gorm:"size:30;not null;comment:'手机'" json:"phone"`
	UserType int    `gorm:"tinyint;default:0;comment:'用户类型'" json:"useType"`
	LeaderId int    `gorm:"default:0;comment:'组长id'" json:"-"`
}

var UserType = map[int]string{
	1: "管理员",
	2: "销售",
	3: "客户服务经理",
	4: "资源分配",
	5: "实施人员",
}

func NewUser(req *forms.ReqUser) *User {
	u := new(User)
	u.Name = req.Name
	u.Email = req.Email
	u.PassWd = util.StringMd5("123456")
	u.Wx = req.Wx
	u.Phone = req.Phone
	u.UserType = req.UserType
	return u
}