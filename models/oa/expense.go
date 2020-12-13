package oa

import (
	"bfimpl/models"
	"bfimpl/services/log"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

var ExcelHeaderArray = []string{"费用发生日期", "费用科目", "费用金额", "备注1", "备注2", "备注3"}

var ExpenseAccountMap = map[string]string{
	"餐补费":     "10001",
	"交通费(市内)": "10002",
	"团队激励":    "10003",
	"活动费":     "10004",
	"办公费":     "10005",
	"招聘费":     "10006",
	"通讯费":     "10007",
	"销售费用":    "10008",
	"充值费用":    "10009",
	"交通费(市外)": "10010",
	"住宿费":     "10011",
	"出差补贴":    "10012",
}

// Expense 报销表
type Expense struct {
	gorm.Model
	EngagementCode  string          `gorm:"size:64;comment:'任务指派编码'" json:"engagement_code"`
	EmpID           int             `gorm:"comment:'报销申请人id'" json:"emp_id"`
	EName           string          `gorm:"size:30;comment:'员工姓名'" json:"e_name"`
	Project         string          `gorm:"size:64;comment:'项目'" json:"project"`
	Status          string          `gorm:"size:20;comment:'申请状态'" json:"status"`
	ExpenseSummary  float64         `gorm:"type:decimal(5,2);comment:'费用总金额'" json:"expense_summary"`
	ApplicationDate time.Time       `gorm:"type:date;comment:'申请日期'" json:"application_date"`
	PaymentDate     *time.Time      `gorm:"type:date;comment:'支付日期'" json:"payment_date"`
	ImportFile      string          `gorm:"size:30;comment:'导入文件'" json:"import_file"`
	LeaderId        int             `gorm:"-" json:"leader_id"`
	ExpenseDetails  []ExpenseDetail `gorm:"ForeignKey:ExpenseID" json:"expense_details"`
}

// ExpenseDetail 报销明细表
type ExpenseDetail struct {
	gorm.Model
	ExpenseID          int             `gorm:"comment:'报销id'" json:"expense_id"`
	ExpenseAccountCode string          `gorm:"size:64;comment:'报销科目编码'" json:"expense_account_code"`
	ExpenseAccount     *ExpenseAccount `json:"expense_account"`
	OcurredDate        models.Date     `gorm:"type:date;comment:'发生日期'" json:"ocurred_date"`
	ExpenseAmount      float64         `gorm:"type:decimal(5,2);comment:'费用金额'" json:"expense_amount"`
	Remarks1           string          `gorm:"size:30;comment:'备注1'" json:"remarks1"`
	Remarks2           string          `gorm:"size:100;comment:'备注2'" json:"remarks2"`
	Remarks3           string          `gorm:"size:100;comment:'备注3'" json:"remarks3"`
}

// ExpenseAccount 报销科目表
type ExpenseAccount struct {
	gorm.Model
	ExpenseAccountCode string `gorm:"size:64;comment:'报销科目编码'" json:"expense_account_code"`
	ExpenseAccountName string `gorm:"size:64;comment:'报销科目名称'" json:"expense_account_name"`
}

// BatchCreateExpenseDetail 批量创建报销明细
func BatchCreateExpenseDetail(db *gorm.DB, expenseID int, details []ExpenseDetail) error {
	log.GLogger.Info("ExpenseDetail: %+v", details[1])
	sql := "INSERT INTO `expense_details` (`created_at`,`updated_at`,`expense_id`,`expense_account_code`,`ocurred_date`,`expense_amount`,`remarks1`" +
		",`remarks2`,`remarks3`) VALUES "
	// 循环data数组,组合sql语句
	for key, value := range details {
		if len(details)-1 == key {
			//最后一条数据 以分号结尾
			sql += fmt.Sprintf("('%s','%s',%d,'%s','%s',%.2f,'%s','%s','%s');", time.Now().Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05"), expenseID, value.ExpenseAccountCode, value.OcurredDate,
				value.ExpenseAmount, value.Remarks1, value.Remarks2, value.Remarks3)
		} else {
			sql += fmt.Sprintf("('%s','%s',%d,'%s','%s',%.2f,'%s','%s','%s'),", time.Now().Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05"), expenseID, value.ExpenseAccountCode, value.OcurredDate,
				value.ExpenseAmount, value.Remarks1, value.Remarks2, value.Remarks3)
		}
	}
	log.GLogger.Info("sql: %s", sql)

	return db.Exec(sql).Error
}
