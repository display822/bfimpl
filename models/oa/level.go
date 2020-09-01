/*
* Auth : acer
* Desc : 级别
* Time : 2020/9/1 22:16
 */

package oa

import "github.com/jinzhu/gorm"

type Level struct {
	gorm.Model
	DepartmentID int     `gorm:"not null;comment:'归属部门'" json:"department_id"`
	LevelName    string  `gorm:"size:50;not null;comment:'级别名称'" json:"level_name"`
	CCRate       float32 `gorm:"type:decimal(5,2);not null;comment:''" json:"cc_rate"`
	OCRate       float32 `gorm:"type:decimal(5,2);not null;comment:''" json:"oc_rate"`
}
