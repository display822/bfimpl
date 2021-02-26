/*
* Auth : acer
* Desc : 员工合同表
* Time : 2020/9/1 21:27
 */

package oa

import (
	"bfimpl/models"

	"github.com/jinzhu/gorm"
)

type EmployeeContract struct {
	gorm.Model
	EmployeeID        int         `json:"employee_id"`
	ContractType      string      `gorm:"size:20;not null;comment:'(劳动合同,劳务合同)'" json:"contract_type"`
	ContractParty     string      `gorm:"size:100;not null;comment:'签约方'" json:"contract_party"`
	ContractMain      string      `gorm:"size:100;not null;comment:'签约主体'" json:"contract_main"`
	ContractStartDate models.Time `gorm:"type:datetime;not null;comment:'合同开始日期'" json:"contract_start_date"`
	ContractEndDate   models.Time `gorm:"type:datetime;not null;comment:'合同结束日期'" json:"contract_end_date"`
	TrialPeriod       int         `gorm:"not null;comment:'试用期(2,6,0)'" json:"trial_period"`
	AnnualLeave       int         `gorm:"not null;comment:'年假'" json:"annual_leave"`
	Status            string      `gorm:"size:30;not null;comment:'Unsigned,Signed'" json:"status"`
	SoftCopy          string      `gorm:"size:200;comment:'合同电子档'" json:"soft_copy"`
	ScannedCopy       string      `gorm:"size:200;comment:'合同扫描件'" json:"scanned_copy"`
}

type ContractSimple struct {
	EndDate       models.Time `gorm:"column:enddate" json:"enddate"`
	EmployeeID    int         `json:"employee_id"`
	ContractMain  string      `json:"contract_main"`
	ContractParty string      `json:"contract_party"`
	ContractType  string      `json:"contract_type"`
}
