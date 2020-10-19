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
	"os"
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
	userToday := make(map[string]bool)
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
			key := row[1] + attendanceTmp.AttendanceDate.String()
			if userToday[key] {
				attendanceTmp.Result = "早退"
			} else {
				attendanceTmp.Result = "迟到"
				userToday[key] = true
			}
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

// @Title 修改临时数据
// @Description 修改临时数据
// @Param	id	query	string	true	"数据id"
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /attendance/create/tmp [post]
func (a *AttendanceController) CreateAttendanceTmp() {
	param := new(oa.AttendanceTmp)
	err := json.Unmarshal(a.Ctx.Input.RequestBody, param)
	if err != nil {
		log.GLogger.Error("parse attendance tmp:%s", err.Error())
		a.ErrorOK(MsgInvalidParam)
	}
	services.Slave().Create(param)
	a.Correct(param)
}

// @Title 查询考勤临时数据
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

// @Title 查询已确认考勤数据
// @Description 查询考勤
// @Param	name	query	string	true	"姓名"
// @Param	year	query	string	false	"年"
// @Param	month	query	string	false	"月"
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /attendance [get]
func (a *AttendanceController) GetUserAttendanceByMonth() {
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
	data := make([]*oa.Attendance, 0)
	services.Slave().Model(oa.Attendance{}).Where("name = ? and attendance_date >= ? and attendance_date <= ?",
		name, startDate, endDate).Find(&data)
	result := make([]*oa.UserAttendanceTmp, 0)
	for _, at := range data {
		tmpData := []*oa.AttendanceTmp{{
			ID:             at.ID,
			EmployeeID:     at.EmployeeID,
			Dept:           at.Dept,
			Name:           at.Name,
			AttendanceDate: at.AttendanceDate,
			CheckTime:      at.CheckIn,
			Status:         at.InStatus,
			Result:         at.InResult,
			LeaveID:        at.LeaveID,
		}, {
			ID:             at.ID,
			EmployeeID:     at.EmployeeID,
			Dept:           at.Dept,
			Name:           at.Name,
			AttendanceDate: at.AttendanceDate,
			CheckTime:      at.CheckOut,
			Status:         at.OutStatus,
			Result:         at.OutResult,
			LeaveID:        at.LeaveID,
		}}
		result = append(result, &oa.UserAttendanceTmp{
			Date: at.AttendanceDate.String(),
			Tmps: tmpData,
		})
	}
	a.Correct(result)
}

// @Title 批量确认考勤
// @Description 查询考勤
// @Param	name	body	string	true	"姓名数组"
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /attendance [post]
func (a *AttendanceController) ConfirmUserAttendance() {
	names := make([]string, 0)
	err := json.Unmarshal(a.Ctx.Input.RequestBody, &names)
	if err != nil {
		log.GLogger.Error("parse names err:%s", err.Error())
		a.ErrorOK(MsgInvalidParam)
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
	tmps := make([]*oa.AttendanceTmp, 0)
	services.Slave().Model(oa.AttendanceTmp{}).Where("attendance_date >= ? and attendance_date <= ? and name in (?)",
		startDate, endDate, names).Find(&tmps)
	data := make([]*oa.AttendanceSimple, 0)
	num := 0
	userDatas := make(map[string]map[string]int)
	for _, row := range tmps {
		//部门，姓名，时间
		ud, uOK := userDatas[row.Name]
		if !uOK {
			ud = make(map[string]int)
		}
		date := row.AttendanceDate
		attendanceIndex, dOK := ud[date.String()]
		if !dOK {
			//新增一条今天的记录,该行数据为签入
			attendance := &oa.AttendanceSimple{
				Dept:           row.Dept,
				Name:           row.Name,
				AttendanceDate: date,
				CheckIn:        row.CheckTime,
				InStatus:       row.Status,
				InResult:       row.Result,
			}
			if row.LeaveID > 0 {
				attendance.LeaveId = row.LeaveID
			}
			data = append(data, attendance)
			ud[date.String()] = num
			num++
		} else {
			// 修改签出时间
			data[attendanceIndex].CheckOut = row.CheckTime
			data[attendanceIndex].OutStatus = row.Status
			data[attendanceIndex].OutResult = row.Result
			if row.LeaveID > 0 {
				data[attendanceIndex].LeaveId = row.LeaveID
			}
		}
		userDatas[row.Name] = ud
	}
	if len(data) == 0 {
		a.ErrorOK("未找到考勤数据")
	}
	//拼接sql
	sql := "insert into attendances(created_at,dept,name,attendance_date,check_in,check_out,in_status,out_status," +
		"in_result,out_result,leave_id) values"
	realData := make([]string, 0)
	now := time.Now().Format(models.TimeFormat)
	for _, d := range data {
		realData = append(realData, d.String(now))
	}
	sql += strings.Join(realData, ",")
	sql += "on duplicate key update updated_at=values(created_at),dept=values(dept),name=values(name),attendance_date=values(attendance_date)" +
		",check_in=values(check_in),check_out=values(check_out),in_status=values(in_status)," +
		"out_status=values(out_status),in_result=values(in_result),out_result=values(out_result),leave_id=values(leave_id);"
	tx := services.Slave().Begin()
	err = tx.Exec(sql).Error
	if err != nil {
		tx.Rollback()
		log.GLogger.Error("考勤sql：%s", err.Error())
		a.ErrorOK(MsgServerErr)
	}
	err = tx.Exec("update attendance_tmp set is_confirm = 1 where attendance_date >= ? and attendance_date <= ?"+
		" and name in (?)", startDate, endDate, names).Error
	if err != nil {
		tx.Rollback()
		log.GLogger.Error("确认考勤sql：%s", err.Error())
		a.ErrorOK(MsgServerErr)
	}
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		log.GLogger.Error("考勤sql：%s", err.Error())
		a.ErrorOK(MsgServerErr)
	}
	a.Correct("")
}

// @Title 导出POS考勤
// @Description 导出POS考勤
// @Param	year	query	string	true	"年"
// @Param	month	query	string	true	"月"
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /attendance/pos [get]
func (a *AttendanceController) ExportPos() {
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
	//查询月份数据
	attendances := make([]*oa.Attendance, 0)
	services.Slave().Model(oa.Attendance{}).Where("attendance_date >= ? and attendance_date <= ?",
		startDate, endDate).Find(&attendances)
	//查询 tapd 账号
	emps := make([]*oa.EmpPos, 0)
	services.Slave().Table("employees").Select("name,tapd").Scan(&emps)
	nameTapd := make(map[string]string)
	for _, emp := range emps {
		nameTapd[emp.Name] = emp.Tapd
	}
	f := excelize.NewFile()
	_ = f.SetSheetRow("Sheet1", "A1", &[]interface{}{"日期时间", "人员编号", "姓名", "姓名"})
	num := 2
	for _, at := range attendances {
		tapd := nameTapd[at.Name]
		_ = f.SetSheetRow("Sheet1", "A"+strconv.Itoa(num), &[]interface{}{
			at.CheckIn.String(), 0, tapd, at.Name})
		num++
		inTime, outTime := strings.Split(at.CheckIn.String(), " "), strings.Split(at.CheckOut.String(), " ")
		if outTime[0] > inTime[0] {
			//下班时间是第二天
			_ = f.SetSheetRow("Sheet1", "A"+strconv.Itoa(num), &[]interface{}{
				inTime[0] + " 23:59:00", 0, tapd, at.Name})
			num++
			//下班时间是07:00
			if outTime[1] == "07:00:00" {
				_ = f.SetSheetRow("Sheet1", "A"+strconv.Itoa(num), &[]interface{}{
					outTime[0] + " 06:59:59", 0, tapd, at.Name})
				num++
			} else if outTime[1] > "07:00:00" {
				_ = f.SetSheetRow("Sheet1", "A"+strconv.Itoa(num), &[]interface{}{
					outTime[0] + " 07:00:00", 0, tapd, at.Name})
				num++
			}
		}
		_ = f.SetSheetRow("Sheet1", "A"+strconv.Itoa(num), &[]interface{}{
			at.CheckOut.String(), 0, tapd, at.Name})
		num++
	}
	fileName := year + "-" + month + "-pos.xlsx"
	f.SaveAs(fileName)
	a.Ctx.Output.Download(fileName, fileName)
	os.Remove(fileName)
}

// @Title 导出考勤请假数据
// @Description 导出考勤请假数据
// @Param	year	query	string	true	"年"
// @Param	month	query	string	true	"月"
// @Success 200 {string} "success"
// @Failure 500 server err
// @router /attendance/data [get]
func (a *AttendanceController) ExportData() {
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
	//查询月份数据
	attendances := make([]*oa.Attendance, 0)
	services.Slave().Model(oa.Attendance{}).Where("attendance_date >= ? and attendance_date <= ?",
		startDate, endDate).Find(&attendances)
	//查询所有请假记录
	leaves := make([]*oa.Leave, 0)
	services.Slave().Model(oa.Leave{}).Where("start_date >= ? and start_date <= ? and status = ?",
		startDate, endDate, models.FlowApproved).Find(&leaves)

	userIndex := make(map[string]int)
	index := 0
	data := make([]*oa.AttendanceExcel, 0)
	for _, at := range attendances {
		i, ok := userIndex[at.Name]
		if !ok {
			i = index
			userIndex[at.Name] = index
			index++
			tmp := &oa.AttendanceExcel{
				Dept: at.Dept,
				Name: at.Name,
			}
			data = append(data, tmp)
		}
		//今天工时
		if !at.CheckIn.IsZero() && !at.CheckOut.IsZero() {
			data[i].Total += at.CheckOut.SubToHour(at.CheckIn)
			if at.InResult == "迟到" {
				data[i].Late++
			}
			if at.OutResult == "早退" {
				data[i].Early++
			}
		} else {
			data[i].Forget += 1
		}
	}
	for _, leave := range leaves {
		//统计请假数据
		if leave.RealDuration != 0 {
			leave.Duration = leave.RealDuration
		}
		i, ok := userIndex[leave.EName]
		if ok {
			if leave.Type == models.LeaveAnnual {
				data[i].Annual += leave.Duration
			} else if leave.Type == models.LeaveShift {
				data[i].Leave += leave.Duration
			} else if leave.Type == models.LeaveSick {
				data[i].Sick += leave.Duration
			}
		}
	}
	//生成excel
	f := excelize.NewFile()
	_ = f.SetSheetRow("Sheet1", "A1", &[]interface{}{"部门", "姓名", "上班总工时", "总调休时长",
		"年假", "病假", "迟到", "早退", "旷工", "忘记打卡"})
	num := 2
	for _, at := range data {
		_ = f.SetSheetRow("Sheet1", "A"+strconv.Itoa(num), &[]interface{}{at.Dept, at.Name, at.Total, at.Leave,
			at.Annual, at.Sick, at.Late, at.Early, at.None, at.Forget})
		num++
	}
	fileName := year + "-" + month + ".xlsx"
	f.SaveAs(fileName)
	a.Ctx.Output.Download(fileName, fileName)
	os.Remove(fileName)
}
