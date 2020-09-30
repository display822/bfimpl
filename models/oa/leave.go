/*
* Auth : acer
* Desc : 请假
* Time : 2020/9/12 22:15
 */

package oa

import (
	"bfimpl/models"

	"github.com/jinzhu/gorm"
)

const (
	PrivateAffair = "PrivateAffair"
	Sick          = "Sick"
	Shift         = "Shift"
	Flexible      = "Flexible"
	Annual        = "Annual"
	Marital       = "Marital"
	Maternity     = "Maternity"
	Funeral       = "Funeral"
)

type Leave struct {
	gorm.Model
	// workday weekend holiday
	EmpID        int         `gorm:"comment:'加班申请人id'" json:"emp_id"`
	EName        string      `gorm:"size:30;comment:'员工姓名'" json:"e_name"`
	Type         string      `gorm:"size:30;comment:'加班类型'" json:"type"`
	Duration     int         `gorm:"comment:'请假时长'" json:"duration"`
	RealDuration int         `gorm:"comment:'实际请假时长'" json:"real_duration"`
	Cause        string      `gorm:"comment:'加班原因'" json:"cause"`
	StartDate    models.Date `gorm:"type:datetime;comment:'开始时间'" json:"start_date"`
	EndDate      models.Date `gorm:"type:datetime;comment:'结束时间'" json:"end_date"`
	Start        string      `gorm:"size:10;comment:'am,pm'" json:"start"`
	End          string      `gorm:"size:10;comment:'am,pm'" json:"end"`
	ReqTime      models.Time `gorm:"type:datetime;comment:'申请时间'" json:"req_time"`
	Status       string      `gorm:"size:20;comment:'申请状态'" json:"status"`
}