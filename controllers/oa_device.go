/**
* @author : yi.zhang
* @description : controllers 描述
* @date   : 2021-01-07 18:20
 */

package controllers

import (
	"bfimpl/models"
	"bfimpl/models/oa"
	"bfimpl/services/log"
	"encoding/json"
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

	d.Correct("")
}

// @Title 设备列表
// @Description 设备列表
// @Success 200 {string} ""
// @Failure 500 server internal err
// @router / [get]
func (d *DeviceController) List() {
	d.Correct("")
}

// @Title 设备详情
// @Description 设备详情
// @Success 200 {string} ""
// @Failure 500 server internal err
// @router / [get]
func (d *DeviceController) Get() {
	d.Correct("")
}

// @Title 设备更新
// @Description 设备更新
// @Success 200 {string} ""
// @Failure 500 server internal err
// @router / [get]
func (d *DeviceController) Put() {
	d.Correct("")
}
