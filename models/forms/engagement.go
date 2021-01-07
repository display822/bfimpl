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

type EngagementResult struct {
	CostSummary  float64 `json:"cost_summary"`
	HourSummary  int     `json:"hour_summary"`
	EmployeeNums int     `json:"employee_nums"`
	List         []E     `json:"list"`
}

type E struct {
	EngagementCode string       `json:"engagement_code"`
	CostSummary    float64      `json:"cost_summary"`
	HourSummary    int          `json:"hour_summary"`
	EmployeeNums   int          `json:"employee_nums"`
	EngagementList []Engagement `json:"engagement_list"`
}
