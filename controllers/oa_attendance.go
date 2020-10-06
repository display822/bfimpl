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
