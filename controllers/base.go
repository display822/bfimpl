package controllers

import (
	"bfimpl/models"
	"bfimpl/services/log"
	"net/http"
	"strconv"

	"strings"

	"bfimpl/services"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
)

const (
	ErrCodeRet      = 1
	MsgInvalidParam = "invalid param"
	MsgServerErr    = "inter server err"
)

type BaseController struct {
	beego.Controller
	valid validation.Validation
}

var AutoLogError bool

// Correct 返回正常的信息，json经过了编码
func (b *BaseController) Correct(data interface{}) {
	b.CorrectEncoding(data, true)
}

// CorrectEncoding 可以json不编码返回
func (b *BaseController) CorrectEncoding(data interface{}, encoding bool) {
	var ret map[string]interface{}
	var ok bool
	if ret, ok = data.(map[string]interface{}); !ok {
		ret = make(map[string]interface{})
		ret["data"] = data
	}
	ret["ret"] = 0
	ret["msg"] = "success"
	b.Data["json"] = ret
	b.ServeJSON(encoding)
	b.StopRun()
}

// Error 最常用的返回错误，返回错误信息即可，错误码为1
func (b *BaseController) Error(msg string) {
	b.ErrorCode(1, http.StatusForbidden, msg)
}

func (b *BaseController) ErrorOK(msg string) {
	b.ErrorCode(1, http.StatusOK, msg)
}

// ErrorCode 可以定制返回错误码
func (b *BaseController) ErrorCode(code int, status int, msg string) {
	b.Data["json"] = struct {
		Ret int    `json:"ret"`
		Msg string `json:"msg"`
	}{code, msg}
	if AutoLogError {
		log.GLogger.Error("requst=%v ret=%v msg=%v", b.Ctx.Request.URL.Path, code, msg)
	}
	b.Ctx.ResponseWriter.WriteHeader(status)
	b.ServeJSON()
	b.StopRun()
}

// ErrorObject 返回错误码，错误提示，和附带一些错误信息
func (b *BaseController) ErrorObject(code int, msg string, obj interface{}) {
	b.Data["json"] = struct {
		Ret  int         `json:"ret"`
		Msg  string      `json:"msg"`
		Data interface{} `json:"data"`
	}{code, msg, obj}
	b.ServeJSON()
	b.StopRun()
}

func (b *BaseController) Prepare() {
	if strings.Contains(b.Ctx.Request.URL.Path, "login") {
		return
	}
	userKey := b.Ctx.Request.Header.Get("Authorization")
	userID, err := services.RedisClient().Get(userKey).Result()
	if err != nil {
		b.ErrorCode(http.StatusUnauthorized, http.StatusOK, "login required")
	}
	//查询用户
	var user models.User
	services.Slave().Take(&user, "id = ?", userID)
	b.Input().Set("userID", userID)
	b.Input().Set("userName", user.Name)
	b.Input().Set("userType", strconv.Itoa(user.UserType))
}
