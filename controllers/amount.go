/*
* Auth : acer
* Desc : 额度 crud
* Time : 2020/6/30 21:52
 */

package controllers

import (
	"bfimpl/models"
	"bfimpl/services"
	"bfimpl/services/log"
	"encoding/json"
	"fmt"
	"time"
)

type AmountController struct {
	BaseController
}

// @Title 新增额度
// @Description 新增服务
// @Param	clientId	body	int	true	"客户id"
// @Param	serviceId	body	int	true	"服务id"
// @Param	amount		body	int	true	"额度"
// @Param	deadline	body	string	true	"到期日期"
// @Param	orderNumber	body	string	true	"订单编号"
// @Param	remark		body	string	true	"备注说明"
// @Success 200 {object} models.Amount
// @Failure 500 server err
// @router / [post]
func (a *AmountController) AddAmount() {
	param := new(models.Amount)
	err := json.Unmarshal(a.Ctx.Input.RequestBody, param)
	if err != nil || param.ID != 0 || param.OrderNumber == "" ||
		param.Amount <= 0 || time.Time(param.Deadline).Before(time.Now()) {
		a.ErrorOK(MsgInvalidParam)
	}
	ok, _ := a.valid.Valid(param)
	if !ok {
		log.GLogger.Error("%s:%s", a.valid.Errors[0].Field, a.valid.Errors[0].Message)
		a.ErrorOK(MsgInvalidParam)
	}
	tx := services.Slave().Begin()
	err = tx.Create(param).Error
	if err != nil {
		tx.Rollback()
		a.ErrorOK(MsgServerErr)
	}
	amountLog := new(models.AmountLog)
	amountLog.AmountId = int(param.ID)
	amountLog.Change = param.Amount
	amountLog.Desc = fmt.Sprintf("订单采买，订单编号:%s", param.OrderNumber)
	amountLog.RealTime = models.Time(time.Now())
	amountLog.Type = "buy"
	tx.Create(amountLog)
	if tx.Commit().Error != nil {
		tx.Rollback()
		a.ErrorOK(MsgServerErr)
	}
	a.Correct(param)
}

// @Title 查询客户的额度列表
// @Description 查询客户的额度列表
// @Param	clientId	query	int	true	"客户id"
// @Success 200 {object} models.RspAmount
// @Failure 500 server err
// @router /list [get]
func (a *AmountController) GetAmounts() {
	clientId, _ := a.GetInt("clientId")
	res := make([]models.RspAmount, 0)
	services.Slave().Raw("select s.service_name, a.amount, a.deadline  "+
		"from amounts a, services s where a.service_id = s.id and a.id = ?", clientId).Scan(&res)

	a.Correct(res)
}

// @Title 查询客户的额度历史
// @Description 查询客户的额度历史
// @Param	amountId	query	int	true	"额度id"
// @Success 200 {object} models.RspAmountLog
// @Failure 500 server err
// @router /log [get]
func (a *AmountController) GetAmountLogs() {
	amountId, _ := a.GetInt("amountId")
	res := make([]models.RspAmountLog, 0)
	services.Slave().Raw("SELECT al.real_time,s.service_name,c.name,al.change,al.desc,"+
		"a.remark FROM amounts a,services s,amount_logs al,clients c WHERE "+
		"a.id = al.amount_id AND a.client_id = c.id AND a.service_id = s.id AND a.id = ?", amountId).Scan(&res)

	a.Correct(res)
}
