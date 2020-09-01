/*
* Auth : acer
* Desc : 工作流，表单元素定义
* Time : 2020/9/1 23:08
 */

package oa

import "github.com/jinzhu/gorm"

type WorkflowFormElementDef struct {
	gorm.Model
	WorkflowDefinitionID int    `gorm:"not null;comment:'所属流程定义ID'"`
	ElementSeq           int    `gorm:"not null;comment:'表单元素序号'"`
	ElementType          string `gorm:"not null;comment:'TextField, TextArea'"`
	ElementName          string `gorm:"not null;comment:'表单元素名称'"`
}
