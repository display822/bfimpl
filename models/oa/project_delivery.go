package oa

import (
	"time"

	"github.com/jinzhu/gorm"
)

// ProjectDelivery 项目交付表
type ProjectDelivery struct {
	gorm.Model
	ProjectID            string    `gorm:"size:64;comment:'项目ID'" json:"project_id"`
	ProjectCategoryCode  int       `gorm:"comment:'项目分类编码'" json:"emp_id"`
	EngagementCode       string    `gorm:"comment:'任务指派编码'" json:"engagement_code"`
	StartDate            time.Time `gorm:"type:date;comment:'交付周期开始日期'" json:"start_date"`
	EndDate              time.Time `gorm:"type:date;comment:'交付周期结束日期'" json:"end_date"`
	MainServiceAmount    int       `gorm:"size:20;comment:'主服务交付数量'" json:"main_service_amount"`
	SubServiceAmount     int       `gorm:"type:decimal(10,2);comment:'自服务交付数量'" json:"sub_service_amount"`
	ProjectDeliveryValue string    `gorm:"type:date;comment:'项目交付价值'" json:"project_delivery_value"`
}
