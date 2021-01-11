// @APIVersion 1.0.0
// @Title Test Api
// @Description WeTest API Project
// @Contact wetest@tencent.com
// @TermsOfServiceUrl http://wetest.qq.com
package routers

import (
	"bfimpl/controllers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
)

func init() {
	ns := beego.NewNamespace("/v1",

		beego.NSNamespace("/login", beego.NSInclude(&controllers.LoginController{})),
		beego.NSNamespace("/user", beego.NSInclude(&controllers.UserController{})),
		beego.NSNamespace("/client", beego.NSInclude(&controllers.ClientController{})),
		beego.NSNamespace("/service", beego.NSInclude(&controllers.ServiceController{})),
		beego.NSNamespace("/amount", beego.NSInclude(&controllers.AmountController{})),
		beego.NSNamespace("/task", beego.NSInclude(&controllers.TaskController{})),

		beego.NSNamespace("/employee", beego.NSInclude(&controllers.EmployeeController{})),
		beego.NSNamespace("/department", beego.NSInclude(&controllers.DepartmentController{})),
		beego.NSNamespace("/file", beego.NSInclude(&controllers.FileController{})),
		beego.NSNamespace("/work", beego.NSInclude(
			&controllers.WorkController{},
			&controllers.AttendanceController{})),
		beego.NSNamespace("/bench", beego.NSInclude(&controllers.BenchController{})),
		beego.NSNamespace("/project", beego.NSInclude(&controllers.ProjectController{})),
		beego.NSNamespace("/expense", beego.NSInclude(&controllers.ExpenseController{})),
		beego.NSNamespace("/engagement", beego.NSInclude(&controllers.EngagementController{})),
		beego.NSNamespace("/low_price_article", beego.NSInclude(&controllers.LowPriceArticleController{})),
		beego.NSNamespace("/low_price_article_requisition", beego.NSInclude(&controllers.LowPriceArticleRequisitionController{})),
		beego.NSNamespace("/device", beego.NSInclude(&controllers.DeviceController{})),
	)

	ips := beego.AppConfig.Strings("bkcors")
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		// AllowAllOrigins:  true,
		AllowOrigins: ips,
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Authorization", "Access-Control-Allow-Origin",
			"Access-Control-Allow-Headers", "Content-Type", "x-csrf-token", "x-requested-with", "projectId"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
	}))
	beego.AddNamespace(ns)
}
