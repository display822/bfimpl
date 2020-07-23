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
	"bfimpl/services/util"
	"encoding/json"
	"math/rand"
	"strconv"
	"time"
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
	//ldap登录成功后，查看信息是否录入
	var user models.User
	err = services.Slave().Where("email = ?", param.UserName+"@broadfun.cn").First(&user).Error
	if err != nil {
		l.ErrorOK("user not found")
	}

	key := strconv.Itoa(rand.Intn(100)) + strconv.FormatInt(time.Now().Unix(), 10)
	user.Session = util.StringMd5(key)
	services.RedisClient().Set(user.Session, user.ID, time.Hour*24)
	l.Correct(user)
}
