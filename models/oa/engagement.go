/**
* @author : yi.zhang
* @description : oa 描述
* @date   : 2020-12-29 17:08
 */

package oa

import (
	"bytes"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

// Engagement 人员管理
type Engagement struct {
	gorm.Model
	EngagementCode string    `gorm:"size:10;comment:'任务指派编码'" json:"engagement_code"`
	EmployeeID     int       `gorm:"comment:'员工id'" json:"employee_id"`
	EmployeeName   string    `gorm:"size:30;comment:'员工姓名'" json:"employee_name"`
	EngagementDate time.Time `gorm:"type:date;comment:'任务日期'" json:"engagement_date"`
	EngagementHour int       `gorm:"comment:'任务耗时'" json:"engagement_hour"`
	EngagementCost float64   `gorm:"type:decimal(10,2);comment:'任务成本'" json:"engagement_cost"`
}

// BatchEngagementCreate 批量插入数据
func BatchEngagementCreate(db *gorm.DB, es []*Engagement) error {
	var buffer bytes.Buffer
	sql := "insert into `engagements` (`created_at`,`updated_at`,`engagement_code`,`employee_id`,`employee_name`,`engagement_date`," +
		"`engagement_hour`,`engagement_cost`) values"
	if _, err := buffer.WriteString(sql); err != nil {
		return err
	}
	for i, pd := range es {
		if i == len(es)-1 {
			buffer.WriteString(fmt.Sprintf("('%s','%s','%s',%d,'%s','%s',%d,%.2f);", time.Now(),
				time.Now(), pd.EngagementCode, pd.EmployeeID, pd.EmployeeName, pd.EngagementDate, pd.EngagementHour,
				pd.EngagementCost,
			))
		} else {
			buffer.WriteString(fmt.Sprintf("('%s','%s','%s',%d,'%s','%s',%d,%.2f),", time.Now(),
				time.Now(), pd.EngagementCode, pd.EmployeeID, pd.EmployeeName, pd.EngagementDate, pd.EngagementHour,
				pd.EngagementCost,
			))
		}
	}
	return db.Exec(buffer.String()).Error
}
