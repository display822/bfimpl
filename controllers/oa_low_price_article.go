/**
* @author : yi.zhang
* @description : controllers 描述
* @date   : 2021-01-11 15:46
 */

package controllers

import (
	"bfimpl/models"
	"bfimpl/models/oa"
	"bfimpl/services"
	"bfimpl/services/log"
	"encoding/json"
	"time"

	"github.com/jinzhu/gorm"
)

type LowPriceArticleController struct {
	BaseController
}

// @Title 创建易耗品
// @Description 创建易耗品
// @Param	body body oa.LowPriceArticle true "易耗品"
// @Success 200 {object} oa.LowPriceArticle
// @Failure 500 server internal err
// @router / [post]
func (l *LowPriceArticleController) Create() {
	// 验证员工身份 (7，8，9)
	userID, _ := l.GetInt("userID")
	userType, _ := l.GetInt("userType")
	if userType != models.UserIT && userType != models.UserFront && userType != models.UserFinance {
		l.ErrorOK("没有权限")
	}

	param := new(oa.LowPriceArticle)
	err := json.Unmarshal(l.Ctx.Input.RequestBody, param)
	if err != nil {
		log.GLogger.Error("parse low_price_article err:%s", err.Error())
		l.ErrorOK(MsgInvalidParam)
	}
	log.GLogger.Info("param :%+v", param)
	param.IngoingOperatorID = userID
	param.IngoingTime = models.Time(time.Now())

	_, ok := oa.LowPriceArticleMap[param.LowPriceArticleCategory]
	if !ok {
		l.ErrorOK("param low_price_article_category error")
	}

	tx := services.Slave().Begin()
	err = tx.Create(&param).Error
	if err != nil {
		log.GLogger.Error("create low_price_article err:%s", err.Error())
		tx.Rollback()
		l.ErrorOK(MsgServerErr)
	}

	requisition := oa.LowPriceArticleRequisition{
		LowPriceArticleID: int(param.ID),
		OperatorID:        userID,
		OperatorCategory:  models.DeviceIngoing,
		Quantity:          param.TotalQuantity,
	}

	err = tx.Create(&requisition).Error
	if err != nil {
		log.GLogger.Error("create low_price_article_requisition err:%s", err.Error())
		tx.Rollback()
		l.ErrorOK(MsgServerErr)
	}
	tx.Commit()

	l.Correct(param)
}

// @Title 易耗品列表
// @Description 易耗品列表
// @Param	pagenum	    query	int	false	"页码"
// @Param	pagesize	query	int	false	"页数"
// @Param	category	query	bool	false	"类别"
// @Param	keyword	query	bool	false	"关键词"
// @Success 200 {object} []oa.LowPriceArticle
// @Failure 500 server internal err
// @router / [get]
func (l *LowPriceArticleController) List() {
	pageNum, _ := l.GetInt("pagenum", 1)
	pageSize, _ := l.GetInt("pagesize", 10)
	category := l.GetString("category")
	keyword := l.GetString("keyword")
	log.GLogger.Info("category:%d;keyword:%s", category, keyword)
	var resp struct {
		Total int                   `json:"total"`
		List  []*oa.LowPriceArticle `json:"list"`
	}
	db := services.Slave()
	if category != "" {
		db = db.Where("low_price_article_category = ?", category)
	}
	if keyword != "" {

	}
	var lpa []*oa.LowPriceArticle
	db.Limit(pageSize).Offset((pageNum - 1) * pageSize).Order("created_at desc").Find(&lpa).Limit(-1).Offset(-1).Count(&resp.Total)
	resp.List = lpa
	l.Correct(resp)
}

// @Title 易耗品详情
// @Description 易耗品详情
// @Success 200 {object} oa.LowPriceArticle
// @Failure 500 server internal err
// @router /:id [get]
func (l *LowPriceArticleController) Get() {
	id, _ := l.GetInt(":id")
	log.GLogger.Info("id:%s", id)
	var lpa oa.LowPriceArticle
	services.Slave().Where("id = ?", id).Find(&lpa)
	l.Correct(lpa)
}

// @Title 员工下易耗品借出列表
// @Description 易耗品借出列表
// @Param	pagenum	    query	int	false	"页码"
// @Param	pagesize	query	int	false	"页数"
// @Success 200 {object} oa.LowPriceArticle
// @Failure 500 server internal err
// @router /employee/outgoing [get]
func (l *LowPriceArticleController) ListOutgoingByEmployee() {
	uID, _ := l.GetInt("userID")

	log.GLogger.Info("uID:%d", uID)

	// 先获取易耗品id
	var articleRequisitions []*oa.LowPriceArticleRequisition
	services.Slave().Where("operator_category = ?", models.DeviceOutgoing).
		Where("is_return = ?", 0).
		Where("associate_employee_id = ?", uID).
		Find(&articleRequisitions)

	var ids []int
	for _, a := range articleRequisitions {
		ids = append(ids, a.LowPriceArticleID)
	}

	// 再查询
	var articles []*oa.LowPriceArticle
	services.Slave().Where("id in (?)", ids).
		Preload("LowPriceArticleRequisitions", func(db *gorm.DB) *gorm.DB {
			return db.Where("operator_category = ?", models.DeviceOutgoing).
				Where("is_return = ?", 0).
				Where("associate_employee_id = ?", uID).Order("created_at desc")
			// return db.
		}).
		Find(&articles)

	l.Correct(articles)
}

// // @Title 员工下易耗品归还列表
// // @Description 员工下易耗品归还列表
// // @Success 200 {string} ""
// // @Failure 500 server internal err
// // @router /employee/return [get]
// func (l *LowPriceArticleController) ListReturnByEmployee() {
// 	eId := l.GetString("eid")
// 	if eId == "" {
// 		l.ErrorOK("need eid")
// 	}
// 	log.GLogger.Info("eId:%d", eId)
//
// 	// 先获取易耗品id
// 	var articleRequisitions []*oa.LowPriceArticleRequisition
// 	var returns []*forms.Return
// 	services.Slave().Where("operator_category = ?", models.DeviceReturn).
// 		Where("associate_employee_id = ?", eId).
// 		Preload("LowPriceArticle").
// 		Find(&articleRequisitions)
//
// 	log.GLogger.Info("articleRequisitions:%d", articleRequisitions)
//
// 	for _, a := range articleRequisitions {
// 		returns = append(returns, &forms.Return{
// 			LowPriceArticleID:   a.LowPriceArticleID,
// 			LowPriceArticleName: a.LowPriceArticle.LowPriceArticleName,
// 			CreatedAt:           a.CreatedAt,
// 		})
// 	}
// 	sort.Sort(forms.ReturnByCreatedAt(returns))
//
// 	l.Correct(returns)
// }
