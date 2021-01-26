/*
* Auth : acer
* Desc : 报销
* Time : 2020/12/4 15:45
 */

package controllers

import (
	"bfimpl/models"
	"bfimpl/models/forms"
	"bfimpl/models/oa"
	"bfimpl/services"
	"bfimpl/services/log"
	"bfimpl/services/util"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/google/uuid"
)

type ExpenseController struct {
	BaseController
}

// @Title 报销列表
// @Description 报销列表
// @Param	pagenum	    query	int	false	"页码"
// @Param	pagesize	query	int	false	"页数"
// @Param	myreq	query	bool	false	"我的报销"
// @Param	mytodo	query	bool	false	"我的审核"
// @Param	status	query	int	false	"状态"
// @Param	todostatus	query	int	false	"0：代办；1：已办"
// @Param	name	query	string	false	"搜索申请人"
// @Param	searchid	query	int	false	"搜索编码"
// @Param	application_date_begin	query	int	false	"费用发生日期开始时间"
// @Param	application_date_end	query	int	false	"费用发生日期结束时间"
// @Success 200 {object} oa.Expense
// @Failure 500 server internal err
// @router / [get]
func (e *ExpenseController) List() {
	pageSize, _ := e.GetInt("pagesize", 10)
	pageNum, _ := e.GetInt("pagenum", 1)
	userType, _ := e.GetInt("userType", 0)
	name := e.GetString("name")
	myReq, _ := e.GetBool("myreq", false)
	myTodo, _ := e.GetBool("mytodo", false)
	status := e.GetString("status")
	userEmail := e.GetString("userEmail")
	searchID := e.GetString("searchid")
	applicationDateBegin := e.GetString("application_date_begin")
	applicationDateEnd := e.GetString("application_date_end")

	log.GLogger.Info("params", userEmail, userType, name, myReq, status, searchID, pageNum, pageSize, applicationDateBegin, applicationDateEnd)

	employee := new(oa.Employee)
	services.Slave().Where("email = ?", userEmail).First(employee)
	log.GLogger.Info("employee: %+v", employee)

	expenses := make([]*oa.Expense, 0)
	query := services.Slave().Debug().Model(oa.Expense{}).Preload("ExpenseDetails")
	if searchID != "" {
		query = query.Where("id like ?", fmt.Sprintf("%%%s%%", searchID))
	}
	if name != "" {
		query = query.Where("e_name like ?", "%"+name+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if applicationDateBegin != "" && applicationDateEnd != "" {
		query = query.Where("application_date > ?", applicationDateBegin).Where("application_date <= ?", applicationDateEnd)
	}

	var resp struct {
		Total int           `json:"total"`
		List  []*oa.Expense `json:"list"`
	}

	if myReq {
		query = query.Where("emp_id = ?", employee.ID)
	}

	eIDs := make([]int, 0)

	if myTodo {
		userID, _ := e.GetInt("userID", 0)
		log.GLogger.Info("userID：%d", userID)
		ids := make([]oa.EntityID, 0)
		var s []string
		if e.GetString("todostatus") != "" {
			if userType == models.UserFinance {
				s = oa.TodoStatusFinanceMap[e.GetString("todostatus")]
			} else {
				s = oa.TodoStatusLeaderMap[e.GetString("todostatus")]
			}
		}

		if len(s) == 0 {
			services.Slave().Debug().Raw("select w.entity_id from workflows w,workflow_nodes wn where w.id = "+
				"wn.workflow_id and w.workflow_definition_id = ? and operator_id = ? and wn.status <> ?"+
				" and wn.node_seq != 1 order by w.entity_id desc", services.GetFlowDefID(services.Expense), userID, models.FlowNA).Scan(&ids)
			// services.Slave().Debug().Raw("select w.entity_id from workflows w,workflow_nodes wn where w.id = "+
			// 	"wn.workflow_id and w.workflow_definition_id = ? and operator_id = ?"+
			// 	" and wn.node_seq != 1 order by w.entity_id desc", services.GetFlowDefID(services.Expense), userID).Scan(&ids)
		} else {
			services.Slave().Debug().Raw("select w.entity_id from workflows w,workflow_nodes wn where w.id = "+
				"wn.workflow_id and w.workflow_definition_id = ? and operator_id = ? and wn.status in (?)"+
				" and wn.node_seq != 1 order by w.entity_id desc", services.GetFlowDefID(services.Expense), userID, s).Scan(&ids)
		}

		for _, eID := range ids {
			eIDs = append(eIDs, eID.EntityID)
		}
		if len(eIDs) != 0 {
			query = query.Where(eIDs)
		}
	}

	query.Limit(pageSize).Offset((pageNum - 1) * pageSize).Order("created_at desc").Find(&expenses).Limit(-1).Offset(-1).Count(&resp.Total)
	resp.List = expenses
	e.Correct(resp)
}

// @Title 申请报销
// @Description 申请报销
// @Param	body	    body	oa.Expense	true	"报销"
// @Success 200 {object} oa.Expense
// @Failure 500 server internal err
// @router / [post]
func (e *ExpenseController) ReqExpense() {
	uID, _ := e.GetInt("userID", 0)
	uEmail := e.GetString("userEmail")
	code := e.GetString("code")
	log.GLogger.Info("ReqExpense query: %d, %s", uID, uEmail)
	t := time.Now()
	validBeginMonthTime := util.GetFirstDateOfMonth(t)
	validEndMonthTime := util.GetLastDateOfMonth(t)
	var otp oa.ExpenseOtp
	if code != "" {
		services.Slave().Where("code = ?", code).Where("emp_id=?", uID).
			Where("created_at < ?", validEndMonthTime).Where("created_at >= ?", validBeginMonthTime).Find(&otp)
		log.GLogger.Info("otp", otp)
		if otp.ID == 0 {
			e.ErrorOK("验证码错误")
		}
	}

	// 获取emp_info
	employee := new(oa.Employee)
	services.Slave().Preload("Department").Preload("Department.Leader").Take(employee, "email = ?", uEmail)
	if employee.ID == 0 {
		e.ErrorOK("未找到员工信息")
	}
	log.GLogger.Info("employee: %+v", employee)

	//查询engagementCode
	engagementCode := new(oa.EngagementCode)
	services.Slave().Model(oa.EngagementCode{}).Where("department_id = ?", employee.DepartmentID).First(engagementCode)
	log.GLogger.Info("engagementCode: %+v", engagementCode)

	param := new(oa.Expense)
	err := json.Unmarshal(e.Ctx.Input.RequestBody, param)
	if err != nil {
		log.GLogger.Error("parse expense info err:%s", err.Error())
		e.ErrorOK(MsgInvalidParam)
	}

	if param.EngagementCode == "" {
		e.ErrorOK("need engagement_code")
	}
	if param.Project == "" {
		e.ErrorOK("need project")
	}
	if param.ImportFile == "" {
		e.ErrorOK("need import_file")
	}
	if param.LeaderId <= 0 {
		e.ErrorOK("need leader_id")
	}
	if len(param.ExpenseDetails) == 0 {
		e.ErrorOK("need expense_details")
	}

	paidCardInfo := e.GetDebitCard(int(employee.ID))
	if paidCardInfo.CardID == "" {
		e.ErrorOK("无合同或银行卡信息")
	}

	param.EmpID = int(employee.ID)
	param.EName = employee.Name
	param.ApplicationDate = time.Now()
	param.Status = models.FlowNA
	// 费用总金额
	var expenseSummary float64

	for _, item := range param.ExpenseDetails {
		if item.ExpenseAmount <= 0 {
			e.ErrorOK("expense_amount error")
		}
		expenseSummary += item.ExpenseAmount
	}
	param.ExpenseSummary = expenseSummary
	log.GLogger.Info("expenseSummary: %v", expenseSummary)

	tx := services.Slave().Begin()
	// 创建报销
	err = tx.Create(param).Error
	if err != nil {
		log.GLogger.Error("create req expense err:%s", err.Error())
		tx.Rollback()
		e.ErrorOK(MsgServerErr)
	}

	log.GLogger.Info("engagementCode.FinanceID: %d", engagementCode.FinanceID)
	if engagementCode.FinanceID <= 0 {
		e.ErrorOK("no finance_id")
	}

	// 执行报销工作流
	err = services.ReqExpense(tx, int(param.ID), uID, param.LeaderId, engagementCode.FinanceID)
	if err != nil {
		log.GLogger.Error("services req expense err:%s", err.Error())
		tx.Rollback()
		e.ErrorOK(MsgServerErr)
	}

	// 创建成功删除验证码
	err = tx.Delete(&otp).Error
	if err != nil {
		log.GLogger.Error("services delete code err:%s", err.Error())
		tx.Rollback()
		e.ErrorOK(MsgServerErr)
	}

	tx.Commit()
	e.Correct(param)
}

// @Title 单条申请报销
// @Description 单条申请报销
// @Param	id	path	int	true	"报销id"
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /:id [get]
func (e *ExpenseController) ExpenseById() {
	eID, _ := e.GetInt(":id", 0)
	expense := new(oa.Expense)
	services.Slave().Debug().Preload("Employee").Preload("Employee.EmployeeBasic").Preload("ExpenseDetails").Preload("ExpenseDetails.ExpenseAccount").Take(expense, "id = ?", eID)
	//oID 查询 workflow
	workflow := new(oa.Workflow)
	services.Slave().Model(oa.Workflow{}).Where("workflow_definition_id = ? and entity_id = ?",
		services.GetFlowDefID(services.Expense), eID).Preload("Nodes").Preload("Nodes.User").
		Preload("Elements").First(workflow)
	if len(workflow.Nodes) != 4 {
		e.ErrorOK("工作流配置错误")
	}
	var resp struct {
		Info     *oa.Expense  `json:"info"`
		WorkFlow *oa.Workflow `json:"work_flow"`
	}
	resp.Info = expense
	resp.WorkFlow = workflow

	e.Correct(resp)
}

// @Title 审批申请报销
// @Description 审批申请报销
// @Param	body	body	forms.ReqApprovalExpense	true
// @Success 200 {string} "success"
// @Failure 500 server err
// @router / [put]
func (e *ExpenseController) ApprovalExpense() {
	param := new(forms.ReqApprovalExpense)
	err := json.Unmarshal(e.Ctx.Input.RequestBody, param)
	if err != nil {
		log.GLogger.Error("parse expense err:%s", err.Error())
		e.ErrorOK(MsgInvalidParam)
	}
	log.GLogger.Info("param :%+v", param)

	expense := new(oa.Expense)
	services.Slave().Debug().Preload("Employee").Take(expense, "id = ?", param.Id)
	log.GLogger.Info("expense:%+v", expense)
	//oID 查询 workflow
	workflow := new(oa.Workflow)
	services.Slave().Model(oa.Workflow{}).Where("workflow_definition_id = ? and entity_id = ?",
		services.GetFlowDefID(services.Expense), param.Id).Preload("Nodes").Preload("Nodes.User").
		Preload("Elements").First(workflow)
	if workflow.Nodes == nil || len(workflow.Nodes) != 4 {
		e.ErrorOK("工作流配置错误")
	}
	isCheck := false
	userID, _ := e.GetInt("userID", 0)

	log.GLogger.Info("userId: %d", userID)
	log.GLogger.Info("expense.Employee.Email:%s", expense.Employee.Email)

	// 负责人，hr审批
	for i, node := range workflow.Nodes {
		log.GLogger.Info("node.OperatorId:%d", node.OperatorID)
		if node.Status == models.FlowProcessing && node.OperatorID == userID {
			isCheck = true
			status := models.FlowRejected
			if param.Status == 1 {
				status = models.FlowApproved
			}
			node.Status = status
			workflow.Elements[i].Value = param.Comment
			if param.Status == 0 {
				workflow.Status = status
				services.Slave().Model(oa.Expense{}).Where("id = ?", param.Id).Updates(map[string]interface{}{
					"status": status,
				})
				if i == 2 {
					otp, _ := util.GenerateOTP(6)
					expenseOtp := oa.ExpenseOtp{
						Code:  otp,
						EmpID: int(expense.Employee.ID),
					}
					services.Slave().Create(&expenseOtp)
					go services.EmailExpenseRejectedUp(expense.Employee.Email, expense.Employee.Name, expense.ApplicationDate, otp)
				}
			} else {
				var nextNodeStatus string
				if i == 1 {
					nextNodeStatus = models.FlowProcessing
				} else if i == 2 {
					services.Slave().Model(oa.Expense{}).Where("id = ?", param.Id).Updates(map[string]interface{}{
						"status": models.FlowUnpaid,
					})
					nextNodeStatus = models.FlowUnpaid
					go services.EmailExpenseApproved(expense.Employee.Email, expense.ID, expense.Employee.Name, expense.ApplicationDate)
				}
				workflow.Nodes[i+1].Status = nextNodeStatus
			}
			services.Slave().Save(workflow)
			break
		}

	}
	if !isCheck {
		e.ErrorOK("您不是当前审批人")
	}
	e.Correct("")
}

// @Title 批量支付报销
// @Description 批量支付报销
// @Param	ids	path	string	true	"报销ids"
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /paid/batch [put]
func (e *ExpenseController) BatchPaidExpense() {
	ids := e.GetString("ids")
	var idList []string
	if ids == "" {
		userType, _ := e.GetInt("userType", 0)
		name := e.GetString("name")
		myReq, _ := e.GetBool("myreq", false)
		myTodo, _ := e.GetBool("mytodo", false)
		status := e.GetString("status")
		userEmail := e.GetString("userEmail")
		searchID := e.GetString("searchid")
		applicationDateBegin := e.GetString("application_date_begin")
		applicationDateEnd := e.GetString("application_date_end")

		employee := new(oa.Employee)
		services.Slave().Where("email = ?", userEmail).First(employee)
		log.GLogger.Info("employee: %+v", employee)

		expenses := make([]*oa.Expense, 0)
		query := services.Slave().Debug().Model(oa.Expense{}).Preload("ExpenseDetails")
		if searchID != "" {
			query = query.Where("id like ?", fmt.Sprintf("%%%s%%", searchID))
		}
		if name != "" {
			query = query.Where("e_name like ?", "%"+name+"%")
		}
		if status != "" {
			query = query.Where("status = ?", status)
		}
		if applicationDateBegin != "" && applicationDateEnd != "" {
			query = query.Where("application_date > ?", applicationDateBegin).Where("application_date <= ?", applicationDateEnd)
		}

		if myReq {
			query = query.Where("emp_id = ?", employee.ID)
		}

		eIDs := make([]int, 0)

		if myTodo {
			userID, _ := e.GetInt("userID", 0)
			log.GLogger.Info("userID：%d", userID)
			ids := make([]oa.EntityID, 0)
			var s []string
			if e.GetString("todostatus") != "" {
				if userType == models.UserFinance {
					s = oa.TodoStatusFinanceMap[e.GetString("todostatus")]
				} else {
					s = oa.TodoStatusLeaderMap[e.GetString("todostatus")]
				}
			}

			if len(s) == 0 {
				services.Slave().Debug().Raw("select w.entity_id from workflows w,workflow_nodes wn where w.id = "+
					"wn.workflow_id and w.workflow_definition_id = ? and operator_id = ? and wn.status <> ?"+
					" and wn.node_seq != 1 order by w.entity_id desc", services.GetFlowDefID(services.Expense), userID, models.FlowNA).Scan(&ids)
				// services.Slave().Debug().Raw("select w.entity_id from workflows w,workflow_nodes wn where w.id = "+
				// 	"wn.workflow_id and w.workflow_definition_id = ? and operator_id = ?"+
				// 	" and wn.node_seq != 1 order by w.entity_id desc", services.GetFlowDefID(services.Expense), userID).Scan(&ids)
			} else {
				services.Slave().Debug().Raw("select w.entity_id from workflows w,workflow_nodes wn where w.id = "+
					"wn.workflow_id and w.workflow_definition_id = ? and operator_id = ? and wn.status in (?)"+
					" and wn.node_seq != 1 order by w.entity_id desc", services.GetFlowDefID(services.Expense), userID, s).Scan(&ids)
			}

			for _, eID := range ids {
				eIDs = append(eIDs, eID.EntityID)
			}
			if len(eIDs) != 0 {
				query = query.Where(eIDs)
			}
		}

		query.Limit(-1).Offset(-1).Order("created_at desc").Find(&expenses)
		for _, expense := range expenses {
			if expense.Status == models.FlowUnpaid {
				idList = append(idList, strconv.Itoa(int(expense.ID)))
			}
		}
	} else {
		idList = strings.Split(ids, ",")
	}
	for _, id := range idList {
		expense := new(oa.Expense)
		services.Slave().Debug().Preload("Employee").Take(expense, "id = ?", id)
		log.GLogger.Info("expense:%+v", expense)
		log.GLogger.Info("expense.Employee.Email:%s", expense.Employee.Email)

		//oID 查询 workflow
		workflow := new(oa.Workflow)
		services.Slave().Model(oa.Workflow{}).Where("workflow_definition_id = ? and entity_id = ?",
			services.GetFlowDefID(services.Expense), id).Preload("Nodes").Preload("Nodes.User").
			Preload("Elements").First(workflow)
		if workflow.Nodes == nil || len(workflow.Nodes) != 4 {
			e.ErrorOK("工作流配置错误")
		}
		userID, _ := e.GetInt("userID", 0)

		if workflow.Nodes[len(workflow.Nodes)-1].OperatorID != userID {
			e.ErrorOK("您不是当前审批人")
		}

		var paymentDate *time.Time

		status := models.FlowPaid
		t := time.Now()
		paymentDate = &t
		paidCardInfo := e.GetDebitCard(expense.EmpID)
		go services.EmailExpensePaid(expense.Employee.Email, expense.Employee.Name, expense.ExpenseSummary, paidCardInfo.CardID, expense.ApplicationDate)

		workflow.Status = status
		services.Slave().Model(oa.Expense{}).Where("id = ?", id).Updates(map[string]interface{}{
			"status":       status,
			"payment_date": paymentDate,
		})

		workflow.Nodes[len(workflow.Nodes)-1].Status = status
		services.Slave().Save(workflow)
	}

	e.Correct("")
}

// @Title 支付报销
// @Description 支付报销
// @Param	body	body	forms.ReqApprovalExpense	true
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /paid [put]
func (e *ExpenseController) PaidExpense() {
	param := new(forms.ReqApprovalExpense)
	err := json.Unmarshal(e.Ctx.Input.RequestBody, param)
	if err != nil {
		log.GLogger.Error("parse expense err:%s", err.Error())
		e.ErrorOK(MsgInvalidParam)
	}
	log.GLogger.Info("param :%+v", param)

	expense := new(oa.Expense)
	services.Slave().Debug().Preload("Employee").Take(expense, "id = ?", param.Id)
	log.GLogger.Info("expense:%+v", expense)
	log.GLogger.Info("expense.Employee.Email:%s", expense.Employee.Email)

	//oID 查询 workflow
	workflow := new(oa.Workflow)
	services.Slave().Model(oa.Workflow{}).Where("workflow_definition_id = ? and entity_id = ?",
		services.GetFlowDefID(services.Expense), param.Id).Preload("Nodes").Preload("Nodes.User").
		Preload("Elements").First(workflow)
	if workflow.Nodes == nil || len(workflow.Nodes) != 4 {
		e.ErrorOK("工作流配置错误")
	}
	userID, _ := e.GetInt("userID", 0)

	if workflow.Nodes[len(workflow.Nodes)-1].OperatorID != userID {
		e.ErrorOK("您不是当前审批人")
	}

	var paymentDate *time.Time
	var status string
	if param.Status == 0 {
		status = models.FlowRejected
		paymentDate = nil
		otp, _ := util.GenerateOTP(6)
		expenseOtp := oa.ExpenseOtp{
			Code:  otp,
			EmpID: int(expense.Employee.ID),
		}
		services.Slave().Create(&expenseOtp)
		go services.EmailExpenseRejectedDown(expense.Employee.Email, expense.Employee.Name, expense.ApplicationDate, otp)
	} else {
		status = models.FlowPaid
		t := time.Now()
		paymentDate = &t
		paidCardInfo := e.GetDebitCard(expense.EmpID)
		go services.EmailExpensePaid(expense.Employee.Email, expense.Employee.Name, expense.ExpenseSummary, paidCardInfo.CardID, expense.ApplicationDate)
	}
	workflow.Status = status
	services.Slave().Model(oa.Expense{}).Where("id = ?", param.Id).Updates(map[string]interface{}{
		"status":       status,
		"payment_date": paymentDate,
	})
	workflow.Nodes[len(workflow.Nodes)-1].Status = status
	services.Slave().Save(workflow)

	e.Correct("")
}

// @Title 解析用户报销内容的excel文件
// @Description 解析用户报销内容的excel文件
// @Param  file form-data binary true "文件"
// @Success 200 {object} forms.ParseExpenseDetailResponse
// @Failure 500 server internal err
// @router /details [post]
func (e *ExpenseController) ParseDetailFile() {
	mf, mfh, err := e.GetFile("file")
	if err != nil {
		log.GLogger.Error("get file err: %s", err.Error())
		e.Error(err.Error())
		return
	}
	defer mf.Close()

	t := time.Now()
	validBeginMonthTime := util.GetFirstDateOfMonth(t)
	validEndMonthTime := util.GetLastDateOfMonth(t)
	log.GLogger.Info("validBeginMonthTime", validBeginMonthTime)
	otpCode := e.GetString("code")
	uID, _ := e.GetInt("userID")
	var validCode bool
	if otpCode != "" {
		var otp oa.ExpenseOtp
		services.Slave().Debug().Where("code = ?", otpCode).Where("emp_id=?", uID).
			Where("created_at < ?", validEndMonthTime).Where("created_at >= ?", validBeginMonthTime).Find(&otp)
		log.GLogger.Info("otp", otp)
		if otp.ID == 0 {
			e.ErrorOK("验证码错误")
		}
		validCode = true
	}

	fs := strings.Split(mfh.Filename, ".")
	ft := fs[len(fs)-1:][0]
	if ft != "xlsx" {
		e.ErrorOK("文件类型错误")
	}
	f, err := excelize.OpenReader(mf)
	if err != nil {
		fmt.Println(err)
		e.ErrorOK(err.Error())
	}
	res, err := Read(f, validCode)
	if err != nil {
		fmt.Println(err)
		e.ErrorOK(err.Error())
	}
	fileName := uuid.New().String() + ".xlsx"

	err = e.SaveToFile("file", "static/"+fileName)
	if err != nil {
		e.ErrorOK("文件保存失败")
	}
	data := forms.ParseExpenseDetailResponse{
		Details:  res,
		FileName: fileName,
	}

	e.Correct(data)
}

func Read(f *excelize.File, validCode bool) ([]*oa.ExpenseDetail, error) {
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return nil, err
	}

	log.GLogger.Info("row len; %d", len(rows))
	if len(rows) < 2 {
		return nil, errors.New("无数据")
	}
	if len(rows[0]) < 6 {
		return nil, errors.New("首行表头字段有误, 无法识别")
	}
	fmt.Println(len(rows[0]))

	for i, v := range rows[0][0:6] {
		if oa.ExcelHeaderArray[i] != v {
			return nil, errors.New("首行表头字段有误, 无法识别")
		}
	}
	var res []*oa.ExpenseDetail
	var errorArray []string
	for i, row := range rows[1:] {
		x := i + 2
		fmt.Println(row)
		var colList [6]string
		for i, colCell := range row {
			colList[i] = colCell
			fmt.Println(colList)
		}

		// 校验费用发生日期
		var ocurredDate models.Date

		if colList[0] == "" {
			errorArray = append(errorArray, fmt.Sprintf("第%d行费用发生日期未填写", x))
		} else {
			log.GLogger.Info("time: %s", colList[0])
			t, err := time.Parse(models.DateFormat, colList[0])
			if err != nil {
				errorArray = append(errorArray, fmt.Sprintf("第%d行费用发生日期格式不正确", x))
				continue
			}
			log.GLogger.Info("ocurredDate: %s", ocurredDate)
			nowTime := time.Now()
			validBeginMonthTime := util.GetFirstDateOfMonth(nowTime)
			log.GLogger.Info("day:%s", nowTime.Day())
			if nowTime.Day() < 11 || validCode {
				if !t.After(validBeginMonthTime.AddDate(0, -1, 0)) {
					errorArray = append(errorArray, fmt.Sprintf("第%d行费用发生日期需要上月内", x))
					continue
				}
			} else {
				if !t.After(validBeginMonthTime) {
					errorArray = append(errorArray, fmt.Sprintf("第%d行费用发生日期需在本月内", x))
					continue
				}
			}
			ocurredDate = models.Date(t)
		}

		// 校验费用科目
		var expenseAccountCode string
		var expenseAccount oa.ExpenseAccount
		if colList[1] == "" {
			errorArray = append(errorArray, fmt.Sprintf("第%d行费用科目未填写", x))
		} else {
			code, ok := oa.ExpenseAccountMap[colList[1]]
			if ok {
				expenseAccountCode = code
				expenseAccount.Code = code
				expenseAccount.ExpenseAccountName = colList[1]
			} else {
				errorArray = append(errorArray, fmt.Sprintf("第%d行费用科目不正确", x))
			}
		}

		// 校验费用金额
		var expenseAmount float64
		if colList[2] == "" {
			errorArray = append(errorArray, fmt.Sprintf("第%d行费用金额未填写", x))
		} else {
			log.GLogger.Info("expenseAmount string：%s", colList[2])
			float, err := strconv.ParseFloat(colList[2], 64)
			if err != nil || float <= 0 {
				errorArray = append(errorArray, fmt.Sprintf("第%d行费用金额格式不正确", x))
			}
			expenseAmount = float
			log.GLogger.Info("float：%v", expenseAmount)
		}

		// 校验备注
		vList := oa.ExpenseAccountValidMap[expenseAccount.ExpenseAccountName]
		for _, v := range vList {
			if colList[v] == "" {
				errorArray = append(errorArray, fmt.Sprintf("第%d行备注%d未填写", x, v-2))
			}
		}

		ed := &oa.ExpenseDetail{
			OcurredDate:        ocurredDate,
			ExpenseAccountCode: expenseAccountCode,
			ExpenseAccount:     &expenseAccount,
			ExpenseAmount:      expenseAmount,
			Remarks1:           colList[3],
			Remarks2:           colList[4],
			Remarks3:           colList[5],
		}
		res = append(res, ed)
	}

	if len(errorArray) > 0 {
		return nil, errors.New(strings.Join(errorArray, "-"))
	}

	return res, nil
}

// @Title engagement_list
// @Description 项目code
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /projects [get]
func (e *ExpenseController) GetProjects() {
	desc := e.GetString("desc")
	uEmail := e.GetString("userEmail")
	//查询部门下项目list
	projects := make([]*oa.EngagementCode, 0)
	//获取emp_info
	employee := new(oa.Employee)
	services.Slave().Preload("Department").Preload("Department.Leader").Take(employee, "email = ?", uEmail)
	if employee.ID == 0 {
		e.ErrorOK("未找到员工信息")
	}
	if employee.Department == nil {
		e.Correct(projects)
	}

	if employee.Department.ID == 0 {
		e.Correct(projects)
	}

	query := services.Slave().Model(oa.EngagementCode{}).Preload("Owner").Where("department_id = ?", employee.Department.ID)
	if desc != "" {
		query = query.Where("engagement_code_desc like ?", "%"+desc+"%")
	}
	query.Find(&projects)

	fmt.Println("employee.Department.Leader", employee.Department.Leader)
	for i := 0; i < len(projects); i++ {
		if projects[i].CodeOwnerID == int(employee.ID) {
			projects[i].Owner = employee.Department.Leader
			projects[i].CodeOwnerID = int(employee.Department.Leader.ID)
			if employee.Department.Leader.ID == employee.ID {
				user := new(models.User)
				services.Slave().Take(user, "id = ?", 2) // 马俊杰
				projects[i].Owner = user
				projects[i].CodeOwnerID = 2 // 马俊杰
			}
		}
	}

	e.Correct(projects)
}

// @Title 报销审批人
// @Description 报销审批人
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /approvals [get]
func (e *ExpenseController) ApprovalUsers() {
	uEmail := e.GetString("userEmail")
	log.GLogger.Info("userEmail:%s", uEmail)
	//获取emp_info
	employee := new(oa.Employee)
	services.Slave().Preload("Department").Preload("Department.Leader").Take(employee, "email = ?", uEmail)
	if employee.ID == 0 {
		e.ErrorOK("未找到员工信息")
	}
	log.GLogger.Info("employee:%s", employee)
	//查询财务 id
	engagementCode := new(oa.EngagementCode)
	services.Slave().Model(oa.EngagementCode{}).Where("department_id = ?", employee.DepartmentID).First(engagementCode)
	log.GLogger.Info("engagementCode:%s", engagementCode)
	u := new(models.User)
	services.Slave().Take(u, "id = ?", engagementCode.FinanceID)
	e.Correct(u.Name)
}

// @Title 支付信息统计
// @Description 支付信息统计
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /paid/info [get]
func (e *ExpenseController) PaidInfo() {
	userID := e.GetString("userID")

	type res struct {
		Sum float64
	}

	t := time.Now().AddDate(0, -1, 0)

	var ExpenseTotal res
	var LastPaidAmount res
	services.Slave().Debug().Raw("select sum(expense_summary) as sum from expenses where emp_id = ? and application_date >= ? and status <> ?;", userID, t, models.FlowRejected).Scan(&ExpenseTotal)
	services.Slave().Debug().Raw("select expenses.expense_summary as sum from expenses where emp_id = ? and status= ? order by created_at desc limit 1;", userID, models.FlowPaid).Scan(&LastPaidAmount)

	data := struct {
		ExpenseTotal   float64 `json:"expense_total"`
		LastPaidAmount float64 `json:"last_paid_amount"`
	}{
		ExpenseTotal:   ExpenseTotal.Sum,
		LastPaidAmount: LastPaidAmount.Sum,
	}
	e.Correct(data)
}

// @Title 员工银行卡
// @Description 员工银行卡
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /:id/debit_card [get]
func (e *ExpenseController) DebitCard() {
	eID, _ := e.GetInt(":id", 0)
	expense := new(oa.Expense)
	services.Slave().Debug().Take(expense, "id = ?", eID)
	if expense.EmpID != 0 {
		e.Correct(e.GetDebitCard(expense.EmpID))
	}
	e.Correct("")
}

// @Title 导出待支付
// @Description 导出待支付
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /export/unpaid [get]
func (e *ExpenseController) ExportUnpaid() {
	ids := e.GetString("ids")
	var idList []string
	if ids == "" {
		userType, _ := e.GetInt("userType", 0)
		name := e.GetString("name")
		myReq, _ := e.GetBool("myreq", false)
		myTodo, _ := e.GetBool("mytodo", false)
		status := e.GetString("status")
		userEmail := e.GetString("userEmail")
		searchID := e.GetString("searchid")
		applicationDateBegin := e.GetString("application_date_begin")
		applicationDateEnd := e.GetString("application_date_end")

		employee := new(oa.Employee)
		services.Slave().Where("email = ?", userEmail).First(employee)
		log.GLogger.Info("employee: %+v", employee)

		expenses := make([]*oa.Expense, 0)
		query := services.Slave().Debug().Model(oa.Expense{}).Preload("ExpenseDetails")
		if searchID != "" {
			query = query.Where("id like ?", fmt.Sprintf("%%%s%%", searchID))
		}
		if name != "" {
			query = query.Where("e_name like ?", "%"+name+"%")
		}
		if status != "" {
			query = query.Where("status = ?", status)
		}
		if applicationDateBegin != "" && applicationDateEnd != "" {
			query = query.Where("application_date > ?", applicationDateBegin).Where("application_date <= ?", applicationDateEnd)
		}

		if myReq {
			query = query.Where("emp_id = ?", employee.ID)
		}

		eIDs := make([]int, 0)

		if myTodo {
			userID, _ := e.GetInt("userID", 0)
			log.GLogger.Info("userID：%d", userID)
			ids := make([]oa.EntityID, 0)
			var s []string
			if e.GetString("todostatus") != "" {
				if userType == models.UserFinance {
					s = oa.TodoStatusFinanceMap[e.GetString("todostatus")]
				} else {
					s = oa.TodoStatusLeaderMap[e.GetString("todostatus")]
				}
			}

			if len(s) == 0 {
				services.Slave().Debug().Raw("select w.entity_id from workflows w,workflow_nodes wn where w.id = "+
					"wn.workflow_id and w.workflow_definition_id = ? and operator_id = ? and wn.status <> ?"+
					" and wn.node_seq != 1 order by w.entity_id desc", services.GetFlowDefID(services.Expense), userID, models.FlowNA).Scan(&ids)
				// services.Slave().Debug().Raw("select w.entity_id from workflows w,workflow_nodes wn where w.id = "+
				// 	"wn.workflow_id and w.workflow_definition_id = ? and operator_id = ?"+
				// 	" and wn.node_seq != 1 order by w.entity_id desc", services.GetFlowDefID(services.Expense), userID).Scan(&ids)
			} else {
				services.Slave().Debug().Raw("select w.entity_id from workflows w,workflow_nodes wn where w.id = "+
					"wn.workflow_id and w.workflow_definition_id = ? and operator_id = ? and wn.status in (?)"+
					" and wn.node_seq != 1 order by w.entity_id desc", services.GetFlowDefID(services.Expense), userID, s).Scan(&ids)
			}

			for _, eID := range ids {
				eIDs = append(eIDs, eID.EntityID)
			}
			if len(eIDs) != 0 {
				query = query.Where(eIDs)
			}
		}

		query.Limit(-1).Offset(-1).Order("created_at desc").Find(&expenses)
		for _, expense := range expenses {
			if expense.Status == models.FlowUnpaid {
				idList = append(idList, strconv.Itoa(int(expense.ID)))
			}
		}
	} else {
		idList = strings.Split(ids, ",")
	}
	log.GLogger.Info("ids: %s", idList)

	var expenses []*oa.Expense
	services.Slave().Where(idList).Find(&expenses)

	log.GLogger.Info("expenses: %s", expenses)

	f := excelize.NewFile()

	f.NewSheet("上海游因")
	f.NewSheet("宁波比孚")
	f.NewSheet("上海品埃")
	f.DeleteSheet("Sheet1")

	// 设置单元格宽度
	f.SetColWidth("上海游因", "A", "A", 20)

	f.SetColWidth("宁波比孚", "A", "A", 6)
	f.SetColWidth("宁波比孚", "B", "D", 10)
	f.SetColWidth("宁波比孚", "B", "D", 10)
	f.SetColWidth("宁波比孚", "E", "E", 17)
	f.SetColWidth("宁波比孚", "F", "F", 20)
	f.SetColWidth("宁波比孚", "G", "G", 25)
	f.SetColWidth("宁波比孚", "H", "K", 17)
	f.SetColWidth("宁波比孚", "L", "L", 20)
	f.SetColWidth("宁波比孚", "M", "M", 15)
	f.SetColWidth("宁波比孚", "M", "M", 15)
	f.SetColWidth("宁波比孚", "N", "Q", 10)
	f.SetColWidth("宁波比孚", "R", "R", 27)
	f.SetColWidth("宁波比孚", "S", "S", 10)

	f.SetColWidth("上海品埃", "A", "A", 6)
	f.SetColWidth("上海品埃", "B", "D", 10)
	f.SetColWidth("上海品埃", "E", "E", 17)
	f.SetColWidth("上海品埃", "F", "F", 20)
	f.SetColWidth("上海品埃", "G", "G", 25)
	f.SetColWidth("上海品埃", "H", "K", 17)
	f.SetColWidth("上海品埃", "L", "L", 20)
	f.SetColWidth("上海品埃", "M", "M", 15)
	f.SetColWidth("上海品埃", "M", "M", 15)
	f.SetColWidth("上海品埃", "N", "Q", 10)
	f.SetColWidth("上海品埃", "R", "R", 27)
	f.SetColWidth("上海品埃", "S", "S", 10)

	_ = f.SetSheetRow("上海游因", "A1", &[]interface{}{"收款帐号", "收款户名", "金额", "开户行", "开户地"})
	_ = f.SetSheetRow("宁波比孚", "A1", &[]interface{}{"币种", "日期", "明细标志", "顺序号", "付款账号开户行",
		"付款账号/卡号", "付款账号名称/卡名称", "收款账号开户行", "收款账号省份", "收款账号地市", "收款账号地区码", "收款账号",
		"收款账号名称", "金额", "汇款用途", "备注信息", "汇款方式", "收款账户短信通知手机号码", "自定义序号"})
	_ = f.SetSheetRow("上海品埃", "A1", &[]interface{}{"币种", "日期", "明细标志", "顺序号", "付款账号开户行",
		"付款账号/卡号", "付款账号名称/卡名称", "收款账号开户行", "收款账号省份", "收款账号地市", "收款账号地区码", "收款账号",
		"收款账号名称", "金额", "汇款用途", "备注信息", "汇款方式", "收款账户短信通知手机号码", "自定义序号"})
	num1 := 2
	num2 := 2
	num3 := 2

	for _, expense := range expenses {
		log.GLogger.Info("empid:%d", expense.EmpID)
		paidCardInfo := e.GetDebitCard(expense.EmpID)
		log.GLogger.Info("paidCardInfo:%s", paidCardInfo)
		issuingBank := paidCardInfo.IssuingBank
		log.GLogger.Info("issuingBank:%s", issuingBank)
		var province, city string
		issuingBankSplit := strings.SplitN(issuingBank, "省", 2)
		if len(issuingBankSplit) == 2 {
			province = issuingBankSplit[0]
			city = issuingBankSplit[1]
		} else {
			city = issuingBank
		}
		if paidCardInfo.PaymentName == "上海游因" {
			err := f.SetSheetRow("上海游因", "A"+strconv.Itoa(num1), &[]interface{}{
				paidCardInfo.CardID, expense.EName, expense.ExpenseSummary, paidCardInfo.BankName, paidCardInfo.BankName,
			})
			log.GLogger.Info("err", err)
			num1++
		} else if paidCardInfo.PaymentName == "宁波比孚" {
			_ = f.SetSheetRow("宁波比孚", "A"+strconv.Itoa(num2), &[]interface{}{
				"RMB" /*币种*/, "" /*日期*/, "" /*明细标志*/, num2 - 1 /*顺序号*/, "工行", /*付款账号开户行*/
				"1001100419005023362" /*付款账号/卡号*/, "宁波比孚信息科技有限公司" /*付款账号名称/卡名称*/, "工行" /*收款账号开户行*/, province, /*收款账号省份*/
				city /*收款账号地市*/, "" /*收款账号地区码*/, paidCardInfo.CardID /*收款账号*/, expense.EName, /*收款账号名称*/
				expense.ExpenseSummary /*金额*/, "报销" /*汇款用途*/, "" /*备注信息*/, "1", /*汇款方式*/
				"" /*收款账户短信通知手机号码*/, "", /*自定义序号*/
			})

			num2++
		} else if paidCardInfo.PaymentName == "上海品埃" {
			_ = f.SetSheetRow("上海品埃", "A"+strconv.Itoa(num3), &[]interface{}{
				"RMB" /*币种*/, "" /*日期*/, "" /*明细标志*/, num3 - 1 /*顺序号*/, "工行", /*付款账号开户行*/
				"1001100409100025962" /*付款账号/卡号*/, "上海品埃信息科技有限公司" /*付款账号名称/卡名称*/, "工行" /*收款账号开户行*/, province, /*收款账号省份*/
				city /*收款账号地市*/, "" /*收款账号地区码*/, paidCardInfo.CardID /*收款账号*/, expense.EName, /*收款账号名称*/
				expense.ExpenseSummary /*金额*/, "报销" /*汇款用途*/, "" /*备注信息*/, "1", /*汇款方式*/
				"" /*收款账户短信通知手机号码*/, "", /*自定义序号*/
			})
			num3++
		}
	}
	f.SetActiveSheet(0)
	f.SaveAs("static/expense.xlsx")
	e.Ctx.Output.Download("static/expense.xlsx", "expense.xlsx")
	os.Remove("static/expense.xlsx")
}

func (e *ExpenseController) GetDebitCard(employeeID int) forms.PaidCardInfo {
	employee := new(oa.Employee)
	services.Slave().Debug().Preload("EmployeeBasic").Take(employee, "id = ?", employeeID)
	log.GLogger.Info("employee:%+v", employee)

	employeeContract := new(oa.EmployeeContract)
	services.Slave().Debug().Where("employee_id = ?", employee.ID).Order("contract_start_date desc").First(employeeContract)
	log.GLogger.Info("employeeContract:%+v", employeeContract)

	var paidCardInfo forms.PaidCardInfo
	if employee.EmployeeBasic == nil {
		return paidCardInfo
	}

	if employeeContract.ContractMain == "上海游因" {
		paidCardInfo.BankName = "招商银行"
		paidCardInfo.CardID = employee.EmployeeBasic.DebitCard2
		paidCardInfo.PaymentName = employeeContract.ContractMain
		paidCardInfo.IssuingBank = employee.EmployeeBasic.IssuingBank2
		return paidCardInfo
	} else if employeeContract.ContractMain == "宁波比孚" {
		paidCardInfo.BankName = "工行"
		paidCardInfo.CardID = employee.EmployeeBasic.DebitCard1
		paidCardInfo.PaymentName = employeeContract.ContractMain
		paidCardInfo.IssuingBank = employee.EmployeeBasic.IssuingBank1
		return paidCardInfo
	} else if employeeContract.ContractMain == "上海品埃" {
		paidCardInfo.BankName = "工行"
		paidCardInfo.CardID = employee.EmployeeBasic.DebitCard1
		paidCardInfo.PaymentName = employeeContract.ContractMain
		paidCardInfo.IssuingBank = employee.EmployeeBasic.IssuingBank1
		return paidCardInfo
	} else {
		return paidCardInfo
	}
}
