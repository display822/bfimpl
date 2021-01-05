/**
* @author : yi.zhang
* @description : controllers 描述
* @date   : 2020-12-29 17:35
 */

package controllers

import (
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

type EngagementController struct {
	BaseController
}

// @Title 人员管理列表
// @Description 人员管理列表
// @Success 200 {string} ""
// @Failure 500 server internal err
// @router / [post]
func (e *EngagementController) List() {
	engagementCodes := e.GetString("engagement_codes")
	beginTime := e.GetString("begin_time")
	endTime := e.GetString("end_time")
	ecs := strings.Split(engagementCodes, ",")
	fmt.Println(ecs)
	fmt.Println(beginTime)
	fmt.Println(endTime)
	var es []oa.Engagement
	services.Slave().Where("engagement_date>=?", beginTime).Where("engagement_date<=?", endTime).
		Where("engagement_code in (?)", ecs).Find(&es)

	m := make(map[string]int)
	m["2020-09-12"] = 1
	res := forms.Engagement{
		EngagementCode: "10001",
		EmployeeName:   "ss",
		DateField:      m,
	}
	e.Correct(res)
}

// @Title 创建人员管理
// @Description 创建人员管理
// @Success 200 {string} ""
// @Failure 500 server internal err
// @router / [post]
func (e *EngagementController) Create() {
	var es []*oa.Engagement
	err := json.Unmarshal(e.Ctx.Input.RequestBody, &es)
	if err != nil {
		log.GLogger.Error("new Engagement err：%s", err.Error())
		e.ErrorOK(MsgInvalidParam)
	}

	log.GLogger.Info("es", es)
	err = oa.BatchEngagementCreate(services.Slave(), es)
	if err != nil {
		e.ErrorOK(MsgServerErr)
	}
	e.Correct("")
}

// @Title 解析人员管理内容的excel文件
// @Description 解析人员管理内容的excel文件
// @Param  file form-data binary true "文件"
// @Success 200 {object}
// @Failure 500 server internal err
// @router /details [post]
func (e *EngagementController) ParseEngagementDetailFile() {
	mf, mfh, err := e.GetFile("file")
	if err != nil {
		log.GLogger.Error("get file err: %s", err.Error())
		e.Error(err.Error())
		return
	}
	defer mf.Close()

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
	res, err := EngagementDetailFile(f)
	if err != nil {
		fmt.Println(err)
		e.ErrorOK(err.Error())
	}

	m := make(map[string]map[string]int)
	for _, item := range res {
		_, ok := m[item.EngagementCode+"-"+item.EmployeeName]
		if !ok {
			mmm := make(map[string]int)
			mmm[item.EngagementDate.Format("2006/01/02")] = item.EngagementHour
			m[item.EngagementCode+"-"+item.EmployeeName] = mmm
		} else {
			m[item.EngagementCode+"-"+item.EmployeeName][item.EngagementDate.Format("2006/01/02")] = item.EngagementHour
		}
	}

	var data []forms.Engagement
	for k, v := range m {
		l := strings.Split(k, "-")
		eng := forms.Engagement{
			EmployeeName:   l[1],
			EngagementCode: l[0],
			DateField:      v,
		}
		data = append(data, eng)
	}

	e.Correct(data)
}

func EngagementDetailFile(f *excelize.File) ([]*oa.Engagement, error) {
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return nil, err
	}

	log.GLogger.Info("row len; %d", len(rows))
	if len(rows) < 2 {
		return nil, errors.New("无数据")
	}
	fmt.Println(len(rows[0]))

	if len(rows[0]) < 9 {
		return nil, errors.New("首行表头字段有误, 无法识别")
	}

	projectCodeIndex := rows[0][0]
	if projectCodeIndex != "项目编号" {
		return nil, errors.New("首行表头字段有误, 无法识别")
	}
	empNameIndex := rows[0][1]
	if empNameIndex != "员工" {
		return nil, errors.New("首行表头字段有误, 无法识别")
	}
	workTimeStringIndex1 := rows[0][2]

	log.GLogger.Info("workTimeStringIndex1", workTimeStringIndex1)
	workTimeIndex1, err := time.Parse("01-02-06", workTimeStringIndex1)
	if err != nil {
		return nil, errors.New("时间无法识别")
	}
	if workTimeIndex1.Weekday() != time.Monday {
		return nil, errors.New("需要从周一起")
	}

	for _, v := range rows[0][3:9] {
		fmt.Println(v)
		index, err := time.Parse("01-02-06", v)
		if err != nil {
			return nil, errors.New("时间无法识别")
		}
		workTimeIndex1 = workTimeIndex1.AddDate(0, 0, 1)
		if !index.Equal(workTimeIndex1) {
			return nil, errors.New("时间不连续")
		}

		fmt.Println(index)
	}

	// TODO 判断是否重复

	log.GLogger.Info("workTimeIndex1", workTimeIndex1)

	var phs []*oa.PublicHoliday
	services.Slave().Where("public_holiday_date >= ?", rows[0][3]).
		Where("public_holiday_date <= ?", rows[0][9]).
		Find(&phs)

	publicHolidayMap := make(map[string]string, len(phs))
	for _, ph := range phs {
		publicHolidayMap[ph.PublicHolidayDate.String()] = ph.HolidayType
	}

	var res []*oa.Engagement
	var errorArray []string
	for i, row := range rows[1:] {
		x := i + 2
		var colList [9]string
		for i, colCell := range row {
			colList[i] = colCell
		}

		// 校验项目编号
		if colList[0] == "" {
			errorArray = append(errorArray, fmt.Sprintf("第%d行项目名称未填写", x))
		}
		engagementCode := new(oa.EngagementCode)
		services.Slave().Model(oa.EngagementCode{}).Where("engagement_code = ?", colList[0]).First(engagementCode)
		if engagementCode.ID == 0 {
			errorArray = append(errorArray, fmt.Sprintf("第%d行项目名称未找到", x))
		}
		fmt.Println(engagementCode)
		// 校验员工
		employee := new(oa.Employee)
		if colList[1] == "" {
			errorArray = append(errorArray, fmt.Sprintf("第%d行员工未填写", x))
		}
		services.Slave().Preload("Level").Take(employee, "name = ?", colList[1])
		fmt.Println(employee)
		fmt.Println(employee.Level)
		if employee.ID == 0 {
			errorArray = append(errorArray, fmt.Sprintf("第%d行员工未找到", x))
		}
		for k, col := range colList[2:9] {
			workTime, _ := time.Parse("01-02-06", rows[0][k+2])

			y := k + 3
			colInt, err := strconv.Atoi(col)
			if err != nil {
				errorArray = append(errorArray, fmt.Sprintf("第%d行%d列工时错误", x, y))
			}

			// 判断是否放假
			ph, ok := publicHolidayMap[rows[0][k+2]]
			if ok {
				if ph == "holiday" { // 放假 不判断

				} else if ph == "workday" { //补假 判断
					if colInt < 8 {
						errorArray = append(errorArray, fmt.Sprintf("第%d行%d列工时小于8小时", x, y))
					}
				}
			} else {
				if k != 5 && k != 6 { // 周末不判断
					if colInt < 8 {
						errorArray = append(errorArray, fmt.Sprintf("第%d行%d列工时小于8小时", x, y))
					}
				}
			}

			engagementCost := float64(((engagementCode.OCRate * employee.Level.OCRate) + (engagementCode.CCRate * employee.Level.CCRate)) * float32(colInt))
			em := &oa.Engagement{
				EngagementCode: colList[0],
				EmployeeID:     int(employee.ID),
				EmployeeName:   employee.Name,
				EngagementDate: workTime,
				EngagementHour: colInt,
				EngagementCost: engagementCost,
			}
			res = append(res, em)
		}

	}

	if len(errorArray) > 0 {
		return nil, errors.New(strings.Join(errorArray, "-"))
	}

	return res, nil
}
