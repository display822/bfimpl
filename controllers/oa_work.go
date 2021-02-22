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
	"fmt"
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
	//查询HRBP id
	engagementCode := new(oa.EngagementCode)
	services.Slave().Model(oa.EngagementCode{}).Where("department_id = ?", employee.DepartmentID).First(engagementCode)
	param := new(oa.Overtime)
	err := json.Unmarshal(w.Ctx.Input.RequestBody, param)
	if err != nil {
		log.GLogger.Error("parse overtime info err:%s", err.Error())
		w.ErrorOK(MsgInvalidParam)
	}
	param.EmpID = int(employee.ID)
	param.EName = employee.Name
	param.ReqTime = models.Time(time.Now())
	param.Status = models.FlowNA
	tx := services.Slave().Begin()
	err = tx.Create(param).Error
	if err != nil {
		log.GLogger.Error("req overtime err:%s", err.Error())
		tx.Rollback()
		w.ErrorOK(MsgServerErr)
	}
	// leaderID := param.LeaderId
	// if employee.Department != nil {
	// 	leaderID = employee.Department.DepartmentLeaderID
	// }
	err = services.ReqOvertime(tx, int(param.ID), uID, param.LeaderId, engagementCode.HRID)
	if err != nil {
		log.GLogger.Error("req overtime err:%s", err.Error())
		tx.Rollback()
		w.ErrorOK(MsgServerErr)
	}
	tx.Commit()
	w.Correct(param)
}

// @Title engagement_list
// @Description 项目code
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /overtime/projects [get]
func (w *WorkController) GetProjects() {
	desc := w.GetString("desc")
	uEmail := w.GetString("userEmail")
	//获取emp_info
	employee := new(oa.Employee)
	services.Slave().Preload("Department").Preload("Department.Leader").Take(employee, "email = ?", uEmail)
	if employee.ID == 0 {
		w.ErrorOK("未找到员工信息")
	}
	//查询部门下项目list
	projects := make([]*oa.EngagementCode, 0)
	query := services.Slave().Model(oa.EngagementCode{}).Preload("Owner").Where("department_id = ?", employee.Department.ID)
	if desc != "" {
		query = query.Where("engagement_code_desc like ?", "%"+desc+"%")
	}
	query.Find(&projects)
	// for _, p :=range projects{
	// 	if p.CodeOwnerID == int(employee.ID){
	// 		p.Owner = employee.Department.Leader
	// 	}
	// }
	fmt.Println("employee.Department.Leader", employee.Department.Leader)
	for i := 0; i < len(projects); i++ {
		fmt.Println("CodeOwnerID", projects[i].CodeOwnerID)
		fmt.Println("int(employee.ID)", int(employee.ID))
		if projects[i].CodeOwnerID == int(employee.ID) {
			projects[i].Owner = employee.Department.Leader
			projects[i].CodeOwnerID = int(employee.Department.Leader.ID)
		}
	}

	w.Correct(projects)
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
	//查询HRBP id
	engagementCode := new(oa.EngagementCode)
	services.Slave().Model(oa.EngagementCode{}).Where("department_id = ?", employee.DepartmentID).First(engagementCode)
	u := new(models.User)
	services.Slave().Take(u, "id = ?", engagementCode.HRID)
	w.Correct(u.Name)
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
	//if len(workflow.Nodes) != 3 {
	//	w.ErrorOK("工作流配置错误")
	//}
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
	myreq, _ := w.GetBool("myreq", false)
	mytodo, _ := w.GetBool("mytodo", false)
	startdate := w.GetString("startdate")
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
	if startdate != "" {
		query = query.Where("start_date >= ?", startdate)
	}
	employee := new(oa.Employee)
	services.Slave().Where("email = ?", userEmail).First(employee)
	if userType != models.UserHR && userType != models.UserLeader {
		//不是hr和部门负责人，只能查自己
		query = query.Where("emp_id = ?", employee.ID)
	} else {
		if name != "" {
			query = query.Where("e_name like ?", "%"+name+"%")
		}
		if myreq {
			//查自己
			query = query.Where("emp_id = ?", employee.ID)
		}
	}
	var resp struct {
		Total int            `json:"total"`
		List  []*oa.Overtime `json:"list"`
	}
	if mytodo {
		//待我审核，查询结点待我审核的id
		qs := make([]string, 0)
		if status == models.FlowNA {
			status = models.FlowProcessing
		}
		if status != "" {
			qs = append(qs, status)
		} else {
			qs = append(qs, models.FlowApproved, models.FlowRejected, models.FlowProcessing)
		}
		userID, _ := w.GetInt("userID", 0)
		ids := make([]*oa.EntityID, 0)
		rsql := "select distinct w.entity_id from workflows w,workflow_nodes wn,overtimes o where w.id = " +
			"wn.workflow_id and w.workflow_definition_id = ? and w.entity_id = o.id and operator_id = ? and wn.status in (?)" +
			" and wn.node_seq != 1"
		p := make([]interface{}, 0)
		p = append(p, services.GetFlowDefID(services.Overtime), userID, qs)
		if name != "" {
			rsql += " and o.e_name like ?"
			p = append(p, "%"+name+"%")
		}
		if workType != "" {
			rsql += " and o.type = ?"
			p = append(p, workType)
		}
		if startdate != "" {
			rsql += " and o.overtime_date >= ?"
			p = append(p, startdate)
		}
		rsql += " order by o.created_at desc"
		services.Slave().Raw(rsql, p...).Scan(&ids)
		resp.Total = len(ids)
		start, end := getPage(resp.Total, pageSize, pageNum)
		eIDs := make([]int, 0)
		for _, eID := range ids[start:end] {
			eIDs = append(eIDs, eID.EntityID)
		}
		services.Slave().Model(oa.Overtime{}).Where(eIDs).Find(&overtimes)
	} else {
		query.Limit(pageSize).Offset((pageNum - 1) * pageSize).Order("created_at desc").Find(&overtimes).Limit(-1).Offset(-1).Count(&resp.Total)
	}
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
	if workflow.Nodes == nil {
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
				//审批通过且类型为weekend,holiday，将加班时长放入leave balance
				if status == models.FlowApproved {
					overtime := new(oa.Overtime)
					services.Slave().Take(overtime, "id = ?", param.Id)
					if overtime.Type == "weekend" || overtime.Type == "holiday" {
						balance := oa.LeaveBalance{
							EmpID:      overtime.EmpID,
							Type:       oa.OverTime,
							Amount:     (overtime.RealDuration) / 8,
							OvertimeID: int(overtime.ID),
						}
						if balance.Amount == 0 {
							balance.Amount = (overtime.Duration) / 8
						}
						services.Slave().Create(&balance)
					}
				}
			}
			break
		}
	}
	if !isCheck {
		w.ErrorOK("您不是当前审批人")
	}
	w.Correct("")
}

// @Title 取消申请加班
// @Description 审批申请加班
// @Param	id	path	int	true	"加班id"
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /overtime/:id [delete]
func (w *WorkController) DeleteOvertime() {
	id, _ := w.GetInt(":id", -1)
	//oID 查询 workflow
	workflow := new(oa.Workflow)
	services.Slave().Model(oa.Workflow{}).Where("workflow_definition_id = ? and entity_id = ?",
		services.GetFlowDefID(services.Overtime), id).Preload("Nodes").Preload("Elements").First(workflow)
	if workflow.Status != models.FlowProcessing {
		w.ErrorOK("流程已结束")
	}
	tx := services.Slave().Begin()
	for _, node := range workflow.Nodes {
		err := tx.Delete(node).Error
		if err != nil {
			tx.Rollback()
			w.ErrorOK(MsgServerErr)
		}
	}
	err := tx.Delete(workflow).Error
	if err != nil {
		tx.Rollback()
		w.ErrorOK(MsgServerErr)
	}
	err = tx.Delete(oa.Overtime{}, "id = ?", id).Error
	if err != nil {
		tx.Rollback()
		w.ErrorOK(MsgServerErr)
	}
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		w.ErrorOK(MsgServerErr)
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
	realDuration, _ := w.GetFloat("real", 0)
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
	//查询HRBP id
	engagementCode := new(oa.EngagementCode)
	services.Slave().Model(oa.EngagementCode{}).Where("department_id = ?", employee.DepartmentID).First(engagementCode)
	param := new(oa.Leave)
	err := json.Unmarshal(w.Ctx.Input.RequestBody, param)
	if err != nil {
		log.GLogger.Error("parse leave info err:%s", err.Error())
		w.ErrorOK(MsgInvalidParam)
	}
	//请年假和周末调休，查询是否有剩余
	if param.Type == oa.Annual || param.Type == oa.Shift {
		data := getRemain(int(employee.ID))
		if param.Type == oa.Annual && (param.Duration)/8 > data.Annual {
			w.ErrorOK("剩余年假不足")
		}
		if param.Type == oa.Shift && (param.Duration)/8 > data.Weekend {
			w.ErrorOK("剩余调休不足")
		}
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
		if leaderID == uID {
			leaderID = 0
		}
	}
	err = services.ReqLeave(tx, int(param.ID), uID, leaderID, engagementCode.HRID, param.Others...)
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
	uID, _ := w.GetInt("userID")
	//获取emp_info
	employee := new(oa.Employee)
	services.Slave().Preload("Department").Preload("Department.Leader").Take(employee, "email = ?", uEmail)
	if employee.ID == 0 {
		w.ErrorOK("未找到员工信息")
	}
	approvalUsers := make([]string, 0)
	if employee.Department != nil && employee.Department.Leader != nil {
		if uID != employee.Department.DepartmentLeaderID {
			approvalUsers = append(approvalUsers, employee.Department.Leader.Name)
		}
	}
	//查询HRBP id
	engagementCode := new(oa.EngagementCode)
	services.Slave().Model(oa.EngagementCode{}).Where("department_id = ?", employee.DepartmentID).First(engagementCode)
	u := new(models.User)
	services.Slave().Take(u, "id = ?", engagementCode.HRID)
	if u.Name != "" {
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
	startdate := w.GetString("startdate")
	myreq, _ := w.GetBool("myreq", false)
	mytodo, _ := w.GetBool("mytodo", false)

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
	if startdate != "" {
		query = query.Where("start_date >= ?", startdate)
	}
	employee := new(oa.Employee)
	services.Slave().Where("email = ?", userEmail).First(employee)
	if userType != models.UserHR && userType != models.UserLeader {
		//不是hr和部门负责人，只能查自己
		query = query.Where("emp_id = ?", employee.ID)
	} else {
		if name != "" {
			query = query.Where("e_name like ?", "%"+name+"%")
		}
		if myreq {
			//查自己
			query = query.Where("emp_id = ?", employee.ID)
		}
	}
	var resp struct {
		Total int         `json:"total"`
		List  []*oa.Leave `json:"list"`
	}
	if mytodo {
		//待我审核，查询结点待我审核的id
		qs := make([]string, 0)
		if status == models.FlowNA {
			status = models.FlowProcessing
		}
		if status != "" {
			qs = append(qs, status)
		} else {
			qs = append(qs, models.FlowApproved, models.FlowRejected, models.FlowProcessing)
		}
		userID, _ := w.GetInt("userID", 0)
		ids := make([]*oa.EntityID, 0)
		rsql := "select distinct w.entity_id from workflows w,workflow_nodes wn,leaves l where w.id = " +
			"wn.workflow_id and w.workflow_definition_id = ? and w.entity_id = l.id and operator_id = ? and wn.status in (?)" +
			" and wn.node_seq != 1"
		p := make([]interface{}, 0)
		p = append(p, services.GetFlowDefID(services.Leave), userID, qs)
		if name != "" {
			rsql += " and l.e_name like ?"
			p = append(p, "%"+name+"%")
		}
		if workType != "" {
			rsql += " and l.type = ?"
			p = append(p, workType)
		}
		if startdate != "" {
			rsql += " and l.start_date >= ?"
			p = append(p, startdate)
		}
		rsql += " order by l.created_at desc"
		services.Slave().Raw(rsql, p...).Scan(&ids)
		resp.Total = len(ids)
		start, end := getPage(resp.Total, pageSize, pageNum)
		eIDs := make([]int, 0)
		for _, eID := range ids[start:end] {
			eIDs = append(eIDs, eID.EntityID)
		}
		services.Slave().Model(oa.Leave{}).Where(eIDs).Find(&leaves)
	} else {
		query.Limit(pageSize).Offset((pageNum - 1) * pageSize).Order("created_at desc").Find(&leaves).Limit(-1).Offset(-1).Count(&resp.Total)
	}
	resp.List = leaves
	w.Correct(resp)
}

func getPage(total, pageSize, pageNum int) (int, int) {
	start, end := 0, total
	if total > pageSize {
		start = (pageNum - 1) * pageSize
		end = start + pageSize
		if start > total {
			start = 0
			end = 0
		} else {
			if end > total {
				end = total
			}
		}
	}
	return start, end
}

// @Title 申请请假列表
// @Description 按日期查询
// @Param	name	query	string	true	"姓名"
// @Param	status	query	string	false	"状态"
// @Param	date	query	string	true	"日期"
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /leavebydate [get]
func (w *WorkController) LeaveListByDate() {
	name := w.GetString("name")
	date := w.GetString("date")
	if name == "" || date == "" {
		w.ErrorOK(MsgInvalidParam)
	}
	leaves := make([]*oa.Leave, 0)
	services.Slave().Model(oa.Leave{}).Where("e_name = ? and start_date <= ? and end_date >= ?",
		name, date, date).Find(&leaves)
	w.Correct(leaves)
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
				//审批通过且类型为weekend,holiday，将加班时长放入leave balance
				if status == models.FlowApproved {
					leave := new(oa.Leave)
					services.Slave().Take(leave, "id = ?", param.Id)
					if leave.Type == "Shift" || leave.Type == "Annual" {
						remain := getRemain(leave.EmpID)
						balance := oa.LeaveBalance{
							EmpID:   leave.EmpID,
							Type:    oa.ShiftLeave,
							Amount:  -(leave.RealDuration) / 8,
							LeaveID: int(leave.ID),
						}
						if balance.Amount == 0 {
							balance.Amount = -(leave.Duration) / 8
						}
						log.GLogger.Info("%f,%f", remain.Weekend, remain.Annual)
						log.GLogger.Info("%f", balance.Amount)
						if leave.Type == "Shift" {
							if remain.Weekend < -balance.Amount {
								w.ErrorOK("剩余调休不足")
							}
						}
						if leave.Type == "Annual" {
							if remain.Annual < -balance.Amount {
								w.ErrorOK("剩余年假不足")
							}
							balance.Type = oa.AnnualLeave
						}
						services.Slave().Create(&balance)
					}
				}
			}
			break
		}
	}
	if !isCheck {
		w.ErrorOK("您不是当前审批人")
	}
	w.Correct("")
}

// @Title 取消申请请假
// @Description 取消申请请假
// @Param	id	path	int	true	"请假id"
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /leave/:id [delete]
func (w *WorkController) DeleteLeave() {
	id, _ := w.GetInt(":id", -1)
	//oID 查询 workflow
	workflow := new(oa.Workflow)
	services.Slave().Model(oa.Workflow{}).Where("workflow_definition_id = ? and entity_id = ?",
		services.GetFlowDefID(services.Leave), id).Preload("Nodes").Preload("Elements").First(workflow)
	if workflow.Status != models.FlowProcessing {
		w.ErrorOK("流程已结束")
	}
	tx := services.Slave().Begin()
	for _, node := range workflow.Nodes {
		err := tx.Delete(node).Error
		if err != nil {
			tx.Rollback()
			w.ErrorOK(MsgServerErr)
		}
	}
	err := tx.Delete(workflow).Error
	if err != nil {
		tx.Rollback()
		w.ErrorOK(MsgServerErr)
	}
	err = tx.Delete(oa.Leave{}, "id = ?", id).Error
	if err != nil {
		tx.Rollback()
		w.ErrorOK(MsgServerErr)
	}
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		w.ErrorOK(MsgServerErr)
	}
	w.Correct("")
}

// @Title 请假时长校验
// @Description 请假时长校验
// @Param	id	path	int	true	"加班id"
// @Param	real	query	string	true	"实际加班时长"
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /leave/:id/check [put]
func (w *WorkController) ValidLeave() {
	lID, _ := w.GetInt(":id", 0)
	realDuration, _ := w.GetFloat("real", 0)
	err := services.Slave().Model(oa.Leave{}).Where("id = ?", lID).Updates(map[string]interface{}{
		"real_duration": realDuration,
	}).Error
	if err != nil {
		w.ErrorOK(MsgServerErr)
	}
	w.Correct("")
}

// @Title 获取剩余年假和周末调休
// @Description 审批申请加班
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /remain/holiday [get]
func (w *WorkController) RemainHoliday() {
	empID, _ := w.GetInt("id", -1)
	if empID == -1 {
		uEmail := w.GetString("userEmail")
		//获取emp_info
		employee := new(oa.Employee)
		services.Slave().Take(employee, "email = ?", uEmail)
		if employee.ID == 0 {
			w.ErrorOK("未找到员工信息")
		}
		empID = int(employee.ID)
	}
	w.Correct(getRemain(empID))
}

func getRemain(empID int) oa.LeaveRemain {
	var remain oa.LeaveRemain
	balances := make([]*oa.LeaveBalance, 0)
	services.Slave().Model(oa.LeaveBalance{}).Where("emp_id = ?", empID).Find(&balances)
	for _, b := range balances {
		switch b.Type {
		case oa.OverTime:
			remain.Weekend += b.Amount
		case oa.Annual:
			remain.Annual += b.Amount
		case oa.ShiftLeave:
			remain.Weekend += b.Amount
		case oa.AnnualLeave:
			remain.Annual += b.Amount
		}
	}
	return remain
}

var months = []int{31, 29, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

//每月28号增加年假
func AddAnnual() {
	emps := make([]*oa.Employee, 0)
	services.Slave().Model(oa.Employee{}).Where("status = 2").Find(&emps)
	m := int(time.Now().Month()) - 1
	for _, emp := range emps {
		//取合同里的年假字段
		contracts := make([]*oa.EmployeeContract, 0)
		services.Slave().Table("employee_contracts").Select("contract_type,contract_party,contract_main,contract_start_date,max(contract_end_date) as contract_end_date,status, employee_id").
			Where("employee_id = ?", emp.ID).Group("employee_id").Scan(&contracts)
		if len(contracts) > 0 {
			annual := float32(contracts[0].AnnualLeave*months[m]) / 365
			if annual > 0 {
				balance := oa.LeaveBalance{
					EmpID:  int(emp.ID),
					Type:   oa.Annual,
					Amount: annual,
				}
				services.Slave().Create(&balance)
			}
		}
	}
}

//每年4月1号清空去年年假
func DeleteAnnual() {
	emps := make([]*oa.Employee, 0)
	services.Slave().Model(oa.Employee{}).Where("status = 2").Find(&emps)
	thisYearTime := time.Date(time.Now().Year(), time.January, 2, 0, 0, 0, 0, time.Local)
	for _, emp := range emps {
		balances := make([]*oa.LeaveBalance, 0)
		services.Slave().Model(oa.LeaveBalance{}).Where("emp_id = ?", emp.ID).Find(&balances)
		var total, thisYear float32
		for _, b := range balances {
			switch b.Type {
			case oa.Annual:
				total += b.Amount
				if b.CreatedAt.After(thisYearTime) {
					thisYear += b.Amount
				}
			case oa.AnnualLeave:
				total += b.Amount
			}
		}
		//如果年假总剩余大于今天增加的年假
		if total > thisYear {
			balance := oa.LeaveBalance{
				EmpID:  int(emp.ID),
				Type:   oa.AnnualLeave,
				Amount: thisYear - total,
			}
			services.Slave().Create(&balance)
		}
	}
}
