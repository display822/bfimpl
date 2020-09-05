/*
* Auth : acer
* Desc : 工作流节点
* Time : 2020/9/1 23:03
 */

package oa

import "github.com/jinzhu/gorm"

type WorkflowNode struct {
	gorm.Model
	WorkflowID int                    `gorm:"not null;comment:'所属工作流'" json:"-"`
	NodeSeq    int                    `gorm:"not null;comment:'节点序号'" json:"node_seq"`
	OperatorID int                    `gorm:"not null;comment:'当前操作人ID,关联UserID'" json:"operator_id"`
	Status     string                 `gorm:"size:20;not null;default:'NA';comment:'NA, Processing, Completed, Approved, Rejected'" json:"status"`
	Elements   []*WorkflowFormElement `json:"elements"`
}
