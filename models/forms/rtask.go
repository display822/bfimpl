/*
* Auth : acer
* Desc : task请求参数
* Time : 2020/7/7 21:51
 */

package forms

import "bfimpl/models"

//创建任务
type ReqTask struct {
	ClientId   int         `json:"clientId"`
	AppName    string      `json:"appName"`
	ServiceId  int         `json:"serviceId"`
	PreAmount  int         `json:"preAmount"`
	PreDate    models.Time `json:"preDate"`
	ExpEndDate models.Time `json:"expEndDate"`
	ManageId   int         `json:"manageId"`
}

//取消任务参数
type ReqCancelTask struct {
	UserId int    `json:"userId"`
	Reason string `json:"reason"`
}
