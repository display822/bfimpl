/*
* Auth : acer
* Desc : 工作流表单元素值
* Time : 2020/9/1 23:08
 */

package oa

import "github.com/jinzhu/gorm"

type WorkflowFormElement struct {
	gorm.Model
	WfElementDefID         int
	WorkflowFormElementDef *WorkflowFormElementDef `gorm:"ForeignKey:WfElementDefID"`
	Value                  string                  `gorm:"not null;comment:'表单元素值'"`
	WorkflowID             int
}
