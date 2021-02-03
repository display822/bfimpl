/*
* Auth : acer
* Desc : 部门
* Time : 2020/9/4 22:45
 */

package controllers

import (
	"bfimpl/models/oa"
	"bfimpl/services"
	"bfimpl/services/log"

	"github.com/jinzhu/gorm"
)

type DepartmentController struct {
	BaseController
}

// @Title 部门列表
// @Description 部门列表
// @Param	a	query	string	true	""
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /list [get]
func (d *DepartmentController) GetDepartments() {
	dps := make([]*oa.Department, 0)
	err := services.Slave().Model(oa.Department{}).Preload("Leader").Find(&dps).Error
	if err != nil {
		log.GLogger.Error("get departments err:%s", err.Error())
		d.ErrorOK("get departments err")
	}
	d.Correct(dps)
}

// @Title 部门下级别
// @Description 部门下级别
// @Param	id	path	int	true	"部门id"
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /level/:id [get]
func (d *DepartmentController) GetLevels() {
	departmentId, _ := d.GetInt(":id", 0)
	levels := make([]*oa.Level, 0)
	err := services.Slave().Model(oa.Level{}).Where("department_id = ?", departmentId).Find(&levels).Error
	if err != nil {
		log.GLogger.Error("get levels err:%s", err.Error())
		d.ErrorOK("get levels err")
	}
	d.Correct(levels)
}

// @Title 部门下服务线
// @Description 部门下服务线
// @Param	id	path	int	true	"部门id"
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /service/:id [get]
func (d *DepartmentController) GetServiceLine() {
	departmentId, _ := d.GetInt(":id", 0)
	s := make([]*oa.ServiceLine, 0)
	services.Slave().Model(oa.ServiceLine{}).Where("department_id = ?", departmentId).Find(&s)
	if len(s) == 0 {
		s = append(s, &oa.ServiceLine{
			Model:        gorm.Model{ID: 0},
			DepartmentID: departmentId,
			ServiceName:  "无",
		})
	}
	d.Correct(s)
}
