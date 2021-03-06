/*
* Auth : acer
* Desc : 考勤
* Time : 2020/9/1 22:44
 */

package oa

import (
	"bfimpl/models"

	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

//考勤确认数据
type Attendance struct {
	ID             uint        `gorm:"primary_key"`
	CreatedAt      time.Time   `json:"-"`
	UpdatedAt      time.Time   `json:"-"`
	DeletedAt      *time.Time  `sql:"index" json:"-"`
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
	OvertimeID     int         `gorm:"size:60;comment:'加班申请记录号'" json:"overtime_id"`
	LeaveID        string      `gorm:"size:255;comment:'休假申请记录号'" json:"leave_id"`
	Shift          float32     `gorm:"comment:'弹性'" json:"-"`
	Overtime       float32     `gorm:"comment:'加班'" json:"-"`
	//ImportFile     string      `gorm:"comment:'导入文件名'" json:"import_file"`
}

//考勤临时数据
type AttendanceTmp struct {
	ID             uint        `gorm:"primary_key"`
	CreatedAt      time.Time   `json:"-"`
	DeletedAt      *time.Time  `sql:"index" json:"-"`
	IsConfirm      int         `gorm:"default:0;not null;comment:'当月是否确认'" json:"-"`
	EmployeeID     int         `gorm:"default:0;not null;comment:'员工ID'" json:"employee_id"`
	Dept           string      `gorm:"index;size:20;comment:'部门'" json:"dept"`
	Name           string      `gorm:"index;size:20;comment:'姓名'" json:"name"`
	AttendanceDate models.Date `gorm:"type:date;not null;comment:'考勤日期'" json:"attendance_date"`
	CheckTime      models.Time `gorm:"type:timestamp;comment:'打卡时间'" json:"check_time"`
	Status         string      `gorm:"size:20;comment:'Normal, Exception'" json:"status"`
	Result         string      `gorm:"size:20;comment:'异常判断'" json:"result"`
	LeaveID        string      `gorm:"size:255;comment:'休假申请记录号'" json:"leave_id"`
}

//节假日补班
type PublicHoliday struct {
	gorm.Model
	HolidayName       string      `gorm:"size:20;comment:'节假日名称'" json:"holiday_name"`
	PublicHolidayDate models.Date `gorm:"type:date;not null;comment:'补班日期'" json:"public_holiday_date"`
	HolidayType       string      `gorm:"comment:'节假日类型'" json:"holiday_type"`
	// workday weekend statutory
}

type DeptUsers struct {
	Dept  string            `json:"dept"`
	Users []*AttendanceUser `json:"users"`
}

type AttendanceUser struct {
	Dept      string `json:"dept"`
	Name      string `json:"name"`
	IsConfirm int    `gorm:"column:is_confirm" json:"is_confirm"`
}

type UserAttendanceTmp struct {
	Date string           `json:"date"`
	Tmps []*AttendanceTmp `json:"tmps"`
}

func (AttendanceTmp) TableName() string {
	return "attendance_tmp"
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
	LeaveId        string
	Shift          float32
	Overtime       float32
}

func (v AttendanceSimple) String(now string) string {
	return fmt.Sprintf("('%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s',%f,%f)", now, v.Dept, v.Name,
		v.AttendanceDate.String(), v.CheckIn.String(), v.CheckOut.String(), v.InStatus, v.OutStatus,
		v.InResult, v.OutResult, v.LeaveId, v.Overtime, v.Shift)
}

func (v AttendanceTmp) String(now string) string {
	return fmt.Sprintf("('%s','%s','%s','%s','%s','%s','%s')", now, v.Dept, v.Name, v.AttendanceDate.String(),
		v.CheckTime.String(), v.Status, v.Result)
}

type DeptUser struct {
	Dept  string            `json:"dept"`
	Users []*UserAttendance `json:"users"`
}

type UserAttendance struct {
	Dept        string        `json:"-"`
	Name        string        `json:"name"`
	Attendances []*Attendance `json:"attendances"`
}

//请假数据统计
type AttendanceExcel struct {
	Dept     string
	Name     string
	Total    int
	Leave    float32
	Annual   float32
	Sick     float32
	Late     int
	Early    int
	Overtime float32
	Shift    float32
	Affair   float32
	None     int
	Forget   int
}
