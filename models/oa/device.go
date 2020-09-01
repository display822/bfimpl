/*
* Auth : acer
* Desc : 设备
* Time : 2020/9/1 22:22
 */

package oa

import (
	"bfimpl/models"

	"github.com/jinzhu/gorm"
)

type Device struct {
	gorm.Model
	DeviceNum         string  `gorm:"size:50;not null;comment:'设备编号'" json:"device_num"`
	DeviceCategory    string  `gorm:"size:20;not null;comment:'(PC,Laptop,iMac,Mobile,Pad,Monitor,Network)'" json:"device_category"`
	Brand             string  `gorm:"size:60;not null;comment:'品牌'" json:"brand"`
	DeviceName        string  `gorm:"size:100;not null;comment:'设备名称'" json:"device_name"`
	DeviceModel       string  `gorm:"size:100;not null;comment:'设备型号'" json:"device_model"`
	Retailer          string  `gorm:"size:100;not null;comment:'零售商'" json:"retailer"`
	PurchasePrice     float32 `gorm:"type:decimal;not null;comment:'购买价格'" json:"purchase_price"`
	VAT               float32 `gorm:"type:decimal;not null;comment:'增值税金额'" json:"vat"`
	WarrantyPeriod    int     `gorm:"type:decimal;not null;comment:'保修期限'" json:"warranty_period"`
	SharedDevice      string  `gorm:"size:10;not null;comment:'公共资产(是,否)'" json:"shared_device"`
	IngoingOperatorID int     `gorm:"size:10;not null;comment:'入库人id'" json:"ingoing_operator_id"`
	Employee1         Employee
	IngoingTime       models.Time `gorm:"type:datetime;not null;comment:'入库时间'" json:"ingoing_time"`
	DeviceStatus      string      `gorm:"size:20;not null;comment:'(在库,出库,报废)'" json:"device_status"`
	CPU               string      `grom:"size:20;comment:'CPU'" json:"cpu" json:"cpu"`
	MEM               string      `grom:"size:20;comment:'内存'" json:"mem"`
	Volume            string      `grom:"size:20;comment:'硬盘容量/运存'" json:"volume"`
	OS                string      `grom:"size:20;comment:'(Windows,Linux,iOS,Android,Mac)'" json:"os"`
	Version           string      `grom:"size:30;comment:'版本'" json:"version"`
	ScreenSize        string      `grom:"size:30;comment:'屏幕尺寸'" json:"screen_size"`
	Resolution        string      `grom:"size:30;comment:'分辨率'" json:"resolution"`
	AspectRatio       string      `grom:"size:20;comment:'屏幕比'" json:"aspect_ratio"`
	MACAddress1       string      `grom:"size:80;comment:'MAC地址1'" json:"mac_address_1"`
	MACAddress2       string      `grom:"size:80;comment:'MAC地址2'" json:"mac_address_2"`
}
