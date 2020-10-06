/*
* Auth : acer
* Desc : 考勤
* Time : 2020/9/1 22:44
 */

package oa

import (
	"bfimpl/models"

	"fmt"

	"github.com/jinzhu/gorm"
)

type Attendance struct {
	gorm.Model
	EmployeeID     int         `gorm:"default:0;not null;comment:'员工ID'" json:"employee_id"`
	Dept           string      `gorm:"index;size:20;comment:'部门'" json:"dept"`
	Name           string      `gorm:"size:20;comment:'姓名'" json:"name"`
	AttendanceDate models.Date `gorm:"type:date;not null;comment:'考勤日期'" json:"attendance_date"`
	CheckIn        models.Time `gorm:"type:datetime;comment:'打卡时间'" json:"check_in"`
	CheckOut       models.Time `gorm:"type:datetime;comment:'打卡时间'" json:"check_out"`
	InStatus       string      `gorm:"size:20;comment:'Normal, Exception'" json:"in_status"`
	OutStatus      string      `gorm:"size:20;comment:'Normal, Exception'" json:"out_status"`
	InResult       string      `gorm:"size:20;comment:'异常判断'" json:"in_result"`
	OutResult      string      `gorm:"size:20;comment:'异常判断'" json:"out_result"`
	OvertimeID     string      `gorm:"size:60;comment:'加班申请记录号'" json:"overtime_id"`
	LeaveID        string      `gorm:"size:60;comment:'休假申请记录号'" json:"leave_id"`
	//ImportFile     string      `gorm:"comment:'导入文件名'" json:"import_file"`
}

type AttendanceSimple struct {
	Dept           string
	Name           string
	AttendanceDate models.Date
	CheckIn        models.Time
	CheckOut       models.Time
	InStatus       string
	OutStatus      string
	InResult       string
	OutResult      string
}

func (v AttendanceSimple) String(now string) string {
	return fmt.Sprintf("('%s','%s','%s','%s','%s','%s','%s','%s','%s','%s')", now, v.Dept, v.Name,
		v.AttendanceDate.String(), v.CheckIn.String(), v.CheckOut.String(), v.InStatus, v.OutStatus,
		v.InResult, v.OutResult)
}
