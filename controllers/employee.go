/*
* Auth : acer
* Desc : 入职离职流程
* Time : 2020/9/4 21:45
 */

package controllers

import (
	"bfimpl/models"
	"bfimpl/models/oa"
	"bfimpl/services"
	"bfimpl/services/log"
	"encoding/json"
	"time"
)

type EmployeeController struct {
	BaseController
}

// @Title hr新建入职
// @Description 新建入职
// @Param	json	body	string	true	"入职员工信息"
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /new [post]
func (e *EmployeeController) NewEmpEntry() {
	userType, _ := e.GetInt("userType", 0)
	if userType != UserHR {
		e.ErrorOK("您不是HR")
	}
	reqEmployee := new(oa.ReqEmployee)
	err := json.Unmarshal(e.Ctx.Input.RequestBody, reqEmployee)
	if err != nil {
		log.GLogger.Error("new employee err：%s", err.Error())
		e.ErrorOK(MsgInvalidParam)
	}

	tx := services.Slave().Begin()
	employee := reqEmployee.ToEmployee()
	err = tx.Create(employee).Error
	if err != nil {
		log.GLogger.Error("create employee err：%s", err.Error())
		tx.Rollback()
		e.ErrorOK(MsgServerErr)
	}
	//创建流程信息
	uID, _ := e.GetInt("userID", 0)
	err = services.CreateEntryWorkflow(tx, int(employee.ID), uID, reqEmployee)
	if err != nil {
		log.GLogger.Error("create entry workflow err：%s", err.Error())
		tx.Rollback()
		e.ErrorOK(MsgServerErr)
	}
	tx.Commit()
	e.Correct(employee)
}

// @Title 入职详情
// @Description 入职详情
// @Param	id	path	int	true	"入职员工id"
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /entry/detail/:id [get]
func (e *EmployeeController) GetEmpEntryInfo() {
	eID, _ := e.GetInt(":id", 0)
	employee := new(oa.Employee)
	services.Slave().Model(oa.Employee{}).Where("id = ?", eID).Preload("Department").
		Preload("Department.Leader").Preload("Level").First(&employee)
	e.Correct(employee)
}

// @Title 入职列表
// @Description 入职列表
// @Param	pagenum	query	int	true	"页数"
// @Param	pagesize	query	int	true	"页大小"
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /entry/list [get]
func (e *EmployeeController) GetEmpEntryList() {

	pageNum, _ := e.GetInt("pagenum", 1)
	pageSize, _ := e.GetInt("pagesize", 10)
	name := e.GetString("name")
	dID, _ := e.GetInt("departmentid", 0)
	status, _ := e.GetInt("status", -1)
	employees := make([]*oa.Employee, 0)
	//未入职和拟入职
	var resp struct {
		Total int            `json:"total"`
		List  []*oa.Employee `json:"list"`
	}
	query := services.Slave().Model(oa.Employee{})
	if dID != 0 {
		query = query.Where("department_id = ?", dID)
	}
	if status != -1 {
		query = query.Where("status = ?", status)
	} else {
		query = query.Where("status in (?)", []int{0, 1})
	}
	if name != "" {
		query = query.Where("name like ?", "%"+name+"%")
	}
	query.Preload("Department").Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&employees).
		Limit(-1).Offset(-1).Count(&resp.Total)
	resp.List = employees
	e.Correct(resp)
}

// @Title 流程信息
// @Description 流程信息
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /workflow/:id [get]
func (e *EmployeeController) GetWorkflowNode() {
	eID, _ := e.GetInt(":id", 0)
	workflow := new(oa.Workflow)
	services.Slave().Model(oa.Workflow{}).Where("workflow_definition_id = ? and entity_id = ?",
		services.GetEntryDef(), eID).Preload("Nodes").Preload("Nodes.User").
		Preload("Elements").Preload("Elements.WorkflowFormElementDef").First(&workflow)

	e.Correct(workflow)
}

// @Title 提交入职流程信息
// @Description 提交入职流程信息
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /workflow/:id [put]
func (e *EmployeeController) CommitWorkflowNode() {
	eID, _ := e.GetInt(":id", 0)
	flowInfo := new(oa.ReqEntryFlow)
	err := json.Unmarshal(e.Ctx.Input.RequestBody, flowInfo)
	if err != nil {
		e.ErrorOK(MsgInvalidParam)
	}
	//修改入职流程信息
	userType, _ := e.GetInt("userType", 0)
	//eID 查询 workflow
	workflow := new(oa.Workflow)
	services.Slave().Model(oa.Workflow{}).Where("workflow_definition_id = ? and entity_id = ?",
		services.GetEntryDef(), eID).Preload("Nodes").Preload("Elements").First(workflow)
	if len(workflow.Nodes) != 3 || len(workflow.Elements) != 3 {
		e.ErrorOK("工作流配置错误")
	}
	workflow.Elements[0].Value = (time.Time(flowInfo.PlanTime)).Format(models.TimeFormat)
	workflow.Elements[1].Value = flowInfo.SeatNumber
	workflow.Elements[2].Value = flowInfo.DeviceReq
	switch userType {
	case UserHR:
		e.ErrorOK("工作流不在当前节点")
	case UserLeader:
		if workflow.Nodes[1].Status != services.FlowProcessing {
			e.ErrorOK("工作流不在当前节点")
		}
		//修改employee信息
		services.Slave().Model(oa.Employee{}).Where("id = ?", eID).Updates(map[string]interface{}{
			"email":      flowInfo.Email,
			"wx_work":    flowInfo.WxWork,
			"tapd":       flowInfo.Tapd,
			"entry_date": workflow.Elements[0].Value,
		})
		//更新节点信息
		workflow.Nodes[1].Status = services.FlowCompleted
		workflow.Nodes[2].Status = services.FlowProcessing
		services.Slave().Save(workflow)
	case UserIT:
		if workflow.Nodes[2].Status != services.FlowProcessing {
			e.ErrorOK("工作流不在当前节点")
		}
		//修改employee信息,最后一个流程，变为已入职
		services.Slave().Model(oa.Employee{}).Where("id = ?", eID).Updates(map[string]interface{}{
			"email":      flowInfo.Email,
			"wx_work":    flowInfo.WxWork,
			"tapd":       flowInfo.Tapd,
			"entry_date": workflow.Elements[0].Value,
			"status":     2,
		})
		workflow.Nodes[2].Status = services.FlowCompleted
		workflow.Status = services.FlowCompleted
		services.Slave().Save(workflow)
	}
	e.Correct("")
}
