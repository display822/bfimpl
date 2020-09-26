/*
* Auth : acer
* Desc : 加班，请假
* Time : 2020/9/12 23:39
 */

package controllers

import (
	"bfimpl/models"
	"bfimpl/models/forms"
	"bfimpl/models/oa"
	"bfimpl/services"
	"bfimpl/services/log"
	"encoding/json"
	"strings"
	"time"
)

type WorkController struct {
	BaseController
}

// @Title 申请加班
// @Description 申请加班
// @Param	json	body	string	true	"加班信息"
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /overtime [post]
func (w *WorkController) ReqOvertime() {
	uID, _ := w.GetInt("userID", 0)
	uEmail := w.GetString("userEmail")
	//获取emp_info
	employee := new(oa.Employee)
	services.Slave().Preload("Department").Preload("Department.Leader").Take(employee, "email = ?", uEmail)
	if employee.ID == 0 {
		w.ErrorOK("未找到员工信息")
	}
	param := new(oa.Overtime)
	err := json.Unmarshal(w.Ctx.Input.RequestBody, param)
	if err != nil {
		log.GLogger.Error("parse overtime info err:%s", err.Error())
		w.ErrorOK(MsgInvalidParam)
	}
	param.EmpID = int(employee.ID)
	param.EName = employee.Name
	param.StartTime = models.Time(time.Now())
	param.EndTime = models.Time(time.Now())
	param.ReqTime = models.Time(time.Now())
	param.Status = models.FlowNA
	tx := services.Slave().Begin()
	err = tx.Create(param).Error
	if err != nil {
		log.GLogger.Error("req overtime err:%s", err.Error())
		tx.Rollback()
		w.ErrorOK(MsgServerErr)
	}
	leaderID := 0
	if employee.Department != nil {
		leaderID = employee.Department.DepartmentLeaderID
	}
	err = services.ReqOvertime(tx, int(param.ID), uID, leaderID)
	if err != nil {
		log.GLogger.Error("req overtime err:%s", err.Error())
		tx.Rollback()
		w.ErrorOK(MsgServerErr)
	}
	tx.Commit()
	w.Correct(param)
}

// @Title 加班审批人
// @Description 加班审批人
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /approvals [get]
func (w *WorkController) ApprovalUsers() {
	uEmail := w.GetString("userEmail")
	//获取emp_info
	employee := new(oa.Employee)
	services.Slave().Preload("Department").Preload("Department.Leader").Take(employee, "email = ?", uEmail)
	if employee.ID == 0 {
		w.ErrorOK("未找到员工信息")
	}
	approvalUsers := make([]string, 0)
	if employee.Department != nil && employee.Department.Leader != nil {
		approvalUsers = append(approvalUsers, employee.Department.Leader.Name)
	}
	u := services.GetWorkUser(models.UserHR)
	if u != nil {
		approvalUsers = append(approvalUsers, u.Name)
	}
	w.Correct(strings.Join(approvalUsers, ";"))
}

// @Title 单条申请加班
// @Description 单条申请加班
// @Param	id	path	int	true	"加班id"
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /overtime/:id [get]
func (w *WorkController) OvertimeById() {
	oID, _ := w.GetInt(":id", 0)
	overtime := new(oa.Overtime)
	services.Slave().Take(overtime, "id = ?", oID)
	//oID 查询 workflow
	workflow := new(oa.Workflow)
	services.Slave().Model(oa.Workflow{}).Where("workflow_definition_id = ? and entity_id = ?",
		services.GetFlowDefID(services.Overtime), oID).Preload("Nodes").Preload("Nodes.User").
		Preload("Elements").First(workflow)
	if len(workflow.Nodes) != 3 {
		w.ErrorOK("工作流配置错误")
	}
	var resp struct {
		Info     *oa.Overtime `json:"info"`
		WorkFlow *oa.Workflow `json:"work_flow"`
	}
	resp.Info = overtime
	resp.WorkFlow = workflow
	w.Correct(resp)
}

// @Title 申请加班列表
// @Description 申请加班列表
// @Param	pagenum	    query	int	false	"分页"
// @Param	pagesize	query	int	false	"分页"
// @Param	name	query	string	false	"姓名"
// @Param	type	query	string	false	"加班类型"
// @Param	status	query	string	false	"状态"
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /overtime [get]
func (w *WorkController) OvertimeList() {
	pageSize, _ := w.GetInt("pagesize", 10)
	pageNum, _ := w.GetInt("pagenum", 1)
	name := w.GetString("name")
	workType := w.GetString("type")
	status := w.GetString("status")

	userType, _ := w.GetInt("userType", 0)
	userEmail := w.GetString("userEmail")
	overtimes := make([]*oa.Overtime, 0)
	query := services.Slave().Model(oa.Overtime{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if workType != "" {
		query = query.Where("type = ?", workType)
	}
	if userType != models.UserHR && userType != models.UserLeader {
		//不是hr和部门负责人，只能查自己
		employee := new(oa.Employee)
		services.Slave().Where("email = ?", userEmail).First(employee)
		query = query.Where("emp_id = ?", employee.ID)
	}
	if name != "" {
		query = query.Where("e_name like ?", "%"+name+"%")
	}
	var resp struct {
		Total int            `json:"total"`
		List  []*oa.Overtime `json:"list"`
	}
	query.Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&overtimes).Limit(-1).Offset(-1).Count(&resp.Total)
	resp.List = overtimes
	w.Correct(resp)
}

// @Title 审批申请加班
// @Description 审批申请加班
// @Param	id	body	int	true	"加班id"
// @Param	comment	body	string	true	"审批意见"
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /overtime [put]
func (w *WorkController) ApprovalOvertime() {
	param := new(forms.ReqApprovalOvertime)
	err := json.Unmarshal(w.Ctx.Input.RequestBody, param)
	if err != nil {
		log.GLogger.Error("parse overtime err:%s", err.Error())
		w.ErrorOK(MsgInvalidParam)
	}
	//oID 查询 workflow
	workflow := new(oa.Workflow)
	services.Slave().Model(oa.Workflow{}).Where("workflow_definition_id = ? and entity_id = ?",
		services.GetFlowDefID(services.Overtime), param.Id).Preload("Nodes").Preload("Nodes.User").
		Preload("Elements").First(workflow)
	if workflow.Nodes == nil || len(workflow.Nodes) != 3 {
		w.ErrorOK("工作流配置错误")
	}
	isCheck := false
	userID, _ := w.GetInt("userID", 0)
	// 负责人，hr审批
	num := len(workflow.Nodes)
	for i, node := range workflow.Nodes {
		if node.Status == models.FlowProcessing && node.OperatorID == userID {
			isCheck = true
			status := models.FlowRejected
			if param.Status == 1 {
				status = models.FlowApproved
			}
			node.Status = status
			workflow.Elements[i].Value = param.Comment
			if i < num-1 {
				//负责人
				if param.Status == 0 {
					workflow.Status = status
					services.Slave().Model(oa.Overtime{}).Where("id = ?", param.Id).Updates(map[string]interface{}{
						"status": status,
					})
				} else {
					workflow.Nodes[i+1].Status = models.FlowProcessing
				}
				services.Slave().Save(workflow)
			} else if i == num-1 {
				//hr
				workflow.Status = status
				services.Slave().Save(workflow)
				services.Slave().Model(oa.Overtime{}).Where("id = ?", param.Id).Updates(map[string]interface{}{
					"status": status,
				})
			}
			break
		}
	}
	if !isCheck {
		w.ErrorOK("您不是当前审批人")
	}
	w.Correct("")
}

// @Title 审批申请加班
// @Description 审批申请加班
// @Param	id	path	int	true	"加班id"
// @Param	real	query	string	true	"实际加班时长"
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /overtime/:id/check [put]
func (w *WorkController) ValidOvertime() {
	oID, _ := w.GetInt(":id", 0)
	realDuration, _ := w.GetInt("real", 0)
	err := services.Slave().Model(oa.Overtime{}).Where("id = ?", oID).Updates(map[string]interface{}{
		"real_duration": realDuration,
	}).Error
	if err != nil {
		w.ErrorOK(MsgServerErr)
	}
	w.Correct("")
}

//==========leave 请假接口==============

// @Title 申请请假
// @Description 申请请假
// @Param	json	body	string	true	"请假信息"
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /leave [post]
func (w *WorkController) ReqLeave() {
	uID, _ := w.GetInt("userID", 0)
	uEmail := w.GetString("userEmail")
	//获取emp_info
	employee := new(oa.Employee)
	services.Slave().Preload("Department").Preload("Department.Leader").Take(employee, "email = ?", uEmail)
	if employee.ID == 0 {
		w.ErrorOK("未找到员工信息")
	}
	param := new(oa.Leave)
	err := json.Unmarshal(w.Ctx.Input.RequestBody, param)
	if err != nil {
		log.GLogger.Error("parse leave info err:%s", err.Error())
		w.ErrorOK(MsgInvalidParam)
	}
	param.EmpID = int(employee.ID)
	param.EName = employee.Name
	param.ReqTime = models.Time(time.Now())
	param.Status = models.FlowNA
	tx := services.Slave().Begin()
	err = tx.Create(param).Error
	if err != nil {
		log.GLogger.Error("req leave err:%s", err.Error())
		tx.Rollback()
		w.ErrorOK(MsgServerErr)
	}
	leaderID := 0
	if employee.Department != nil {
		leaderID = employee.Department.DepartmentLeaderID
	}
	err = services.ReqLeave(tx, int(param.ID), uID, leaderID)
	if err != nil {
		log.GLogger.Error("req leave err:%s", err.Error())
		tx.Rollback()
		w.ErrorOK(MsgServerErr)
	}
	tx.Commit()
	w.Correct(param)
}

// @Title 请假审批人
// @Description 请假审批人
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /leave/approvals [get]
func (w *WorkController) LeaveApprovalUsers() {
	uEmail := w.GetString("userEmail")
	//获取emp_info
	employee := new(oa.Employee)
	services.Slave().Preload("Department").Preload("Department.Leader").Take(employee, "email = ?", uEmail)
	if employee.ID == 0 {
		w.ErrorOK("未找到员工信息")
	}
	approvalUsers := make([]string, 0)
	if employee.Department != nil && employee.Department.Leader != nil {
		approvalUsers = append(approvalUsers, employee.Department.Leader.Name)
	}
	u := services.GetWorkUser(models.UserHR)
	if u != nil {
		approvalUsers = append(approvalUsers, u.Name)
	}
	w.Correct(strings.Join(approvalUsers, ";"))
}

// @Title 单条申请请假
// @Description 单条申请请假
// @Param	id	path	int	true	"请假id"
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /leave/:id [get]
func (w *WorkController) LeaveById() {
	lID, _ := w.GetInt(":id", 0)
	leave := new(oa.Leave)
	services.Slave().Take(leave, "id = ?", lID)
	//oID 查询 workflow
	workflow := new(oa.Workflow)
	services.Slave().Model(oa.Workflow{}).Where("workflow_definition_id = ? and entity_id = ?",
		services.GetFlowDefID(services.Leave), lID).Preload("Nodes").Preload("Nodes.User").
		Preload("Elements").First(workflow)
	var resp struct {
		Info     *oa.Leave    `json:"info"`
		WorkFlow *oa.Workflow `json:"work_flow"`
	}
	resp.Info = leave
	resp.WorkFlow = workflow
	w.Correct(resp)
}

// @Title 申请请假列表
// @Description 申请请假列表
// @Param	pagenum	    query	int	false	"分页"
// @Param	pagesize	query	int	false	"分页"
// @Param	name	query	string	false	"姓名"
// @Param	type	query	string	false	"请假类型"
// @Param	status	query	string	false	"状态"
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /leave [get]
func (w *WorkController) LeaveList() {
	pageSize, _ := w.GetInt("pagesize", 10)
	pageNum, _ := w.GetInt("pagenum", 1)
	name := w.GetString("name")
	workType := w.GetString("type")
	status := w.GetString("status")

	userType, _ := w.GetInt("userType", 0)
	userEmail := w.GetString("userEmail")
	leaves := make([]*oa.Leave, 0)
	query := services.Slave().Model(oa.Leave{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if workType != "" {
		query = query.Where("type = ?", workType)
	}
	if userType != models.UserHR && userType != models.UserLeader {
		//不是hr和部门负责人，只能查自己
		employee := new(oa.Employee)
		services.Slave().Where("email = ?", userEmail).First(employee)
		query = query.Where("emp_id = ?", employee.ID)
	}
	if name != "" {
		query = query.Where("e_name like ?", "%"+name+"%")
	}
	var resp struct {
		Total int         `json:"total"`
		List  []*oa.Leave `json:"list"`
	}
	query.Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&leaves).Limit(-1).Offset(-1).Count(&resp.Total)
	resp.List = leaves
	w.Correct(resp)
}

// @Title 审批请假
// @Description 审批请假
// @Param	id	body	int	true	"请假id"
// @Param	comment	body	string	true	"审批意见"
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /leave [put]
func (w *WorkController) ApprovalLeave() {
	param := new(forms.ReqApprovalOvertime)
	err := json.Unmarshal(w.Ctx.Input.RequestBody, param)
	if err != nil {
		log.GLogger.Error("parse overtime err:%s", err.Error())
		w.ErrorOK(MsgInvalidParam)
	}
	//oID 查询 workflow
	workflow := new(oa.Workflow)
	services.Slave().Model(oa.Workflow{}).Where("workflow_definition_id = ? and entity_id = ?",
		services.GetFlowDefID(services.Leave), param.Id).Preload("Nodes").Preload("Nodes.User").
		Preload("Elements").First(workflow)

	isCheck := false
	userID, _ := w.GetInt("userID", 0)
	// 负责人，hr审批
	num := len(workflow.Nodes)
	for i, node := range workflow.Nodes {
		if node.Status == models.FlowProcessing && node.OperatorID == userID {
			isCheck = true
			status := models.FlowRejected
			if param.Status == 1 {
				status = models.FlowApproved
			}
			node.Status = status
			workflow.Elements[i].Value = param.Comment
			if i < num-1 {
				//负责人
				if param.Status == 0 {
					//中间节点拒绝
					workflow.Status = status
					services.Slave().Model(oa.Leave{}).Where("id = ?", param.Id).Updates(map[string]interface{}{
						"status": status,
					})
				} else {
					workflow.Nodes[i+1].Status = models.FlowProcessing
				}
				services.Slave().Save(workflow)
			} else if i == num-1 {
				//hr
				workflow.Status = status
				services.Slave().Save(workflow)
				services.Slave().Model(oa.Leave{}).Where("id = ?", param.Id).Updates(map[string]interface{}{
					"status": status,
				})
			}
			break
		}
	}
	if !isCheck {
		w.ErrorOK("您不是当前审批人")
	}
	w.Correct("")
}

// @Title 审批申请加班
// @Description 审批申请加班
// @Param	id	path	int	true	"加班id"
// @Param	real	query	string	true	"实际加班时长"
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /leave/:id/check [put]
func (w *WorkController) ValidLeave() {
	lID, _ := w.GetInt(":id", 0)
	realDuration, _ := w.GetInt("real", 0)
	err := services.Slave().Model(oa.Leave{}).Where("id = ?", lID).Updates(map[string]interface{}{
		"real_duration": realDuration,
	}).Error
	if err != nil {
		w.ErrorOK(MsgServerErr)
	}
	w.Correct("")
}
