/*
* Auth : acer
* Desc : task请求参数
* Time : 2020/7/7 21:51
 */

package forms

import (
	"bfimpl/models"
)

//创建任务
type ReqTask struct {
	ClientId   int         `json:"clientId"`
	AppName    string      `json:"appName"`
	ServiceId  int         `json:"serviceId"`
	PreAmount  int         `json:"preAmount"`
	PreDate    models.Time `json:"preDate"`
	ExpEndDate models.Time `json:"expEndDate"`
	ManageId   int         `json:"manageId"`
}

//取消任务参数
type ReqCancelTask struct {
	UserId int    `json:"userId"`
	Reason string `json:"reason"`
}

// 任务详细信息录入
type ReqTaskDetail struct {
	RealServiceId   int         `json:"realServiceId"`
	RealAmount      int         `json:"realAmount"`
	ExpDeliverTime  models.Time `json:"expDeliverTime"`
	ExpEndTime      models.Time `json:"expEndTime"`
	Version         string      `json:"version"`
	PkgAddress      string      `json:"pkgAddress"`
	TestType        string      `json:"testType"`
	TestExtInfo     string      `json:"testExtInfo"`
	WhiteList       string      `json:"whiteList"`
	TestAccountType string      `json:"testAccountType"`
	AccountReUse    string      `json:"reUse"`
	AccountAddress  string      `json:"accountAddress"`
	AccountNum      int         `json:"accountNum"`
	PhoneNum        int         `json:"phoneNum"`
	ConcurrentNum   int         `json:"concurrentNum"`
	ReqPhone        string      `json:"reqPhone"`
	ExtReq          string      `json:"extReq"`
	InstanceTxt     string      `json:"instanceTxt"`
	InstanceMv      string      `json:"instanceMv"`
	ChangeLog       string      `json:"changeLog"`
}

func (d *ReqTaskDetail) GetTaskDetail() *models.TaskDetail {
	return &models.TaskDetail{
		Version:         d.Version,
		PkgAddress:      d.PkgAddress,
		TestType:        d.TestType,
		TestExtInfo:     d.TestExtInfo,
		WhiteList:       d.WhiteList,
		TestAccountType: d.TestAccountType,
		AccountReUse:    d.AccountReUse,
		AccountAddress:  d.AccountAddress,
		AccountNum:      d.AccountNum,
		PhoneNum:        d.PhoneNum,
		ConcurrentNum:   d.ConcurrentNum,
		ReqPhone:        d.ReqPhone,
		ExtReq:          d.ExtReq,
		InstanceTxt:     d.InstanceTxt,
		InstanceMv:      d.InstanceMv,
		ChangeLog:       d.ChangeLog,
	}
}

func (d *ReqTaskDetail) GetTaskHistory() *models.TaskHistory {
	return &models.TaskHistory{
		ExpDeliverTime:  d.ExpDeliverTime,
		ExpEndTime:      d.ExpEndTime,
		Version:         d.Version,
		PkgAddress:      d.PkgAddress,
		TestType:        d.TestType,
		TestExtInfo:     d.TestExtInfo,
		WhiteList:       d.WhiteList,
		TestAccountType: d.TestAccountType,
		AccountReUse:    d.AccountReUse,
		AccountAddress:  d.AccountAddress,
		AccountNum:      d.AccountNum,
		PhoneNum:        d.PhoneNum,
		ConcurrentNum:   d.ConcurrentNum,
		ReqPhone:        d.ReqPhone,
		ExtReq:          d.ExtReq,
		InstanceTxt:     d.InstanceTxt,
		InstanceMv:      d.InstanceMv,
		ChangeLog:       d.ChangeLog,
	}
}

//指派任务参数
type ReqAssignTask struct {
	ExeUserId   int `json:"exeUserId"`
	AssignSrvId int `json:"serviceId"`
	Amount      int `json:"amount"`
}

// 完成任务参数
type ReqFinishTask struct {
	UsedTime     int    `json:"usedTime"`
	ExecuteBatch int    `json:"executeBatch"`
	ExecuteTai   int    `json:"executeTai"`
	DelayTime    int    `json:"delayTime"`
	Desc         string `json:"desc"`
	Tags         []int  `json:"tags"`
}

// 结单任务参数
type ReqCommentTask struct {
	RealTime       models.Time `json:"realTime"`
	ReExecuteTimes int         `json:"reExeTimes"`
	Score          int         `json:"score"`
	Other          string      `json:"other"`
	CommentType    int         `json:"commentType"`
}

// 退次任务参数
type ReqBackAmount struct {
	Amount int    `json:"amount"`
	Remark string `json:"remark"`
}

type QueryTaskSum struct {
	Id          int    `json:"id"`
	ServiceName string `json:"service_name"`
	Status      string `json:"status"`
}

type RspTaskSum struct {
	Name   string `json:"name"`
	ReqNum int    `json:"reqNum"`
	ImpNum int    `json:"impNum"`
}

//额度延期参数
type ReqAmountDelay struct {
	Deadline models.Time `json:"deadline"`
	Remark   string      `json:"remark"`
}
