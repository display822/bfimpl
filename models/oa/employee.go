/*
* Auth : acer
* Desc : 员工主表
* Time : 2020/9/1 21:27
 */

package oa

import (
	"bfimpl/models"

	"time"

	"github.com/jinzhu/gorm"
)

type Employee struct {
	gorm.Model
	Name             string         `gorm:"size:60;not null;default:'';comment:'姓名'" json:"name"`
	Gender           string         `gorm:"size:10;not null;default:'';comment:'性别'" json:"gender"`
	Status           int            `gorm:"default:0;comment:'未入职，拟入职，已入职，已离职, 已解约'" json:"status"`
	Mobile           string         `gorm:"size:20;not null;default:'';comment:'手机'" json:"mobile"`
	EmpNo            string         `gorm:"size:20;not null;default:'';comment:'员工编号'" json:"emp_no"`
	IDCard           string         `gorm:"size:20;not null;default:'';comment:'身份证'" json:"id_card"`
	Nation           string         `gorm:"size:20;not null;default:'';comment:'民族'" json:"nation"`
	Age              int            `gorm:"not null;default:0;comment:'年龄'" json:"age"`
	PoliticStatus    string         `gorm:"size:20;not null;default:'';comment:'政治面貌'" json:"politic_status"`
	InterviewComment string         `gorm:"not null;default:'';comment:'面试评价'" json:"interview_comment"`
	Resume           string         `gorm:"not null;default:'';comment:'简历地址'" json:"resume"`
	Email            string         `gorm:"unique_index;size:60;not null;default:'';comment:'邮箱'" json:"email"`
	PersonalEmail    string         `gorm:"size:60;not null;default:'';comment:'私人邮箱'" json:"personal_email"`
	WxWork           string         `gorm:"size:60;not null;default:'';comment:'企业微信'" json:"wx_work"`
	Tapd             string         `gorm:"size:60;not null;default:'';comment:'tapd账号'" json:"tapd"`
	ServiceLine      string         `gorm:"size:60;not null;default:'';comment:'服务线'" json:"service_line"`
	DepartmentID     int            `gorm:"not null;default:0;comment:'部门'" json:"department_id"`
	Department       *Department    `json:"department"`
	LevelID          int            `gorm:"not null;default:0;comment:'级别'" json:"level_id"`
	Level            *Level         `json:"level"`
	CreateTime       models.Time    `gorm:"type:datetime" json:"create_time"`
	CreatorId        int            `gorm:"comment:'创建人id'" json:"-"`
	Annual           int            `gorm:"comment:'年假'" json:"annual"`
	Position         string         `gorm:"size:30;default:'';comment:'岗位'" json:"position"`
	PlanDate         models.Time    `gorm:"type:datetime;comment:'计划入职日期'" json:"plan_date"`
	EntryDate        models.Time    `gorm:"type:datetime;comment:'入职日期'" json:"entry_date"`
	ResignationDate  models.Time    `gorm:"type:datetime;comment:'离职日期'" json:"resignation_date"`
	ConfirmDate      models.Time    `gorm:"type:datetime;comment:'转正日期'" json:"confirm_date"`
	Reason           string         `gorm:"not null;default:'';comment:'离职原因'" json:"reason"`
	ReqUser          string         `gorm:"not null;default:'';comment:'离职申请人'" json:"req_user"`
	EmployeeBasic    *EmployeeBasic `json:"employee_basic"`
	//Contracts        []*EmployeeContract `json:"contracts"`
}

type EIDStatus struct {
	EntityID int    `gorm:"column:entity_id"`
	Status   string `gorm:"column:status"`
}

type SimpleEmp struct {
	ID   int
	Name string
}

type EmpPos struct {
	Name string `gorm:"column:name"`
	Tapd string `gorm:"column:tapd"`
}

type ReqEmployee struct {
	Name             string      `json:"name"`
	Gender           string      `json:"gender"`
	Status           int         `json:"status"`
	Mobile           string      `json:"mobile"`
	IDCard           string      `json:"id_card"`
	InterviewComment string      `json:"interview_comment"`
	Resume           string      `json:"resume"`
	Email            string      `json:"email"`
	WxWork           string      `json:"wx_work"`
	Tapd             string      `json:"tapd"`
	ServiceLine      string      `json:"service_line"`
	DepartmentID     int         `json:"department_id"`
	LeaderID         int         `json:"leader_id"`
	LevelID          int         `json:"level_id"`
	Position         string      `json:"position"`
	EntryDate        models.Time `json:"entry_date"`
	PlanDate         models.Time `json:"plan_date"`
	SeatNumber       string      `json:"seat_number"`
	DeviceReq        string      `json:"device_req"`
}

func (r *ReqEmployee) ToEmployee() *Employee {
	return &Employee{
		Name:             r.Name,
		Gender:           r.Gender,
		Status:           r.Status,
		Mobile:           r.Mobile,
		IDCard:           r.IDCard,
		InterviewComment: r.InterviewComment,
		Resume:           r.Resume,
		Email:            r.Email,
		WxWork:           r.WxWork,
		Tapd:             r.Tapd,
		ServiceLine:      r.ServiceLine,
		DepartmentID:     r.DepartmentID,
		LevelID:          r.LevelID,
		Position:         r.Position,
		PlanDate:         r.PlanDate,
		CreateTime:       models.Time(time.Now()),
	}
}

type SocialSecurity struct {
	gorm.Model
	Name        string `gorm:"size:60;not null;default:'-';comment:'文件名'" json:"name"`
	DownloadUrl string `gorm:"size:200;not null;default:'';comment:'下载地址'" json:"download_url"`
}
