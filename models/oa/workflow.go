/*
* Auth : acer
* Desc : 工作流
* Time : 2020/9/1 22:58
 */

package oa

import (
	"bfimpl/models"

	"github.com/jinzhu/gorm"
)

type Workflow struct {
	gorm.Model
	WorkflowDefinitionID int                    `json:"-"`
	WorkflowDefinition   *WorkflowDefinition    `json:"definition"`
	Status               string                 `gorm:"size:20;not null;comment:'状态(Processing, Completed, Approved, Rejected)'"`
	EntityID             int                    `gorm:"not null;comment:'实体ID'"`
	Nodes                []*WorkflowNode        `json:"nodes"`
	Elements             []*WorkflowFormElement `json:"elements"`
}

//入职流程信息
type ReqEntryFlow struct {
	Email      string      `json:"email"`
	WxWork     string      `json:"wx_work"`
	Tapd       string      `json:"tapd"`
	PlanTime   models.Time `json:"plan_date"`
	SeatNumber string      `json:"seat_number"`
	DeviceReq  string      `json:"device_req"`
}

type WorkflowId struct {
	Id int `json:"id"`
}
