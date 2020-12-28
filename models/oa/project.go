package oa

import (
	"bytes"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

var ProjectExcelHeaderArray = []string{"项目名称", "项目编码", "主服务交付数量", "子服务交付数量"}

// ProjectDelivery 项目表
type Project struct {
	gorm.Model
	ProjectName string `gorm:"size:64;comment:'项目名称'" json:"project_name"`
}

type ProjectCategory struct {
	gorm.Model
	ProjectCategoryCode  string  `gorm:"comment:'项目分类编码'" json:"project_category_code"`
	ServiceCategoryDesc  string  `gorm:"comment:'项目分类描述'" json:"service_category_desc"`
	MainServiceQuotation float64 `gorm:"type:decimal(10,2);;comment:'主服务报价'" json:"main_service_quotation"`
	SubServiceQuotation  float64 `gorm:"type:decimal(10,2);;comment:'子服务报价'" json:"sub_service_quotation"`
}

// ProjectDelivery 项目交付表
type ProjectDelivery struct {
	gorm.Model
	ProjectID            int        `gorm:"size:64;comment:'项目ID'" json:"project_id"`
	ProjectCategoryCode  string     `gorm:"comment:'项目分类编码'" json:"project_category_code"`
	EngagementCode       string     `gorm:"comment:'任务指派编码'" json:"engagement_code"`
	StartDate            *time.Time `gorm:"type:date;comment:'交付周期开始日期'" json:"start_date"`
	EndDate              *time.Time `gorm:"type:date;comment:'交付周期结束日期'" json:"end_date"`
	MainServiceAmount    int        `gorm:"comment:'主服务交付数量'" json:"main_service_amount"`
	SubServiceAmount     int        `gorm:"comment:'子服务交付数量'" json:"sub_service_amount"`
	ProjectDeliveryValue float64    `gorm:"type:decimal(10,2);;comment:'项目交付价值'" json:"project_delivery_value"`
}

// BatchProjectCreate 批量插入数据
func BatchProjectCreate(db *gorm.DB, pds []*ProjectDelivery) error {
	var buffer bytes.Buffer
	sql := "insert into `employees` (`created_at`,`updated_at`,`project_id`,`project_category_code`,`engagement_code`," +
		"`start_date`,`end_date`,`main_service_amount`,`sub_service_amount`,`project_delivery_value`) values"
	if _, err := buffer.WriteString(sql); err != nil {
		return err
	}
	for i, _ := range pds {
		if i == len(pds)-1 {
			buffer.WriteString(fmt.Sprintf("('%s','%s','%d','%s','%s','%s','%s',%d,%d,%f);"))
		} else {
			buffer.WriteString(fmt.Sprintf("('%s','%s','%d','%s','%s','%s','%s',%d,%d,%f),"))
		}
	}
	return db.Exec(buffer.String()).Error
}
