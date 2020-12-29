/**
* @author : yi.zhang
* @description : forms 描述
* @date   : 2020-12-30 16:32
 */

package forms

type Engagement struct {
	EngagementCode string
	EmployeeCount  int
	EmployeeHour   int
	EngagementCost float64
	DataField      []string
	Employees      []Employee
}

type Employee struct {
	EmployeeName   string
	EmployeeCount  int
	EmployeeHour   int
	EngagementCost float64
	Employees      []Employee
}
