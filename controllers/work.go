/*
* Auth : acer
* Desc : 加班，请假
* Time : 2020/9/12 23:39
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
	err = services.ReqOvertime(tx, int(param.ID), uID, int(employee.Department.Leader.ID))
	if err != nil {
		log.GLogger.Error("req overtime err:%s", err.Error())
		tx.Rollback()
		w.ErrorOK(MsgServerErr)
	}
	tx.Commit()
	w.Correct(param)
}
