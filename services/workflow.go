/*
* Auth : acer
* Desc : 工作流
* Time : 2020/9/5 9:42
 */

package services

import (
	"bfimpl/models/oa"

	"bfimpl/models"

	"errors"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

const (
	EmployeeEntry  = iota
	FlowNA         = "NA"
	FlowProcessing = "Processing"
	FlowCompleted  = "Completed"
	FlowApproved   = "Approved"
	FlowRejected   = "Rejected"
)

var WorkFlowDef map[int]int
var itUserID int

func init() {
	WorkFlowDef = make(map[int]int)
	var tmp struct {
		ID int `json:"id"`
	}
	//查询入职流程定义id
	Slave().Raw("select id from workflow_definitions where workflow_purpose = ? ",
		"EmployeeEntry").Scan(&tmp)
	WorkFlowDef[EmployeeEntry] = tmp.ID

	//查询it user_id
	var u models.User
	Slave().Model(models.User{}).Where("user_type = ?", 7).First(&u)
	itUserID = int(u.ID)
}

func CreateEntryWorkflow(db *gorm.DB, eID, uID int, reqEmployee *oa.ReqEmployee) error {
	//工作流
	workflow := oa.Workflow{
		WorkflowDefinitionID: WorkFlowDef[EmployeeEntry],
		Status:               oa.Processing,
		EntityID:             eID,
	}
	err := db.Create(&workflow).Error
	if err != nil {
		return err
	}
	//查询elements
	eleDef := make([]*oa.WorkflowFormElementDef, 0)
	db.Model(oa.WorkflowFormElementDef{}).Where("workflow_definition_id = ?", WorkFlowDef[EmployeeEntry]).Find(&eleDef)
	if len(eleDef) != 3 {
		return errors.New("wrong workflow elements")
	}
	//元素填写
	ele1 := oa.WorkflowFormElement{
		WfElementDefID: int(eleDef[0].ID),
		Value:          reqEmployee.EntryDate.String(),
		WorkflowID:     workflow.ID,
	}
	ele2 := oa.WorkflowFormElement{
		WfElementDefID: int(eleDef[1].ID),
		Value:          reqEmployee.SeatNumber,
		WorkflowID:     workflow.ID,
	}
	ele3 := oa.WorkflowFormElement{
		WfElementDefID: int(eleDef[2].ID),
		Value:          reqEmployee.DeviceReq,
		WorkflowID:     workflow.ID,
	}
	err = db.Create(&ele1).Error
	err = db.Create(&ele2).Error
	err = db.Create(&ele3).Error
	if err != nil {
		return err
	}
	//节点1，hr录入
	nodeHr := oa.WorkflowNode{
		WorkflowID: int(workflow.ID),
		NodeSeq:    1,
		OperatorID: uID,
		Status:     FlowCompleted,
	}
	err = db.Create(&nodeHr).Error
	if err != nil {
		return err
	}
	//节点2，部门负责人填写
	nodeLeader := oa.WorkflowNode{
		WorkflowID: int(workflow.ID),
		NodeSeq:    2,
		OperatorID: reqEmployee.LeaderID,
	}
	err = db.Create(&nodeLeader).Error
	if err != nil {
		return err
	}
	//节点3，it填写
	nodeIT := oa.WorkflowNode{
		WorkflowID: int(workflow.ID),
		NodeSeq:    3,
		OperatorID: itUserID,
	}
	err = db.Create(&nodeIT).Error
	if err != nil {
		return err
	}

	return nil
}
