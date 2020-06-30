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
	SaleId       int    `gorm:"not null;comment:'销售id'" json:"-"`
	MainManageId int    `gorm:"not null;comment:'主客户服务经理id'" json:"-"`
	SubManageId  int    `gorm:"not null;comment:'副客户服务经理id'" json:"-"`
}

// 客户额度
type Amount struct {
	gorm.Model
	ServiceId string `gorm:"size:60;not null;comment:'服务id'" json:"serviceId"`
	Deadline  Time   `gorm:"type:datetime;comment:'到期时间'" json:"deadline"`
	Amount    int    `gorm:"not null;comment:'剩余额度'" json:"amount"`
	ClientId  int    `gorm:"not null;comment:'客户id'" json:"-"`
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
