/**
* @author : chen lie
* @description : 财务code
* @date   : 2020-09-30 15:58
 */

package oa

import (
	"bfimpl/models"

	"github.com/jinzhu/gorm"
)

type EngagementCode struct {
	gorm.Model
	EngagementCode     string       `gorm:"size:10;comment:'任务指派编码'" json:"engagement_code"`
	EngagementCodeDesc string       `gorm:"comment:'任务指派编码说明'" json:"engagement_code_desc"`
	Category           string       `gorm:"comment:'Deliverable,NonDeliverable,Operation'" json:"category"`
	CCRate             float32      `gorm:"type:decimal(5,2);comment:''" json:"cc_rate"`
	OCRate             float32      `gorm:"type:decimal(5,2);comment:''" json:"oc_rate"`
	DepartmentID       int          `gorm:"comment:'部门id'" json:"department_id"`
	CodeOwnerID        int          `gorm:"comment:'负责人id'" json:"code_owner_id"`
	Owner              *models.User `gorm:"ForeignKey:CodeOwnerID" json:"owner"`
	TeamLeaderID       int          `gorm:"comment:'组长'" json:"team_leader_id"`
	HRID               int          `gorm:"comment:'人事'" json:"hr_id"`
	FinanceID          int          `gorm:"comment:'财务'" json:"finance_id"`
}
