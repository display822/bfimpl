/*
* Auth : acer
* Desc : 设备
* Time : 2020/9/1 22:22
 */

package oa

import (
	"bfimpl/models"
	"bfimpl/services/log"
	"fmt"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
)

var DeviceCategoryMap = map[string]struct{}{
	"PC":      {},
	"Laptop":  {},
	"iMac":    {},
	"Mobile":  {},
	"Pad":     {},
	"Monitor": {},
	"Network": {},
}

var DeviceTodoStatusLeaderMap = map[string][]string{
	"0": {models.FlowProcessing},                    // 待审核
	"1": {models.FlowApproved, models.FlowRejected}, // 审核完成
}

// Device 设备表
type Device struct {
	gorm.Model
	DeviceCode         string               `gorm:"size:50;not null;comment:'设备编码'" json:"device_code"`
	DeviceCategory     string               `gorm:"size:20;not null;comment:'(PC,Laptop,iMac,Mobile,Pad,Monitor,Network)'" json:"device_category"`
	Brand              string               `gorm:"size:60;not null;comment:'品牌'" json:"brand"`
	DeviceName         string               `gorm:"size:100;not null;comment:'设备名称'" json:"device_name"`
	DeviceModel        string               `gorm:"size:100;not null;comment:'设备型号'" json:"device_model"`
	SharedDevice       string               `gorm:"size:10;not null;comment:'公共资产(是,否)'" json:"shared_device"`
	IngoingOperatorID  int                  `gorm:"size:10;not null;comment:'入库人id'" json:"ingoing_operator_id"`
	IngoingTime        models.Time          `gorm:"type:datetime;not null;comment:'入库时间'" json:"ingoing_time"`
	DeviceStatus       string               `gorm:"size:20;not null;comment:'(free:空闲,occupy:占用,fix:修理,scrap:报废)'" json:"device_status"`
	CPU                string               `grom:"size:20;comment:'CPU'" json:"cpu"`
	GPU                string               `grom:"size:20;comment:'GPU'" json:"gpu"`
	MEM                string               `grom:"size:20;comment:'内存'" json:"mem"`
	Volume             string               `grom:"size:20;comment:'存储容量'" json:"volume"`
	OS                 string               `grom:"size:20;comment:'(Windows,Linux,iOS,Android,Mac)'" json:"os"`
	Core               string               `grom:"size:20;comment:'核心'" json:"core"`
	Version            string               `grom:"size:30;comment:'版本'" json:"version"`
	ScreenSize         string               `grom:"size:30;comment:'屏幕尺寸'" json:"screen_size"`
	Resolution         string               `grom:"size:30;comment:'分辨率'" json:"resolution"`
	AspectRatio        string               `grom:"size:20;comment:'屏幕比'" json:"aspect_ratio"`
	MACAddress1        string               `grom:"size:80;comment:'MAC地址1'" json:"mac_address_1"`
	MACAddress2        string               `grom:"size:80;comment:'MAC地址2'" json:"mac_address_2"`
	Retailer           string               `gorm:"size:100;not null;comment:'零售商'" json:"retailer"`
	PurchasePrice      float64              `gorm:"type:decimal(10,2);not null;comment:'购买价格'" json:"purchase_price"`
	PurchaseDate       models.Time          `gorm:"type:datetime;comment:'购买日期'" json:"purchase_date"`
	VAT                float64              `gorm:"type:decimal(10,2);not null;comment:'增值税金额'" json:"vat"`
	WarrantyPeriod     int                  `gorm:"not null;comment:'保修期限'" json:"warranty_period"`
	Site               string               `gorm:"size:100;not null;comment:'位置'" json:"site"`
	DeviceApplyID      int                  `gorm:"size:50;comment:'申请单id'" json:"device_apply_id"`
	DeviceApply        *DeviceApply         `gorm:"ForeignKey:DeviceApplyID" json:"device_apply"`
	IsApply            int                  `gorm:"size:10;not null;comment:'是否可申领'" json:"is_apply"`
	DeviceRequisitions []*DeviceRequisition `json:"device_requisitions"`
	DeviceApplys       []*DeviceApply       `json:"device_applys"`
	CanApply           bool                 `gorm:"-" json:"can_apply"`
	Depreciate         float64              `gorm:"-" json:"depreciate"`
}

// DeviceApply 设备申请表
type DeviceApply struct {
	gorm.Model
	DeviceID             int         `gorm:"size:50;not null;comment:'设备ID'" json:"device_id"`
	Device               *Device     `gorm:"ForeignKey:DeviceID" json:"device"`
	EngagementCode       string      `gorm:"size:64;comment:'任务指派编码'" json:"engagement_code"`
	EmpID                int         `gorm:"comment:'报销申请人id'" json:"emp_id"`
	Employee             *Employee   `gorm:"ForeignKey:EmpID" json:"employee"`
	EName                string      `gorm:"size:30;comment:'员工姓名'" json:"e_name"`
	Status               string      `gorm:"size:20;comment:'申请状态'" json:"status"`
	Type                 int         `gorm:"size:10;default:0;comment:'类型(0: 自申请，1: 管理员分配)'" json:"type"`
	Project              string      `gorm:"size:64;comment:'项目'" json:"project"`
	ApplicationDate      time.Time   `gorm:"type:datetime;comment:'申请日期'" json:"application_date"`
	ReceiveDate          time.Time   `gorm:"type:datetime;comment:'领用日期'" json:"receive_date"`
	IsReturn             int         `gorm:"not null;comment:'是否归还'" json:"is_return"`
	OutgoingOperatorID   int         `gorm:"not null;comment:'出库操作人ID，关联EmployeeID'" json:"outgoing_operator_id"`
	OutgoingOperatorName string      `gorm:"not null;comment:'出库操作人Name，关联EmployeeName'" json:"outgoing_operator_name"`
	OutgoingTime         models.Time `gorm:"type:datetime;comment:'出库时间'" json:"outgoing_time"`
	LeaderId             int         `gorm:"-" json:"leader_id"`
	CanReceive           bool        `gorm:"-" json:"can_receive"` // 是否能领用
	CanRevoke            bool        `gorm:"-" json:"can_revoke"`  // 是否能撤销
}

// DeviceApplyInfo 申请设备基本信息
type DeviceApplyInfo struct {
	Employee       *Employee `json:"employee"`
	Device         *Device   `json:"device"`
	CollectDevices string    `json:"collect_devices"`
}

// SetDepreciate 设置折旧信息
func (d *Device) SetDepreciate() {
	t := time.Time(d.IngoingTime)
	ingoingTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	log.GLogger.Info("ingoingTime: %s", ingoingTime)

	ingoingTimeAfter1Month := ingoingTime.AddDate(0, 1, 0)
	ingoingTimeAfter35Month := ingoingTime.AddDate(0, 35, 0)
	ingoingTimeAfter36Month := ingoingTime.AddDate(0, 36, 0)
	log.GLogger.Info("ingoingTimeAfter1Month: %s", ingoingTimeAfter1Month)
	log.GLogger.Info("ingoingTimeAfter35Month: %s", ingoingTimeAfter35Month)
	log.GLogger.Info("ingoingTimeAfter36Month: %s", ingoingTimeAfter36Month)

	nowTime := time.Now()
	log.GLogger.Info("nowTime: %s", nowTime)

	if nowTime.After(ingoingTime) && nowTime.Before(ingoingTimeAfter1Month) { // 第1个月
		d.Depreciate = 0
	} else if nowTime.After(ingoingTimeAfter36Month) { // 第36个月之后
		d.Depreciate = 0
	} else if nowTime.After(ingoingTimeAfter1Month) && nowTime.Before(ingoingTimeAfter36Month) { // 第2-36个月
		// 计算每个月的折旧
		depreciate := d.PurchasePrice / 36
		log.GLogger.Info("depreciate: %s", depreciate)
		monthDepreciate, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", depreciate), 64)
		if nowTime.After(ingoingTimeAfter35Month) {
			log.GLogger.Info("nowTime After ingoingTimeAfter35Month")
			d.Depreciate = d.PurchasePrice - monthDepreciate*35 // 第36个月
		} else {
			d.Depreciate = depreciate
		}
	}
}
