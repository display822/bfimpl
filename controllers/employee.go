/*
* Auth : acer
* Desc : 入职离职流程
* Time : 2020/9/4 21:45
 */

package controllers

type EmployeeController struct {
	BaseController
}

// @Title hr新建入职
// @Description 新建入职
// @Param	json	body	string	true	"入职员工信息"
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /new [post]
func (s *EmployeeController) NewEmpEntry() {

}
