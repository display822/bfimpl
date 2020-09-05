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
	EmployeeEntry  = "EmployeeEntry"
	EmployeeLeave  = "EmployeeLeave"
	FlowNA         = "NA"
	FlowProcessing = "Processing"
	FlowCompleted  = "Completed"
	FlowApproved   = "Approved"
	FlowRejected   = "Rejected"
)

var WorkFlowDef map[string]int
var itUserID int
var financeUserID int
var frontUserID int

func init() {
	WorkFlowDef = make(map[string]int)
	//查询流程定义id
	flowDefs := make([]*oa.WorkflowDefinition, 0)
	Slave().Model(oa.WorkflowDefinition{}).Find(&flowDefs)
	for _, def := range flowDefs {
		WorkFlowDef[def.WorkflowPurpose] = int(def.ID)
	}
	//查询hr6  it7  caiwu8
	users := make([]*models.User, 0)
	Slave().Model(models.User{}).Where("user_type in (?)", []int{6, 7, 8}).Find(&users)
	for _, u := range users {
		if u.UserType == models.UserIT {
			itUserID = int(u.ID)
		} else if u.UserType == models.UserFinance {
			financeUserID = int(u.ID)
		} else if u.UserType == models.UserFront {
			frontUserID = int(u.ID)
		}
	}
}

func GetEntryDef() int {
	return WorkFlowDef[EmployeeEntry]
}

func GetLeaveDef() int {
	return WorkFlowDef[EmployeeLeave]
}

//入职流程工作流
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
		Name:           eleDef[0].ElementName,
		Value:          reqEmployee.EntryDate.String(),
		WorkflowID:     workflow.ID,
	}
	ele2 := oa.WorkflowFormElement{
		WfElementDefID: int(eleDef[1].ID),
		Name:           eleDef[1].ElementName,
		Value:          reqEmployee.SeatNumber,
		WorkflowID:     workflow.ID,
	}
	ele3 := oa.WorkflowFormElement{
		WfElementDefID: int(eleDef[2].ID),
		Name:           eleDef[2].ElementName,
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
		Status:     FlowProcessing,
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

func CreateLeaveWorkflow(db *gorm.DB, eID, uID int) error {
	//工作流
	workflow := oa.Workflow{
		WorkflowDefinitionID: WorkFlowDef[EmployeeLeave],
		Status:               oa.Processing,
		EntityID:             eID,
	}
	err := db.Create(&workflow).Error
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
	//节点2，it填写
	nodeIT := oa.WorkflowNode{
		WorkflowID: int(workflow.ID),
		NodeSeq:    2,
		OperatorID: itUserID,
		Status:     FlowProcessing,
	}
	err = db.Create(&nodeIT).Error
	if err != nil {
		return err
	}
	//节点3，财务填写
	nodeFinance := oa.WorkflowNode{
		WorkflowID: int(workflow.ID),
		NodeSeq:    3,
		OperatorID: financeUserID,
	}
	err = db.Create(&nodeFinance).Error
	if err != nil {
		return err
	}
	//节点3，财务填写
	nodeFront := oa.WorkflowNode{
		WorkflowID: int(workflow.ID),
		NodeSeq:    4,
		OperatorID: frontUserID,
	}
	err = db.Create(&nodeFront).Error
	if err != nil {
		return err
	}
	return nil
}
