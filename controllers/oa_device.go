/**
* @author : yi.zhang
* @description : controllers 描述
* @date   : 2021-01-07 18:20
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
	"time"

	"github.com/jinzhu/gorm"
)

type DeviceController struct {
	BaseController
}

// @Title 创建设备
// @Description 创建设备
// @Param	body body oa.Device true "设备"
// @Success 200 {object} oa.Device
// @Failure 500 server internal err
// @router / [post]
func (d *DeviceController) Create() {
	// 验证员工身份 (7，8，9)
	userType, _ := d.GetInt("userType")
	if userType != models.UserIT && userType != models.UserFront && userType != models.UserFinance {
		d.ErrorOK("没有权限")
	}
	userID, _ := d.GetInt("userID")
	userName := d.GetString("userName")

	param := new(oa.Device)
	err := json.Unmarshal(d.Ctx.Input.RequestBody, param)
	if err != nil {
		log.GLogger.Error("parse device err:%s", err.Error())
		d.ErrorOK(MsgInvalidParam)
	}

	param.IngoingOperatorID = userID
	param.IngoingTime = models.Time(time.Now())

	log.GLogger.Info("param :%+v", param)

	tx := services.Slave().Begin()
	err = tx.Create(param).Error
	if err != nil {
		log.GLogger.Error("create device err：%s", err.Error())
		tx.Rollback()
		d.ErrorOK(MsgServerErr)
	}

	deviceRequisition := oa.DeviceRequisition{
		DeviceID:         int(param.ID),
		OperatorCategory: models.DeviceIngoing,
		OperatorID:       userID,
		OperatorName:     userName,
	}
	err = tx.Create(&deviceRequisition).Error
	if err != nil {
		log.GLogger.Error("create deviceRequisition err：%s", err.Error())
		tx.Rollback()
		d.ErrorOK(MsgServerErr)
	}
	tx.Commit()

	d.Correct("")
}

// @Title 设备列表
// @Description 设备列表
// @Param	pagenum	    query	int	false	"页码"
// @Param	pagesize	query	int	false	"页数"
// @Param	device_category	query	string	false	"设备分类"
// @Param	device_status	query	string	false	"设备状态"
// @Param	keyword	query	string	false	"搜索关键词"
// @Success 200 {object} []oa.Device
// @Failure 500 server internal err
// @router / [get]
func (d *DeviceController) List() {
	pageSize, _ := d.GetInt("pagesize", 10)
	pageNum, _ := d.GetInt("pagenum", 1)
	deviceCategory := d.GetString("device_category")
	deviceStatus := d.GetString("device_status")
	keyword := d.GetString("keyword")

	var list []*oa.Device
	db := services.Slave()
	if deviceCategory != "" {
		db = db.Where("device_category =?", deviceCategory)
	}
	if deviceStatus != "" {
		db = db.Where("device_status =?", deviceStatus)
	}
	if keyword != "" {
		k := fmt.Sprintf("%%%s%%", keyword)
		db = db.Where("device_name like ?", k)
		db = db.Or("device_num like ?", k)
		// 更多模糊查询...
	}
	var resp struct {
		Total int          `json:"total"`
		List  []*oa.Device `json:"list"`
	}
	db.Limit(pageSize).Offset((pageNum-1)*pageSize).Order("created_at desc").
		Preload("DeviceApply", func(db *gorm.DB) *gorm.DB {
			return db.Where("status = ?", "")
		}).
		Find(&list).Limit(-1).Offset(-1).Count(&resp.Total)
	resp.List = list

	d.Correct(resp)
}

// @Title 设备详情
// @Description 设备详情
// @Success 200 {object} oa.Device
// @Failure 500 server internal err
// @router /:id [get]
func (d *DeviceController) Get() {
	dID, _ := d.GetInt(":id", 0)
	var device oa.Device
	services.Slave().Where("id = ?", dID).
		Preload("DeviceRequisitions").
		Find(&device)
	d.Correct(device)
}

// @Title 设备更新
// @Description 设备更新
// @Param	body body oa.Device true "设备"
// @Success 200 {string} ""
// @Failure 500 server internal err
// @router / [put]
func (d *DeviceController) Put() {
	userType, _ := d.GetInt("userType")
	if userType != models.UserIT && userType != models.UserFront && userType != models.UserFinance {
		d.ErrorOK("没有权限")
	}

	param := new(oa.Device)
	err := json.Unmarshal(d.Ctx.Input.RequestBody, param)
	if err != nil {
		log.GLogger.Error("parse expense err:%s", err.Error())
		d.ErrorOK(MsgInvalidParam)
	}
	log.GLogger.Info("param :%+v", param)

	services.Slave().Save(param)

	d.Correct("")
}

// @Title 申请设备项目
// @Description 申请设备项目
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /projects [get]
func (d *DeviceController) GetProjects() {
	desc := d.GetString("desc")
	uEmail := d.GetString("userEmail")
	//获取emp_info
	employee := new(oa.Employee)
	services.Slave().Preload("Department").Preload("Department.Leader").Take(employee, "email = ?", uEmail)
	if employee.ID == 0 {
		d.ErrorOK("未找到员工信息")
	}
	//查询部门下项目list
	projects := make([]*oa.EngagementCode, 0)
	query := services.Slave().Model(oa.EngagementCode{}).Preload("Owner").Where("department_id = ?", employee.Department.ID)
	if desc != "" {
		query = query.Where("engagement_code_desc like ?", "%"+desc+"%")
	}
	query.Find(&projects)

	fmt.Println("employee.Department.Leader", employee.Department.Leader)
	for i := 0; i < len(projects); i++ {
		if projects[i].CodeOwnerID == int(employee.ID) {
			projects[i].Owner = employee.Department.Leader
			projects[i].CodeOwnerID = int(employee.Department.Leader.ID)
			if employee.Department.Leader.ID == employee.ID {
				user := new(models.User)
				services.Slave().Take(user, "id = ?", 2) // 马俊杰
				projects[i].Owner = user
				projects[i].CodeOwnerID = 2 // 马俊杰
			}
		}
	}

	d.Correct(projects)
}

// @Title 申请设备列表
// @Description 申请设备列表
// @Param	pagenum	    query	int	false	"页码"
// @Param	pagesize	query	int	false	"页数"
// @Param	myreq	query	bool	false	"我的报销"
// @Param	mytodo	query	bool	false	"我的审核"
// @Param	status	query	int	false	"状态"
// @Success 200 {object} []oa.DeviceApply
// @Failure 500 server internal err
// @router /apply [get]
func (d *DeviceController) ApplyList() {
	pageSize, _ := d.GetInt("pagesize", 10)
	pageNum, _ := d.GetInt("pagenum", 1)
	userType, _ := d.GetInt("userType", 0)
	name := d.GetString("name")
	myReq, _ := d.GetBool("myreq", false)
	myTodo, _ := d.GetBool("mytodo", false)
	status := d.GetString("status")
	userEmail := d.GetString("userEmail")
	searchID := d.GetString("searchid")
	applicationDateBegin := d.GetString("application_date_begin")
	applicationDateEnd := d.GetString("application_date_end")

	log.GLogger.Info("params", userEmail, userType, name, myReq, status, searchID, pageNum, pageSize, applicationDateBegin, applicationDateEnd)

	employee := new(oa.Employee)
	services.Slave().Where("email = ?", userEmail).First(employee)
	log.GLogger.Info("employee: %+v", employee)

	deviceApplys := make([]*oa.DeviceApply, 0)
	query := services.Slave().Debug().Model(oa.DeviceApply{})
	//if searchID != "" {
	//	query = query.Where("id like ?", fmt.Sprintf("%%%s%%", searchID))
	//}
	//if name != "" {
	//	query = query.Where("e_name like ?", "%"+name+"%")
	//}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	//if applicationDateBegin != "" && applicationDateEnd != "" {
	//	query = query.Where("application_date > ?", applicationDateBegin).Where("application_date <= ?", applicationDateEnd)
	//}

	var resp struct {
		Total int               `json:"total"`
		List  []*oa.DeviceApply `json:"list"`
	}

	if myReq {
		query = query.Where("emp_id = ?", employee.ID)
		query.Limit(pageSize).Offset((pageNum - 1) * pageSize).Preload("device").Order("created_at desc").Find(&deviceApplys).Limit(-1).Offset(-1).Count(&resp.Total)
	}

	if myTodo {
		userID, _ := d.GetInt("userID", 0)
		log.GLogger.Info("userID：%d", userID)
		ids := make([]oa.EntityID, 0)
		if status == "" {
			services.Slave().Debug().Raw("select w.entity_id from workflows w,workflow_nodes wn where w.id = "+
				"wn.workflow_id and w.workflow_definition_id = ? and operator_id = ? and wn.status <> ?"+
				" and wn.node_seq != 1 order by w.entity_id desc", services.GetFlowDefID(services.Device), userID, models.FlowHide).Scan(&ids)
		} else {
			services.Slave().Debug().Raw("select w.entity_id from workflows w,workflow_nodes wn where w.id = "+
				"wn.workflow_id and w.workflow_definition_id = ? and operator_id = ? and wn.status = ?"+
				" and wn.node_seq != 1 order by w.entity_id desc", services.GetFlowDefID(services.Device), userID, status).Scan(&ids)
		}

		resp.Total = len(ids)
		log.GLogger.Info("resp.Total: %d", len(ids))
		log.GLogger.Info("ids", ids)
		start, end := getPage(resp.Total, pageSize, pageNum)
		eIDs := make([]int, 0)
		for _, eID := range ids[start:end] {
			eIDs = append(eIDs, eID.EntityID)
		}
		services.Slave().Debug().Model(&oa.DeviceApply{}).Preload("Device").Order("created_at desc").Where(eIDs).Find(&deviceApplys)
	}

	resp.List = deviceApplys
	d.Correct(resp)
}

// @Title 已申请员工
// @Description 已申请员工
// @Success 200 {object} oa.Device
// @Failure 500 server internal err
// @router /apply/employee [get]
func (d *DeviceController) ApplyEmployee() {
	//dID := d.GetString("did")
}

// @Title 撤回申请设备
// @Description 撤回申请设备
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /apply/:id/recall [post]
func (d *DeviceController) RecallDevice() {
	id, _ := d.GetInt(":id")
	userID, _ := d.GetInt("userID")
	deviceApply := new(oa.DeviceApply)
	services.Slave().Debug().Preload("Employee").Take(deviceApply, "id = ?", id)
	log.GLogger.Info("deviceApply:%+v", deviceApply)

	if deviceApply.EmpID != userID {
		d.ErrorOK("没有权限")
	}
	//oID 查询 workflow
	workflow := new(oa.Workflow)
	services.Slave().Model(oa.Workflow{}).Where("workflow_definition_id = ? and entity_id = ?",
		services.GetFlowDefID(services.Device), id).Preload("Nodes").Preload("Nodes.User").
		Preload("Elements").First(workflow)
	if workflow.Nodes == nil || len(workflow.Nodes) != 2 {
		d.ErrorOK("工作流配置错误")
	}

	log.GLogger.Info("userId: %d", userID)
	log.GLogger.Info("expense.Employee.Email:%s", deviceApply.Employee.Email)

	if workflow.Status == models.FlowProcessing {
		workflow.Nodes[1].Status = models.FlowHide
	}

	workflow.Status = models.FlowRevoked
	services.Slave().Save(workflow)
	services.Slave().Model(oa.DeviceApply{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status": models.FlowRevoked,
	})

	d.Correct("")
}

// @Title 审批申请设备
// @Description 审批申请设备
// @Param	body	body	forms.ReqApprovalDevice	true
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /apply [put]
func (d *DeviceController) ApprovalDevice() {
	param := new(forms.ReqApprovalDevice)
	err := json.Unmarshal(d.Ctx.Input.RequestBody, param)
	if err != nil {
		log.GLogger.Error("parse ReqApprovalDevice err:%s", err.Error())
		d.ErrorOK(MsgInvalidParam)
	}
	log.GLogger.Info("param :%+v", param)

	deviceApply := new(oa.DeviceApply)
	services.Slave().Debug().Preload("Employee").Take(deviceApply, "id = ?", param.Id)
	log.GLogger.Info("deviceApply:%+v", deviceApply)
	//oID 查询 workflow
	workflow := new(oa.Workflow)
	services.Slave().Model(oa.Workflow{}).Where("workflow_definition_id = ? and entity_id = ?",
		services.GetFlowDefID(services.Device), param.Id).Preload("Nodes").Preload("Nodes.User").
		Preload("Elements").First(workflow)
	if workflow.Nodes == nil || len(workflow.Nodes) != 2 {
		d.ErrorOK("工作流配置错误")
	}
	isCheck := false
	userID, _ := d.GetInt("userID", 0)

	log.GLogger.Info("userId: %d", userID)
	log.GLogger.Info("expense.Employee.Email:%s", deviceApply.Employee.Email)

	// 负责人，hr审批
	for i, node := range workflow.Nodes {
		log.GLogger.Info("node.OperatorId:%d", node.OperatorID)
		if node.Status == models.FlowProcessing && node.OperatorID == userID {
			isCheck = true
			status := models.FlowRejected
			if param.Status == 1 {
				status = models.FlowApproved
			}
			node.Status = status
			workflow.Elements[i].Value = param.Comment

			workflow.Status = status
			services.Slave().Model(oa.DeviceApply{}).Where("id = ?", param.Id).Updates(map[string]interface{}{
				"status": status,
			})

			services.Slave().Save(workflow)
			break
		}

	}
	if !isCheck {
		d.ErrorOK("您不是当前审批人")
	}
	d.Correct("")
}

// @Title 设备借出
// @Description 设备借出
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /apply [put]
func (d *DeviceController) BorrowDevice() {
	userType, _ := d.GetInt("userType")
	if userType != models.UserIT && userType != models.UserFront && userType != models.UserFinance {
		d.ErrorOK("没有权限")
	}
	//dID := d.GetString("did")
	//eID := d.GetString("eid")

}

// @Title 申请设备基本信息
// @Description 申请设备基本信息
// @Success 200 {object} oa.Device
// @Failure 500 server internal err
// @router /apply/info [get]
func (d *DeviceController) ApplyInfo() {
	uID, _ := d.GetInt("userID", 0)
	uEmail := d.GetString("userEmail")
	dID := d.GetString("id")
	log.GLogger.Info("ReqExpense query: %d, %s", uID, uEmail)
	// 获取emp_info
	employee := new(oa.Employee)
	services.Slave().Preload("Department").Take(employee, "email = ?", uEmail)
	if employee.ID == 0 {
		d.ErrorOK("未找到员工信息")
	}

	device := new(oa.Device)
	services.Slave().Where("id =?", dID).Find(device)

	var collectDevices []*oa.Device
	//services.Slave().Where("id =?", dID).Find(device)

	res := oa.DeviceApplyInfo{
		Employee:       employee,
		Device:         device,
		CollectDevices: collectDevices,
	}

	d.Correct(res)
}

// @Title 申请设备
// @Description 申请设备
// @Param	body	    body	oa.DeviceApply	true	"设备申请"
// @Success 200 {object} ooa.DeviceApply
// @Failure 500 server internal err
// @router /apply [post]
func (d *DeviceController) ReqDevice() {
	uID, _ := d.GetInt("userID", 0)
	uEmail := d.GetString("userEmail")
	log.GLogger.Info("ReqExpense query: %d, %s", uID, uEmail)
	// 获取emp_info
	employee := new(oa.Employee)
	services.Slave().Preload("Department").Preload("Department.Leader").Take(employee, "email = ?", uEmail)
	if employee.ID == 0 {
		d.ErrorOK("未找到员工信息")
	}
	log.GLogger.Info("employee: %+v", employee)

	param := new(oa.DeviceApply)
	err := json.Unmarshal(d.Ctx.Input.RequestBody, param)
	if err != nil {
		log.GLogger.Error("parse DeviceApply info err:%s", err.Error())
		d.ErrorOK(MsgInvalidParam)
	}

	if param.DeviceID <= 0 {
		d.ErrorOK("need DeviceID")
	}
	if param.EngagementCode == "" {
		d.ErrorOK("need EngagementCode")
	}
	if param.Project == "" {
		d.ErrorOK("need project")
	}
	if param.LeaderId <= 0 {
		d.ErrorOK("need leader_id")
	}

	param.EmpID = int(employee.ID)
	param.EName = employee.Name
	param.ApplicationDate = time.Now()
	param.Status = models.FlowNA

	tx := services.Slave().Begin()
	// 创建申请记录
	err = tx.Create(param).Error
	if err != nil {
		log.GLogger.Error("create req device err:%s", err.Error())
		tx.Rollback()
		d.ErrorOK(MsgServerErr)
	}

	// 执行设备工作流
	err = services.ReqDeviceApply(tx, int(param.ID), uID, param.LeaderId)
	if err != nil {
		log.GLogger.Error("services req device err:%s", err.Error())
		tx.Rollback()
		d.ErrorOK(MsgServerErr)
	}
	tx.Commit()
	d.Correct(param)
}
