/*
* Auth : acer
* Desc : 工作台
* Time : 2020/9/26 16:18
 */

package controllers

import (
	"bfimpl/models"
	"bfimpl/models/oa"
	"bfimpl/services"
)

type BenchController struct {
	BaseController
}

// @Title 工作台我的审批
// @Description 工作台
// @Param	type	query	string	true	"待办，已办"
// @Param	pagenum	    query	int	false	"分页"
// @Param	pagesize	query	int	false	"分页"
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /myapprove [get]
func (b *BenchController) GetMyApprove() {
	pageSize, _ := b.GetInt("pagesize", 10)
	pageNum, _ := b.GetInt("pagenum", 1)
	userID, _ := b.GetInt("userID", 0)
	status := b.GetString("type", "todo")

	//待办
	flowStatus := []string{models.FlowProcessing}
	if status != "todo" {
		//已办
		flowStatus = []string{models.FlowApproved, models.FlowCompleted}
	}
	var resp struct {
		Total int            `json:"total"`
		List  []*oa.Workflow `json:"list"`
	}
	nodes := make([]*oa.WorkflowNode, 0)
	services.Slave().Model(oa.WorkflowNode{}).Where("node_seq != 1 and operator_id = ? and status in (?)", userID,
		flowStatus).Order("created_at desc").Find(&nodes)
	flowIds := make([]int, 0)
	ex := make(map[int]bool)
	for _, n := range nodes {
		if _, ok := ex[n.WorkflowID]; !ok {
			ex[n.WorkflowID] = true
			flowIds = append(flowIds, n.WorkflowID)
		}
	}
	resp.Total = len(flowIds)
	start, end := getPage(resp.Total, pageSize, pageNum)
	//查询flow
	flowIds = flowIds[start:end]
	services.Slave().Model(oa.Workflow{}).Preload("WorkflowDefinition").Preload("Nodes").
		Preload("Nodes.User").Where(flowIds).Find(&resp.List)
	b.Correct(resp)
}

// @Title 工作台我的申请
// @Description 工作台
// @Param	type	query	string	true	"待办，已办"
// @Param	pagenum	    query	int	false	"分页"
// @Param	pagesize	query	int	false	"分页"
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /myreq [get]
func (b *BenchController) GetMyRequest() {
	pageSize, _ := b.GetInt("pagesize", 10)
	pageNum, _ := b.GetInt("pagenum", 1)
	userID, _ := b.GetInt("userID", 0)
	status := b.GetString("type", "todo")

	var resp struct {
		Total int            `json:"total"`
		List  []*oa.Workflow `json:"list"`
	}
	//待办
	flowStatus := []string{models.FlowProcessing}
	if status != "todo" {
		//已办
		flowStatus = []string{models.FlowRejected, models.FlowApproved, models.FlowCompleted}
	}
	//我的申请
	flows := make([]*oa.WorkflowId, 0)
	services.Slave().Raw("select w.id from workflows w,workflow_nodes wn where w.id = wn.workflow_id and "+
		"operator_id = ? and node_seq = 1 and w.status in (?)", userID, flowStatus).Scan(&flows)
	resp.Total = len(flows)
	if resp.Total > pageSize {
		start := (pageNum - 1) * pageSize
		end := start + pageSize
		if start > resp.Total {
			start = 0
			end = 0
		} else {
			if end > resp.Total {
				end = resp.Total
			}
		}
		flows = flows[start:end]
	}
	ids := make([]int, 0)
	for _, f := range flows {
		ids = append(ids, f.Id)
	}
	//查询flow
	services.Slave().Model(oa.Workflow{}).Preload("WorkflowDefinition").Preload("Nodes").
		Preload("Nodes.User").Where(ids).Find(&resp.List)
	b.Correct(resp)
}
