/*
* Auth : acer
* Desc : 工作流表单元素值
* Time : 2020/9/1 23:08
 */

package oa

import "github.com/jinzhu/gorm"

type WorkflowFormElement struct {
	gorm.Model
	WfElementDefID         int                     `json:"-"`
	WorkflowFormElementDef *WorkflowFormElementDef `gorm:"ForeignKey:WfElementDefID" json:"elementDef"`
	Name                   string                  `gorm:"size:60;not null;default:'';comment:'表单元素名称'" json:"name"`
	Value                  string                  `gorm:"not null;comment:'表单元素值'" json:"value"`
	WorkflowID             uint                    `json:"-"`
}
