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
	AccountReUse    int         `json:"reUse"`
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
	}
}

//指派任务参数
type ReqAssignTask struct {
	ExeUserId int `json:"exeUserId"`
	Amount    int `json:"amount"`
}
