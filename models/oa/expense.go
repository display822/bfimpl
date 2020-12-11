package oa

import (
	"bfimpl/models"

	"github.com/jinzhu/gorm"
)

var ExcelHeaderArray = []string{"费用发生日期"}

// Expense 报销表
type Expense struct {
	gorm.Model
	EngagementCode  string      `gorm:"size:64;comment:'任务指派编码'" json:"engagement_code"`
	EmpID           int         `gorm:"comment:'报销申请人id'" json:"emp_id"`
	EName           string      `gorm:"size:30;comment:'员工姓名'" json:"e_name"`
	Status          string      `gorm:"size:20;comment:'申请状态'" json:"status"`
	ExpenseSummary  int         `gorm:"comment:'费用总金额'" json:"expense_summary"`
	ApplicationDate models.Date `gorm:"type:date;comment:'申请日期'" json:"application_date"`
	PaymentDate     models.Date `gorm:"type:date;comment:'支付日期'" json:"payment_date"`
	ImportFile      string      `gorm:"size:30;comment:'导入文件'" json:"import_file"`
}

// ExpenseDetail 报销明细表
type ExpenseDetail struct {
	gorm.Model
	ExpenseID          int         `gorm:"comment:'报销id'" json:"expense_id"`
	ExpenseAccountCode string      `gorm:"size:64;comment:'报销科目编码'" json:"expense_account_code"`
	OcurredDate        models.Date `gorm:"type:date;comment:'申请日期'" json:"ocurred_date"`
	ExpenseAmount      int         `gorm:"comment:'费用金额'" json:"expense_amount"`
	Remarks1           string      `gorm:"size:30;comment:'备注1'" json:"remarks1"`
	Remarks2           string      `gorm:"size:100;comment:'备注2'" json:"remarks2"`
	Remarks3           string      `gorm:"size:100;comment:'备注3'" json:"remarks3"`
}

// ExpenseAccount 报销科目表
type ExpenseAccount struct {
	gorm.Model
	ExpenseAccountCode string `gorm:"size:64;comment:'报销科目编码'" json:"expense_account_code"`
	ExpenseAccountName string `gorm:"size:64;comment:'报销科目名称'" json:"expense_account_name"`
}
