/*
* Auth : acer
* Desc : 任务类
* Time : 2020/7/7 21:48
 */

package controllers

import (
	"bfimpl/models"
	"bfimpl/models/forms"
	"bfimpl/services"
	"bfimpl/services/util"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type TaskController struct {
	BaseController
}

// @Title 提测
// @Description 提测
// @Param	clientId	body	int	true	"客户id"
// @Param	appName		body	string	true	"应用名称"
// @Param	serviceId	body	int		true	"服务id"
// @Param	preAmount	body	int	true	"预计额度"
// @Param	preDate		body	string		true	"预计测试日期"
// @Param	expEndDate	body	int	string	"期望结单日期"
// @Param	manageId	body	int	true	"客户服务经理id"
// @Success 200 {object} models.Task
// @Failure 500 server err
// @router / [post]
func (t *TaskController) NewTask() {
	param := new(forms.ReqTask)
	err := json.Unmarshal(t.Ctx.Input.RequestBody, param)
	if err != nil {
		t.ErrorOK(MsgInvalidParam)
	}
	//额度大于0， 测试不早于当天，结单不早于测试当天
	today, _ := time.Parse(models.DateFormat, time.Now().Format(models.DateFormat))
	if param.PreAmount <= 0 ||
		time.Time(param.PreDate).Before(today) ||
		time.Time(param.ExpEndDate).Before(time.Time(param.PreDate)) {
		t.ErrorOK(MsgInvalidParam)
	}
	//查询状态启用,可实施,amount>0的 amounts
	aOut := make([]*models.AmountSimple, 0)
	services.Slave().Raw("select a.id,a.amount,a.order_number from amounts a,services s where a.service_id = s.id "+
		"and a.client_id = ? and s.id = ? and a.deadline > ? and s.state=0 and s.use !=2 and a.amount >0 "+
		"order by deadline", param.ClientId, param.ServiceId, time.Now().Format(models.DateFormat)).Scan(&aOut)
	outSum := 0
	for _, t := range aOut {
		outSum += t.Amount
	}
	if outSum < param.PreAmount {
		t.ErrorOK("额度不足")
	}

	task := new(models.Task)
	task.ClientId = param.ClientId
	task.AppName = param.AppName
	task.ServiceId = param.ServiceId
	task.PreAmount = param.PreAmount
	task.ManageId = param.ManageId
	task.PreDate = param.PreDate
	task.ExpEndDate = param.ExpEndDate
	task.RealTime = models.Time(time.Now())
	//任务编号
	task.Serial = time.Now().Format("20060102") + "_" + util.StringMd5(strconv.FormatInt(time.Now().Unix(), 10))

	tx := services.Slave().Begin()
	err = tx.Create(task).Error
	if err != nil {
		tx.Rollback()
		t.ErrorOK(MsgServerErr)
	}
	// 扣除额度
	// 转出aOut
	remain := param.PreAmount
	var msg string
	for _, o := range aOut {
		msg = fmt.Sprintf("%s，任务编号:%s", "实施消耗", o.OrderNumber)
		if o.Amount < remain {
			//转换完break
			err = tx.Model(models.Amount{}).Where("id = ?", o.Id).Update("amount", 0).Error
			createAmountLogSimple(tx, o, msg, models.Amount_Use, "", task.Serial)
			if err != nil {
				tx.Rollback()
				t.ErrorOK(MsgServerErr)
			}
			remain -= o.Amount
			continue
		}
		err = tx.Model(models.Amount{}).Where("id = ?", o.Id).Update("amount", o.Amount-remain).Error
		createAmountLogSimple(tx, o, msg, models.Amount_Use, "", task.Serial)
		if err != nil {
			tx.Rollback()
			t.ErrorOK(MsgServerErr)
		}
		break
	}
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		t.ErrorOK(MsgServerErr)
	}
	t.Correct(task)
}
