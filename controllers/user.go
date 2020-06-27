package controllers

import (
	"bfimpl/models"
	"bfimpl/models/forms"
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
// @Success 200 {object} models.User
// @Failure 500 server err
// @router / [post]
func (u *UserController) AddUser() {
	reqUser := new(forms.ReqUser)

	err := json.Unmarshal(u.Ctx.Input.RequestBody, reqUser)
	if err != nil {
		u.ErrorOK(MsgInvalidParam)
	}
	b, e := u.valid.Valid(reqUser)
	if !b {
		u.ErrorOK(e.Error())
	}
	user := models.NewUser(reqUser)
	err = user.Create()
	if err != nil {
		u.ErrorOK(err.Error())
	}
	u.Correct(user)
}

// @Title 资源分配人员列表
// @Description 无
// @Success 200 {object} []models.User
// @Failure 500 server err
// @router /leaders [get]
func (u *UserController) GroupLeaders() {
	//userType = 4
	users, e := models.GetLeaders()
	if e != nil {
		u.ErrorOK(e.Error())
	}
	u.Correct(users)
}
