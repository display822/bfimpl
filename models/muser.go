package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Name     string `gorm:"size:30;not null;comment:'姓名'" json:"name" valid:"Required"`
	Email    string `gorm:"unique_index;size:50;not null;comment:'邮箱'" json:"email" valid:"Required"`
	Wx       string `gorm:"size:50;not null;comment:'企业微信'" json:"wx" valid:"Required"`
	Phone    string `gorm:"size:30;not null;comment:'手机'" json:"phone" valid:"Required"`
	UserType int    `gorm:"tinyint;default:0;comment:'用户类型'" json:"userType" valid:"Required"`
	LeaderId int    `gorm:"default:0;comment:'组长id'" json:"leaderId"`
}

var UserType = map[int]string{
	1: "admin",
	2: "sale",
	3: "manager",
	4: "tm",
	5: "implement",
}

type ReqLogin struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}
