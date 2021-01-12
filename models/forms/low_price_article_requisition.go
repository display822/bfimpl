/**
* @author : yi.zhang
* @description : forms 描述
* @date   : 2021-01-12 12:45
 */

package forms

import (
	"time"
)

type ReturnByCreatedAt []*Return

func (p ReturnByCreatedAt) Len() int           { return len(p) }
func (p ReturnByCreatedAt) Less(i, j int) bool { return p[i].CreatedAt.After(p[j].CreatedAt) }
func (p ReturnByCreatedAt) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type Return struct {
	LowPriceArticleID   int       `json:"low_price_article_id"`
	LowPriceArticleName string    `json:"low_price_article_name"`
	CreatedAt           time.Time `json:"created_at"`
}
