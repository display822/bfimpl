/**
* @author : chen lie
* @description : 客户 crud
* @date   : 2020-06-30 17:47
 */

package controllers

import (
	"bfimpl/models"
	"bfimpl/services"
	"bfimpl/services/log"
	"encoding/json"

	"github.com/jinzhu/gorm"
)

type ClientController struct {
	BaseController
}

// @Title 新增客户
// @Description 新增客户
// @Param	name		body	string	true	"客户名称"
// @Param	number		body	string	true	"编号"
// @Param	type		body	int		true	"0内部1外部"
// @Param	level		body	string	true	"S,A,B"
// @Param	saleId		body	int		true	"销售id"
// @Param	mainManageId	body	int	true	"主客户服务经理id"
// @Param	subManageId		body	int	true	"副客户服务经理id"
// @Success 200 {object} models.Client
// @Failure 500 server err
// @router / [post]
func (c *ClientController) AddClient() {
	param := new(models.Client)

	err := json.Unmarshal(c.Ctx.Input.RequestBody, param)
	if err != nil {
		c.ErrorOK(MsgInvalidParam)
	}
	ok, _ := c.valid.Valid(param)
	if !ok {
		log.GLogger.Error("%s:%s", c.valid.Errors[0].Field, c.valid.Errors[0].Message)
		c.ErrorOK(MsgInvalidParam)
	}
	if param.Type != 0 && param.Type != 1 {
		c.ErrorOK("invalid client type")
	}
	tmp := new(models.Client)
	err = services.Slave().Where("number =?", param.Number).First(tmp).Error
	if err != gorm.ErrRecordNotFound {
		c.ErrorOK("客户编号重复")
	}
	err = services.Slave().Create(param).Error
	if err != nil {
		c.ErrorOK(err.Error())
	}
	c.Correct(param)
}

// @Title 客户列表
// @Description 客户列表
// @Success 200 {object} []models.Client
// @Failure 500 server err
// @router /list [get]
func (c *ClientController) GetClients() {
	clients := make([]models.Client, 0)
	services.Slave().Model(models.Client{}).Find(&clients)

	c.Correct(clients)
}
