/*
* Auth : acer
* Desc : 工作流定义
* Time : 2020/9/1 22:53
 */

package oa

import "github.com/jinzhu/gorm"

type WorkflowDefinition struct {
	gorm.Model
	WorkflowType    string `gorm:"size:20;not null;comment:'流程类别(Approval, Business)'" json:"workflow_type"`
	WorkflowPurpose string `gorm:"size:30;not null;comment:'工作流目的(EmployeeEntry, Overtime, Leave, Expense, DeviceRequisition)'" json:"workflow_purpose"`
	WorkflowEntity  string `gorm:"size:30;not null;comment:'工作流所用实体(Employee, DeviceRequisition, Overtime, Leave, Expense)'" json:"workflow_entity"`
}
