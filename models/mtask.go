/*
* Auth : acer
* Desc : 任务结构体
* Time : 2020/6/29 15:00
 */

package models

import "github.com/jinzhu/gorm"

// 若是客户服务经理提测，则状态标记为对接中，提测时间点和需求对接起始点保持一致
// 需求变更 任务状态变为deliver，涉及额度和任务类型，变为 cancel
type Task struct {
	gorm.Model

	ClientName  string `gorm:"size:30;not null;comment:'客户名称'" json:"clientName"`
	AppName     string `gorm:"size:50;not null;comment:'游戏名称'" json:"appName"`
	ServiceName string `gorm:"size:60;not null;comment:'服务名称'" json:"serviceName"`
	PreAmount   int    `gorm:"not null;comment:'预计额度'" json:"preAmount"`
	PreDate     Time   `gorm:"type:time;comment:'预计提测日期'" json:"preDate"`
	ExpDate     Time   `gorm:"type:time;comment:'期望结单日期'" json:"expDate"`
	ManageId    int    `gorm:"comment:'客户服务经理id'" json:"manageId"`
	RealTime    Time   `gorm:"type:time;not null;comment:'创建时间'" json:"realTime"`
	FrozenTime  Time   `gorm:"type:time;not null;comment:'冻结时间'" json:"frozenTime"`
	PauseTime   Time   `gorm:"type:time;not null;comment:'暂停时间'" json:"pauseTime"`
	Status      int    `gorm:"not null;comment:'任务状态'" json:"status"`
	Serial      string `gorm:"not null;comment:'任务编号'" json:"serial"`
	Reason      string `gorm:"default:'';comment:'任务取消原因'" json:"reason"`
	RealAmount  string `gorm:"comment:'交付评估额度'" json:"realAmount"`
	ExeUserId   int    `gorm:"index;comment:'被指派人员id'" json:"exeUserId"`
	ChangeLog   string `gorm:"size:256;comment:'变更说明'" json:"changeLog"`
}

var TaskStatus = map[string]int{
	"create":   1,
	"cancel":   2,
	"confirm":  3,
	"frozen":   4,
	"deliver":  5,
	"execute":  6,
	"finish":   7,
	"archived": 8,
}

// 任务详细信息
type TaskDetail struct {
	gorm.Model
	TaskID          int    `gorm:"index;comment:'任务id'" json:"taskId"`
	Version         string `gorm:"size:30;comment:'测试版本'" json:"version"`
	PkgAddress      string `gorm:"size:256;comment:'安装包地址'" json:"pkgAddress"`
	TestType        int    `gorm:"comment:'环境类型'" json:"testType"`
	TestExtInfo     string `gorm:"size:256;comment:'测试环境补充信息'" json:"testExtInfo"`
	WhiteListType   int    `gorm:"type:tinyint;comment:'白名单类型'" json:"whiteListType"`
	WhiteList       string `gorm:"size:512;comment:'白名单'" json:"whiteList"`
	TestAccountType string `gorm:"size:40;comment:'测试账号类型'" json:"testAccountType"`
	AccountReUse    int    `gorm:"type:tinyint;default:0;comment:'账号是否重复使用'" json:"reUse"`
	AccountAddress  string `gorm:"size:256;comment:'账号文件地址'" json:"accountAddress"`

	AccountNum    int    `gorm:"comment:'账号数量'" json:"accountNum"`
	PhoneNum      int    `gorm:"comment:'手机号/微信数量'" json:"phoneNum"`
	ConcurrentNum int    `gorm:"comment:'并发数'" json:"concurrentNum"`
	ReqPhone      string `gorm:"size:256;comment:'机型需求'" json:"reqPhone"`
	ExtReq        string `gorm:"size:256;comment:'其他需求'" json:"extReq"`
	InstanceTxt   string `gorm:"size:256;comment:'文字用例内网地址'" json:"instanceTxt"`
	InstanceMv    string `gorm:"size:256;comment:'视频用例内网地址'" json:"instanceMv"`
}

//var WLType = map[int]string{
//	0: "无",
//	1: "IP地址",
//	2: "账号白名单",
//}

//任务执行信息
type TaskExeInfo struct {
	gorm.Model
	TaskID       int    `gorm:"index;comment:'任务id'" json:"taskId"`
	UsedTime     int    `gorm:"comment:'任务执行时长'" json:"usedTime"`
	ExecuteBatch int    `gorm:"comment:'任务执行批次'" json:"executeBatch"`
	ExecuteTai   int    `gorm:"comment:'任务执行台次'" json:"executeTai"`
	DelayTime    int    `gorm:"comment:'外部延误时常'" json:"delayTime"`
	Desc         string `gorm:"comment:'执行说明'" json:"desc"`
	Tag          int    `gorm:"comment:'执行用到的tag'" json:"tag"`
}

type TaskComment struct {
	gorm.Model
	TaskID         int    `gorm:"index;comment:'任务id'" json:"taskId"`
	CommentType    int    `gorm:"comment:'0实施评价1客户经理评价'" json:"commentType"`
	RealTime       Time   `gorm:"type:time;not null;comment:'交付时间'" json:"realTime"`
	ReExecuteTimes int    `gorm:"default:0;comment:'返工次数'" json:"reExeTimes"`
	Score          int    `gorm:"comment:'评分'" json:"score"`
	Other          string `gorm:"size:256;comment:'其他信息'" json:"other"`
}
