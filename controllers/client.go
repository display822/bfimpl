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
		c.ErrorOK("客户编号已存在")
	}
	err = services.Slave().Create(param).Error
	if err != nil {
		c.ErrorOK(err.Error())
	}
	c.Correct(param)
}

// @Title 客户列表
// @Description 客户列表
// @Param	saleId		query	int		false	"销售id"
// @Param	manageId	query	int		false	"客服经理id"
// @Success 200 {object} []models.Client
// @Failure 500 server err
// @router /list [get]
func (c *ClientController) GetClients() {

	uID, _ := c.GetInt("userID", 0)
	if uID == 0 {
		c.ErrorOK("need user id")
	}
	userType, _ := c.GetInt("userType", 0)
	clients := make([]models.Client, 0)
	db := services.Slave().Model(models.Client{})
	if userType == 2 {
		//销售
		db = db.Where("sale_id = ?", uID)
	} else if userType == 3 {
		db = db.Where("main_manage_id = ? or sub_manage_id = ?", uID, uID)
	}
	db.Find(&clients)

	c.Correct(clients)
}

// @Title 修改客户
// @Description 修改客户
// @Param	name		body	string	true	"客户名称"
// @Param	number		body	string	true	"编号"
// @Param	type		body	int		true	"0内部1外部"
// @Param	level		body	string	true	"S,A,B"
// @Param	saleId		body	int		true	"销售id"
// @Param	mainManageId	body	int	true	"主客户服务经理id"
// @Param	subManageId		body	int	true	"副客户服务经理id"
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /:id [put]
func (c *ClientController) UpdateClient() {
	param := new(models.Client)
	id, _ := c.GetInt(":id")
	err := json.Unmarshal(c.Ctx.Input.RequestBody, param)
	if err != nil {
		c.ErrorOK(MsgInvalidParam)
	}
	services.Slave().Model(models.Client{}).Where("id = ?", id).Updates(map[string]interface{}{
		"name":         param.Name,
		"number":       param.Number,
		"type":         param.Type,
		"level":        param.Level,
		"saleId":       param.SaleId,
		"mainManageId": param.MainManageId,
		"subManageId":  param.SubManageId,
	})
	c.Correct("")
}

// @Title 获取客户
// @Description 获取客户
// @Param	id		path	int	true	"客户id"
// @Success 200 {string} models.Client
// @Failure 500 server err
// @router /:id [get]
func (c *ClientController) GetClient() {
	id, _ := c.GetInt(":id")
	var client models.Client
	err := services.Slave().Table("clients").Take(&client, "id = ?", id).Error
	if err != nil {
		c.ErrorOK(MsgServerErr)
	}
	users := make([]models.User, 0)
	services.Slave().Table("users").Where("id in (?,?,?)",
		client.SaleId, client.MainManageId, client.SubManageId).Find(&users)
	result := models.RspClient{
		Client: client,
	}
	for i := range users {
		if users[i].ID == uint(client.SaleId) {
			result.Sale = users[i]
		}
		if users[i].ID == uint(client.MainManageId) {
			result.Manager = users[i]
		}
		if users[i].ID == uint(client.SubManageId) {
			result.SubManager = users[i]
		}
	}

	c.Correct(result)
}
