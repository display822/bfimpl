package controllers

import (
	"bfimpl/models"
	"bfimpl/services"
	"bfimpl/services/log"
	"encoding/json"
)

type UserController struct {
	BaseController
}

// @Title 新增用户
// @Description 新增用户
// @Param	name	body	string	true	"姓名"
// @Param	email	body	string	true	"邮箱"
// @Param	wx		body	string	true	"企业微信"
// @Param	phone	body	string	true	"手机"
// @Param	userType	body	int	true	"用户类型"
// @Param	leaderId	body	int	false	"组长id"
// @Success 200 {object} models.User
// @Failure 500 server err
// @router / [post]
func (u *UserController) AddUser() {
	reqUser := new(models.User)

	err := json.Unmarshal(u.Ctx.Input.RequestBody, reqUser)
	if err != nil {
		u.ErrorOK(MsgInvalidParam)
	}
	b, _ := u.valid.Valid(reqUser)
	if !b {
		log.GLogger.Error("%s:%s", u.valid.Errors[0].Field, u.valid.Errors[0].Message)
		u.ErrorOK(MsgInvalidParam)
	}

	err = services.Slave().Create(reqUser).Error
	if err != nil {
		u.ErrorOK(err.Error())
	}
	u.Correct(reqUser)
}

// @Title 资源分配人员列表
// @Description 无
// @Success 200 {object} []models.User
// @Failure 500 server err
// @router /leaders [get]
func (u *UserController) GroupLeaders() {
	//userType = 4
	users := make([]*models.User, 0)
	err := services.Slave().Where("user_type = ?", 4).Find(&users).Error
	if err != nil {
		u.ErrorOK(err.Error())
	}
	u.Correct(users)
}

// @Title 人员列表
// @Description 按类型查询
//  1: "admin",
//	2: "sale",
//	3: "manager",
//	4: "tm",
//	5: "implement",
// @Param	type		query	int		true	"类型"
// @Param	pageSize	query	int		true	"每页条数"
// @Param	pageNum		query	int		true	"页数"
// @Success 200 {object} []models.User
// @Failure 500 server err
// @router /list [get]
func (u *UserController) UserList() {
	userType, _ := u.GetInt("type", 0)
	pageSize, _ := u.GetInt("pageSize", 10)
	pageNum, _ := u.GetInt("pageNum", 1)

	users := make([]*models.User, 0)
	query := services.Slave().Model(models.User{})
	if userType != 0 {
		query = query.Where("user_type = ?", userType)
	}
	total := 0
	err := query.Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&users).Limit(-1).Offset(-1).Count(&total).Error
	if err != nil {
		u.ErrorOK(err.Error())
	}
	var res = struct {
		Total int            `json:"total"`
		Users []*models.User `json:"users"`
	}{
		Total: total,
		Users: users,
	}

	u.Correct(res)
}
