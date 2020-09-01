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
	Name            string              `gorm:"size:60;not null;comment:'姓名'" json:"name"`
	Mobile          string              `gorm:"size:20;not null;comment:'手机'" json:"mobile"`
	Email           string              `gorm:"size:60;not null;comment:'邮箱'" json:"email"`
	WxWork          string              `gorm:"size:60;not null;comment:'企业微信'" json:"wx_work"`
	DepartmentID    int                 `gorm:"not null;comment:'部门'" json:"department_id"`
	ServiceLine     string              `gorm:"size:60;not null;comment:'服务线'" json:"service_line"`
	LevelID         int                 `gorm:"not null;comment:'级别'" json:"level_id"`
	Level           *Level              `json:"level"`
	EntryDate       models.Time         `gorm:"type:datetime;comment:'入职日期'" json:"entry_date"`
	ResignationDate models.Time         `gorm:"type:datetime;comment:'离职日期'" json:"resignation_date"`
	Status          int                 `gorm:"default:0;comment:'拟入职，已入职，未入职，已离职, 已解约'" json:"status"`
	EmployeeBasic   *EmployeeBasic      `json:"employee_basic"`
	Contracts       []*EmployeeContract `json:"contracts"`
}
