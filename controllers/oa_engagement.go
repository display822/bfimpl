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

// @Title 人员管理周数据列表
// @Description 人员管理周数据列表
// @Param	pagenum	    query	int	false	"分页"
// @Param	pagesize	query	int	false	"分页"
// @Param	period_time	query	string	false	"周期时间"
// @Success 200 {string} ""
// @Failure 500 server internal err
// @router /period [get]
func (e *EngagementController) PeriodList() {
	userID, _ := e.GetInt("userID")
	periodTime := e.GetString("period_time")
	var es []oa.Engagement
	db := services.Slave().Where("department_id =?", userID)
	if periodTime != "" {
		db = db.Where("period_time = ?", periodTime)
	}
	db.Order("created_at desc").Group("period_time").Find(&es)

	e.Correct(es)
}

// @Title 人员管理周数据详细
// @Description 人员管理周数据详细
// @Param	period_time	query	string	true	"周期时间"
// @Success 200 {string} ""
// @Failure 500 server internal err
// @router /period/detail [get]
func (e *EngagementController) PeriodDetail() {
	userID, _ := e.GetInt("userID")
	periodTime := e.GetString("period_time")
	var es []oa.Engagement
	services.Slave().Where("department_id =?", userID).Where("period_time = ?", periodTime).
		Order("created_at").Find(&es)

	m := make(map[string]map[string]int)
	for _, item := range es {
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

// @Title 人员管理信息
// @Description 人员管理信息
// @Param	engagement_codes query string	false	"项目编号"
// @Param	begin_time query string	false	"开始时间"
// @Param	end_time query string	false	"结束时间"
// @Success 200 {string} ""
// @Failure 500 server internal err
// @router / [get]
func (e *EngagementController) List() {
	engagementCodes := e.GetString("engagement_codes")
	if engagementCodes == "" {
		e.ErrorOK("need engagement_codes")
	}
	beginTime := e.GetString("begin_time")
	if beginTime == "" {
		e.ErrorOK("need begin_time")
	}
	endTime := e.GetString("end_time")
	if endTime == "" {
		e.ErrorOK("need end_time")
	}

	ecs := strings.Split(engagementCodes, ",")
	fmt.Println(ecs)
	var ers forms.EngagementResult
	for _, ec := range ecs {
		var ee forms.E
		fmt.Println(ec)

		var es []oa.Engagement
		services.Slave().Debug().Where("engagement_date >= ?", beginTime).Where("engagement_date <=?", endTime).
			Where("engagement_code = ?", ec).Find(&es)
		if len(es) == 0 {
			fmt.Println("0")
			continue
		}

		for _, item := range es {
			ee.CostSummary += item.EngagementCost
			ee.HourSummary += item.EngagementHour
		}

		m := make(map[string]map[string]int)
		for _, item := range es {
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
		ee.EngagementCode = ec
		ee.EmployeeNums = len(data)
		ee.EngagementList = data

		ers.CostSummary += ee.CostSummary
		ers.HourSummary += ee.HourSummary
		ers.EmployeeNums += ee.EmployeeNums
		ers.List = append(ers.List, ee)
	}

	e.Correct(ers)
}

// @Title 创建人员管理
// @Description 创建人员管理
// @Param  file form-data binary true "文件"
// @Success 200 {string} ""
// @Failure 500 server internal err
// @router / [post]
func (e *EngagementController) Create() {
	userID, _ := e.GetInt("userID")
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
		e.ErrorOK(err.Error())
	}

	res, err := EngagementDetailFile(f, userID)
	if err != nil {
		e.ErrorOK(err.Error())
	}

	var eng []oa.Engagement
	err = services.Slave().Where("period_time = ?", res[0].PeriodTime).Where("department_id = ?", userID).Find(&eng).Error
	if err == nil {
		log.GLogger.Info("Exist")
		err = services.Slave().Delete(&eng).Error
		if err != nil {
			e.ErrorOK(MsgServerErr)
		}
	}

	err = oa.BatchEngagementCreate(services.Slave(), res)
	if err != nil {
		e.ErrorOK(MsgServerErr)
	}

	e.Correct("")
}

// @Title 解析人员管理内容的excel文件
// @Description 解析人员管理内容的excel文件
// @Param  file form-data binary true "文件"
// @Success 200 {object} ""
// @Failure 500 server internal err
// @router /details [post]
func (e *EngagementController) ParseEngagementDetailFile() {
	userID, _ := e.GetInt("userID")
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
		e.ErrorOK(err.Error())
	}

	res, err := EngagementDetailFile(f, userID)
	if err != nil {
		e.ErrorOK(err.Error())
	}

	log.GLogger.Info("res:%s", res)

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

func EngagementDetailFile(f *excelize.File, departmentID int) ([]*oa.Engagement, error) {
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return nil, err
	}

	log.GLogger.Info("row len; %d", len(rows))
	if len(rows) < 2 {
		return nil, errors.New("无数据")
	}
	log.GLogger.Info("len(rows[0]): %d", len(rows[0]))

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

	workTimeIndex, err := time.Parse("01-02-06", rows[0][2])
	workStartTime := workTimeIndex
	if err != nil {
		return nil, errors.New("时间无法识别")
	}
	if workTimeIndex.Weekday() != time.Monday {
		return nil, errors.New("需要从周一起")
	}

	for _, v := range rows[0][3:9] {
		index, err := time.Parse("01-02-06", v)
		if err != nil {
			return nil, errors.New("时间无法识别")
		}
		workTimeIndex = workTimeIndex.AddDate(0, 0, 1)
		if !index.Equal(workTimeIndex) {
			return nil, errors.New("时间不连续")
		}
	}

	// 重复map
	existMap := make(map[string][]int)

	log.GLogger.Info("workTimeIndex", workTimeIndex)
	log.GLogger.Info("startTime", workStartTime)
	log.GLogger.Info("endTime", workTimeIndex)
	var phs []oa.PublicHoliday
	services.Slave().Where("public_holiday_date >= ?", workStartTime).
		Where("public_holiday_date <= ?", workTimeIndex).
		Find(&phs)
	log.GLogger.Info("phs", phs)
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
			break
		}
		// 校验员工
		employee := new(oa.Employee)
		if colList[1] == "" {
			errorArray = append(errorArray, fmt.Sprintf("第%d行员工未填写", x))
		}
		services.Slave().Preload("Level").Take(employee, "name = ?", colList[1])

		if employee.ID == 0 {
			errorArray = append(errorArray, fmt.Sprintf("第%d行员工未找到", x))
			break
		}

		// 计算重复
		v, ok := existMap[colList[0]+colList[1]]
		if ok {
			existMap[colList[0]+colList[1]] = append(v, x)
		} else {
			l := []int{x}
			existMap[colList[0]+colList[1]] = l
		}

		for k, col := range colList[2:9] {
			workTime, _ := time.Parse("01-02-06", rows[0][k+2])

			y := k + 3
			colInt, err := strconv.Atoi(col)
			if err != nil {
				errorArray = append(errorArray, fmt.Sprintf("第%d行%d列工时错误", x, y))
			}

			// 判断是否放假
			ph, ok := publicHolidayMap[workTime.Format("2006-01-02")]
			if ok {
				if ph == "holiday" { // 放假 不判断

				} else if ph == "workday" { //补假 判断
					if colInt < 8 {
						errorArray = append(errorArray, fmt.Sprintf("第%d行第%d列工时小于8小时", x, y))
					}
				}
			} else {
				if k != 5 && k != 6 { // 周末不判断
					if colInt < 8 {
						errorArray = append(errorArray, fmt.Sprintf("第%d行第%d列工时小于8小时", x, y))
					}
				}
			}

			engagementCost := float64(((engagementCode.OCRate * employee.Level.OCRate) + (engagementCode.CCRate * employee.Level.CCRate)) * float32(colInt))
			em := &oa.Engagement{
				EngagementCode: colList[0],
				EmployeeID:     int(employee.ID),
				DepartmentID:   departmentID,
				EmployeeName:   employee.Name,
				PeriodTime:     workStartTime.Format("2006/01/02") + "-" + workTimeIndex.Format("2006/01/02"),
				EngagementDate: workTime,
				EngagementHour: colInt,
				EngagementCost: engagementCost,
			}
			res = append(res, em)
		}

		// 检测重复
		for _, v := range existMap {
			if len(v) > 1 {
				var msg []string
				for _, index := range v {
					msg = append(msg, fmt.Sprintf("第%d行", index))
				}
				errorArray = append(errorArray, fmt.Sprintf("%s重复", strings.Join(msg, ",")))
			}
		}
	}

	if len(errorArray) > 0 {
		return nil, errors.New(strings.Join(errorArray, "-"))
	}

	return res, nil
}
