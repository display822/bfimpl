/*
* Auth : acer
* Desc : 考勤
* Time : 2020/9/1 22:44
 */

package oa

import (
	"bfimpl/models"

	"github.com/jinzhu/gorm"
)

type Attendance struct {
	gorm.Model
	EmployeeID          int         `gorm:"not null;comment:'员工ID'" json:"employee_id"`
	AttendanceDate      models.Time `gorm:"type:datetime;not null;comment:'考勤日期'" json:"attendance_date"`
	CheckinTime         models.Time `gorm:"type:datetime;comment:'签入时间'" json:"checkin_time"`
	CheckoutTime        models.Time `gorm:"type:datetime;comment:'签出时间'" json:"checkout_time"`
	AttendanceStatus    string      `gorm:"size:20;comment:'System, Normal, Pending, Confirmed'" json:"attendance_status"`
	PendingCause        string      `gorm:"size:20;comment:'Flexible,HomeOffice, Overnight, Late, Absent, Mistake'" json:"pending_cause"`
	OvertimeReferenceID string      `gorm:"size:60;comment:'加班申请记录号'" json:"overtime_reference_id"`
	LeaveReferenceID    string      `gorm:"size:60;comment:'休假申请记录号'" json:"leave_reference_id"`
	ImportFile          string      `gorm:"comment:'导入文件名'" json:"import_file"`
}
