/*
* Auth : acer
* Desc : 员工主表
* Time : 2020/9/1 21:27
 */

package oa

import (
	"bfimpl/models"

	"github.com/jinzhu/gorm"
)

type Employee struct {
	gorm.Model
	Name             string      `gorm:"size:60;not null;default:'';comment:'姓名'" json:"name"`
	Gender           string      `gorm:"size:10;not null;default:'';comment:'性别'" json:"gender"`
	Status           int         `gorm:"default:0;comment:'拟入职，已入职，未入职，已离职, 已解约'" json:"status"`
	Mobile           string      `gorm:"size:20;not null;default:'';comment:'手机'" json:"mobile"`
	IDCard           string      `gorm:"size:20;not null;default:'';comment:'身份证'" json:"id_card"`
	InterviewComment string      `gorm:"not null;default:'';comment:'面试评价'" json:"interview_comment"`
	Resume           string      `gorm:"not null;default:'';comment:'简历地址'" json:"resume"`
	Email            string      `gorm:"size:60;not null;default:'';comment:'邮箱'" json:"email"`
	WxWork           string      `gorm:"size:60;not null;default:'';comment:'企业微信'" json:"wx_work"`
	Tapd             string      `gorm:"size:60;not null;default:'';comment:'tapd账号'" json:"tapd"`
	ServiceLine      string      `gorm:"size:60;not null;default:'';comment:'服务线'" json:"service_line"`
	DepartmentID     int         `gorm:"not null;default:0;comment:'部门'" json:"department_id"`
	Department       *Department `json:"department"`
	LevelID          int         `gorm:"not null;default:0;comment:'级别'" json:"level_id"`
	Level            *Level      `json:"level"`
	Position         string      `gorm:"size:30;default:'';comment:'岗位'" json:"position"`
	EntryDate        models.Time `gorm:"type:datetime;comment:'入职日期'" json:"entry_date"`
	ResignationDate  models.Time `gorm:"type:datetime;comment:'离职日期'" json:"resignation_date"`
	//EmployeeBasic    *EmployeeBasic      `json:"employee_basic"`
	//Contracts        []*EmployeeContract `json:"contracts"`
}
