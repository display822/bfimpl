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

		beego.NSNamespace("/user", beego.NSInclude(&controllers.UserController{})),
		beego.NSNamespace("/client", beego.NSInclude(&controllers.ClientController{})),
		beego.NSNamespace("/service", beego.NSInclude(&controllers.ServiceController{})),
		beego.NSNamespace("/amount", beego.NSInclude(&controllers.AmountController{})),
		beego.NSNamespace("/task", beego.NSInclude(&controllers.TaskController{})),
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
