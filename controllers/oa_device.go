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
	"sort"
	"strconv"
	"strings"
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

	_, ok := oa.DeviceCategoryMap[param.DeviceCategory]
	if !ok {
		d.ErrorOK("param device_category error")
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

	d.Correct(param)
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
	userID, _ := d.GetInt("userID")
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
		Preload("DeviceApplys", func(db *gorm.DB) *gorm.DB {
			return db.Where("status = ?", models.FlowApproved).Or("status =?", models.FlowNA)
		}).
		Preload("DeviceApply", func(db *gorm.DB) *gorm.DB {
			return db.Where("status = ?", models.FlowReceived)
		}).
		Find(&list).Limit(-1).Offset(-1).Count(&resp.Total)

	for _, item := range list {
		item.CanApply = true
		if item.IsApply == 1 {
			item.CanApply = false
		}
		if item.DeviceStatus != models.DeviceFree {
			item.CanApply = false
		}
		for _, a := range item.DeviceApplys {
			// if a.Status != models.FlowReceived && a.Status != models.FlowRevoked {
			// 	item.CanApply = false
			// }
			if a.EmpID == userID {
				item.CanApply = false
				break
			}
		}

		// var deviceApplys []*oa.DeviceApply
		// for _, a := range item.DeviceApplys {
		// 	if a.Status == models.FlowApproved || a.Status == models.FlowNA {
		// 		deviceApplys = append(deviceApplys, a)
		// 	}
		// }
		//
		// item.DeviceApplys = deviceApplys
	}
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

	_, ok := oa.DeviceCategoryMap[param.DeviceCategory]
	if !ok {
		d.ErrorOK("param device_category error")
	}

	var device oa.Device
	err = services.Slave().Where("id = ?", param.ID).Find(&device).Error
	if err != nil {
		log.GLogger.Error("get device err:%s", err.Error())
		d.ErrorOK(MsgServerErr)
	}

	// 占用中不能改
	if device.DeviceStatus == models.DevicePossessed {
		d.ErrorOK("设备占用中")
	}
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
	projects := make([]*oa.EngagementCode, 0)
	//获取emp_info
	employee := new(oa.Employee)
	services.Slave().Preload("Department").Preload("Department.Leader").Take(employee, "email = ?", uEmail)
	if employee.ID == 0 {
		d.ErrorOK("未找到员工信息")
	}

	if employee.Department.ID == 0 {
		d.Correct(projects)
	}
	//查询部门下项目list

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

// @Title 领用设备项目
// @Description 领用设备项目
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /distribution/projects [get]
func (d *DeviceController) GetDistributionProjects() {
	desc := d.GetString("desc")
	uid := d.GetString("user_id")
	projects := make([]*oa.EngagementCode, 0)
	//获取emp_info
	employee := new(oa.Employee)
	services.Slave().Preload("Department").Preload("Department.Leader").Take(employee, "id = ?", uid)
	if employee.ID == 0 {
		d.ErrorOK("未找到员工信息")
	}

	if employee.Department.ID == 0 {
		d.Correct(projects)
	}
	//查询部门下项目list

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
// @Param	category	query	string	false	"类别"
// @Param	search	query	string	false	"search"
// @Success 200 {object} []oa.DeviceApply
// @Failure 500 server internal err
// @router /apply [get]
func (d *DeviceController) ListApply() {
	pageSize, _ := d.GetInt("pagesize", 10)
	pageNum, _ := d.GetInt("pagenum", 1)
	userType, _ := d.GetInt("userType", 0)
	myReq, _ := d.GetBool("myreq", false)
	myTodo, _ := d.GetBool("mytodo", false)
	status := d.GetString("status")
	userEmail := d.GetString("userEmail")
	category := d.GetString("category")
	search := d.GetString("search")
	log.GLogger.Info("params", userEmail, userType, myReq, status, pageNum, pageSize)

	employee := new(oa.Employee)
	services.Slave().Where("email = ?", userEmail).First(employee)
	log.GLogger.Info("employee: %+v", employee)

	deviceApplys := make([]*oa.DeviceApply, 0)
	query := services.Slave().Debug().Select("device_applies.*").
		Joins("LEFT JOIN devices ON device_applies.device_id = devices.id and devices.deleted_at is NULL")

	if category != "" {
		query = query.Where("devices.device_category = ?", category)
	}
	if status != "" {
		if status == models.FlowApproved {
			query = query.Where("device_applies.status in (?)", []string{status, models.FlowDistributed})
		} else {
			query = query.Where("device_applies.status = ?", status)
		}
	}
	if search != "" {
		query = query.Where("devices.device_category like ?", fmt.Sprintf("%%%s%%", search))
		query = query.Or("devices.device_code like ?", fmt.Sprintf("%%%s%%", search))
		query = query.Or("device_applies.id like ?", fmt.Sprintf("%%%s%%", search))
		query = query.Or("devices.brand like ?", fmt.Sprintf("%%%s%%", search))
		query = query.Or("devices.device_model like ?", fmt.Sprintf("%%%s%%", search))
		query = query.Or("device_applies.project like ?", fmt.Sprintf("%%%s%%", search))
	}

	var resp struct {
		Total int               `json:"total"`
		List  []*oa.DeviceApply `json:"list"`
	}

	if myReq {
		query = query.Where("device_applies.emp_id = ?", employee.ID)
	}

	eIDs := make([]int, 0)

	if myTodo {
		userID, _ := d.GetInt("userID", 0)
		log.GLogger.Info("userID：%d", userID)
		ids := make([]oa.EntityID, 0)
		var s []string
		if d.GetString("todostatus") != "" {
			s = oa.DeviceTodoStatusLeaderMap[d.GetString("todostatus")]
		}
		if len(s) == 0 {
			services.Slave().Debug().Raw("select w.entity_id from workflows w,workflow_nodes wn where w.id = "+
				"wn.workflow_id and w.workflow_definition_id = ? and operator_id = ? and wn.status <> ?"+
				" and wn.node_seq != 1 order by w.entity_id desc", services.GetFlowDefID(services.Device), userID, models.FlowHide).Scan(&ids)
		} else {
			services.Slave().Debug().Raw("select w.entity_id from workflows w,workflow_nodes wn where w.id = "+
				"wn.workflow_id and w.workflow_definition_id = ? and operator_id = ? and wn.status in (?)"+
				" and wn.node_seq != 1 order by w.entity_id desc", services.GetFlowDefID(services.Device), userID, s).Scan(&ids)
		}
		for _, eID := range ids {
			eIDs = append(eIDs, eID.EntityID)
		}
		query = query.Where(eIDs)
		log.GLogger.Info("eid:%s", eIDs)
	}
	query.Limit(pageSize).Offset((pageNum - 1) * pageSize).Preload("Device").Order("device_applies.created_at desc").Find(&deviceApplys).Limit(-1).Offset(-1).Count(&resp.Total)

	for _, deviceApply := range deviceApplys {
		if (deviceApply.Status == models.FlowApproved && deviceApply.Device.DeviceApplyID == int(deviceApply.ID)) || deviceApply.Status == models.FlowDistributed {
			deviceApply.CanReceive = true
		}

		if deviceApply.Status == models.FlowNA || deviceApply.Status == models.FlowApproved {
			deviceApply.CanRevoke = true
		}
	}
	resp.List = deviceApplys
	d.Correct(resp)
}

// @Title 申请设备列表
// @Description 申请设备列表
// @Success 200 {object} oa.DeviceApply
// @Failure 500 server internal err
// @router /apply/:id [get]
func (d *DeviceController) GetApply() {
	id, _ := d.GetInt(":id", 0)
	deviceApply := new(oa.DeviceApply)
	services.Slave().Debug().Preload("Device").Take(deviceApply, "id = ?", id)
	// oID 查询 workflow
	workflow := new(oa.Workflow)
	services.Slave().Model(oa.Workflow{}).Where("workflow_definition_id = ? and entity_id = ?",
		services.GetFlowDefID(services.Device), id).Preload("Nodes").Preload("Nodes.User").
		Preload("Elements").First(workflow)
	log.GLogger.Info("workflow", workflow)
	if len(workflow.Nodes) != 2 {
		d.ErrorOK("工作流配置错误")
	}
	var resp struct {
		Info     *oa.DeviceApply `json:"info"`
		WorkFlow *oa.Workflow    `json:"work_flow"`
	}
	resp.Info = deviceApply
	resp.WorkFlow = workflow

	d.Correct(resp)
}

// @Title 设备已申请列表
// @Description 设备已申请列表
// @Success 200 {object} oa.Device
// @Failure 500 server internal err
// @router /:id/apply [get]
func (d *DeviceController) ListDeviceApply() {
	id, _ := d.GetInt(":id")
	var deviceApplys []*oa.DeviceApply
	services.Slave().Where("device_id = ?", id).Where("status = ?", models.FlowApproved).
		Find(&deviceApplys)
	d.Correct(deviceApplys)
}

// @Title 撤回申请设备
// @Description 撤回申请设备
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /apply/:id/revoke [put]
func (d *DeviceController) RevokeDevice() {
	id, _ := d.GetInt(":id")
	userID, _ := d.GetInt("userID")
	deviceApply := new(oa.DeviceApply)
	services.Slave().Debug().Preload("Device").Preload("Employee").Take(deviceApply, "id = ?", id)
	log.GLogger.Info("deviceApply:%+v", deviceApply)

	if deviceApply.EmpID != userID {
		d.ErrorOK("没有权限")
	}

	if deviceApply.Status != models.FlowNA && deviceApply.Status != models.FlowApproved {
		d.ErrorOK("不可撤销")
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
// @router /apply/approval [put]
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

// @Title 设备领用
// @Description 设备领用
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /:id/receive [put]
func (d *DeviceController) ReceiveDevice() {
	userID, _ := d.GetInt("userID")
	userName := d.GetString("userName")
	id, _ := d.GetInt(":id")
	tx := services.Slave().Begin()
	var device oa.Device
	err := tx.Where("id = ?", id).Preload("DeviceApply").Find(&device).Error
	if err != nil {
		log.GLogger.Error("get device err:%s", err.Error())
		tx.Rollback()
		d.ErrorOK(MsgServerErr)
	}
	log.GLogger.Info("device", device)
	if device.DeviceApply.EmpID != userID {
		d.ErrorOK("没有领用权限")
	}

	if device.DeviceApply.Status == models.FlowReceived {
		d.ErrorOK("已领用")
	}

	device.DeviceStatus = models.DevicePossessed
	device.DeviceApply.Status = models.FlowReceived
	device.DeviceApply.ReceiveDate = time.Now()

	err = tx.Save(&device).Error
	if err != nil {
		log.GLogger.Error("save device err:%s", err.Error())
		tx.Rollback()
		d.ErrorOK(MsgServerErr)
	}

	// 添加记录
	deviceRequisition := oa.DeviceRequisition{
		DeviceID:              id,
		AssociateEmployeeID:   userID,
		AssociateEmployeeName: userName,
		OperatorCategory:      models.DeviceOutgoing,
		OperatorID:            device.DeviceApply.OutgoingOperatorID,
		OperatorName:          device.DeviceApply.OutgoingOperatorName,
		Comment:               "",
	}
	err = tx.Create(&deviceRequisition).Error
	if err != nil {
		log.GLogger.Error("create deviceRequisition err:%s", err.Error())
		tx.Rollback()
		d.ErrorOK(MsgServerErr)
	}

	tx.Commit()
	d.Correct("")
}

// @Title 设备归还
// @Description 设备归还
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /:id/return [put]
func (d *DeviceController) ReturnDevice() {
	userID, _ := d.GetInt("userID")
	userName := d.GetString("userName")
	userType, _ := d.GetInt("userType")
	status := d.GetString("status")

	id, _ := d.GetInt(":id")
	if userType != models.UserIT && userType != models.UserFront && userType != models.UserFinance {
		d.ErrorOK("没有权限")
	}
	tx := services.Slave().Begin()
	var device oa.Device
	err := tx.Where("id = ?", id).Preload("DeviceApply").Find(&device).Error
	if err != nil {
		log.GLogger.Error("get device err:%s", err.Error())
		tx.Rollback()
		d.ErrorOK(MsgServerErr)
	}
	log.GLogger.Info("device", device)

	if device.DeviceStatus != models.DevicePossessed {
		d.ErrorOK("设备状态错误")
	}

	var s string
	if status == "0" {
		s = models.DeviceFree
	} else if status == "1" {
		s = models.DeviceFixing
	} else if status == "2" {
		s = models.DeviceScrap
	}
	device.DeviceStatus = s
	device.DeviceApplyID = 0
	device.DeviceApply.Status = models.DeviceReturn

	err = tx.Save(&device).Error
	if err != nil {
		log.GLogger.Error("save device err:%s", err.Error())
		tx.Rollback()
		d.ErrorOK(MsgServerErr)
	}

	// 添加记录
	deviceRequisition := oa.DeviceRequisition{
		DeviceID:              id,
		AssociateEmployeeID:   device.DeviceApply.EmpID,
		AssociateEmployeeName: device.DeviceApply.EName,
		OperatorCategory:      models.DeviceReturn,
		OperatorID:            userID,
		OperatorName:          userName,
		Comment:               "",
	}
	err = tx.Create(&deviceRequisition).Error
	if err != nil {
		log.GLogger.Error("create deviceRequisition err:%s", err.Error())
		tx.Rollback()
		d.ErrorOK(MsgServerErr)
	}

	tx.Commit()
	d.Correct("")
}

// @Title 设备分配
// @Description 设备分配
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /:id/distribution [put]
func (d *DeviceController) DistributionDevice() {
	userID, _ := d.GetInt("userID")
	userName := d.GetString("userName")
	userType, _ := d.GetInt("userType")
	id, _ := d.GetInt(":id")

	employeeID, _ := d.GetInt("employee_id")
	employeeName := d.GetString("employee_name")
	engagementCode := d.GetString("engagement_code")
	project := d.GetString("project")
	if employeeID <= 0 {
		d.ErrorOK("need employee_id")
	}
	if employeeName == "" {
		d.ErrorOK("need employee_name")
	}
	if engagementCode == "" {
		d.ErrorOK("need engagement_code")
	}
	if project == "" {
		d.ErrorOK("need project")
	}
	if userType != models.UserIT && userType != models.UserFront && userType != models.UserFinance {
		d.ErrorOK("没有权限")
	}
	tx := services.Slave().Begin()

	var device oa.Device
	err := tx.Where("id = ?", id).Find(&device).Error
	if err != nil {
		log.GLogger.Error("get device err:%s", err.Error())
		tx.Rollback()
		d.ErrorOK(MsgServerErr)
	}
	if device.DeviceStatus != models.DeviceFree {
		d.ErrorOK("设备不为空闲，不可借出")
	}

	if device.DeviceApplyID != 0 {
		d.ErrorOK("不可重复借出")
	}

	// 创建申请单
	deviceApply := oa.DeviceApply{
		DeviceID:             id,
		EngagementCode:       engagementCode,
		EmpID:                employeeID,
		EName:                employeeName,
		Status:               models.FlowDistributed,
		Project:              project,
		ApplicationDate:      time.Now(),
		OutgoingOperatorID:   userID,
		OutgoingOperatorName: userName,
		OutgoingTime:         models.Time(time.Now()),
	}

	err = tx.Create(&deviceApply).Error
	if err != nil {
		log.GLogger.Error("create device_apply err:%s", err.Error())
		tx.Rollback()
		d.ErrorOK(MsgServerErr)
	}

	device.DeviceApplyID = int(deviceApply.ID)

	err = tx.Save(&device).Error
	if err != nil {
		log.GLogger.Error("save device err:%s", err.Error())
		tx.Rollback()
		d.ErrorOK(MsgServerErr)
	}

	tx.Commit()
	d.Correct("")
}

// @Title 设备借出
// @Description 设备借出
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /:id/outgoing [put]
func (d *DeviceController) OutgoingDevice() {
	userID, _ := d.GetInt("userID")
	userName := d.GetString("userName")
	userType, _ := d.GetInt("userType")
	if userType != models.UserIT && userType != models.UserFront && userType != models.UserFinance {
		d.ErrorOK("没有权限")
	}
	id, _ := d.GetInt(":id")
	deviceApplyID, _ := d.GetInt("device_apply_id")

	if deviceApplyID <= 0 {
		d.ErrorOK("need device_apply_id")
	}

	// 判断设备状态
	tx := services.Slave().Begin()
	var device oa.Device
	err := tx.Where("id = ?", id).Find(&device).Error
	if err != nil {
		log.GLogger.Error("get device err:%s", err.Error())
		tx.Rollback()
		d.ErrorOK(MsgServerErr)
	}
	if device.DeviceStatus != models.DeviceFree {
		d.ErrorOK("设备不为空闲，不可借出")
	}

	if device.DeviceApplyID != 0 {
		d.ErrorOK("不可重复借出")
	}
	// 自己申请 将设备关联申请单
	device.DeviceApplyID = deviceApplyID

	err = tx.Save(&device).Error
	if err != nil {
		log.GLogger.Error("get device err:%s", err.Error())
		tx.Rollback()
		d.ErrorOK(MsgServerErr)
	}

	var deviceApply oa.DeviceApply
	err = tx.Where("id = ?", deviceApplyID).Find(&deviceApply).Error
	if err != nil {
		log.GLogger.Error("get device err:%s", err.Error())
		tx.Rollback()
		d.ErrorOK(MsgServerErr)
	}
	deviceApply.OutgoingOperatorID = userID
	deviceApply.OutgoingOperatorName = userName
	deviceApply.OutgoingTime = models.Time(time.Now())
	err = tx.Save(&deviceApply).Error
	if err != nil {
		log.GLogger.Error("get device err:%s", err.Error())
		tx.Rollback()
		d.ErrorOK(MsgServerErr)
	}
	tx.Commit()
	d.Correct("")
}

// @Title 申请设备基本信息
// @Description 申请设备基本信息
// @Success 200 {object} oa.Device
// @Failure 500 server internal err
// @router /apply/info [get]
func (d *DeviceController) ApplyInfo() {
	uID, _ := d.GetInt("userID", 0)
	deviceID, _ := d.GetInt("device_id")
	employeeID, _ := d.GetInt("employee_id")
	log.GLogger.Info("ReqExpense query: %d, %d", uID, employeeID)
	// 获取emp_info
	employee := new(oa.Employee)
	services.Slave().Preload("Department").Take(employee, "id = ?", employeeID)
	if employee.ID == 0 {
		d.ErrorOK("未找到员工信息")
	}

	device := new(oa.Device)
	services.Slave().Where("id =?", deviceID).Find(device)

	var deviceApplys []*oa.DeviceApply
	services.Slave().Where("status =?", models.FlowReceived).
		Where("emp_id =?", employeeID).
		Preload("Device").
		Find(&deviceApplys)

	var collectDevices []string
	for _, item := range deviceApplys {
		if item.Device != nil {
			collectDevices = append(collectDevices, item.Device.DeviceName)
		}
	}
	res := oa.DeviceApplyInfo{
		Employee:       employee,
		Device:         device,
		CollectDevices: strings.Join(collectDevices, ","),
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

// @Title 员工下设备借出列表
// @Description 易耗品借出列表
// @Param	pagenum	    query	int	false	"页码"
// @Param	pagesize	query	int	false	"页数"
// @Success 200 {object} oa.Device
// @Failure 500 server internal err
// @router /employee/outgoing [get]
func (d *DeviceController) ListOutgoingByEmployee() {
	uID, _ := d.GetInt("userID")
	log.GLogger.Info("uID:%d", uID)

	var deviceApplys []*oa.DeviceApply
	services.Slave().Where("status =?", models.FlowReceived).
		Where("emp_id =?", uID).
		Preload("Device").
		Order("receive_date desc").
		Find(&deviceApplys)

	d.Correct(deviceApplys)
}

// @Title 员工下归还列表
// @Description 员工下归还列表
// @Success 200 {string} ""
// @Failure 500 server internal err
// @router /employee/return [get]
func (d *DeviceController) ListReturnByEmployee() {
	uID, _ := d.GetInt("userID")
	log.GLogger.Info("uID:%d", uID)
	var returns []*forms.Return

	// 先获取易耗品id
	var deviceApplys []*oa.DeviceApply
	services.Slave().Where("status = ?", models.DeviceReturn).
		Where("emp_id = ?", uID).
		Preload("Device").
		Find(&deviceApplys)

	log.GLogger.Info("deviceApplys:%d", deviceApplys)
	for _, a := range deviceApplys {
		returns = append(returns, &forms.Return{
			ID:        a.Device.DeviceCode,
			Name:      a.Device.DeviceName,
			CreatedAt: a.CreatedAt,
		})
	}

	// 先获取易耗品id
	var articleRequisitions []*oa.LowPriceArticleRequisition
	services.Slave().Where("operator_category = ?", models.DeviceReturn).
		Where("associate_employee_id = ?", uID).
		Preload("LowPriceArticle").
		Find(&articleRequisitions)

	log.GLogger.Info("articleRequisitions:%d", articleRequisitions)

	for _, a := range articleRequisitions {
		returns = append(returns, &forms.Return{
			ID:        strconv.Itoa(a.LowPriceArticleID),
			Name:      a.LowPriceArticle.LowPriceArticleName,
			CreatedAt: a.CreatedAt,
		})
	}

	sort.Sort(forms.ReturnByCreatedAt(returns))

	d.Correct(returns)
}
