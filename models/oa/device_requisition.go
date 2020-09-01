/*
* Auth : acer
* Desc : 申请设备
* Time : 2020/9/1 22:37
 */

package oa

import (
	"bfimpl/models"

	"github.com/jinzhu/gorm"
)

type DeviceRequisition struct {
	gorm.Model
	DeviceID           int         `gorm:"not null;comment:'设备ID'" json:"device_id"`
	UserID             int         `gorm:"not null;comment:'借用人ID，关联EmployeeID'" json:"user_id"`
	Status             string      `gorm:"not null;comment:'(已领用,已归还,已灭失)'" json:"status"`
	OutgoingOperatorID int         `gorm:"not null;comment:'出库操作人ID，关联EmployeeID'" json:"outgoing_operator_id"`
	OutgoingTime       models.Time `gorm:"type:datetime;comment:'出库时间'" json:"outgoing_time"`
	IngoingOperatorID  int         `gorm:"not null;comment:'入库操作人ID，关联EmployeeID'" json:"ingoing_operator_id"`
	IngoingTime        models.Time `gorm:"type:datetime;comment:'入库时间'" json:"ingoing_time"`
}
