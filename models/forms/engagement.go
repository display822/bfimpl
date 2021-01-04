/**
* @author : yi.zhang
* @description : forms 描述
* @date   : 2020-12-30 16:32
 */

package forms

type Engagement struct {
	EngagementCode string         `json:"engagement_code"`
	EmployeeName   string         `json:"employee_name"`
	DateField      map[string]int `json:"date_field"`
}
