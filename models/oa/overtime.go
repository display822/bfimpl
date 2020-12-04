/*
* Auth : acer
* Desc : 加班
* Time : 2020/9/12 22:15
 */

package oa

import (
	"bfimpl/models"

	"github.com/jinzhu/gorm"
)

type Overtime struct {
	gorm.Model
	EngagementCode string `gorm:"size:64;comment:'任务指派编码'" json:"engagement_code"`
	// workday weekend holiday
	EmpID        int         `gorm:"comment:'加班申请人id'" json:"emp_id"`
	EName        string      `gorm:"size:30;comment:'员工姓名'" json:"e_name"`
	Type         string      `gorm:"size:30;comment:'加班类型'" json:"type"`
	Project      string      `gorm:"size:64;comment:'项目'" json:"project"`
	Duration     int         `gorm:"comment:'加班时长'" json:"duration"`
	RealDuration int         `gorm:"comment:'实际加班时长'" json:"real_duration"`
	Cause        string      `gorm:"comment:'加班原因'" json:"cause"`
	OvertimeDate models.Date `gorm:"type:date;comment:'加班日期'" json:"overtime_date"`
	StartTime    models.Time `gorm:"type:datetime;comment:'开始时间'" json:"start_time"`
	EndTime      models.Time `gorm:"type:datetime;comment:'结束时间'" json:"end_time"`
	ReqTime      models.Time `gorm:"type:datetime;comment:'申请时间'" json:"req_time"`
	LeaderId     int         `gorm:"-" json:"leader_id"`
	Status       string      `gorm:"size:20;comment:'申请状态'" json:"status"`
}
