/*
* Auth : acer
* Desc : 登录
* Time : 2020/7/10 9:49
 */

package controllers

import (
	"bfimpl/models"
	"bfimpl/services"
	"bfimpl/services/log"
	"encoding/json"
)

type LoginController struct {
	BaseController
}

// @Title 登录
// @Description 登录
// @Param	username	body	string	true	"用户名"
// @Param	password	body	string	true	"密码"
// @Success 200 {object} models.User
// @Failure 500 server err
// @router / [post]
func (l *LoginController) Login() {
	param := new(models.ReqLogin)
	err := json.Unmarshal(l.Ctx.Input.RequestBody, param)
	if err != nil {
		l.ErrorOK(MsgInvalidParam)
	}

	b, e := services.LdapService().Login(param.UserName, param.Password)
	if e != nil {
		log.GLogger.Error(e.Error())
	}
	if !b {
		l.ErrorOK("login fail")
	}
	l.Correct("")
}
