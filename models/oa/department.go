/*
* Auth : acer
* Desc : 部门
* Time : 2020/9/1 22:13
 */

package oa

import (
	"bfimpl/models"

	"github.com/jinzhu/gorm"
)

type Department struct {
	gorm.Model
	DepartmentName     string        `gorm:"size:50;not null;comment:'部门名称'" json:"department_name"`
	DepartmentLeaderID int           `gorm:"not null;comment:'领导id'" json:"department_leader_id"`
	PID                int           `gorm:"default:0;comment:'父部门id'" json:"-"`
	Leader             *models.User  `gorm:"ForeignKey:DepartmentLeaderID" json:"leader"`
	Children           []*Department `json:"children"`
}

type ServiceLine struct {
	gorm.Model
	DepartmentID int `gorm:"not null;comment:'归属部门'" json:"department_id"`
	ServiceName  int `gorm:"not null;comment:'服务线名称'" json:"service_name"`
}
