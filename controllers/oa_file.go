/*
* Auth : acer
* Desc : 文件上传
* Time : 2020/9/6 14:51
 */

package controllers

import (
	"bfimpl/models"
	"bfimpl/models/oa"
	"bfimpl/services"
	"bfimpl/services/log"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/minio/minio-go"
)

type FileController struct {
	BaseController
}

// @Title 上传文件
// @Description 上传的文件
// @Param  file form-data binary true "文件"
// @Success 200 {string} "success"
// @Failure 500 server internal err
// @router /upload [post]
func (m *FileController) Upload() {
	_ = m.Ctx.Input.ParseFormOrMulitForm(100 << 20)
	file, header, err := m.GetFile("file")
	bucket := m.GetString("bucket", "default")
	if err != nil {
		log.GLogger.Error(err.Error())
		m.Error("need file")
	}
	defer file.Close()
	//创建path
	dir := "static/" + bucket
	f, e1 := os.Open(dir)
	if e1 != nil {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			m.Error(fmt.Sprintf("create %s err.", dir))
		}
	}
	f.Close()
	fileName := strconv.FormatInt(time.Now().Unix(), 10) + "_" + header.Filename
	tmpFile, err := os.Create(dir + "/" + fileName)
	if err != nil {
		log.GLogger.Error("upload err:%s", err.Error())
		m.Error(err.Error())
	}
	defer tmpFile.Close()
	_, err = io.Copy(tmpFile, file)
	if err != nil {
		m.Error(err.Error())
	}
	go upload(fileName, bucket)
	m.Correct("/" + bucket + "/" + fileName)
}
func upload(filename, bucket string) {
	contentType := "application/octet-stream"
	// 查询bucket 是否存在
	services.InitBucket(bucket)
	// 使用FPutObject上传一个文件。
	filePath := "static/" + bucket + "/" + filename
	n, err := services.MinIOClient().FPutObject(bucket, filename, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.GLogger.Error(err.Error())
	} else {
		log.GLogger.Info("Successfully uploaded %s of size %d\n", filename, n)
	}
	os.Remove(filePath)
}

// @Title 社保文件列表
// @Description 社保文件列表
// @Success 200 {string} "success"
// @Failure 500 server internal err
// @router /sslist [post]
func (m *FileController) SocialSecurityList() {
	ss := make([]*oa.SocialSecurity, 0)
	services.Slave().Model(oa.SocialSecurity{}).Find(&ss)
	m.Correct(ss)
}

//生成社保信息
func GeneraSheBao() {
	//生成excel
	f := excelize.NewFile()
	row := 1
	_ = f.SetCellStr("Sheet1", "A"+strconv.Itoa(row), "在职")
	row++
	_ = f.SetSheetRow("Sheet1", "A"+strconv.Itoa(row), &[]interface{}{"主体", "员工姓名", "状态", "入职日期",
		"离职日期", "身份证号", "户籍性质", "公积金号"})
	row++
	//查询在职员工
	existEmp := make([]*oa.Employee, 0)
	services.Slave().Where("status = 2").Preload("EmployeeBasic").Find(&existEmp)
	existIds := make([]uint, 0)
	for _, emp := range existEmp {
		existIds = append(existIds, emp.ID)
	}
	//查询本月入职和离职流程
	inAndOutEmp := make([]*oa.Employee, 0)
	now := time.Now()
	end := fmt.Sprintf("%d-%2d-15", now.Year(), now.Month())
	pre := now.AddDate(0, -1, 0)
	start := fmt.Sprintf("%d-%02d-16", pre.Year(), pre.Month())
	services.Slave().Where("(entry_date >= ? and entry_date <= ?) or (resignation_date >= ? and resignation_date <= ?)",
		start, end, start, end).Find(&inAndOutEmp)
	userInIds := make([]int, 0)
	userOutIds := make([]int, 0)
	except := make(map[uint]bool)
	for _, emp := range inAndOutEmp {
		if emp.ResignationDate.IsZero() {
			//入职
			except[emp.ID] = true
			userInIds = append(userInIds, int(emp.ID))
		} else {
			//离职
			userOutIds = append(userOutIds, int(emp.ID))
		}
	}
	//在职信息=============================
	//查询合同
	existMain := make([]*oa.ContractSimple, 0)
	services.Slave().Table("employee_contracts").Select("contract_main,max(contract_end_date) as enddate, employee_id").
		Where("employee_id in (?)", existIds).Group("employee_id").Scan(&existMain)
	empContract := make(map[int]*oa.ContractSimple)
	for i := range existMain {
		// eid -> 合同
		empContract[existMain[i].EmployeeID] = existMain[i]
	}
	userIn := make(map[int]*oa.Employee)
	for i, emp := range existEmp {
		if except[emp.ID] {
			userIn[int(emp.ID)] = existEmp[i]
			continue
		}
		var contractMain, huji, fund string
		if emp.EmployeeBasic != nil {
			huji = emp.EmployeeBasic.HujiType
			fund = emp.EmployeeBasic.PublicFund
		}
		if m, ok := empContract[int(emp.ID)]; ok {
			contractMain = m.ContractMain
		}
		_ = f.SetSheetRow("Sheet1", "A"+strconv.Itoa(row), &[]interface{}{contractMain, emp.Name, models.EmpStatus[emp.Status],
			emp.EntryDate, "-", emp.IDCard, huji, fund})
		row++
	}
	//新入职信息=============================
	_ = f.SetCellStr("Sheet1", "A"+strconv.Itoa(row), "新入职")
	row++
	for _, emp := range userIn {
		var contractMain, huji, fund string
		if emp.EmployeeBasic != nil {
			huji = emp.EmployeeBasic.HujiType
			fund = emp.EmployeeBasic.PublicFund
		}
		if m, ok := empContract[int(emp.ID)]; ok {
			contractMain = m.ContractMain
		}
		_ = f.SetSheetRow("Sheet1", "A"+strconv.Itoa(row), &[]interface{}{contractMain, emp.Name, models.EmpStatus[emp.Status],
			emp.EntryDate, "-", emp.IDCard, huji, fund})
		row++
	}
	_ = f.SetCellStr("Sheet1", "A"+strconv.Itoa(row), "已离职")
	row++
	//离职信息=============================
	//离职员工
	leaveEmp := make([]*oa.Employee, 0)
	services.Slave().Where(userOutIds).Preload("EmployeeBasic").Find(&leaveEmp)
	leaveMain := make([]*oa.ContractSimple, 0)
	services.Slave().Table("employee_contracts").Select("contract_main,max(contract_end_date) as enddate, employee_id").
		Where("employee_id in (?)", userOutIds).Group("employee_id").Scan(&leaveMain)
	leaveContract := make(map[int]*oa.ContractSimple)
	for i := range leaveMain {
		// eid -> 合同
		leaveContract[leaveMain[i].EmployeeID] = leaveMain[i]
	}
	for _, emp := range leaveEmp {
		var contractMain, huji, fund string
		if emp.EmployeeBasic != nil {
			huji = emp.EmployeeBasic.HujiType
			fund = emp.EmployeeBasic.PublicFund
		}
		if m, ok := leaveContract[int(emp.ID)]; ok {
			contractMain = m.ContractMain
		}
		_ = f.SetSheetRow("Sheet1", "A"+strconv.Itoa(row), &[]interface{}{contractMain, emp.Name, models.EmpStatus[emp.Status],
			emp.EntryDate, emp.ResignationDate, emp.IDCard, huji, fund})
		row++
	}
	saveFile := new(oa.SocialSecurity)
	fileName := fmt.Sprintf("%s至%s信息表.xlsx", start, end)
	saveFile.Name = fileName
	saveFile.DownloadUrl = "/socialsecurity/" + fileName
	f.SaveAs("static/" + saveFile.DownloadUrl)
	//保存数据库信息
	upload(fileName, "socialsecurity")
	services.Slave().Create(saveFile)
	os.Remove("static/" + saveFile.DownloadUrl)
}
