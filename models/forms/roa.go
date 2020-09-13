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
