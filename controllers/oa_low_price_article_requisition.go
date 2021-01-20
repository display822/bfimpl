/**
* @author : yi.zhang
* @description : controllers 描述
* @date   : 2021-01-11 16:30
 */

package controllers

import (
	"bfimpl/models"
	"bfimpl/models/oa"
	"bfimpl/services"
	"bfimpl/services/log"
	"encoding/json"
	"strings"
)

type LowPriceArticleRequisitionController struct {
	BaseController
}

// @Title 易耗品操作记录列表
// @Description 易耗品列表
// @Param	pagenum	    query	int	false	"页码"
// @Param	pagesize	query	int	false	"页数"
// @Param	low_price_article_id	query	bool	false	"易耗品id"
// @Param	employee_name	query	bool	false	"员工姓名"
// @Success 200 {object} []oa.LowPriceArticleRequisition
// @Failure 500 server internal err
// @router / [get]
func (l *LowPriceArticleRequisitionController) List() {
	pageNum, _ := l.GetInt("pagenum", 1)
	pageSize, _ := l.GetInt("pagesize", 10)
	lowPriceArticleID := l.GetString("low_price_article_id")
	category := l.GetString("category")
	employeeName := l.GetString("employee_name")
	var resp struct {
		Total int                              `json:"total"`
		List  []*oa.LowPriceArticleRequisition `json:"list"`
	}
	var lpar []*oa.LowPriceArticleRequisition
	db := services.Slave()
	if category != "" {
		db = db.Where("operator_category = ?", category)
	}
	if employeeName != "" {
		db = db.Where("associate_employee_name = ?", employeeName)
	}
	db.Where("low_price_article_id =?", lowPriceArticleID).Limit(pageSize).Offset((pageNum - 1) * pageSize).Order("created_at desc").Find(&lpar).Limit(-1).Offset(-1).Count(&resp.Total)
	resp.List = lpar
	l.Correct(resp)
}

// @Title 易耗品借出
// @Description 易耗品借出
// @Param	body body oa.LowPriceArticleRequisition true "易耗品记录"
// @Success 200 {string} ""
// @Failure 500 server internal err
// @router /outgoing [post]
func (l *LowPriceArticleRequisitionController) Outgoing() {
	uID, _ := l.GetInt("userID", 0)
	uName := l.GetString("userName")
	lpar := new(oa.LowPriceArticleRequisition)
	err := json.Unmarshal(l.Ctx.Input.RequestBody, lpar)
	if err != nil {
		log.GLogger.Error("parse low_price_article_requisition err:%s", err.Error())
		l.ErrorOK(MsgInvalidParam)
	}
	if lpar.Quantity <= 0 {
		l.ErrorOK("数量需要大于0")
	}
	tx := services.Slave().Begin()

	var lowPriceArticle oa.LowPriceArticle
	err = tx.Model(&oa.LowPriceArticle{}).Where("id = ?", lpar.LowPriceArticleID).Find(&lowPriceArticle).Error
	if err != nil {
		log.GLogger.Error("get low_price_article_requisition err:%s", err.Error())
		tx.Rollback()
		l.ErrorOK(MsgServerErr)
	}

	if lpar.Quantity > (lowPriceArticle.TotalQuantity - lowPriceArticle.ScrapQuantity - lowPriceArticle.OutgoingQuantity) {
		l.ErrorOK("数量不足")
	}

	lpar.OperatorID = uID
	lpar.OperatorName =

		log.GLogger.Info("low_price_article_requisition :%+v", lpar)
	var lpars []*oa.LowPriceArticleRequisition
	for i := 0; i < lpar.Quantity; i++ {
		lpars = append(lpars, lpar)
	}
	// 添加借出记录
	err = oa.BatchRequisitionOutGoing(tx, lpars)
	if err != nil {
		log.GLogger.Error("create low_price_article_requisition err:%s", err.Error())
		tx.Rollback()
		l.ErrorOK(MsgServerErr)
	}

	lowPriceArticle.OutgoingQuantity += lpar.Quantity

	// 主题的借出数量添加
	err = tx.Save(&lowPriceArticle).Error
	if err != nil {
		log.GLogger.Error("save low_price_article_requisition err:%s", err.Error())
		tx.Rollback()
		l.ErrorOK(MsgServerErr)
	}
	tx.Commit()

	l.Correct("")
}

// @Title 易耗品归还
// @Description 易耗品归还
// @Param	id	    query	string	true	"需归还易耗品id"
// @Param	status	query	string	true	"归还状态(ok,scrap)"
// @Success 200 {string} ""
// @Failure 500 server internal err
// @router /:id/return [post]
func (l *LowPriceArticleRequisitionController) Return() {
	uID, _ := l.GetInt("userID", 0)
	uName := l.GetString("userName")
	id := l.GetString(":id")
	if id == "" {
		l.ErrorOK("need id")
	}
	status := l.GetString("status")
	if status == "" {
		l.ErrorOK("need status")
	}
	tx := services.Slave().Begin()
	var lpar oa.LowPriceArticleRequisition
	err := services.Slave().Where(id).Find(&lpar).Error
	if err != nil {
		log.GLogger.Error("get low_price_article_requisition err:%s", err.Error())
		tx.Rollback()
		l.ErrorOK(MsgServerErr)
	}
	log.GLogger.Info("lpars", lpar)
	if lpar.OperatorCategory != models.DeviceOutgoing {
		l.ErrorOK("类别错误, 无法归还")
	}
	if lpar.IsReturn == 1 {
		l.ErrorOK("已归还")
	}

	var lowPriceArticle oa.LowPriceArticle
	err = services.Slave().Where(lpar.LowPriceArticleID).Find(&lowPriceArticle).Error
	if err != nil {
		log.GLogger.Error("get low_price_article err:%s", err.Error())
		tx.Rollback()
		l.ErrorOK(MsgServerErr)
	}

	log.GLogger.Info("lowPriceArticle", lowPriceArticle)

	// 判断是否有需要归还的
	if lowPriceArticle.NeedReturn == 0 {
		l.ErrorOK("无需归还")
	}

	newLowPriceArticleRequisition := &oa.LowPriceArticleRequisition{
		LowPriceArticleID:     lpar.LowPriceArticleID,
		OperatorID:            uID,
		OperatorName:          uName,
		AssociateEmployeeID:   lpar.AssociateEmployeeID,
		AssociateEmployeeName: lpar.AssociateEmployeeName,
		OperatorCategory:      models.DeviceReturn,
		Quantity:              1,
		Comment:               lpar.Comment,
	}

	// 添加一条归还记录
	err = tx.Create(&newLowPriceArticleRequisition).Error
	if err != nil {
		log.GLogger.Error("create low_price_article_requisition err:%s", err.Error())
		tx.Rollback()
		l.ErrorOK(MsgServerErr)
	}

	// 更新为已归还
	lpar.IsReturn = 1
	err = tx.Save(&lpar).Error
	if err != nil {
		log.GLogger.Error("save low_price_article_requisition err:%s", err.Error())
		tx.Rollback()
		l.ErrorOK(MsgServerErr)
	}

	if status == models.DeviceScrap {
		// 添加一条报废记录
		baofei := &oa.LowPriceArticleRequisition{
			LowPriceArticleID: lpar.LowPriceArticleID,
			OperatorID:        uID,
			OperatorName:      uName,
			OperatorCategory:  models.DeviceScrap,
			Quantity:          1,
			Comment:           lpar.Comment,
		}
		// 添加一条报废记录
		err = tx.Create(&baofei).Error
		if err != nil {
			log.GLogger.Error("create low_price_article_requisition err:%s", err.Error())
			tx.Rollback()
			l.ErrorOK(MsgServerErr)
		}

		// 主体的报废数量+1
		lowPriceArticle.ScrapQuantity += 1
	}

	// 主体的借出数量-1
	lowPriceArticle.OutgoingQuantity -= 1
	err = tx.Save(&lowPriceArticle).Error
	if err != nil {
		log.GLogger.Error("save low_price_article err:%s", err.Error())
		tx.Rollback()
		l.ErrorOK(MsgServerErr)
	}

	tx.Commit()

	l.Correct("")
}

// @Title 易耗品批量归还
// @Description 易耗品批量归还
// @Param	ids	    query	string	true	"需归还易耗品ids"
// @Success 200 {string} ""
// @Failure 500 server internal err
// @router /return/batch [post]
func (l *LowPriceArticleRequisitionController) BatchReturn() {
	uID, _ := l.GetInt("userID", 0)
	uName := l.GetString("userName")
	ids := l.GetString("ids")
	idList := strings.Split(ids, ",")
	log.GLogger.Info("uID :%d", uID)
	log.GLogger.Info("idList :%s", idList)
	tx := services.Slave().Begin()
	var lpars []*oa.LowPriceArticleRequisition
	var lpa oa.LowPriceArticle
	for _, id := range idList {
		var lpar oa.LowPriceArticleRequisition
		err := services.Slave().Where(id).Find(&lpar).Error
		if err != nil {
			log.GLogger.Error("get low_price_article_requisition err:%s", err.Error())
			tx.Rollback()
			l.ErrorOK(MsgServerErr)
		}
		log.GLogger.Info("lpars", lpar)
		if lpar.OperatorCategory != models.DeviceOutgoing {
			l.ErrorOK("归还易耗品中含有类别错误, 无法归还")
		}
		if lpar.IsReturn == 1 {
			l.ErrorOK("归还易耗品中含有已归还")
		}

		var lowPriceArticle oa.LowPriceArticle
		err = services.Slave().Where(lpar.LowPriceArticleID).Find(&lowPriceArticle).Error
		if err != nil {
			log.GLogger.Error("get low_price_article err:%s", err.Error())
			tx.Rollback()
			l.ErrorOK(MsgServerErr)
		}

		log.GLogger.Info("lowPriceArticle", lowPriceArticle)

		// 判断是否有需要归还的
		if lowPriceArticle.NeedReturn == 0 {
			l.ErrorOK("归还易耗品中含有无需归还")
		}
		lpa = lowPriceArticle
		lpars = append(lpars, &lpar)
	}

	log.GLogger.Info("lpars", lpars)

	// 批量插入归还记录
	err := oa.BatchRequisitionReturn(tx, uID, uName, lpars)
	if err != nil {
		log.GLogger.Error("BatchRequisitionReturn err:%s", err.Error())
		tx.Rollback()
		l.ErrorOK(MsgServerErr)
	}

	// 批量更新为已归还
	err = tx.Model(&oa.LowPriceArticleRequisition{}).Where("id in (?)", ids).Update("is_return", 1).Error
	if err != nil {
		log.GLogger.Error("update is_return err:%s", err.Error())
		tx.Rollback()
		l.ErrorOK(MsgServerErr)
	}
	// 主体的借出数量-1
	lpa.OutgoingQuantity -= len(idList)
	err = tx.Save(&lpa).Error
	if err != nil {
		log.GLogger.Error("save low_price_article err:%s", err.Error())
		tx.Rollback()
		l.ErrorOK(MsgServerErr)
	}

	tx.Commit()
	l.Correct("")
}

// @Title 易耗品报废
// @Description 易耗品报废
// @Param	body body oa.LowPriceArticleRequisition true "易耗品记录"
// @Success 200 {string} ""
// @Failure 500 server internal err
// @router /scrap [post]
func (l *LowPriceArticleRequisitionController) Scrap() {
	// 添加一条归还记录
	uID, _ := l.GetInt("userID", 0)
	uName := l.GetString("userName")
	lpar := new(oa.LowPriceArticleRequisition)
	err := json.Unmarshal(l.Ctx.Input.RequestBody, lpar)
	if err != nil {
		log.GLogger.Error("parse low_price_article_requisition err:%s", err.Error())
		l.ErrorOK(MsgInvalidParam)
	}
	log.GLogger.Info("lpars", lpar)

	if lpar.Quantity < 0 {
		l.ErrorOK("数量需要大于0")
	}
	tx := services.Slave().Begin()

	var lowPriceArticle oa.LowPriceArticle
	err = services.Slave().Where(lpar.LowPriceArticleID).Find(&lowPriceArticle).Error
	if err != nil {
		log.GLogger.Error("get low_price_article err:%s", err.Error())
		tx.Rollback()
		l.ErrorOK(MsgServerErr)
	}

	log.GLogger.Info("lowPriceArticle", lowPriceArticle)

	if lpar.Quantity > (lowPriceArticle.TotalQuantity - lowPriceArticle.OutgoingQuantity - lowPriceArticle.ScrapQuantity) {
		l.ErrorOK("报废数量不足")
	}

	// 判断报废数量
	newLowPriceArticleRequisition := &oa.LowPriceArticleRequisition{
		LowPriceArticleID: lpar.LowPriceArticleID,
		OperatorID:        uID,
		OperatorName:      uName,
		OperatorCategory:  models.DeviceScrap,
		Quantity:          lpar.Quantity,
		Comment:           lpar.Comment,
	}

	err = tx.Create(&newLowPriceArticleRequisition).Error
	if err != nil {
		log.GLogger.Error("create low_price_article_requisition err:%s", err.Error())
		tx.Rollback()
		l.ErrorOK(MsgServerErr)
	}

	// 主体的报废数量
	lowPriceArticle.ScrapQuantity += lpar.Quantity
	err = tx.Save(&lowPriceArticle).Error
	if err != nil {
		log.GLogger.Error("save low_price_article err:%s", err.Error())
		tx.Rollback()
		l.ErrorOK(MsgServerErr)
	}

	tx.Commit()
	l.Correct("")
}
