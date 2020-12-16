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
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

type ExpenseController struct {
	BaseController
}

// @Title 报销列表
// @Description 报销列表
// @Param	pagenum	    query	int	false	"页码"
// @Param	pagesize	query	int	false	"页数"
// @Param	myreq	query	bool	false	"我的报销"
// @Param	status	query	int	false	"状态"
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

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if applicationDateBegin != "" && applicationDateEnd != "" {
		query = query.Where("application_date > ?", applicationDateBegin).Where("application_date <= ?", applicationDateEnd)
	}
	if userType != models.UserHR && userType != models.UserLeader {
		//不是hr和部门负责人，只能查自己
		query = query.Where("emp_id = ?", employee.ID)
	} else {
		if name != "" {
			query = query.Where("e_name like ?", "%"+name+"%")
		}
		if myReq {
			// 查自己
			query = query.Where("emp_id = ?", employee.ID)
		}
	}
	var resp struct {
		Total int           `json:"total"`
		List  []*oa.Expense `json:"list"`
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
	log.GLogger.Info("ReqExpense query: %d, %s", uID, uEmail)
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

	param.EmpID = int(employee.ID)
	param.EName = employee.Name
	param.ApplicationDate = time.Now()
	param.Status = models.FlowNA
	// 费用总金额
	var expenseSummary float64

	for _, item := range param.ExpenseDetails {
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

	//log.GLogger.Info("param.ExpenseDetails: %v", param.ExpenseDetails)
	//// 批量创建报销详情
	//err = oa.BatchCreateExpenseDetail(tx, int(param.ID), param.ExpenseDetails)
	//if err != nil {
	//	log.GLogger.Error("batch create expense detail err:%s", err.Error())
	//	tx.Rollback()
	//	e.ErrorOK(MsgServerErr)
	//}

	// 执行报销工作流
	err = services.ReqExpense(tx, int(param.ID), uID, param.LeaderId, engagementCode.FinanceID)
	if err != nil {
		log.GLogger.Error("services req expense err:%s", err.Error())
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
// @router /expense/:id [get]
func (e *ExpenseController) ExpenseById() {
	eID, _ := e.GetInt(":id", 0)
	expense := new(oa.Expense)
	services.Slave().Preload("ExpenseDetails").Take(expense, "id = ?", eID)
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
// @Param	id	body	int	true	"报销id"
// @Param	comment	body	string	true	"审批意见"
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /expense [put]
func (e *ExpenseController) ApprovalExpense() {
	param := new(forms.ReqApprovalExpense)
	err := json.Unmarshal(e.Ctx.Input.RequestBody, param)
	if err != nil {
		log.GLogger.Error("parse expense err:%s", err.Error())
		e.ErrorOK(MsgInvalidParam)
	}
	//oID 查询 workflow
	workflow := new(oa.Workflow)
	services.Slave().Model(oa.Workflow{}).Where("workflow_definition_id = ? and entity_id = ?",
		services.GetFlowDefID(services.Expense), param.Id).Preload("Nodes").Preload("Nodes.User").
		Preload("Elements").First(workflow)
	if workflow.Nodes == nil || len(workflow.Nodes) != 4 {
		e.ErrorOK("工作流配置错误")
	}
	//isCheck := false
	//userID, _ := e.GetInt("userID", 0)
	// 负责人，hr审批
	//num := len(workflow.Nodes)
	//for i, node := range workflow.Nodes {
	//	if node.Status == models.FlowProcessing && node.OperatorID == userID {
	//		isCheck = true
	//		status := models.FlowRejected
	//		if param.Status == 1 {
	//			status = models.FlowApproved
	//		}
	//		node.Status = status
	//		workflow.Elements[i].Value = param.Comment
	//		if i < num-1 {
	//			//负责人
	//			if param.Status == 0 {
	//				workflow.Status = status
	//				services.Slave().Model(oa.Overtime{}).Where("id = ?", param.Id).Updates(map[string]interface{}{
	//					"status": status,
	//				})
	//			} else {
	//				workflow.Nodes[i+1].Status = models.FlowProcessing
	//			}
	//			services.Slave().Save(workflow)
	//		} else if i == num-1 {
	//			//hr
	//			workflow.Status = status
	//			services.Slave().Save(workflow)
	//			services.Slave().Model(oa.Overtime{}).Where("id = ?", param.Id).Updates(map[string]interface{}{
	//				"status": status,
	//			})
	//			//审批通过且类型为weekend,holiday，将加班时长放入leave balance
	//			if status == models.FlowApproved {
	//				overtime := new(oa.Overtime)
	//				services.Slave().Take(overtime, "id = ?", param.Id)
	//				if overtime.Type == "weekend" || overtime.Type == "holiday" {
	//					balance := oa.LeaveBalance{
	//						EmpID:      overtime.EmpID,
	//						Type:       oa.OverTime,
	//						Amount:     (overtime.RealDuration) / 8,
	//						OvertimeID: int(overtime.ID),
	//					}
	//					if balance.Amount == 0 {
	//						balance.Amount = (overtime.Duration) / 8
	//					}
	//					services.Slave().Create(&balance)
	//				}
	//			}
	//		}
	//		break
	//	}
	//}
	//if !isCheck {
	//	e.ErrorOK("您不是当前审批人")
	//}
	e.Correct("")
}

// @Title 支付报销
// @Description 支付报销
// @Param	id	body	int	true	"报销id"
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /expense/paid [put]
func (e *ExpenseController) PaidExpense() {

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
	fmt.Println(mf, mfh)
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
	res, err := Read(f)
	if err != nil {
		fmt.Println(err)
		e.ErrorOK(err.Error())
	}
	data := forms.ParseExpenseDetailResponse{
		Details:  res,
		FileName: mfh.Filename,
	}

	e.Correct(data)
}

func Read(f *excelize.File) ([]*oa.ExpenseDetail, error) {
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
			t, err := time.Parse(models.DateFormat, colList[0])
			if err != nil {
				errorArray = append(errorArray, fmt.Sprintf("第%d行费用发生日期格式不正确", x))
			}
			log.GLogger.Info("ocurredDate: %s", ocurredDate)
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
				expenseAccount.ExpenseAccountCode = code
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
	//获取emp_info
	employee := new(oa.Employee)
	services.Slave().Preload("Department").Preload("Department.Leader").Take(employee, "email = ?", uEmail)
	if employee.ID == 0 {
		e.ErrorOK("未找到员工信息")
	}
	//查询部门下项目list
	projects := make([]*oa.EngagementCode, 0)
	query := services.Slave().Model(oa.EngagementCode{}).Preload("Owner").Where("department_id = ?", employee.Department.ID)
	if desc != "" {
		query = query.Where("engagement_code_desc like ?", "%"+desc+"%")
	}
	query.Find(&projects)
	// for _, p :=range projects{
	// 	if p.CodeOwnerID == int(employee.ID){
	// 		p.Owner = employee.Department.Leader
	// 	}
	// }
	fmt.Println("employee.Department.Leader", employee.Department.Leader)
	for i := 0; i < len(projects); i++ {
		fmt.Println("CodeOwnerID", projects[i].CodeOwnerID)
		fmt.Println("int(employee.ID)", int(employee.ID))
		if projects[i].CodeOwnerID == int(employee.ID) {
			projects[i].Owner = employee.Department.Leader
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
