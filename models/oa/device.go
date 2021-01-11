/*
* Auth : acer
* Desc : 设备
* Time : 2020/9/1 22:22
 */

package oa

import (
	"bfimpl/models"
	"github.com/jinzhu/gorm"
	"time"
)

// Device 设备表
type Device struct {
	gorm.Model
	DeviceNum         string       `gorm:"size:50;not null;comment:'设备编号'" json:"device_num"`
	DeviceCategory    string       `gorm:"size:20;not null;comment:'(PC,Laptop,iMac,Mobile,Pad,Monitor,Network)'" json:"device_category"`
	Brand             string       `gorm:"size:60;not null;comment:'品牌'" json:"brand"`
	DeviceName        string       `gorm:"size:100;not null;comment:'设备名称'" json:"device_name"`
	DeviceModel       string       `gorm:"size:100;not null;comment:'设备型号'" json:"device_model"`
	SharedDevice      string       `gorm:"size:10;not null;comment:'公共资产(是,否)'" json:"shared_device"`
	IngoingOperatorID int          `gorm:"size:10;not null;comment:'入库人id'" json:"ingoing_operator_id"`
	IngoingTime       models.Time  `gorm:"type:datetime;not null;comment:'入库时间'" json:"ingoing_time"`
	DeviceStatus      string       `gorm:"size:20;not null;comment:'(free:空闲,occupy:占用,fix:修理,scrap:报废)'" json:"device_status"`
	CPU               string       `grom:"size:20;comment:'CPU'" json:"cpu"`
	GPU               string       `grom:"size:20;comment:'GPU'" json:"gpu"`
	MEM               string       `grom:"size:20;comment:'内存'" json:"mem"`
	Volume            string       `grom:"size:20;comment:'存储容量'" json:"volume"`
	OS                string       `grom:"size:20;comment:'(Windows,Linux,iOS,Android,Mac)'" json:"os"`
	Core              string       `grom:"size:20;comment:'核心'" json:"core"`
	Version           string       `grom:"size:30;comment:'版本'" json:"version"`
	ScreenSize        string       `grom:"size:30;comment:'屏幕尺寸'" json:"screen_size"`
	Resolution        string       `grom:"size:30;comment:'分辨率'" json:"resolution"`
	AspectRatio       string       `grom:"size:20;comment:'屏幕比'" json:"aspect_ratio"`
	MACAddress1       string       `grom:"size:80;comment:'MAC地址1'" json:"mac_address_1"`
	MACAddress2       string       `grom:"size:80;comment:'MAC地址2'" json:"mac_address_2"`
	Retailer          string       `gorm:"size:100;not null;comment:'零售商'" json:"retailer"`
	PurchasePrice     float32      `gorm:"type:decimal;not null;comment:'购买价格'" json:"purchase_price"`
	PurchaseDate      models.Time  `gorm:"type:datetime;not null;comment:'购买日期'" json:"purchase_date"`
	VAT               float32      `gorm:"type:decimal;not null;comment:'增值税金额'" json:"vat"`
	WarrantyPeriod    int          `gorm:"type:decimal;not null;comment:'保修期限'" json:"warranty_period"`
	Site              string       `gorm:"size:100;not null;comment:'位置'" json:"site"`
	IsApply           int          `gorm:"size:10;not null;comment:'是否可申领'" json:"is_apply"`
	DeviceLogs        []*DeviceLog `json:"device_logs"`
}

// DeviceLog 设备活动记录
type DeviceLog struct {
	gorm.Model
	DeviceID   int    `gorm:"size:50;not null;comment:'设备ID'" json:"device_id"`
	OperatorID int    `gorm:"size:10;not null;comment:'操作人id'" json:"operator_id"`
	Type       string `gorm:"size:10;not null;comment:'(入库,借出,归还)'" json:"type"`
	EID        int    `gorm:"size:10;not null;comment:'关联员工'" json:"eid"`
	EName      string `gorm:"size:10;not null;comment:'关联员工'" json:"ename"`
	Comment    string `gorm:"size:10;not null;comment:'备注'" json:"comment"`
}

// DeviceApply 设备申请表
type DeviceApply struct {
	gorm.Model
	DeviceID           int         `gorm:"size:50;not null;comment:'设备ID'" json:"device_id"`
	EngagementCode     string      `gorm:"size:64;comment:'任务指派编码'" json:"engagement_code"`
	EmpID              int         `gorm:"comment:'报销申请人id'" json:"emp_id"`
	Employee           *Employee   `gorm:"ForeignKey:EmpID" json:"employee"`
	EName              string      `gorm:"size:30;comment:'员工姓名'" json:"e_name"`
	Device             Device      `json:"device"`
	Status             string      `gorm:"size:20;comment:'申请状态'" json:"status"`
	Project            string      `gorm:"size:64;comment:'项目'" json:"project"`
	ApplicationDate    time.Time   `gorm:"type:date;comment:'申请日期'" json:"application_date"`
	OutgoingOperatorID int         `gorm:"not null;comment:'出库操作人ID，关联EmployeeID'" json:"outgoing_operator_id"`
	OutgoingTime       models.Time `gorm:"type:datetime;comment:'出库时间'" json:"outgoing_time"`
	LeaderId           int         `gorm:"-" json:"leader_id"`
}

// DeviceApplyInfo 申请设备基本信息
type DeviceApplyInfo struct {
	Employee       *Employee
	Device         *Device
	CollectDevices []*Device
}
