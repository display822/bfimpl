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
	if userType != models.UserHR {
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

// @Title employee详情
// @Description employee详情
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
	number := e.GetString("emp_no")
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
	if number != "" {
		//员工编号
		query = query.Where("emp_no like ?", "%"+number+"%")
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
	case models.UserHR:
		e.ErrorOK("工作流不在当前节点")
	case models.UserLeader:
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
	case models.UserIT:
		if workflow.Nodes[2].Status != services.FlowProcessing {
			e.ErrorOK("工作流不在当前节点")
		}
		//修改employee信息,最后一个流程，变为已入职
		services.Slave().Model(oa.Employee{}).Where("id = ?", eID).Updates(map[string]interface{}{
			"email":     flowInfo.Email,
			"wx_work":   flowInfo.WxWork,
			"tapd":      flowInfo.Tapd,
			"plan_date": workflow.Elements[0].Value,
			"status":    2,
		})
		workflow.Nodes[2].Status = services.FlowCompleted
		workflow.Status = services.FlowCompleted
		services.Slave().Save(workflow)
	}
	e.Correct("")
}

// @Title hr新建离职
// @Description hr新建离职
// @Param	json	body	string	true	"离职员工信息"
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /leave/:id [post]
func (e *EmployeeController) NewEmpLeave() {
	eID, _ := e.GetInt(":id", 0)
	userType, _ := e.GetInt("userType", 0)
	operator := e.GetString("userName")
	if userType != models.UserHR {
		e.ErrorOK("您不是HR")
	}
	reqLeave := new(oa.ReqQuitFlow)
	err := json.Unmarshal(e.Ctx.Input.RequestBody, reqLeave)
	if err != nil {
		log.GLogger.Error("req leave err：%s", err.Error())
		e.ErrorOK(MsgInvalidParam)
	}

	tx := services.Slave().Begin()
	quitInfo := reqLeave.ToEntity()
	quitInfo.EmployeeID = eID
	err = tx.Create(quitInfo).Error
	if err != nil {
		log.GLogger.Error("employee leave err：%s", err.Error())
		tx.Rollback()
		e.ErrorOK(MsgServerErr)
	}
	//更新employee信息
	err = tx.Model(oa.Employee{}).Where("id = ?", eID).Updates(map[string]interface{}{
		"req_user":         operator,
		"reason":           reqLeave.Reason,
		"status":           3,
		"resignation_date": reqLeave.ResignationDate.String(),
	}).Error
	if err != nil {
		log.GLogger.Error("employee leave err：%s", err.Error())
		tx.Rollback()
		e.ErrorOK(MsgServerErr)
	}
	//创建流程信息
	uID, _ := e.GetInt("userID", 0)
	err = services.CreateLeaveWorkflow(tx, eID, uID)
	if err != nil {
		log.GLogger.Error("create leave workflow err：%s", err.Error())
		tx.Rollback()
		e.ErrorOK(MsgServerErr)
	}
	tx.Commit()
	e.Correct("")
}

// @Title 离职流程信息
// @Description 离职流程信息
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /leave/:id [get]
func (e *EmployeeController) GetLeaveInfo() {
	eID, _ := e.GetInt(":id", 0)
	flowInfo := new(oa.QuitFlowInfo)
	services.Slave().Model(oa.QuitFlowInfo{}).Where("employee_id = ?", eID).First(flowInfo)
	e.Correct(flowInfo)
}

// @Title 提交离职流程信息
// @Description 提交离职流程信息
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /leave/:id [put]
func (e *EmployeeController) CommitLeaveInfoNode() {
	eID, _ := e.GetInt(":id", 0)
	flowInfo := new(oa.QuitFlowInfo)
	err := json.Unmarshal(e.Ctx.Input.RequestBody, flowInfo)
	if err != nil {
		e.ErrorOK(MsgInvalidParam)
	}
	//修改入职流程信息
	userType, _ := e.GetInt("userType", 0)
	//eID 查询 workflow
	workflow := new(oa.Workflow)
	services.Slave().Model(oa.Workflow{}).Where("workflow_definition_id = ? and entity_id = ?",
		services.GetLeaveDef(), eID).Preload("Nodes").Preload("Elements").First(workflow)
	if len(workflow.Nodes) != 4 {
		e.ErrorOK("工作流配置错误")
	}
	switch userType {
	case models.UserLeader:
		e.ErrorOK("工作流不在当前节点")
	case models.UserHR:
		e.ErrorOK("工作流不在当前节点")
	case models.UserIT:
		if workflow.Nodes[1].Status != services.FlowProcessing {
			e.ErrorOK("工作流不在当前节点")
		}
		//修改 flowInfo 信息
		services.Slave().Save(flowInfo)
		//更新节点信息
		workflow.Nodes[1].Status = services.FlowCompleted
		workflow.Nodes[2].Status = services.FlowProcessing
		services.Slave().Save(workflow)
	case models.UserFinance:
		if workflow.Nodes[2].Status != services.FlowProcessing {
			e.ErrorOK("工作流不在当前节点")
		}
		//修改 flowInfo 信息
		services.Slave().Save(flowInfo)
		workflow.Nodes[2].Status = services.FlowCompleted
		workflow.Nodes[3].Status = services.FlowProcessing
		services.Slave().Save(workflow)
	case models.UserFront:
		//前台
		if workflow.Nodes[3].Status != services.FlowProcessing {
			e.ErrorOK("工作流不在当前节点")
		}
		//修改 flowInfo 信息,最后一个流程，变为已离职
		services.Slave().Save(flowInfo)
		workflow.Nodes[3].Status = services.FlowCompleted
		workflow.Status = services.FlowCompleted
		services.Slave().Save(workflow)
		services.Slave().Model(oa.Employee{}).Where("id = ?", eID).Updates(map[string]interface{}{
			"status": 4,
		})
	}
	e.Correct("")
}

// @Title 保存员工所有信息
// @Description 保存员工所有信息
// @Param	json	body	string	true	"员工所有信息"
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /save [put]
func (e *EmployeeController) SaveEmpInfo() {
	employee := new(oa.Employee)
	err := json.Unmarshal(e.Ctx.Input.RequestBody, employee)
	if err != nil {
		log.GLogger.Error("parse employee info err:%s", err.Error())
		e.ErrorOK(MsgInvalidParam)
	}
	services.Slave().Save(employee)
	e.Correct("")
}

// @Title employee详情
// @Description employee详情
// @Param	id	path	int	true	"员工id"
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /detail/:id [get]
func (e *EmployeeController) GetEmpInfo() {
	eID, _ := e.GetInt(":id", 0)
	employee := new(oa.Employee)
	services.Slave().Model(oa.Employee{}).Where("id = ?", eID).Preload("Department").
		Preload("Department.Leader").Preload("Level").Preload("EmployeeBasic").First(&employee)
	e.Correct(employee)
}
