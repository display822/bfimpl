/*
* Auth : acer
* Desc : oa请求
* Time : 2020/9/13 11:13
 */

package forms

//加班审批信息
type ReqApprovalOvertime struct {
	Id      int    `json:"id"`
	Status  int    `json:"status"`
	Comment string `json:"comment"`
}

//报销审批信息
type ReqApprovalExpense struct {
	Id      int    `json:"id"`
	Status  int    `json:"status"`
	Comment string `json:"comment"`
}

// 设备审批信息
type ReqApprovalDevice struct {
	Id      int    `json:"id"`
	Status  int    `json:"status"`
	Comment string `json:"comment"`
}
