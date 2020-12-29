package controllers

import (
	"bfimpl/models/oa"
	"bfimpl/services"
	"bfimpl/services/log"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/jinzhu/gorm"
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
	pageSize, _ := p.GetInt("pagesize", 10)
	pageNum, _ := p.GetInt("pagenum", 1)
	periodTime := p.GetString("period_time")
	projectName := p.GetString("project_name")

	var resp struct {
		Total              int                   `json:"total"`
		TotalDeliveryValue float64               `json:"total_delivery_value"`
		List               []*oa.ProjectDelivery `json:"list"`
	}

	db := services.Slave()
	if periodTime != "" {
		db = db.Where("period_time = ?", periodTime)
	}
	if projectName != "" {
		db = db.Where("project_name=?", projectName)
	}
	var totalDeliveryValue struct {
		N float64
	}

	db.Table("project_deliveries").Select("sum(project_delivery_value) as n").Scan(&totalDeliveryValue)
	resp.TotalDeliveryValue = totalDeliveryValue.N

	db.Limit(pageSize).Offset((pageNum - 1) * pageSize).Order("created_at desc").Find(&resp.List).Limit(-1).Offset(-1).Count(&resp.Total)

	p.Correct(resp)
}

// @Title 项目交付过滤字段
// @Description 项目交付过滤字段
// @Success 200 {string} ""
// @Failure 500 server internal err
// @router /filter [get]
func (p *ProjectController) FilterField() {
	var resp struct {
		ProjectName []string `json:"project_name"`
		PeriodTime  []string `json:"period_time"`
	}
	var projectName []*oa.ProjectDelivery
	services.Slave().Group("project_name").Order("created_at desc").Find(&projectName)
	var periodTime []*oa.ProjectDelivery
	services.Slave().Group("period_time").Order("created_at desc").Find(&periodTime)

	for _, item := range projectName {
		resp.ProjectName = append(resp.ProjectName, item.ProjectName)
	}
	for _, item := range periodTime {
		resp.PeriodTime = append(resp.PeriodTime, item.PeriodTime)
	}

	p.Correct(resp)
}

// @Title 创建项目交付
// @Description 创建项目交付
// @Success 200 {string} ""
// @Failure 500 server internal err
// @router / [post]
func (p *ProjectController) Create() {
	var ps []*oa.ProjectDelivery
	err := json.Unmarshal(p.Ctx.Input.RequestBody, &ps)
	if err != nil {
		log.GLogger.Error("new ProjectDelivery err：%s", err.Error())
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
func (p *ProjectController) ParseProjectDetailFile() {
	periodTime := p.GetString("period_time")
	if periodTime == "" {
		p.ErrorOK("need period_time")
	}
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
	res, err := ReadProjectFile(f, periodTime)
	if err != nil {
		fmt.Println(err)
		p.ErrorOK(err.Error())
	}

	p.Correct(res)
}

func ReadProjectFile(f *excelize.File, periodTime string) ([]*oa.ProjectDelivery, error) {
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
	for i, row := range rows[1:] {
		x := i + 2
		fmt.Println(row)
		var colList [4]string
		for i, colCell := range row {
			colList[i] = colCell
			fmt.Println(colList)
		}

		// 校验项目名称
		if colList[0] == "" {
			errorArray = append(errorArray, fmt.Sprintf("第%d行项目名称未填写", x))
		}

		// 校验项目编码
		var projectCategory oa.ProjectCategory
		if colList[1] == "" {
			errorArray = append(errorArray, fmt.Sprintf("第%d行项目编码未填写", x))
		} else {
			err = services.Slave().Where("project_category_code=?", colList[1]).First(&projectCategory).Error
			if err != nil || err == gorm.ErrRecordNotFound {
				errorArray = append(errorArray, fmt.Sprintf("第%d行项目编码未找到", x))
			}
		}

		// 校验主服务交付数量
		var mainServiceAmount int
		if colList[2] == "" {
			errorArray = append(errorArray, fmt.Sprintf("第%d行主服务交付数量未填写", x))
		} else {
			log.GLogger.Info("mainServiceAmount string：%s", colList[2])
			mainServiceAmount, err = strconv.Atoi(colList[2])
			if err != nil {
				errorArray = append(errorArray, fmt.Sprintf("第%d行主服务交付数量不正确", x))
			}
			log.GLogger.Info("mainServiceAmount：%v", mainServiceAmount)
		}

		// 校验子服务交付数量
		var subServiceAmount int
		if colList[3] == "" {
			errorArray = append(errorArray, fmt.Sprintf("第%d行子服务交付数量未填写", x))
		} else {
			log.GLogger.Info("subServiceAmount string：%s", colList[3])
			subServiceAmount, err = strconv.Atoi(colList[3])
			if err != nil {
				errorArray = append(errorArray, fmt.Sprintf("第%d行子服务交付数量不正确", x))
			}
			log.GLogger.Info("subServiceAmount：%v", subServiceAmount)
		}

		projectDeliveryValue := projectCategory.MainServiceQuotation*float64(mainServiceAmount) +
			projectCategory.SubServiceQuotation*float64(subServiceAmount)
		ed := &oa.ProjectDelivery{
			ProjectName:          colList[0],
			ProjectCategoryCode:  colList[1],
			PeriodTime:           periodTime,
			MainServiceAmount:    mainServiceAmount,
			SubServiceAmount:     subServiceAmount,
			ProjectDeliveryValue: projectDeliveryValue,
		}
		res = append(res, ed)
	}

	if len(errorArray) > 0 {
		return nil, errors.New(strings.Join(errorArray, "-"))
	}

	return res, nil
}
