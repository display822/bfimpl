/**
* @author : yi.zhang
* @description : oa 描述
* @date   : 2021-01-11 15:19
 */

package oa

import (
	"bfimpl/models"

	"github.com/jinzhu/gorm"
)

var LowPriceArticleMap = map[string]struct{}{
	"Office": {},
	"IT":     {},
}

// LowPriceArticle 低值易耗品表
type LowPriceArticle struct {
	gorm.Model
	LowPriceArticleCategory     string                        `gorm:"size:20;not null;comment:'(Office,IT)'" json:"low_price_article_category"`
	LowPriceArticleName         string                        `gorm:"size:100;not null;comment:'物品名称'" json:"low_price_article_name"`
	Brand                       string                        `gorm:"size:60;not null;comment:'品牌'" json:"brand"`
	Retailer                    string                        `gorm:"size:100;not null;comment:'零售商'" json:"retailer"`
	Site                        string                        `gorm:"size:100;not null;comment:'位置'" json:"site"`
	PurchasePrice               float64                       `gorm:"type:decimal(10,2);not null;comment:'购买价格'" json:"purchase_price"`
	IngoingOperatorID           int                           `gorm:"size:10;not null;comment:'入库人id'" json:"ingoing_operator_id"`
	IngoingOperatorName         string                        `gorm:"not null;comment:'入库人Name'" json:"ingoing_operator_name"`
	IngoingTime                 models.Time                   `gorm:"type:datetime;not null;comment:'入库时间'" json:"ingoing_time"`
	TotalQuantity               int                           `gorm:"not null;comment:'物品数量'" json:"total_quantity"`
	OutgoingQuantity            int                           `gorm:"not null;comment:'物品借出数量'" json:"outgoing_quantity"`
	ScrapQuantity               int                           `gorm:"not null;comment:'物品报废数量'" json:"scrap_quantity"`
	NeedReturn                  int                           `gorm:"size:10;not null;comment:'需归还'" json:"need_return"`
	Comment                     string                        `gorm:"size:2000;not null;comment:'备注'" json:"comment"`
	LowPriceArticleRequisitions []*LowPriceArticleRequisition `json:"low_price_article_requisitions"`
}
