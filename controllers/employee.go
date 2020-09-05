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
