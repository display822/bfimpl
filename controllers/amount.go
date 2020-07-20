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
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
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
		log.GLogger.Error(err.Error())
		a.ErrorOK(MsgInvalidParam)
	}
	ok, _ := a.valid.Valid(param)
	if !ok {
		log.GLogger.Error("%s:%s", a.valid.Errors[0].Field, a.valid.Errors[0].Message)
		a.ErrorOK(MsgInvalidParam)
	}
	param.OrderNumber = strconv.FormatInt(time.Now().Unix(), 10)
	tx := services.Slave().Begin()
	err = tx.Create(param).Error
	if err != nil {
		tx.Rollback()
		a.ErrorOK(MsgServerErr)
	}
	createAmountLog(tx, param, "订单采买", models.Amount_Buy, param.Amount)
	if tx.Commit().Error != nil {
		tx.Rollback()
		a.ErrorOK(MsgServerErr)
	}
	a.Correct(param)
}

// @Title 查询客户的额度列表
// @Description 查询客户的额度列表
// @Param	clientId	query	int	true	"客户id"
// @Param	deadline	query	string	false	"过期时间，默认所有"
// @Param	use			query	int	false	"服务类型，不传返回所有，1可实施2可转换3可实施可转换"
// @Success 200 {object} models.RspAmount
// @Failure 500 server err
// @router /list [get]
func (a *AmountController) GetAmounts() {
	clientId, _ := a.GetInt("clientId")
	deadline := a.GetString("deadline")
	use, _ := a.GetInt("use", 0)
	if deadline == "" {
		deadline = "2020-07-01"
	}
	query := "select a.id id, a.amount amount,al.`change`,a.deadline deadline, a.service_id,al.type,s.service_name from amounts a," +
		"amount_logs al,services s where a.id = al.amount_id and a.client_id = ? and a.service_id = s.id and a.deadline > ? "
	if use != 0 {
		query += "and s.use = ? "
	}
	query += "order by a.service_id,a.deadline"
	// 查询额度log
	res := make([]models.ClientAmount, 0)
	if use != 0 {
		services.Slave().Raw(query, clientId, deadline, use).Scan(&res)
	} else {
		services.Slave().Raw(query, clientId, deadline).Scan(&res)
	}
	//统计log
	//service_id为key表示是否统计
	find := make(map[int]*models.RspAmount)
	ids := make([]int, 0)
	for _, ca := range res {
		if amount, ok := find[ca.ServiceId]; ok {
			amount.CalData(ca)
			find[amount.ServiceId] = amount
		} else {
			ids = append(ids, ca.ServiceId)
			rspAmount := new(models.RspAmount)
			rspAmount.ServiceId = ca.ServiceId
			rspAmount.ServiceName = ca.ServiceName
			rspAmount.Deadline = ca.Deadline
			rspAmount.CalData(ca)
			find[ca.ServiceId] = rspAmount
		}
	}
	data := make([]*models.RspAmount, 0)
	for _, i := range ids {
		data = append(data, find[i])
	}
	a.Correct(data)
}

// @Title 查询客户的额度历史
// @Description 查询客户的额度历史
// @Param	clientId	query	int	true	"客户id"
// @Param	serviceId	query	int	true	"服务id"
// @Success 200 {object} models.RspAmountLog
// @Failure 500 server err
// @router /log [get]
func (a *AmountController) GetAmountLogs() {
	clientId, _ := a.GetInt("clientId")
	serviceId, _ := a.GetInt("serviceId")
	res := make([]models.RspAmountLog, 0)
	services.Slave().Raw("SELECT al.real_time,s.service_name,a.order_number,al.change,al.desc,"+
		"al.remark,al.type FROM amounts a,amount_logs al,services s WHERE a.client_id = ? AND a.service_id = ?"+
		" AND a.id = al.amount_id AND a.service_id = s.id order by al.real_time desc", clientId, serviceId).Scan(&res)

	a.Correct(res)
}

// @Title 客户额度转换
// @Description 客户额度转换
// @Param	clientId	json	int	true	"客户id"
// @Param	sOutId		json	int	true	"转出服务id"
// @Param	sOutNum		json	int	true	"转出服务额度"
// @Param	sInId		json	int	true	"转入服务id"
// @Param	sInNum		json	int	true	"转入服务额度"
// @Param	remark		json	string	true	"备注说明"
// @Success 200 {string} ""
// @Failure 500 server err
// @router /switch [put]
func (a *AmountController) SwitchAmount() {
	param := new(models.ReqSwitchAmount)
	err := json.Unmarshal(a.Ctx.Input.RequestBody, param)
	if err != nil {
		a.ErrorOK(err.Error())
	}
	//转出转入额度小于0
	if param.SOutNum <= 0 || param.SInNum <= 0 {
		a.ErrorOK("额度填写错误")
	}
	if param.SOutId == param.SInId {
		a.ErrorOK("相同服务不可转换")
	}
	deadline := time.Now().Format(models.DateFormat)
	//查询状态启用,可转换,amount>0的 amounts
	aOut := make([]*models.AmountSimple, 0)
	services.Slave().Raw("select a.id,a.amount,a.order_number from amounts a,services s where a.service_id = s.id "+
		"and a.client_id = ? and s.id = ? and a.deadline > ? and s.state=0 and s.use >1 and a.amount >0 "+
		"order by deadline", param.ClientId, param.SOutId, deadline).Scan(&aOut)
	outSum := 0
	for _, t := range aOut {
		outSum += t.Amount
	}
	if outSum < param.SOutNum {
		a.ErrorOK("可转出额度不足")
	}
	// 查询转入amounts
	aIn := make([]*models.AmountSimple, 0)
	services.Slave().Raw("select a.id,a.amount,a.order_number from amounts a,services s where a.service_id = s.id "+
		"and a.client_id = ? and s.id = ? and a.deadline > ? and s.state=0 and s.use >1 "+
		"order by deadline", param.ClientId, param.SInId, deadline).Scan(&aIn)
	if len(aIn) == 0 {
		a.ErrorOK("不存在可转入的服务")
	}
	// ==============start convert===============
	// 额度转换关联字段
	refer := strconv.FormatInt(time.Now().UnixNano(), 10)
	tx := services.Slave().Begin()
	//转入aIn[0]
	err = tx.Model(models.Amount{}).Where("id = ?", aIn[0].Id).Update("amount", aIn[0].Amount+param.SInNum).Error
	msg := fmt.Sprintf("%s，订单编号:%s", "额度转换", aIn[0].OrderNumber)
	createAmountLogSimple(tx, aIn[0], msg, models.Amount_ConvIn, param.Remark, refer, param.SInNum)
	if err != nil {
		tx.Rollback()
		a.ErrorOK(MsgServerErr)
	}
	//转出aOut
	remain := param.SOutNum

	for _, o := range aOut {
		msg = fmt.Sprintf("%s，订单编号:%s", "额度转换", o.OrderNumber)
		if o.Amount < remain {
			//转换完break
			err = tx.Model(models.Amount{}).Where("id = ?", o.Id).Update("amount", 0).Error
			createAmountLogSimple(tx, o, msg, models.Amount_ConvOut, param.Remark, refer, o.Amount)
			if err != nil {
				tx.Rollback()
				a.ErrorOK(MsgServerErr)
			}
			remain -= o.Amount
			continue
		}
		err = tx.Model(models.Amount{}).Where("id = ?", o.Id).Update("amount", o.Amount-remain).Error
		createAmountLogSimple(tx, o, msg, models.Amount_ConvOut, param.Remark, refer, remain)
		if err != nil {
			tx.Rollback()
			a.ErrorOK(MsgServerErr)
		}
		break
	}
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		a.ErrorOK(MsgServerErr)
	}
	a.Correct("")
}

func createAmountLog(db *gorm.DB, param *models.Amount, msg, t string, amount int) error {
	amountLog := new(models.AmountLog)
	amountLog.AmountId = int(param.ID)
	amountLog.Change = amount * models.AmountChange[t]
	amountLog.Desc = fmt.Sprintf("%s，订单编号:%s", msg, param.OrderNumber)
	amountLog.RealTime = models.Time(time.Now())
	amountLog.Type = t
	amountLog.Remark = param.Remark
	return db.Create(amountLog).Error
}

func createAmountLogSimple(db *gorm.DB, param *models.AmountSimple, msg, t, r, refer string, amount int) error {
	amountLog := new(models.AmountLog)
	amountLog.AmountId = int(param.Id)
	amountLog.Change = amount * models.AmountChange[t]
	amountLog.Desc = msg
	amountLog.RealTime = models.Time(time.Now())
	amountLog.Type = t
	amountLog.Remark = r
	amountLog.Refer = refer
	return db.Create(amountLog).Error
}

func AmountDelayOut() {
	log.GLogger.Info("start amount delay...")
	//查询到期日期小于今天的amount
	amounts := make([]*models.Amount, 0)
	err := services.Slave().Where("amount > 0 and deadline < ?", time.Now().Format(models.DateFormat)).
		Find(&amounts).Error
	if err != nil {
		log.GLogger.Error("amount delay task err:%s", err.Error())
		return
	}
	tx := services.Slave().Begin()
	for _, a := range amounts {
		//创建log
		err = createAmountLog(tx, a, "额度过期", models.Amount_Delay_Out, a.Amount)
		if err != nil {
			tx.Rollback()
			break
		}
		//更新额度为0
		err = tx.Model(&a).Update("amount", 0).Error
		if err != nil {
			tx.Rollback()
			break
		}
	}
	tx.Commit()
	log.GLogger.Info("amount delay finish")
}
