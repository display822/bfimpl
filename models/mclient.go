/*
* Auth : acer
* Desc : 客户结构体
* Time : 2020/6/28 21:16
 */

package models

import "github.com/jinzhu/gorm"

// 客户
type Client struct {
	gorm.Model
	Name         string `gorm:"size:30;not null;comment:'名称'" json:"name"`
	Number       string `gorm:"unique_index;size:50;not null;comment:'编号'" json:"number"`
	Type         int    `gorm:"type:tinyint;default:0;comment:'0内部1外部'" json:"type"`
	Level        string `gorm:"size:5;not null;comment:'级别S,A,B'" json:"level"`
	SaleId       int    `gorm:"comment:'销售id'" json:"saleId"`
	MainManageId int    `gorm:"comment:'主客户服务经理id'" json:"mainManageId"`
	SubManageId  int    `gorm:"comment:'副客户服务经理id'" json:"subManageId"`
}

// 客户额度
type Amount struct {
	gorm.Model
	ClientId    int    `gorm:"index;not null;comment:'客户id'" json:"clientId"`
	ServiceId   int    `gorm:"index;not null;comment:'服务id'" json:"serviceId"`
	Amount      int    `gorm:"not null;comment:'剩余额度'" json:"amount"`
	Deadline    Time   `gorm:"type:date;comment:'到期时间'" json:"deadline"`
	OrderNumber string `gorm:"size:100;comment:'订单编号'" json:"orderNumber"`
	Remark      string `gorm:"size:100;comment:'备注'" json:"remark"`
}

// 额度变动
type AmountLog struct {
	gorm.Model
	AmountId int    `gorm:"not null;comment:'额度id'" json:"-"`
	Change   int    `gorm:"not null;comment:'额度变动'" json:"change"`
	Desc     string `gorm:"size:100;comment:'事项说明'" json:"desc"`
	RealTime Time   `gorm:"type:datetime;comment:'发生时间'" json:"realTime"`
	Refer    string `gorm:"size:100;comment:'额度转换关联'" json:"-"`
	Type     string `gorm:"comment:'变动类型delay,convert'" json:"-"`
	Remark   string `gorm:"size:100;comment:'备注'" json:"remark"`
	TaskId   int    `gorm:"comment:'任务退次关联'" json:"-"`
}

// 服务
type Service struct {
	gorm.Model
	ServiceName string `gorm:"size:60;not null;comment:'服务名称'" json:"serviceName"`
	State       int    `gorm:"type:tinyint;comment:'0启用1禁用'" json:"state"`
	Use         int    `gorm:"not null;comment:'1可实施2可转换'" json:"use"`
	Sort        int    `gorm:"index;comment:'排序字段'" json:"sort"`
}

type ClientAmount struct {
	ServiceName string `json:"service_name"`
	Type        string `json:"type"`
	ServiceId   int    `json:"service_id"`
	Deadline    Time   `json:"deadline"`
	Change      int    `json:"change"`
}

// 额度列表
type RspAmount struct {
	ServiceName string `json:"service_name"`
	//总数
	Amount    int  `json:"amount"`
	Used      int  `json:"used"`
	Remain    int  `json:"remain"`
	Delay     int  `json:"delay"`
	Deadline  Time `json:"deadline"`
	ServiceId int  `json:"service_id"`
}

func (amount *RspAmount) CalData(ca ClientAmount) {
	amount.Remain += ca.Change
	if ca.Type == Amount_ConvOut || ca.Type == Amount_Use {
		amount.Used += ca.Change * AmountChange[Amount_Use]
	} else if ca.Type == Amount_Delay {
		amount.Delay += ca.Change * AmountChange[Amount_Delay]
	} else if ca.Type == Amount_Buy || ca.Type == Amount_ConvIn {
		amount.Amount += ca.Change
	} else if ca.Type == Amount_Cancel {
		amount.Used -= ca.Change
	}
}

// 额度历史
type RspAmountLog struct {
	RealTime    Time   `json:"real_time"`
	ServiceName string `json:"service_name"`
	OrderNumber string `json:"order_number"`
	Change      int    `json:"change"`
	Desc        string `json:"desc"`
	Remark      string `json:"remark"`
}

type ReqSwitchAmount struct {
	ClientId int    `json:"clientId"`
	SOutId   int    `json:"sOutId"`
	SOutNum  int    `json:"sOutNum"`
	SInId    int    `json:"sInId"`
	SInNum   int    `json:"sInNum"`
	Remark   string `json:"remark"`
}

type AmountSimple struct {
	Id          int    `json:"id"`
	Amount      int    `json:"amount"`
	OrderNumber string `json:"order_number"`
}

type RspClient struct {
	Client
	Sale       User `json:"sale"`
	Manager    User `json:"manager"`
	SubManager User `json:"subManager"`
}
