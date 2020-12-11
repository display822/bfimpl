/*
* Auth : acer
* Desc : 报销
* Time : 2020/12/4 15:45
 */

package controllers

import (
	"bfimpl/models"
	"bfimpl/models/oa"
	"bfimpl/services/log"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

type ExpenseController struct {
	BaseController
}

// @Title 解析用户报销内容的excel文件
// @Description 解析用户报销内容的excel文件
// @Param  file form-data binary true "文件"
// @Success 200 {string} "success"
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
	err = Read(f)
	if err != nil {
		fmt.Println(err)
		e.ErrorOK(err.Error())
	}
}

func Read(f *excelize.File) error {
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return err
	}

	log.GLogger.Info("row len; %d", len(rows))
	if len(rows) <= 2 {
		return errors.New("无数据")
	}
	if len(rows[0]) < 6 {
		return errors.New("首行表头字段有误, 无法识别")
	}
	fmt.Println(len(rows[0]))

	for i, v := range rows[0][0:6] {
		if oa.ExcelHeaderArray[i] != v {
			return errors.New("首行表头字段有误, 无法识别")
		}
	}
	var errorArray []string
	for i, row := range rows[1:] {
		fmt.Println(row)
		var colList [6]string
		for i, colCell := range row {
			colList[i] = colCell
			fmt.Println(colList)
		}

		if colList[0] == "" {
			errorArray = append(errorArray, fmt.Sprintf("第%d行费用发生日期未填写", i))
		} else {
			_, err := time.Parse(models.TimeFormat, colList[0])
			if err != nil {
				errorArray = append(errorArray, fmt.Sprintf("第%d行费用发生日期格式不正确", i))
			}
		}

	}
	if len(errorArray) > 0 {
		return errors.New(strings.Join(errorArray, "-"))
	}
	return nil
}
