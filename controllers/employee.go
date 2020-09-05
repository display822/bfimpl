/*
* Auth : acer
* Desc : 入职离职流程
* Time : 2020/9/4 21:45
 */

package controllers

import (
	"bfimpl/models/oa"
	"bfimpl/services"
	"bfimpl/services/log"
	"encoding/json"
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
	if userType != 6 {
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
