/*
* Auth : acer
* Desc : 申请设备
* Time : 2020/9/1 22:37
 */

package oa

import (
	"github.com/jinzhu/gorm"
)

// DeviceRequisition DeviceRequisition
type DeviceRequisition struct {
	gorm.Model
	DeviceID              int     `gorm:"not null;comment:'设备ID'" json:"device_id"`
	Device                *Device `json:"device"`
	AssociateEmployeeID   int     `gorm:"not null;comment:'关联EmployeeID'" json:"associate_employee_id"`
	AssociateEmployeeName string  `gorm:"not null;comment:'关联EmployeeName'" json:"associate_employee_name"`
	OperatorCategory      string  `gorm:"size:20;not null;comment:'操作类别(入库,报废,借出,归还)'" json:"operator_category"`
	OperatorID            int     `gorm:"not null;comment:'出库操作人ID，关联EmployeeID'" json:"outgoing_operator_id"`
	OperatorName          string  `gorm:"not null;comment:'出库操作人Name，关联EmployeeName'" json:"outgoing_operator_name"`
	Comment               string  `gorm:"size:2000;not null;comment:'备注'" json:"comment"`
}
