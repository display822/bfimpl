/*
* Auth : acer
* Desc : 离职流程信息
* Time : 2020/9/5 17:08
 */

package oa

import (
	"bfimpl/models"

	"github.com/jinzhu/gorm"
)

type QuitFlowInfo struct {
	gorm.Model
	EmployeeID int
	Account    string `gorm:"size:60;comment:'账号'" json:"account"`
	Computer   string `gorm:"comment:'电脑'" json:"computer"`
	Phone      string `gorm:"comment:'手机'" json:"phone"`
	Expense    string `gorm:"size:60;comment:'报销'" json:"expense"`
	DeviceReq  string `gorm:"comment:'物品领用归还'" json:"device_req"`
	WorkDay    string `gorm:"size:20;comment:'实际出勤天数'" json:"work_day"`
	OffDay     string `gorm:"size:20;comment:'旷工天数'" json:"off_day"`
	HalfDay    string `gorm:"size:20;comment:'病假'" json:"half_day"`
	ChangeDay  string `gorm:"size:20;comment:'调休'" json:"change_day"`
	Others     string `gorm:"comment:'其他结算'" json:"others"`
	LateDay    string `gorm:"size:20;comment:'迟到'" json:"late_day"`
	ThingsDay  string `gorm:"size:20;comment:'事假'" json:"things_day"`
	SalaryDay  string `gorm:"size:20;comment:'带薪假'" json:"salary_day"`
	AnnualDay  string `gorm:"size:20;comment:'年假'" json:"annual_day"`
}

type ReqQuitFlow struct {
	Account         string      `json:"account"`
	Computer        string      `json:"computer"`
	Phone           string      `json:"phone"`
	Expense         string      `json:"expense"`
	DeviceReq       string      `json:"device_req"`
	WorkDay         string      `json:"work_day"`
	OffDay          string      `json:"off_day"`
	HalfDay         string      `json:"half_day"`
	ChangeDay       string      `json:"change_day"`
	Others          string      `json:"others"`
	LateDay         string      `json:"late_day"`
	ThingsDay       string      `json:"things_day"`
	SalaryDay       string      `json:"salary_day"`
	AnnualDay       string      `json:"annual_day"`
	ResignationDate models.Time `json:"resignation_date"`
	Reason          string      `json:"reason"`
}

func (r *ReqQuitFlow) ToEntity() *QuitFlowInfo {
	return &QuitFlowInfo{
		Account:   r.Account,
		Computer:  r.Computer,
		Phone:     r.Phone,
		Expense:   r.Expense,
		DeviceReq: r.DeviceReq,
		WorkDay:   r.WorkDay,
		OffDay:    r.OffDay,
		HalfDay:   r.HalfDay,
		ChangeDay: r.ChangeDay,
		Others:    r.Others,
		LateDay:   r.LateDay,
		ThingsDay: r.ThingsDay,
		SalaryDay: r.SalaryDay,
		AnnualDay: r.AnnualDay,
	}
}
