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
	"strconv"
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
	employee.Email = strconv.Itoa(int(time.Now().UnixNano()))
	employee.CreatorId, _ = e.GetInt("userID")
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

// @Title employee删除
// @Description employee删除
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /:id [delete]
func (e *EmployeeController) DeleteEmp() {
	eID, _ := e.GetInt(":id", 0)
	userType, _ := e.GetInt("userType", 0)
	if userType != models.UserHR {
		e.ErrorOK("您不是HR")
	}
	services.Slave().Delete(oa.Employee{}, "id = ?", eID)
	e.Correct("")
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
	createID, _ := e.GetInt("creator", -1)
	userType, _ := e.GetInt("userType", 0)
	userID, _ := e.GetInt("userID", 0)
	flow, _ := e.GetInt("flow", 1)
	employees := make([]*oa.Employee, 0)
	//未入职和拟入职
	var resp struct {
		Total int            `json:"total"`
		List  []*oa.Employee `json:"list"`
	}
	query := services.Slave().Model(oa.Employee{})

	//查流程表，得到员工id列表
	if userType != models.UserHR {
		ids := make([]*oa.EntityID, 0)
		//如果是IT,只显示流程在自己这入职
		services.Slave().Raw("select w.entity_id from workflows w, workflow_nodes wn where w.id = wn.workflow_id"+
			" and w.workflow_definition_id = ? and wn.operator_id = ? and wn.status = ? limit ?,?",
			flow, userID, models.FlowProcessing, (pageNum-1)*pageSize, pageSize).Scan(&ids)
		resp.Total = len(ids)
		start, end := getPage(resp.Total, pageSize, pageNum)
		eIDs := make([]int, 0)
		for _, eID := range ids[start:end] {
			eIDs = append(eIDs, eID.EntityID)
		}
		query.Where(eIDs).Order("updated_at desc").Find(&employees)
		resp.List = employees
		e.Correct(resp)
	}
	if dID != 0 {
		query = query.Where("department_id = ?", dID)
	}
	if status != -1 {
		query = query.Where("status = ?", status)
	} else {
		if flow == 1 {
			//入职
			query = query.Where("status in (?)", []int{0, 1, 2, 4})
		} else {
			query = query.Where("status in (?)", []int{3})
		}
	}
	if createID != -1 {
		query = query.Where("creator_id = ?", createID)
	}
	if name != "" {
		query = query.Where("name like ?", "%"+name+"%")
	}
	if number != "" {
		//员工编号
		query = query.Where("emp_no like ?", "%"+number+"%")
	}
	query.Preload("Department").Limit(pageSize).Offset((pageNum - 1) * pageSize).Order("created_at desc").
		Find(&employees).Limit(-1).Offset(-1).Count(&resp.Total)
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
	flowType := e.GetString("type", "entry")
	var flowDefId int
	if flowType == "entry" {
		flowDefId = services.GetEntryDef()
	} else if flowType == "leave" {
		flowDefId = services.GetLeaveDef()
	}
	workflow := new(oa.Workflow)
	services.Slave().Model(oa.Workflow{}).Where("workflow_definition_id = ? and entity_id = ?",
		flowDefId, eID).Preload("Nodes").Preload("Nodes.User").
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
		if workflow.Nodes[1].Status != models.FlowProcessing {
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
		workflow.Nodes[1].Status = models.FlowCompleted
		workflow.Nodes[2].Status = models.FlowProcessing
		services.Slave().Save(workflow)
	case models.UserIT:
		if workflow.Nodes[2].Status != models.FlowProcessing {
			e.ErrorOK("工作流不在当前节点")
		}
		//修改employee信息,最后一个流程，变为已入职
		services.Slave().Model(oa.Employee{}).Where("id = ?", eID).Updates(map[string]interface{}{
			"email":     flowInfo.Email,
			"wx_work":   flowInfo.WxWork,
			"tapd":      flowInfo.Tapd,
			"plan_date": workflow.Elements[0].Value,
			//"status":    2,
		})
		workflow.Nodes[2].Status = models.FlowCompleted
		workflow.Status = models.FlowCompleted
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

	quitInfo := reqLeave.ToEntity()
	quitInfo.EmployeeID = eID
	err = services.Slave().Create(quitInfo).Error
	if err != nil {
		log.GLogger.Error("employee leave err：%s", err.Error())
		e.ErrorOK(MsgServerErr)
	}
	//更新employee信息
	err = services.Slave().Model(oa.Employee{}).Where("id = ?", eID).Updates(map[string]interface{}{
		"req_user":         operator,
		"reason":           reqLeave.Reason,
		"status":           3,
		"resignation_date": reqLeave.ResignationDate.String(),
	}).Error
	if err != nil {
		log.GLogger.Error("employee leave err：%s", err.Error())
		e.ErrorOK(MsgServerErr)
	}
	//创建流程信息
	uID, _ := e.GetInt("userID", 0)
	err = services.CreateLeaveWorkflow(services.Slave(), eID, uID)
	if err != nil {
		log.GLogger.Error("create leave workflow err：%s", err.Error())
		e.ErrorOK(MsgServerErr)
	}
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
	flowInfo.EmployeeID = eID
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
		if workflow.Nodes[1].Status != models.FlowProcessing {
			e.ErrorOK("工作流不在当前节点")
		}
		//修改 flowInfo 信息
		services.Slave().Save(flowInfo)
		//更新节点信息
		workflow.Nodes[1].Status = models.FlowCompleted
		workflow.Nodes[2].Status = models.FlowProcessing
		services.Slave().Save(workflow)
	case models.UserFinance:
		if workflow.Nodes[2].Status != models.FlowProcessing {
			e.ErrorOK("工作流不在当前节点")
		}
		//修改 flowInfo 信息
		services.Slave().Save(flowInfo)
		workflow.Nodes[2].Status = models.FlowCompleted
		workflow.Nodes[3].Status = models.FlowProcessing
		services.Slave().Save(workflow)
	case models.UserFront:
		//前台
		if workflow.Nodes[3].Status != models.FlowProcessing {
			e.ErrorOK("工作流不在当前节点")
		}
		//修改 flowInfo 信息,最后一个流程，变为已离职
		services.Slave().Save(flowInfo)
		workflow.Nodes[3].Status = models.FlowCompleted
		workflow.Status = models.FlowCompleted
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

// @Title employee入职状态
// @Description employee入职状态
// @Param	status	path	int	true	"入职状态"
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /status/:id [put]
func (e *EmployeeController) UpdateEmpStatus() {
	eID, _ := e.GetInt(":id", 0)
	status, _ := e.GetInt("status", 0)
	entryDate := e.GetString("entry_date")
	m := make(map[string]interface{})
	m["status"] = status
	if entryDate != "" {
		m["entry_date"] = entryDate
	}
	err := services.Slave().Model(oa.Employee{}).Where("id = ?", eID).Updates(m).Error
	if err != nil {
		log.GLogger.Error("update entry status:%s", err.Error())
		e.ErrorOK(MsgServerErr)
	}
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

// @Title employee新建合同
// @Description 新建合同
// @Param	id	path	int	true	"员工id"
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /contract/:id [post]
func (e *EmployeeController) CreateEmpContract() {
	eID, _ := e.GetInt(":id", 0)
	contract := new(oa.EmployeeContract)
	err := json.Unmarshal(e.Ctx.Input.RequestBody, contract)
	if err != nil {
		log.GLogger.Error("parse contract err:%s", err.Error())
		e.ErrorOK(MsgInvalidParam)
	}
	contract.EmployeeID = eID
	err = services.Slave().Create(contract).Error
	if err != nil {
		log.GLogger.Error("create contract err:%s", err.Error())
		e.ErrorOK(MsgServerErr)
	}

	e.Correct(contract)
}

// @Title 合同列表
// @Description 合同列表
// @Param	pagesize	query	int	true	"页大小"
// @Param	pagenum	query	int	true	"页数"
// @Param	name	query	string	true	"姓名"
// @Param	status	query	string	true	"状态"
// @Param	number 	query	string	true	"编号"
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /contracts [get]
func (e *EmployeeController) GetContracts() {
	pageSize, _ := e.GetInt("pagesize", 10)
	pageNum, _ := e.GetInt("pagenum", 1)
	name := e.GetString("name")
	mainObj := e.GetString("main")
	status := e.GetString("status")
	number := e.GetString("number")
	contracts := make([]*oa.EmployeeContract, 0)
	query := services.Slave().Model(oa.EmployeeContract{})
	var resp struct {
		Total int                    `json:"total"`
		List  []*oa.EmployeeContract `json:"list"`
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if number != "" {
		//查询 eID
		eIDs := make([]int, 0)
		services.Slave().Model(oa.Employee{}).Where("emp_no like ?", "%"+number+"%").Pluck("ID", &eIDs)
		query = query.Where("employee_id in (?)", eIDs)
	}
	if name != "" {
		query = query.Where("contract_party like ?", "%"+name+"%")
	}
	if mainObj != "" {
		query = query.Where("contract_main like ?", "%"+mainObj+"%")
	}
	query.Limit(pageSize).Offset((pageNum - 1) * pageSize).Order("created_at desc").Find(&contracts).
		Limit(-1).Offset(-1).Count(&resp.Total)
	resp.List = contracts

	e.Correct(resp)
}

// @Title 员工合同列表
// @Description 员工合同列表
// @Param	id	path	int	true	"员工id"
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /contracts/:id [get]
func (e *EmployeeController) GetEmpContracts() {
	eID, _ := e.GetInt(":id", 0)
	contracts := make([]*oa.EmployeeContract, 0)
	services.Slave().Model(oa.EmployeeContract{}).Where("employee_id = ?", eID).Find(&contracts)
	e.Correct(contracts)
}

// @Title 修改合同
// @Description 修改合同
// @Param	id	path	int	true	"员工id"
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /contract [put]
func (e *EmployeeController) UpdateContract() {
	contract := new(oa.EmployeeContract)
	err := json.Unmarshal(e.Ctx.Input.RequestBody, contract)
	if err != nil {
		log.GLogger.Error("parse contract err:%s", err.Error())
		e.ErrorOK(MsgInvalidParam)
	}
	services.Slave().Save(contract)
	e.Correct("")
}

// @Title 获取合同信息
// @Description 获取合同信息
// @Param	id	path	int	true	"合同id"
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /contract/:id [get]
func (e *EmployeeController) GetContract() {
	cID, _ := e.GetInt(":id", 0)
	contract := new(oa.EmployeeContract)
	services.Slave().Take(contract, "id = ?", cID)
	e.Correct(contract)
}

// @Title 删除合同
// @Description 删除合同
// @Param	id	path	int	true	"员工id"
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /contract/:id [delete]
func (e *EmployeeController) DelContract() {
	cID, _ := e.GetInt(":id", 0)
	services.Slave().Delete(oa.EmployeeContract{}, "id = ?", cID)
	e.Correct("")
}

// @Title 合同即将到期
// @Description 合同即将到期
// @Param	pagesize	query	int	true	"页大小"
// @Param	pagenum	query	int	true	"页数"
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /contract/continue [get]
func (e *EmployeeController) MoreContract() {
	pageSize, _ := e.GetInt("pagesize", 10)
	pageNum, _ := e.GetInt("pagenum", 1)
	deadLine := time.Now().AddDate(0, 0, 45)

	var resp struct {
		Total int                  `json:"total"`
		List  []*oa.ContractSimple `json:"list"`
	}
	services.Slave().Table("employee_contracts").Select("contract_party,max(contract_end_date) as enddate, employee_id").
		Group("employee_id").Having("enddate < ?", deadLine.Format(models.DateFormat)).Count(&resp.Total).
		Limit(pageSize).Offset((pageNum - 1) * pageSize).Scan(&resp.List)
	e.Correct(resp)
}
