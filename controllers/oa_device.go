/**
* @author : yi.zhang
* @description : controllers 描述
* @date   : 2021-01-07 18:20
 */

package controllers

import (
	"bfimpl/models"
	"bfimpl/models/oa"
	"bfimpl/services"
	"bfimpl/services/log"
	"encoding/json"
	"fmt"
	"time"
)

type DeviceController struct {
	BaseController
}

// @Title 创建设备
// @Description 创建设备
// @Success 200 {string} ""
// @Failure 500 server internal err
// @router / [post]
func (d *DeviceController) Create() {
	// 验证员工身份 (7，8，9)
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

	param.IngoingTime = models.Time(time.Now())
	err = services.Slave().Create(param).Error
	if err != nil {
		log.GLogger.Error("create employee err：%s", err.Error())
		d.ErrorOK(MsgServerErr)
	}

	// TODO 添加活动记录

	d.Correct("")
}

// @Title 设备列表
// @Description 设备列表
// @Param	pagenum	    query	int	false	"页码"
// @Param	pagesize	query	int	false	"页数"
// @Param	device_category	query	string	false	"设备分类"
// @Param	device_status	query	string	false	"设备状态"
// @Param	keyword	query	string	false	"搜索关键词"
// @Success 200 {string} ""
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
	db.Limit(pageSize).Offset((pageNum - 1) * pageSize).Order("created_at desc").Find(&list).Limit(-1).Offset(-1).Count(&resp.Total)
	resp.List = list

	d.Correct(resp)
}

// @Title 设备详情
// @Description 设备详情
// @Success 200 {string} ""
// @Failure 500 server internal err
// @router / [get]
func (d *DeviceController) Get() {
	dID, _ := d.GetInt(":id", 0)
	var device oa.Device
	services.Slave().Where("id = ?", dID).Find(&device)
	d.Correct(device)
}

// @Title 设备更新
// @Description 设备更新
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
// @Param	body	    body	oa.Device	true	"设备"
// @Success 200 {object} oa.Device
// @Failure 500 server internal err
// @router /apply [get]
func (d *DeviceController) ApplyList() {

}

// @Title 已申请员工
// @Description 已申请员工
// @Param	body	    body	oa.Device	true	"设备"
// @Success 200 {object} oa.Device
// @Failure 500 server internal err
// @router /apply/employee [get]
func (d *DeviceController) ApplyEmployee() {
	//dID := d.GetString("did")
}

// @Title 审批申请设备
// @Description 审批申申请设备
// @Param	body	body	oa.Device	true
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /apply [put]
func (d *DeviceController) ApprovalDevice() {

}

// @Title 设备借出
// @Description 设备借出
// @Param	body	body	oa.Device	true
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
// @Param	body	    body	oa.Device	true	"设备"
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
// @Param	body	    body	oa.Device	true	"设备"
// @Success 200 {object} oa.Device
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
