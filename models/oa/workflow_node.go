/*
* Auth : acer
* Desc : 工作流节点
* Time : 2020/9/1 23:03
 */

package oa

import "github.com/jinzhu/gorm"

type WorkflowNode struct {
	gorm.Model
	WorkflowID int    `gorm:"not null;comment:'所属工作流'"`
	NodeSeq    int    `gorm:"not null;comment:'节点序号'"`
	OperatorID int    `gorm:"not null;comment:'当前操作人ID,关联EmployeeID'"`
	Status     string `gorm:"size:20;not null;comment:'NA, Processing, Completed, Approved, Rejected'"`
}
