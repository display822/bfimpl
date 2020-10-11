/*
* Auth : acer
* Desc : 考勤
* Time : 2020/10/6 16:47
 */

package controllers

import (
	"bfimpl/models"
	"bfimpl/models/oa"
	"bfimpl/services/log"
	"strings"
	"time"

	"bfimpl/services"

	"strconv"

	"fmt"

	"encoding/json"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

const (
	excelTime = "1/2/06 15:04"
	Normal    = "Normal"
	Exception = "Exception"
)

type AttendanceController struct {
	BaseController
}

// @Title 上传考勤
// @Description 上传考勤
// @Param	file	body	binary	true	"excel文件"
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /attendance [post]
func (a *AttendanceController) UploadAttendance() {
	f, _, err := a.GetFile("file")
	if err != nil {
		log.GLogger.Error("上传考勤：%s", err.Error())
		a.ErrorOK(MsgInvalidParam)
	}
	//解析 xlsx
	defer f.Close()
	file, err := excelize.OpenReader(f)
	if err != nil {
		log.GLogger.Error("读取考勤：%s", err.Error())
		a.ErrorOK(MsgServerErr)
	}
	rows, err := file.GetRows("Sheet1")
	if err != nil {
		log.GLogger.Error("读取Sheet1：%s", err.Error())
		a.ErrorOK(MsgServerErr)
	}
	users := make([]string, 0)
	userDatas := make(map[string]map[string]*oa.AttendanceSimple)
	for _, row := range rows[1:] {
		if len(row) < 3 {
			continue
		}
		//部门，姓名，时间
		ud, uOK := userDatas[row[1]]
		if !uOK {
			users = append(users, row[1])
			ud = make(map[string]*oa.AttendanceSimple)
		}
		checkTime, err := time.Parse(excelTime, row[2])
		if err != nil {
			log.GLogger.Error(err.Error())
			continue
		}
		date := models.Date(checkTime)
		attendance, dOK := ud[date.String()]
		if !dOK {
			//新增一条今天的记录,该行数据为签入
			attendance = &oa.AttendanceSimple{
				Dept:           row[0],
				Name:           row[1],
				AttendanceDate: date,
				CheckIn:        models.Time(checkTime),
				InStatus:       Normal,
			}
			if strings.Split(checkTime.String(), " ")[1] > "09:45" {
				//迟到
				attendance.InStatus = Exception
				attendance.InResult = "迟到"
			}
		} else {
			// 修改签出时间
			attendance.CheckOut = models.Time(checkTime)
			attendance.OutStatus = Normal
			if strings.Split(checkTime.String(), " ")[1] < "18:30" {
				//迟到
				attendance.OutStatus = Exception
				attendance.OutResult = "早退"
			}
		}
		ud[date.String()] = attendance
		userDatas[row[1]] = ud
	}
	//拼接sql
	sql := "insert into attendances(created_at,dept,name,attendance_date,check_in,check_out,in_status,out_status," +
		"in_result,out_result) values"
	realData := make([]string, 0)
	now := time.Now().Format(models.TimeFormat)
	for _, u := range users {
		for _, v := range userDatas[u] {
			realData = append(realData, v.String(now))
		}
	}
	sql += strings.Join(realData, ",")
	sql += "on duplicate key update updated_at=values(created_at),dept=values(dept),name=values(name),attendance_date=values(attendance_date)" +
		",check_in=values(check_in),check_out=values(check_out),in_status=values(in_status)," +
		"out_status=values(out_status),in_result=values(in_result),out_result=values(out_result);"
	err = services.Slave().Exec(sql).Error
	if err != nil {
		log.GLogger.Error("考勤sql：%s", err.Error())
		a.ErrorOK(MsgServerErr)
	}
	a.Correct("")
}

// @Title 查询考勤
// @Description 查询考勤
// @Param	name	query	string	true	"姓名"
// @Param	year	query	string	false	"年"
// @Param	month	query	string	false	"月"
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /attendance [get]
func (a *AttendanceController) GetAttendances() {
	name := a.GetString("name")
	year := a.GetString("year")
	month := a.GetString("month")
	imonth, _ := a.GetInt("month", -1)
	if year == "" || month == "" {
		now := time.Now()
		year = strconv.Itoa(now.Year())
		imonth = int(now.Month())
		month = strconv.Itoa(imonth)
		if len(month) == 1 {
			month = "0" + month
		}
	}
	startDate := strings.Join([]string{year, month, "01"}, "-")
	endDate := fmt.Sprintf("%s-%s-%d", year, month, models.Months[imonth])
	query := services.Slave().Model(oa.Attendance{}).Where("attendance_date >= ? and attendance_date <= ?",
		startDate, endDate)
	if name != "" {
		query = query.Where("name like ?", "%"+name+"%")
	}
	data := make([]*oa.Attendance, 0)
	query.Order("attendance_date").Find(&data)

	order := make(map[string]int)
	userAttendances := make([]*oa.UserAttendance, 0)
	userNum := 0
	for i, attendance := range data {
		uaIndex, ok := order[attendance.Name]
		if !ok {
			order[attendance.Name] = userNum
			userNum++
			tmp := &oa.UserAttendance{
				Dept:        attendance.Dept,
				Name:        attendance.Name,
				Attendances: []*oa.Attendance{data[i]},
			}
			userAttendances = append(userAttendances, tmp)
		} else {
			userAttendances[uaIndex].Attendances = append(userAttendances[uaIndex].Attendances, data[i])
		}
	}
	deptUser := make([]*oa.DeptUser, 0)
	deptNum := 0
	for i, u := range userAttendances {
		duIndex, ok := order[u.Dept]
		if !ok {
			order[u.Dept] = deptNum
			deptNum++
			tmp := &oa.DeptUser{
				Dept:  u.Dept,
				Users: []*oa.UserAttendance{userAttendances[i]},
			}
			deptUser = append(deptUser, tmp)
		} else {
			deptUser[duIndex].Users = append(deptUser[duIndex].Users, userAttendances[i])
		}
	}
	a.Correct(deptUser)
}

// @Title 修改考勤
// @Description 修改考勤
// @Param	json	body	json	true	"考勤数据"
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /attendance [put]
func (a *AttendanceController) UpdateAttendance() {
	attendance := new(oa.Attendance)
	err := json.Unmarshal(a.Ctx.Input.RequestBody, attendance)
	if err != nil {
		log.GLogger.Error("update attendance:%s", err.Error())
		a.ErrorOK(MsgInvalidParam)
	}
	tmp := new(oa.Attendance)
	err = services.Slave().Take(tmp, "id = ?", attendance.ID).Error
	if err != nil {
		a.ErrorOK(MsgInvalidParam)
	}
	attendance.CreatedAt = tmp.CreatedAt
	attendance.UpdatedAt = tmp.UpdatedAt
	services.Slave().Save(attendance)
	a.Correct("")
}

// @Title 上传考勤到临时表
// @Description 上传考勤到临时表
// @Param	file	body	binary	true	"excel文件"
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /attendance/tmp [post]
func (a *AttendanceController) UploadAttendanceTmp() {
	f, _, err := a.GetFile("file")
	if err != nil {
		log.GLogger.Error("上传考勤：%s", err.Error())
		a.ErrorOK(MsgInvalidParam)
	}
	//解析 xlsx
	defer f.Close()
	file, err := excelize.OpenReader(f)
	if err != nil {
		log.GLogger.Error("读取考勤：%s", err.Error())
		a.ErrorOK(MsgServerErr)
	}
	rows, err := file.GetRows("Sheet1")
	if err != nil {
		log.GLogger.Error("读取Sheet1：%s", err.Error())
		a.ErrorOK(MsgServerErr)
	}
	users := make([]string, 0)
	userDatas := make(map[string][]*oa.AttendanceTmp)
	for _, row := range rows[1:] {
		if len(row) < 3 {
			continue
		}
		//部门，姓名，时间
		ud, uOK := userDatas[row[1]]
		if !uOK {
			users = append(users, row[1])
			ud = make([]*oa.AttendanceTmp, 0)
		}
		checkTime, err := time.Parse(excelTime, row[2])
		if err != nil {
			log.GLogger.Error(err.Error())
			continue
		}
		attendanceTmp := &oa.AttendanceTmp{
			Dept:           row[0],
			Name:           row[1],
			AttendanceDate: models.Date(checkTime),
			CheckTime:      models.Time(checkTime),
			Status:         Normal,
		}
		t := strings.Split(attendanceTmp.CheckTime.String(), " ")
		if len(t) == 2 && t[1] > "09:45" && t[1] < "18:30" {
			attendanceTmp.Status = Exception
		}
		userDatas[row[1]] = append(ud, attendanceTmp)
	}
	//拼接sql
	sql := "insert into attendance_tmp(created_at,dept,name,attendance_date,check_time,status) values"
	realData := make([]string, 0)
	now := time.Now().Format(models.TimeFormat)
	for _, u := range users {
		for _, v := range userDatas[u] {
			realData = append(realData, v.String(now))
		}
	}
	sql += strings.Join(realData, ",")
	sql += "on duplicate key update dept=values(dept),name=values(name),attendance_date=values(attendance_date)" +
		",check_time=values(check_time),status=values(status);"
	err = services.Slave().Exec(sql).Error
	if err != nil {
		log.GLogger.Error("考勤sql：%s", err.Error())
		a.ErrorOK(MsgServerErr)
	}
	a.Correct("")
}

// @Title 查询考勤人员
// @Description 查询考勤
// @Param	name	query	string	true	"姓名"
// @Param	year	query	string	false	"年"
// @Param	month	query	string	false	"月"
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /attendance/deptuser [get]
func (a *AttendanceController) GetAttendanceUserByDept() {
	name := a.GetString("name")
	year := a.GetString("year")
	month := a.GetString("month")
	imonth, _ := a.GetInt("month", -1)
	if year == "" || month == "" {
		now := time.Now()
		year = strconv.Itoa(now.Year())
		imonth = int(now.Month())
		month = strconv.Itoa(imonth)
		if len(month) == 1 {
			month = "0" + month
		}
	}
	startDate := strings.Join([]string{year, month, "01"}, "-")
	endDate := fmt.Sprintf("%s-%s-%d", year, month, models.Months[imonth])
	data := make([]*oa.AttendanceUser, 0)
	services.Slave().Raw("select dept,name,min(is_confirm) is_confirm from attendance_tmp where "+
		"attendance_date >= ? and attendance_date <= ? and name like ? group by dept,name", startDate, endDate,
		"%"+name+"%").Scan(&data)
	result := make([]*oa.DeptUsers, 0)
	order := make(map[string]int)
	num := 0
	for i, d := range data {
		dIndex, ok := order[d.Dept]
		if !ok {
			order[d.Dept] = num
			num++
			deptUser := &oa.DeptUsers{
				Dept:  d.Dept,
				Users: []*oa.AttendanceUser{data[i]},
			}
			result = append(result, deptUser)
		} else {
			result[dIndex].Users = append(result[dIndex].Users, data[i])
		}
	}
	a.Correct(result)
}

// @Title 删除临时数据
// @Description 删除临时数据
// @Param	id	query	string	true	"数据id"
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /attendance/tmp/:id [delete]
func (a *AttendanceController) DeleteAttendanceTmp() {
	id, _ := a.GetInt(":id", -1)
	err := services.Slave().Delete(oa.AttendanceTmp{}, "id = ?", id).Error
	if err != nil {
		log.GLogger.Error("delete attendance tmp:%s", err.Error())
		a.ErrorOK(MsgServerErr)
	}
	a.Correct("")
}

// @Title 修改临时数据
// @Description 修改临时数据
// @Param	id	query	string	true	"数据id"
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /attendance/tmp [put]
func (a *AttendanceController) UpdateAttendanceTmp() {
	param := new(oa.AttendanceTmp)
	err := json.Unmarshal(a.Ctx.Input.RequestBody, param)
	if err != nil {
		log.GLogger.Error("parse attendance tmp:%s", err.Error())
		a.ErrorOK(MsgInvalidParam)
	}
	tmp := new(oa.AttendanceTmp)
	err = services.Slave().Take(tmp, "id = ?", param.ID).Error
	if err != nil {
		log.GLogger.Error("query attendance tmp:%s", err.Error())
		a.ErrorOK(MsgServerErr)
	}
	param.CreatedAt = tmp.CreatedAt
	services.Slave().Save(param)
	a.Correct(param)
}

// @Title 查询考勤
// @Description 查询考勤
// @Param	name	query	string	true	"姓名"
// @Param	year	query	string	false	"年"
// @Param	month	query	string	false	"月"
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /attendance/tmp [get]
func (a *AttendanceController) GetUserAttendanceTmps() {
	name := a.GetString("name")
	if name == "" {
		a.ErrorOK("请选择员工姓名")
	}
	year := a.GetString("year")
	month := a.GetString("month")
	imonth, _ := a.GetInt("month", -1)
	if year == "" || month == "" {
		now := time.Now()
		year = strconv.Itoa(now.Year())
		imonth = int(now.Month())
		month = strconv.Itoa(imonth)
		if len(month) == 1 {
			month = "0" + month
		}
	}
	startDate := strings.Join([]string{year, month, "01"}, "-")
	endDate := fmt.Sprintf("%s-%s-%d", year, month, models.Months[imonth])
	tmpData := make([]*oa.AttendanceTmp, 0)
	services.Slave().Model(oa.AttendanceTmp{}).Where("name = ? and attendance_date >= ? and attendance_date <= ?",
		name, startDate, endDate).Find(&tmpData)

	result := make([]*oa.UserAttendanceTmp, 0)
	order := make(map[string]int)
	num := 0
	for i, tmp := range tmpData {
		date := tmp.AttendanceDate.String()
		dIndex, ok := order[date]
		if !ok {
			order[date] = num
			num++
			ut := &oa.UserAttendanceTmp{
				Date: date,
				Tmps: []*oa.AttendanceTmp{tmpData[i]},
			}
			result = append(result, ut)
		} else {
			result[dIndex].Tmps = append(result[dIndex].Tmps, tmpData[i])
		}
	}
	a.Correct(result)
}
