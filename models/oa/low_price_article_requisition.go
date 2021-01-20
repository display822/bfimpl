/**
* @author : yi.zhang
* @description : oa 描述
* @date   : 2021-01-11 15:41
 */

package oa

import (
	"bfimpl/models"
	"bytes"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

type LowPriceArticleRequisition struct {
	gorm.Model
	LowPriceArticleID     int              `gorm:"not null;comment:'设备ID'" json:"low_price_article_id"`
	LowPriceArticle       *LowPriceArticle `json:"low_price_article"`
	OperatorID            int              `gorm:"not null;comment:'操作人ID，关联EmployeeID'" json:"operator_id"`
	OperatorName          string           `gorm:"not null;comment:'操作人Name，关联EmployeeName'" json:"operator_name"`
	AssociateEmployeeID   int              `gorm:"not null;comment:'关联EmployeeID'" json:"associate_employee_id"`
	AssociateEmployeeName string           `gorm:"not null;comment:'关联EmployeeName'" json:"associate_employee_name"`
	OperatorCategory      string           `gorm:"size:20;not null;comment:'操作类别(入库,报废,借出,归还)'" json:"operator_category"`
	Quantity              int              `gorm:"not null;comment:'数量'" json:"quantity"`
	Comment               string           `gorm:"size:2000;not null;comment:'备注'" json:"comment"`
	IsReturn              int              `gorm:"not null;comment:'是否归还'" json:"is_return"`
}

// BatchRequisitionOutGoing 批量插入数据
func BatchRequisitionOutGoing(db *gorm.DB, lpars []*LowPriceArticleRequisition) error {
	var buffer bytes.Buffer
	sql := "insert into `low_price_article_requisitions` (`created_at`,`updated_at`,`low_price_article_id`,`operator_id`," +
		"`operator_name`, `associate_employee_id`,`associate_employee_name`, `operator_category`,`quantity`,`comment`) values"
	if _, err := buffer.WriteString(sql); err != nil {
		return err
	}
	for i, lpar := range lpars {
		if i == len(lpars)-1 {
			buffer.WriteString(fmt.Sprintf("('%s','%s',%d, %d, '%s', %d,'%s','%s',%d,'%s');", time.Now(),
				time.Now(), lpar.LowPriceArticleID, lpar.OperatorID, lpar.OperatorName, lpar.AssociateEmployeeID, lpar.AssociateEmployeeName,
				models.DeviceOutgoing, 1, lpar.Comment,
			))
		} else {
			buffer.WriteString(fmt.Sprintf("('%s','%s',%d, %d, '%s', %d,'%s','%s',%d,'%s'),", time.Now(),
				time.Now(), lpar.LowPriceArticleID, lpar.OperatorID, lpar.OperatorName, lpar.AssociateEmployeeID, lpar.AssociateEmployeeName,
				models.DeviceOutgoing, 1, lpar.Comment,
			))
		}
	}
	return db.Exec(buffer.String()).Error
}

// BatchRequisitionReturn 批量插入数据
func BatchRequisitionReturn(db *gorm.DB, operatorID int, operatorName string, lpars []*LowPriceArticleRequisition) error {
	var buffer bytes.Buffer
	sql := "insert into `low_price_article_requisitions` (`created_at`,`updated_at`,`low_price_article_id`,`operator_id`," +
		"`operator_name`, `associate_employee_id`,`associate_employee_name`, `operator_category`,`quantity`,`comment`) values"
	if _, err := buffer.WriteString(sql); err != nil {
		return err
	}
	for i, lpar := range lpars {
		if i == len(lpars)-1 {
			buffer.WriteString(fmt.Sprintf("('%s','%s',%d, %d, '%s', %d,'%s','%s',%d,'%s');", time.Now(),
				time.Now(), lpar.LowPriceArticleID, operatorID, operatorName, lpar.AssociateEmployeeID, lpar.AssociateEmployeeName,
				models.DeviceReturn, 1, "",
			))
		} else {
			buffer.WriteString(fmt.Sprintf("('%s','%s',%d, %d, '%s', %d,'%s','%s',%d,'%s'),", time.Now(),
				time.Now(), lpar.LowPriceArticleID, operatorID, operatorName, lpar.AssociateEmployeeID, lpar.AssociateEmployeeName,
				models.DeviceReturn, 1, "",
			))
		}
	}
	return db.Exec(buffer.String()).Error
}
