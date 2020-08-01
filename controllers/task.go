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
	"bfimpl/services/log"
	"bfimpl/services/util"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type TaskController struct {
	BaseController
}

// @Title 任务看板
// @Description 任务看板
// @Success 200 {object} []forms.RspTaskSum
// @Failure 500 server err
// @router /dashboard [get]
func (t *TaskController) TaskDashboard() {
	data := make([]forms.QueryTaskSum, 0)
	services.Slave().Raw("select s.id,s.service_name, t.status from services s LEFT JOIN tasks t on s.id = t.service_id").Scan(&data)
	var sId []int
	services.Slave().Model(models.Service{}).Order("sort").Pluck("id", &sId)
	result := make([]forms.RspTaskSum, len(sId))
	m := make(map[int]int)
	for i, v := range sId {
		m[v] = i
	}
	for _, d := range data {
		index := m[d.Id]
		result[index].Name = d.ServiceName
		if d.Status == "" {
			continue
		} else if d.Status == models.TaskCreate || d.Status == models.TaskConfirm {
			result[index].ReqNum += 1
		} else if d.Status != models.TaskEnd && d.Status != models.TaskCancel {
			result[index].ImpNum += 1
		}
	}
	t.Correct(result)
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

	uID, _ := t.GetInt("userID", 0)
	if uID == 0 {
		t.ErrorOK("need user id")
	}
	userType, _ := t.GetInt("userType", 0)
	// 管理员，销售，经理
	if userType != 1 && userType != 2 && userType != 3 {
		t.ErrorOK("无提测权限")
	}

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

	client := new(models.Client)
	services.Slave().Take(client, "id = ?", param.ClientId)
	task := new(models.Task)
	task.ClientId = param.ClientId
	task.AppName = param.AppName
	task.ServiceId = param.ServiceId
	task.PreAmount = param.PreAmount
	task.ManageId = param.ManageId
	task.PreDate = param.PreDate
	task.ExpEndDate = param.ExpEndDate
	if client.Level == "A" {
		task.ClientLevel = 3
	} else if client.Level == "B" {
		task.ClientLevel = 6
	}
	task.Status = models.TaskCreate
	task.RealTime = models.Time(time.Now())
	if userType == 3 {
		//经理提测直接显示对接中
		task.Status = models.TaskConfirm
		task.TMAcceptTime = task.RealTime
	}
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
	msg := fmt.Sprintf("%s，任务编号:%s", "实施消耗", task.Serial)
	for _, o := range aOut {
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
	//task log
	taskLog := models.TaskLog{
		TaskID:     task.ID,
		CreateTime: models.Time(time.Now()),
		Title:      t.GetString("userName") + "创建任务",
	}
	err = tx.Create(&taskLog).Error
	if err != nil {
		tx.Rollback()
		t.ErrorOK(MsgServerErr)
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
		t.ErrorOK("need id")
	}
	var task models.Task
	err := services.Slave().Preload("Client").Preload("Client.Sale").Preload("Manage").
		Preload("Service").Preload("RealService").Preload("ExeUser").
		Preload("TaskDetail").Preload("Logs").Take(&task, "id = ?", id).Error
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
	uID, _ := t.GetInt("userID", 0)
	if uID == 0 {
		t.ErrorOK("need user id")
	}
	pageSize, _ := t.GetInt("pageSize", 10)
	pageNum, _ := t.GetInt("pageNum", 1)

	tasks := make([]models.Task, 0)
	total := 0
	query := services.Slave().Model(models.Task{}).Where("status = ?", status)
	if status == models.TaskExecute {
		query = query.Or("status = ?", models.TaskPause)
	}
	//查询用户
	userType, _ := t.GetInt("userType", 0)
	switch userType {
	case 1:
	case 2:
		clientIds := make([]int, 0)
		services.Slave().Model(models.Client{}).Where("sale_id = ?", uID).Pluck("id", &clientIds)
		query = query.Where("client_id in (?)", clientIds)
	case 3:
		query = query.Where("manage_id = ?", uID)
	case 4:
		if status != models.TaskCreate && status != models.TaskConfirm && status != models.TaskFrozen {
			exeIds := make([]int, 0)
			services.Slave().Model(models.User{}).Where("leader_id = ?", uID).Pluck("id", &exeIds)
			query = query.Where("exe_user_id in (?)", exeIds)
		}
	case 5:
		query = query.Where("exe_user_id = ?", uID)
	default:
		t.ErrorOK("invalid user type")
	}
	query = query.Preload("Client").Preload("Service").Preload("RealService").Preload("Manage").
		Preload("ExeUser").Limit(pageSize).Offset((pageNum - 1) * pageSize)

	if status == models.TaskConfirm || status == models.TaskCreate {
		query = query.Order("exp_end_date, client_level")
	} else if status == models.TaskCancel {
		query = query.Order("cancel_time desc, client_level")
	} else if status == models.TaskEnd {
		query = query.Order("end_time desc, client_level")
	} else {
		query = query.Order("exp_end_time, client_level")
	}
	err := query.Find(&tasks).Limit(-1).Offset(-1).Count(&total).Error
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

/*
1. 期望结单时间<=当天
2. 期望结单时间-当时<=2但未进入执行中
3. 期望结单日期-当天<=3但未进入需求对接中（仅对管理员、销售、客户服务经理）
4. 任务处于执行中：暂停、待重启执行
5. 按期望结单时间升序排列(仅有结单日期的时间视做23：59：59)，客户级别降序排列
*/

// @Title 亟需关注
// @Description 亟需关注
// @Success 200 {object} []models.Task
// @Failure 500 server err
// @router /high [get]
func (t *TaskController) TaskImportant() {
	uID, _ := t.GetInt("userID", 0)
	if uID == 0 {
		t.ErrorOK("need user id")
	}
	userType, _ := t.GetInt("userType", 0)
	// 查询任务
	tasks := make([]models.Task, 0)
	today := time.Now().AddDate(0, 0, 1).Format(models.DateFormat)
	nextTwo := time.Now().AddDate(0, 0, 2).Format(models.TimeFormat)
	nextThree := time.Now().AddDate(0, 0, 3).Format(models.DateFormat)
	query := services.Slave().Model(models.Task{}).Where("status = ? ", models.TaskPause).
		Or("exp_end_time < ?", today).
		Or("status != ? and exp_end_time <= ?", models.TaskExecute, nextTwo)
	switch userType {
	case 1:
		//管理员
		query = query.Or("status = ? and exp_end_date <= ?", models.TaskCreate, nextThree)
	case 2:
		//销售,自己客户信息
		clientIds := make([]int, 0)
		services.Slave().Model(models.Client{}).Where("sale_id = ?", uID).Pluck("id", &clientIds)
		query = query.Or("status = ? and exp_end_date <= ? ", models.TaskCreate, nextThree).
			Where("client_id in (?)", clientIds)
	case 3:
		//客户服务经理
		query = query.Or("status = ? and exp_end_date <= ? ", models.TaskCreate, nextThree).
			Where("manage_id = ?", uID)
	case 4:
		//组长
		exeIds := make([]int, 0)
		services.Slave().Model(models.User{}).Where("leader_id = ?", uID).Pluck("id", &exeIds)
		query = query.Where("exe_user_id in (?)", exeIds)
	case 5:
		//实施
		query = query.Where("exe_user_id = ?", uID)
	default:
		t.ErrorOK("invalid user type")
	}
	query = query.Preload("Client").Preload("Service").Preload("RealService").Order("exp_end_time, client_level")
	query.Find(&tasks)
	//去掉end,cancel
	result := make([]models.Task, 0)
	for _, t := range tasks {
		if t.Status != models.TaskEnd && t.Status != models.TaskCancel {
			result = append(result, t)
		}
	}
	t.Correct(result)
}

// @Title 今日结单
// @Description 今日结单
// @Param	type	query	int 	true	"默认今天，1明天"
// @Success 200 {object} []models.Task
// @Failure 500 server err
// @router /focus [get]
func (t *TaskController) TaskToday() {
	qType, _ := t.GetInt("type", 0)
	if qType != 0 {
		qType = 1
	}
	uID, _ := t.GetInt("userID", 0)
	if uID == 0 {
		t.ErrorOK("need user id")
	}
	userType, _ := t.GetInt("userType", 0)
	// 查询任务
	tasks := make([]models.Task, 0)
	query := services.Slave().Model(models.Task{}).Where("status != ? and status != ?", models.TaskCancel, models.TaskEnd)
	today := time.Now().AddDate(0, 0, qType).Format(models.DateFormat)
	morrow := time.Now().AddDate(0, 0, qType+1).Format(models.DateFormat)
	switch userType {
	case 1:
		//管理员
		query = query.Where("exp_end_time >= ? and exp_end_time < ?", today, morrow)
	case 2:
		//销售,自己客户信息
		clientIds := make([]int, 0)
		services.Slave().Model(models.Client{}).Where("sale_id = ?", uID).Pluck("id", &clientIds)
		query = query.Where("exp_end_time >= ? and exp_end_time < ? and client_id in (?)",
			today, morrow, clientIds)
	case 3:
		//客户服务经理
		query = query.Where("exp_end_time >= ? and exp_end_time < ? and manage_id = ?",
			today, morrow, uID)
	case 4:
		//组长
		exeIds := make([]int, 0)
		services.Slave().Model(models.User{}).Where("leader_id = ?", uID).Pluck("id", &exeIds)
		query = query.Where("exp_end_time >= ? and exp_end_time < ? and exe_user_id in (?)",
			today, morrow, exeIds)
	case 5:
		//实施
		query = query.Where("exp_deliver_time >= ? and exp_deliver_time < ? and exe_user_id = ?",
			today, morrow, uID)
	default:
		t.ErrorOK("invalid user type")
	}
	query.Preload("Client").Preload("Service").Preload("RealService").Find(&tasks)
	t.Correct(tasks)
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

	//task log
	taskLog := models.TaskLog{
		TaskID:     uint(id),
		CreateTime: models.Time(time.Now()),
		Title:      t.GetString("userName") + "接受任务",
	}
	tx := services.Slave().Begin()
	//更新任务状态和 确认时间
	err := tx.Model(models.Task{}).Where("id = ?", id).Updates(map[string]interface{}{
		"tm_accept_time": models.Time(time.Now()),
		"status":         models.TaskConfirm,
	}).Error
	if err != nil {
		tx.Rollback()
		t.ErrorOK(MsgServerErr)
	}
	err = tx.Create(&taskLog).Error
	if err != nil {
		tx.Rollback()
		t.ErrorOK(MsgServerErr)
	}
	tx.Commit()
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
	uID, _ := t.GetInt("userID", 0)
	if uID == 0 {
		t.ErrorOK("need user id")
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
		"cancel_user_id": uID,
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
	err = createAmountLogSimpleT(tx, int(amount.ID), "任务取消", models.Amount_Cancel, "", task.Serial, task.PreAmount)
	if err != nil {
		tx.Rollback()
		t.ErrorOK("add amount fail")
	}
	taskLog := models.TaskLog{
		TaskID:     uint(id),
		CreateTime: models.Time(time.Now()),
		Title:      t.GetString("userName") + "取消任务",
	}
	err = tx.Create(&taskLog).Error
	if err != nil {
		tx.Rollback()
		t.ErrorOK(MsgServerErr)
	}
	tx.Commit()
	t.Correct("")
}

// @Title 任务信息录入
// @Description 任务信息录入
// @Param	id		path	int		true	"任务id"
// @Param   json    body   forms.ReqTaskDetail true "任务详情"
// @Success 200 {"string"} success
// @Failure 500 server err
// @router /save/:id [post]
func (t *TaskController) SaveTaskDetail() {
	id, _ := t.GetInt(":id", 0)
	param := new(forms.ReqTaskDetail)
	err := json.Unmarshal(t.Ctx.Input.RequestBody, param)
	if err != nil {
		log.GLogger.Error(err.Error())
		t.ErrorOK(MsgInvalidParam)
	}
	var task models.Task
	err = services.Slave().Take(&task, "id = ?", id).Error
	if err != nil {
		t.ErrorOK("invalid taskId")
	}
	//更新task
	services.Slave().Model(&task).Updates(map[string]interface{}{
		"exp_deliver_time": param.ExpDeliverTime,
		"exp_end_time":     param.ExpEndTime,
		"real_amount":      param.RealAmount,
		"real_service_id":  param.RealServiceId,
	})
	logContent := "信息录入"
	if param.ChangeLog != "" {
		logContent = "变更需求"
	}
	// 更新detail
	tmp := new(models.TaskDetail)
	services.Slave().Where("task_id = ?", id).First(tmp)
	taskDetail := param.GetTaskDetail()
	taskDetail.TaskID = id
	if tmp.ID == 0 {
		services.Slave().Create(taskDetail)
	} else {
		taskDetail.ID = tmp.ID
		services.Slave().Save(taskDetail)
	}
	//创建日志
	taskLog := models.TaskLog{
		TaskID:     uint(id),
		CreateTime: models.Time(time.Now()),
		Title:      t.GetString("userName") + logContent,
		Desc:       param.ChangeLog,
	}
	services.Slave().Create(&taskLog)
	//创建历史版本
	history := param.GetTaskHistory()
	history.TaskLogID = int(taskLog.ID)
	services.Slave().Create(&history)
	t.Correct("")
}

// @Title 任务信息历史
// @Description 任务信息历史
// @Param	id		path	int		true	"任务id"
// @Success 200 {object} []models.TaskLog
// @Failure 500 server err
// @router /history/:id [get]
func (t *TaskController) TaskHistory() {

	id, _ := t.GetInt(":id", 0)
	taskLogs := make([]models.TaskLog, 0)
	services.Slave().Model(models.TaskLog{}).Where("task_id = ?", id).Preload("TaskHistory").
		Joins("inner join task_histories on task_logs.id = task_histories.task_log_id").
		Order("task_logs.create_time desc").Find(&taskLogs)

	t.Correct(taskLogs)
}

// @Title 任务冻结
// @Description 任务冻结
// @Param	id		path	int		true	"任务id"
// @Success 200 {"string"} success
// @Failure 500 server err
// @router /frozen/:id [put]
func (t *TaskController) FrozenTask() {
	id, _ := t.GetInt(":id", 0)
	if id == 0 {
		t.ErrorOK(MsgInvalidParam)
	}

	var task models.Task
	err := services.Slave().Take(&task, "id = ?", id).Error
	if task.ID == 0 || err != nil {
		t.ErrorOK("invalid taskID")
	}
	//冻结时服务类型变更
	if task.ServiceId != task.RealServiceId {
		//查询状态启用,可实施,amount>0的 amounts
		aOut := make([]*models.AmountSimple, 0)
		services.Slave().Raw("select a.id,a.amount,a.order_number from amounts a,services s where a.service_id = s.id "+
			"and a.client_id = ? and s.id = ? and a.deadline > ? and s.state=0 and s.use !=2 and a.amount >0 "+
			"order by deadline", task.ClientId, task.RealServiceId, time.Now().Format(models.DateFormat)).Scan(&aOut)
		outSum := 0
		for _, a := range aOut {
			outSum += a.Amount
		}
		if outSum < task.RealAmount {
			t.ErrorOK("额度不足")
		}
		tx := services.Slave().Begin()
		//原额度反冲，类似取消
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
		err = createAmountLogSimpleT(tx, int(amount.ID), "冻结变更", models.Amount_Frozen_In, "", task.Serial, task.PreAmount)
		if err != nil {
			tx.Rollback()
			t.ErrorOK("add amount fail")
		}
		//消耗新额度
		// 转出aOut
		remain := task.RealAmount
		msg := fmt.Sprintf("%s，任务编号:%s", "实施消耗", task.Serial)
		for _, o := range aOut {
			if o.Amount < remain {
				//转换完break
				err = tx.Model(models.Amount{}).Where("id = ?", o.Id).Update("amount", 0).Error
				err = createAmountLogSimple(tx, o, msg, models.Amount_Use, "", task.Serial, o.Amount)
				if err != nil {
					tx.Rollback()
					t.ErrorOK(MsgServerErr)
				}
				remain -= o.Amount
				continue
			}
			err = tx.Model(models.Amount{}).Where("id = ?", o.Id).Update("amount", o.Amount-remain).Error
			err = createAmountLogSimple(tx, o, msg, models.Amount_Use, "", task.Serial, remain)
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
	} else {
		if task.RealAmount < task.PreAmount {
			//反冲
			tx := services.Slave().Begin()
			//原额度反冲，类似取消
			var amount models.Amount
			err = tx.Model(models.Amount{}).Where("client_id = ? and service_id = ?", task.ClientId, task.ServiceId).
				Order("deadline desc").First(&amount).Error
			if err != nil {
				tx.Rollback()
				t.ErrorOK("add amount fail")
			}
			err = tx.Model(&amount).UpdateColumn("amount", amount.Amount+task.PreAmount-task.RealAmount).Error
			if err != nil {
				tx.Rollback()
				t.ErrorOK("add amount fail")
			}
			// 反冲记录
			err = createAmountLogSimpleT(tx, int(amount.ID), "冻结变更", models.Amount_Frozen_In, "", task.Serial, task.PreAmount-task.RealAmount)
			if err != nil {
				tx.Rollback()
				t.ErrorOK("add amount fail")
			}
		} else if task.RealAmount > task.PreAmount {
			//消耗
			aOut := make([]*models.AmountSimple, 0)
			services.Slave().Raw("select a.id,a.amount,a.order_number from amounts a,services s where a.service_id = s.id "+
				"and a.client_id = ? and s.id = ? and a.deadline > ? and s.state=0 and s.use !=2 and a.amount >0 "+
				"order by deadline", task.ClientId, task.ServiceId, time.Now().Format(models.DateFormat)).Scan(&aOut)
			outSum := 0
			for _, a := range aOut {
				outSum += a.Amount
			}
			remain := task.RealAmount - task.PreAmount
			if outSum < remain {
				t.ErrorOK("额度不足")
			}
			tx := services.Slave().Begin()
			msg := fmt.Sprintf("%s，任务编号:%s", "实施消耗", task.Serial)
			for _, o := range aOut {
				if o.Amount < remain {
					//转换完break
					err = tx.Model(models.Amount{}).Where("id = ?", o.Id).Update("amount", 0).Error
					err = createAmountLogSimple(tx, o, msg, models.Amount_Use, "", task.Serial, o.Amount)
					if err != nil {
						tx.Rollback()
						t.ErrorOK(MsgServerErr)
					}
					remain -= o.Amount
					continue
				}
				err = tx.Model(models.Amount{}).Where("id = ?", o.Id).Update("amount", o.Amount-remain).Error
				err = createAmountLogSimple(tx, o, msg, models.Amount_Use, "", task.Serial, remain)
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
		}
	}

	//更新任务状态和 冻结时间
	err = services.Slave().Model(&task).Updates(map[string]interface{}{
		"frozen_time": models.Time(time.Now()),
		"status":      models.TaskFrozen,
	}).Error
	if err != nil {
		t.ErrorOK(MsgServerErr)
	}
	taskLog := models.TaskLog{
		TaskID:     uint(id),
		CreateTime: models.Time(time.Now()),
		Title:      t.GetString("userName") + "冻结需求",
	}
	services.Slave().Create(&taskLog)
	t.Correct("")
}

// @Title 任务指派
// @Description 任务指派
// @Param	id		path	int		true	"任务id"
// @Param	json	body	forms.ReqAssignTask		true	"body参数"
// @Success 200 {"string"} success
// @Failure 500 server err
// @router /assign/:id [put]
func (t *TaskController) AssignTask() {
	id, _ := t.GetInt(":id", 0)
	param := new(forms.ReqAssignTask)
	err := json.Unmarshal(t.Ctx.Input.RequestBody, param)
	if err != nil {
		t.ErrorOK(MsgInvalidParam)
	}
	var task models.Task
	err = services.Slave().Take(&task, "id = ?", id).Error
	if err != nil {
		t.ErrorOK("invalid taskId")
	}
	//更新任务状态，取消时间，原因，取消人id
	err = services.Slave().Model(&task).Updates(map[string]interface{}{
		"assign_time":    models.Time(time.Now()),
		"status":         models.TaskAssign,
		"exe_user_id":    param.ExeUserId,
		"deliver_amount": param.Amount,
	}).Error
	if err != nil {
		t.ErrorOK(MsgServerErr)
	}
	taskLog := models.TaskLog{
		TaskID:     uint(id),
		CreateTime: models.Time(time.Now()),
		Title:      t.GetString("userName") + "资源指派",
	}
	services.Slave().Create(&taskLog)
	t.Correct("")
}

// @Title 任务执行
// @Description 任务启动执行
// @Param	id	path	int	true	"任务id"
// @Success 200 {"string"} success
// @Failure 500 server err
// @router /execute/:id [put]
func (t *TaskController) ExecuteTask() {
	id, _ := t.GetInt(":id", 0)
	if id == 0 {
		t.ErrorOK(MsgInvalidParam)
	}

	//更新任务状态和 执行时间
	err := services.Slave().Model(models.Task{}).Where("id = ?", id).Updates(map[string]interface{}{
		"execute_time": models.Time(time.Now()),
		"status":       models.TaskExecute,
	}).Error
	if err != nil {
		t.ErrorOK(MsgServerErr)
	}
	taskLog := models.TaskLog{
		TaskID:     uint(id),
		CreateTime: models.Time(time.Now()),
		Title:      t.GetString("userName") + "启动执行",
	}
	services.Slave().Create(&taskLog)
	t.Correct("")
}

// @Title 任务暂停
// @Description 任务暂停
// @Param	id	path	int	true	"任务id"
// @Success 200 {"string"} success
// @Failure 500 server err
// @router /pause/:id [put]
func (t *TaskController) PauseTask() {
	id, _ := t.GetInt(":id", 0)
	if id == 0 {
		t.ErrorOK(MsgInvalidParam)
	}

	//更新任务状态和 暂停时间
	err := services.Slave().Model(models.Task{}).Where("id = ?", id).Updates(map[string]interface{}{
		"pause_time": models.Time(time.Now()),
		"status":     models.TaskPause,
		"is_pause":   1,
	}).Error
	if err != nil {
		t.ErrorOK(MsgServerErr)
	}
	taskLog := models.TaskLog{
		TaskID:     uint(id),
		CreateTime: models.Time(time.Now()),
		Title:      t.GetString("userName") + "暂停任务",
	}
	services.Slave().Create(&taskLog)
	t.Correct("")
}

// @Title 任务变更完成
// @Description 任务变更完成
// @Param	id	path	int	true	"任务id"
// @Success 200 {"string"} success
// @Failure 500 server err
// @router /change/:id [put]
func (t *TaskController) ChangeFinish() {
	id, _ := t.GetInt(":id", 0)
	if id == 0 {
		t.ErrorOK(MsgInvalidParam)
	}
	//更新任务状态和 暂停时间
	err := services.Slave().Model(models.Task{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_pause": 0,
	}).Error
	if err != nil {
		t.ErrorOK(MsgServerErr)
	}
	t.Correct("")
}

// @Title 任务执行标签
// @Description 标签列表
// @Success 200 {object} []models.Tag
// @Failure 500 server err
// @router /tags [get]
func (t *TaskController) TaskTags() {
	data := make([]models.Tag, 0)
	services.Slave().Model(models.Tag{}).Find(&data)
	t.Correct(data)
}

// @Title 任务完成
// @Description 任务完成
// @Param	id	path	int	true	"任务id"
// @Param	json	body	forms.ReqFinishTask	true	"body请求参数"
// @Success 200 {"string"} success
// @Failure 500 server err
// @router /finish/:id [put]
func (t *TaskController) FinishTask() {
	id, _ := t.GetInt(":id", 0)
	if id == 0 {
		t.ErrorOK(MsgInvalidParam)
	}
	param := new(forms.ReqFinishTask)
	err := json.Unmarshal(t.Ctx.Input.RequestBody, param)
	if err != nil {
		log.GLogger.Error(err.Error())
		t.ErrorOK(MsgInvalidParam)
	}

	tags := make([]*models.Tag, 0)
	err = services.Slave().Where(param.Tags).Find(&tags).Error
	if err != nil {
		t.ErrorOK("need tag id")
	}
	tx := services.Slave().Begin()
	//创建执行任务信息
	taskExeInfo := models.TaskExeInfo{
		TaskID:       id,
		UsedTime:     param.UsedTime,
		ExecuteBatch: param.ExecuteBatch,
		ExecuteTai:   param.ExecuteTai,
		DelayTime:    param.DelayTime,
		Desc:         param.Desc,
		Tags:         tags,
	}
	err = tx.Create(&taskExeInfo).Error
	if err != nil {
		tx.Rollback()
		t.Error(MsgServerErr)
	}
	//更新任务状态和 完成时间
	err = tx.Model(models.Task{}).Where("id = ?", id).Updates(map[string]interface{}{
		"finish_time": models.Time(time.Now()),
		"status":      models.TaskFinish,
	}).Error
	if err != nil {
		tx.Rollback()
		t.ErrorOK(MsgServerErr)
	}
	taskLog := models.TaskLog{
		TaskID:     uint(id),
		CreateTime: models.Time(time.Now()),
		Title:      t.GetString("userName") + "执行完成",
	}
	err = tx.Create(&taskLog).Error
	if err != nil {
		tx.Rollback()
		t.ErrorOK(MsgServerErr)
	}
	tx.Commit()
	t.Correct("")
}

// @Title 任务执行信息
// @Description 任务执行信息
// @Param	id	path	int	true	"任务id"
// @Success 200 {object} models.TaskExeInfo
// @Failure 500 server err
// @router /exeinfo/:id [get]
func (t *TaskController) GetTaskExeInfo() {
	id, _ := t.GetInt(":id", 0)
	if id == 0 {
		t.ErrorOK("invalid id")
	}
	var info models.TaskExeInfo
	services.Slave().Where("task_id = ?", id).Preload("Tags").First(&info)
	t.Correct(info)
}

// @Title 任务结单
// @Description 任务结单
// @Param	id	path	int	true	"任务id"
// @Success 200 {"string"} success
// @Failure 500 server err
// @router /end/:id [put]
func (t *TaskController) EndTask() {
	id, _ := t.GetInt(":id", 0)
	if id == 0 {
		t.ErrorOK(MsgInvalidParam)
	}
	//更新任务状态和 完成时间
	err := services.Slave().Model(models.Task{}).Where("id = ?", id).Updates(map[string]interface{}{
		"end_time": models.Time(time.Now()),
		"status":   models.TaskEnd,
	}).Error
	if err != nil {
		log.GLogger.Error("EndTask:%s", err.Error())
		t.ErrorOK(MsgServerErr)
	}
	taskLog := models.TaskLog{
		TaskID:     uint(id),
		CreateTime: models.Time(time.Now()),
		Title:      t.GetString("userName") + "任务结单",
	}
	services.Slave().Create(&taskLog)
	t.Correct("")
}

// @Title 任务评价
// @Description 结单任务评价，销售评价客户经理
// @Param	id	path	int	true	"任务id"
// @Param	json	body	forms.ReqCommentTask	true	"body请求参数"
// @Success 200 {"string"} success
// @Failure 500 server err
// @router /comment/:id [put]
func (t *TaskController) CommentTask() {
	id, _ := t.GetInt(":id", 0)
	if id == 0 {
		t.ErrorOK(MsgInvalidParam)
	}
	var task models.Task
	services.Slave().Take(&task, "id = ?", id)
	//if task.Status != models.TaskEnd {
	//	t.ErrorOK("任务未结单")
	//}
	param := new(forms.ReqCommentTask)
	err := json.Unmarshal(t.Ctx.Input.RequestBody, param)
	if err != nil {
		log.GLogger.Error(err.Error())
		t.ErrorOK(MsgInvalidParam)
	}

	//创建执行任务信息
	taskComment := models.TaskComment{
		TaskID:         id,
		CommentType:    param.CommentType,
		RealTime:       param.RealTime,
		ReExecuteTimes: param.ReExecuteTimes,
		Score:          param.Score,
		Other:          param.Other,
	}
	err = services.Slave().Create(&taskComment).Error
	if err != nil {
		t.Error(MsgServerErr)
	}
	t.Correct("")
}

// @Title 获取任务评价信息
// @Description commentType 0实施评价1客服经理评价
// @Param	id	path	int	true	"任务id"
// @Success 200 {object} []models.TaskComment
// @Failure 500 server err
// @router /comment/:id [get]
func (t *TaskController) TaskComments() {
	id, _ := t.GetInt(":id", 0)
	if id == 0 {
		t.ErrorOK(MsgInvalidParam)
	}
	comments := make([]models.TaskComment, 0)
	query := services.Slave().Where("task_id = ?", id)
	userType, _ := t.GetInt("userType", 0)
	if userType == 3 {
		//经理
		query = query.Where("comment_type = ?", 0)
	} else if userType == 5 {
		//实施
		query = query.Where("comment_type = ?", 10)
	}
	query.Find(&comments)
	t.Correct(comments)
}

// @Title 任务退次
// @Description 任务退次,退次额度小于实际额度 realAmount
// @Param	id		path	int	true	"任务serial"
// @Param	json	body	forms.ReqBackAmount	true	"退次参数"
// @Success 200 {"string"} success
// @Failure 500 server err
// @router /backamount/:id [put]
func (t *TaskController) TaskBackAmount() {
	taskSerial := t.GetString(":id")
	param := new(forms.ReqBackAmount)
	err := json.Unmarshal(t.Ctx.Input.RequestBody, param)
	if err != nil {
		log.GLogger.Error(err.Error())
		t.ErrorOK(MsgInvalidParam)
	}
	if param.Amount <= 0 {
		t.ErrorOK(MsgInvalidParam)
	}
	var task models.Task
	err = services.Slave().Where("serial = ?", taskSerial).First(&task).Error
	if err != nil {
		t.ErrorOK("invalid taskId")
	}
	//检查真实提测额度
	if param.Amount > task.RealAmount {
		t.ErrorOK("退次额度超过提测额度")
	}
	tx := services.Slave().Begin()
	err = tx.Model(&task).UpdateColumn("real_amount", task.RealAmount-param.Amount).Error
	if err != nil {
		tx.Rollback()
		t.ErrorOK(MsgServerErr)
	}
	//反冲额度
	var amount models.Amount
	err = tx.Model(models.Amount{}).Where("client_id = ? and service_id = ?", task.ClientId, task.RealServiceId).
		Order("deadline desc").First(&amount).Error
	if err != nil {
		tx.Rollback()
		t.ErrorOK("add amount fail")
	}
	err = tx.Model(&amount).UpdateColumn("amount", amount.Amount+param.Amount).Error
	if err != nil {
		tx.Rollback()
		t.ErrorOK("add amount fail")
	}
	// 反冲记录
	err = createAmountLogSimpleT(tx, int(amount.ID), "任务退次", models.Amount_Back, param.Remark, task.Serial, param.Amount)
	if err != nil {
		tx.Rollback()
		t.ErrorOK("add amount fail")
	}
	tx.Commit()
	t.Correct("")
}
