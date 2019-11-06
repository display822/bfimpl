package controllers

import (
	"api-project-go/models"
	"api-project-go/models/forms"
	"api-project-go/services"
	"github.com/astaxie/beego/httplib"
	"net/http"
	"time"
)

type DemoController struct {
	BaseController
}

// 展示路由和注释的使用

// @Title 获取用户
// @Description 通过用户ID获取用户信息
// @Param	uid	path	string	true	"用户id"
// @Success 200 {object} models.User
// @Failure 403 invlaid param
// @Failure 404 user not found
// @router /user/:uid [get]
func (d *DemoController) GetUser() {
	uid := d.GetString(":uid")
	if uid == "" {
		d.Error("invlaid param")
		return
	}
	u, err := models.GetUser(uid)
	if err != nil {
		d.ErrorCode(http.StatusNotFound, ErrCodeRet, "user not found")
		return
	}
	d.Correct(u)
}

// @Title 获取所有用户
// @Description 拉取所有用户
// @Success 200 {object} []models.User
// @router /users [get]
func (d *DemoController) GetAllUsers() {
	users := models.GetAllUsers()
	d.Correct(users)
}

// 展示Redis的使用，以及表单的验证. GET/POST/DELETE请求

// @Title 设置Redis的值
// @Description 设置Redis的值
// @Param key formData string true "redis设置的key"
// @Param value formData string true "redis设置的value"
// @Success 200 {string} success
// @Failure 403 invlaid param
// @router /redis/string [post]
func (d *DemoController) SetRedis() {
	var req forms.RedisStringReq
	if err := d.ParseForm(&req); err != nil {
		d.Error("invlaid param")
		return
	}
	_, _ = d.valid.Valid(&req)
	if d.valid.HasErrors() {
		d.Correct(req)
		return
	}
	err := services.RedisClient().Set(req.Key, req.Value, -1).Err()
	if err != nil {
		d.Error(err.Error())
		return
	}
	d.Correct(req)
}

// @Title 获取Redis的值
// @Description 设置Redis的值
// @Param key path string true "redis的key"
// @Success 200 {string} data
// @Failure 403 invlaid param
// @router /redis/string/:key [get]
func (d *DemoController) GetRedis() {
	key := d.GetString(":key")
	if key == "" {
		d.Error("invlaid param")
		return
	}
	value, err := services.RedisClient().Get(key).Result()
	if err != nil {
		d.Error(err.Error())
		return
	}
	d.Correct(value)
}

// @Title 获取数据库表信息
// @Description 获取数据库表信息
// @Success 200 {object} []models.Table
// @router /mysql/tables [get]
func (d *DemoController) GetMySQLInfo() {
	tables := new(models.Table).GetAllTables()
	d.Correct(tables)
}

// @Title 示例http请求如何发送
// @Description 示例http请求如何发送
// @Success 200 {object}
// @router /http/get [get]
func (d *DemoController) HttpRequestDemo() {
	var req struct {
		Ret int    `json:"ret"`
		Msg string `json:"msg"`
	}
	err := httplib.NewBeegoRequest("http://localhost:8080", http.MethodGet).
		SetTimeout(time.Second, time.Second).
		ToJSON(&req)
	if err != nil {
		d.Error("req error=" + err.Error())
		return
	}
	d.Correct(req)
}
