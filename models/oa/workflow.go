/*
* Auth : acer
* Desc : 工作流
* Time : 2020/9/1 22:58
 */

package oa

import "github.com/jinzhu/gorm"

type Workflow struct {
	gorm.Model
	WorkflowDefinitionID int                    `json:"-"`
	WorkflowDefinition   *WorkflowDefinition    `json:"-"`
	Status               string                 `gorm:"size:20;not null;comment:'状态(Processing, Completed, Approved, Rejected)'"`
	EntityID             int                    `gorm:"not null;comment:'实体ID'"`
	Nodes                []*WorkflowNode        `json:"nodes"`
	Elements             []*WorkflowFormElement `json:"elements"`
}
