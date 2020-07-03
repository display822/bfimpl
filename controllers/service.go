/**
* @author : chen lie
* @description : service crud
* @date   : 2020-06-30 17:48
 */

package controllers

import (
	"bfimpl/models"
	"bfimpl/services"
	"bfimpl/services/log"
	"encoding/json"
	"strconv"
	"strings"
)

type ServiceController struct {
	BaseController
}

// @Title 新增服务
// @Description 新增服务
// @Param	serviceName	body	string	true	"服务名"
// @Param	state	body	string	true	"0启用1禁用"
// @Param	use		body	string	true	"1可实施2可转换"
// @Success 200 {object} models.Service
// @Failure 500 server err
// @router / [post]
func (s *ServiceController) AddService() {
	param := new(models.Service)

	err := json.Unmarshal(s.Ctx.Input.RequestBody, param)
	if err != nil {
		s.ErrorOK(MsgInvalidParam)
	}
	b, _ := s.valid.Valid(param)
	if !b {
		log.GLogger.Error("%s:%s", s.valid.Errors[0].Field, s.valid.Errors[0].Message)
		s.ErrorOK(MsgInvalidParam)
	}
	//查询数据条数
	var count int
	services.Slave().Model(models.Service{}).Count(&count)
	param.Sort = count + 1
	err = services.Slave().Create(param).Error
	if err != nil {
		s.ErrorOK(err.Error())
	}
	s.Correct(param)
}

// @Title 修改服务
// @Description 修改服务
// @Param	serviceName	body	string	true "服务名"
// @Param	state	body	string	true	"0启用1禁用"
// @Param	use		body	string	true	"1可实施2可转换"
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /:id [put]
func (s *ServiceController) UpdateService() {
	param := new(models.Service)
	id, _ := s.GetInt(":id")
	err := json.Unmarshal(s.Ctx.Input.RequestBody, param)
	if err != nil {
		s.ErrorOK(MsgInvalidParam)
	}
	services.Slave().Model(models.Service{}).Where("id = ?", id).Updates(map[string]interface{}{
		"service_name": param.ServiceName,
		"use":          param.Use,
		"state":        param.State,
	})
	s.Correct("")
}

// @Title 移动服务
// @Description 交换两个服务位置
// @Param	ids	query	string	true	"两个服务id，逗号隔开"
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /switch [put]
func (s *ServiceController) SwitchService() {
	param := s.GetString("ids")
	ids := strings.Split(param, ",")
	if len(ids) < 2 {
		s.ErrorOK(MsgInvalidParam)
	}
	id1, e1 := strconv.Atoi(ids[0])
	id2, e2 := strconv.Atoi(ids[1])
	if e1 != nil || e2 != nil {
		s.ErrorOK(MsgInvalidParam)
	}
	tx := services.Slave().Begin()
	var service1, service2 models.Service
	tx.Take(&service1, "id = ?", id1)
	tx.Take(&service2, "id = ?", id2)
	sort1, sort2 := service1.Sort, service2.Sort
	tx.Table("services").Where("id = ?", id1).UpdateColumn("sort", sort2)
	tx.Table("services").Where("id = ?", id2).UpdateColumn("sort", sort1)
	err := tx.Commit().Error
	if err != nil {
		tx.Rollback()
		s.ErrorOK(err.Error())
	}
	s.Correct("")
}

// @Title 服务列表
// @Description 服务列表
// @Success 200 {object} []models.Service
// @Failure 500 server err
// @router /list [get]
func (s *ServiceController) GetServices() {
	srvs := make([]models.Service, 0)
	services.Slave().Model(models.Service{}).Find(&srvs)
	s.Correct(srvs)
}
