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

// @Title 任务提测
// @Description 任务提测
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
	task.Status = models.TaskCreate
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
			createAmountLogSimple(tx, o, msg, models.Amount_Use, "", task.Serial, o.Amount)
			if err != nil {
				tx.Rollback()
				t.ErrorOK(MsgServerErr)
			}
			remain -= o.Amount
			continue
		}
		err = tx.Model(models.Amount{}).Where("id = ?", o.Id).Update("amount", o.Amount-remain).Error
		createAmountLogSimple(tx, o, msg, models.Amount_Use, "", task.Serial, remain)
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

// @Title 单个任务
// @Description 单个任务
// @Param	id		path	int		true	"任务id"
// @Success 200 {object} models.Task
// @Failure 500 server err
// @router /:id [get]
func (t *TaskController) Task() {
	id, _ := t.GetInt(":id", 0)
	if id == 0 {
		t.GetInt("need id")
	}
	var task models.Task
	err := services.Slave().Preload("Client").Preload("TaskDetail").Take(&task, "id = ?", id).Error
	if err != nil {
		t.ErrorOK("invalid id")
	}
	t.Correct(task)
}

// @Title 任务列表
// @Description 任务列表
// @Param	status		query	string	true	"状态"
// @Param	pageSize	query	int		false	"条数"
// @Param	pageNum		query	int		false	"页数"
// @Success 200 {object} []models.Task
// @Failure 500 server err
// @router /list [get]
func (t *TaskController) TaskList() {
	status := t.GetString("status")
	if status == "" {
		t.ErrorOK("need status")
	}
	pageSize, _ := t.GetInt("pageSize", 10)
	pageNum, _ := t.GetInt("pageNum", 1)

	tasks := make([]models.Task, 0)
	total := 0
	err := services.Slave().Model(models.Task{}).Where("status = ?", status).Preload("Client").
		Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&tasks).Limit(-1).Offset(-1).Count(&total).Error
	if err != nil {
		t.ErrorOK(MsgServerErr)
	}
	var resp struct {
		Total int           `json:"total"`
		List  []models.Task `json:"list"`
	}
	resp.Total = total
	resp.List = tasks
	t.Correct(resp)
}

// @Title 任务确认
// @Description 任务确认
// @Param	id	path	int	true	"任务id"
// @Success 200 {"string"} success
// @Failure 500 server err
// @router /confirm/:id [put]
func (t *TaskController) ConfirmTask() {
	id, _ := t.GetInt(":id", 0)
	if id == 0 {
		t.ErrorOK(MsgInvalidParam)
	}

	//更新任务状态和 确认时间
	err := services.Slave().Model(models.Task{}).Where("id = ?", id).Updates(map[string]interface{}{
		"tm_accept_time": models.Time(time.Now()),
		"status":         models.TaskConfirm,
	}).Error
	if err != nil {
		t.ErrorOK(MsgServerErr)
	}
	t.Correct("")
}

// @Title 任务取消
// @Description 任务取消
// @Param	id		path	int		true	"任务id"
// @Param	userId	body	int		true	"操作人id"
// @Param	reason	body	string	true	"取消原因"
// @Success 200 {"string"} success
// @Failure 500 server err
// @router /cancel/:id [put]
func (t *TaskController) CancelTask() {
	id, _ := t.GetInt(":id", 0)
	if id == 0 {
		t.ErrorOK(MsgInvalidParam)
	}
	param := new(forms.ReqCancelTask)
	err := json.Unmarshal(t.Ctx.Input.RequestBody, param)
	if err != nil {
		t.ErrorOK(MsgInvalidParam)
	}

	tx := services.Slave().Begin()
	var task models.Task
	err = tx.Model(models.Task{}).Take(&task, "id = ?", id).Error
	if err != nil {
		t.ErrorOK("task not found")
	}
	//更新任务状态，取消时间，原因，取消人id
	err = tx.Model(&task).Updates(map[string]interface{}{
		"cancel_time":    models.Time(time.Now()),
		"status":         models.TaskCancel,
		"cancel_user_id": param.UserId,
		"reason":         param.Reason,
	}).Error
	if err != nil {
		tx.Rollback()
		t.ErrorOK(MsgServerErr)
	}
	// 额度反冲
	// 查询额度
	var amount models.Amount
	err = tx.Model(models.Amount{}).Where("client_id = ? and service_id = ?", task.ClientId, task.ServiceId).
		Order("deadline desc").First(&amount).Error
	if err != nil {
		tx.Rollback()
		t.ErrorOK("add amount fail")
	}
	err = tx.Model(&amount).UpdateColumn("amount", amount.Amount+task.PreAmount).Error
	if err != nil {
		tx.Rollback()
		t.ErrorOK("add amount fail")
	}
	// 反冲记录
	err = createAmountLog(tx, &amount, "任务取消", models.Amount_Cancel, task.PreAmount)
	if err != nil {
		tx.Rollback()
		t.ErrorOK("add amount fail")
	}
	tx.Commit()
	t.Correct("")
}
