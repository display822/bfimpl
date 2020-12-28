package controllers

import (
	"bfimpl/models/oa"
	"bfimpl/services"
	"bfimpl/services/log"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

type ProjectController struct {
	BaseController
}

// @Title 项目交付列表
// @Description 项目交付列表
// @Success 200 {string} ""
// @Failure 500 server internal err
// @router /details [post]
func (p *ProjectController) List() {

	p.Correct("")
}

// @Title 创建项目交付
// @Description 创建项目交付
// @Success 200 {string} ""
// @Failure 500 server internal err
// @router /details [post]
func (p *ProjectController) Create() {
	var ps []*oa.ProjectDelivery
	err := json.Unmarshal(p.Ctx.Input.RequestBody, &ps)
	if err != nil {
		log.GLogger.Error("new employee err：%s", err.Error())
		p.ErrorOK(MsgInvalidParam)
	}

	log.GLogger.Info("ps", ps)
	err = oa.BatchProjectCreate(services.Slave(), ps)
	if err != nil {
		p.ErrorOK(MsgServerErr)
	}
	p.Correct("")
}

// @Title 解析项目交付内容的excel文件
// @Description 解析项目交付内容的excel文件
// @Param  file form-data binary true "文件"
// @Success 200 {object} []oa.ProjectDelivery
// @Failure 500 server internal err
// @router /details [post]
func (p *ProjectController) ParseDetailFile() {
	mf, mfh, err := p.GetFile("file")
	if err != nil {
		log.GLogger.Error("get file err: %s", err.Error())
		p.Error(err.Error())
		return
	}
	defer mf.Close()

	fs := strings.Split(mfh.Filename, ".")
	ft := fs[len(fs)-1:][0]
	if ft != "xlsx" {
		p.ErrorOK("文件类型错误")
	}
	f, err := excelize.OpenReader(mf)
	if err != nil {
		fmt.Println(err)
		p.ErrorOK(err.Error())
	}
	res, err := ReadProjectFile(f)
	if err != nil {
		fmt.Println(err)
		p.ErrorOK(err.Error())
	}

	p.Correct(res)
}

func ReadProjectFile(f *excelize.File) ([]*oa.ProjectDelivery, error) {
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return nil, err
	}

	log.GLogger.Info("row len; %d", len(rows))
	if len(rows) < 2 {
		return nil, errors.New("无数据")
	}
	if len(rows[0]) < 3 {
		return nil, errors.New("首行表头字段有误, 无法识别")
	}
	fmt.Println(len(rows[0]))

	for i, v := range rows[0][0:3] {
		if oa.ProjectExcelHeaderArray[i] != v {
			return nil, errors.New("首行表头字段有误, 无法识别")
		}
	}
	var res []*oa.ProjectDelivery
	var errorArray []string
	for _, row := range rows[1:] {
		//x := i + 2
		fmt.Println(row)
		var colList [6]string
		for i, colCell := range row {
			colList[i] = colCell
			fmt.Println(colList)
		}

		// 校验数量
		//var ocurredDate models.Date
		//if colList[0] == "" {
		//	errorArray = append(errorArray, fmt.Sprintf("第%d行费用发生日期未填写", x))
		//} else {
		//	log.GLogger.Info("time: %s", colList[0])
		//	t, err := time.Parse(models.DateFormat, colList[0])
		//	if err != nil {
		//		errorArray = append(errorArray, fmt.Sprintf("第%d行费用发生日期格式不正确", x))
		//	}
		//	log.GLogger.Info("ocurredDate: %s", ocurredDate)
		//	ocurredDate = models.Date(t)
		//}

		// 校验费用科目
		//var expenseAccountCode string
		//var expenseAccount oa.ExpenseAccount
		//if colList[1] == "" {
		//	errorArray = append(errorArray, fmt.Sprintf("第%d行费用科目未填写", x))
		//} else {
		//	code, ok := oa.ExpenseAccountMap[colList[1]]
		//	if ok {
		//		expenseAccountCode = code
		//		expenseAccount.Code = code
		//		expenseAccount.ExpenseAccountName = colList[1]
		//	} else {
		//		errorArray = append(errorArray, fmt.Sprintf("第%d行费用科目不正确", x))
		//	}
		//}

		//// 校验费用金额
		//var expenseAmount float64
		//if colList[2] == "" {
		//	errorArray = append(errorArray, fmt.Sprintf("第%d行费用金额未填写", x))
		//} else {
		//	log.GLogger.Info("expenseAmount string：%s", colList[2])
		//	float, err := strconv.ParseFloat(colList[2], 64)
		//	if err != nil || float <= 0 {
		//		errorArray = append(errorArray, fmt.Sprintf("第%d行费用金额格式不正确", x))
		//	}
		//	expenseAmount = float
		//	log.GLogger.Info("float：%v", expenseAmount)
		//}

		//// 校验备注
		//vList := oa.ExpenseAccountValidMap[expenseAccount.ExpenseAccountName]
		//for _, v := range vList {
		//	if colList[v] == "" {
		//		errorArray = append(errorArray, fmt.Sprintf("第%d行备注%d未填写", x, v))
		//	}
		//}

		ed := &oa.ProjectDelivery{
			ProjectCategoryCode: colList[1],
			StartDate:           nil,
			EndDate:             nil,
			MainServiceAmount:   1,
			SubServiceAmount:    2,
		}
		res = append(res, ed)
	}

	if len(errorArray) > 0 {
		return nil, errors.New(strings.Join(errorArray, "-"))
	}

	return res, nil
}
