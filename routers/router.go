// @APIVersion 1.0.0
// @Title Test Api
// @Description WeTest API Project
// @Contact wetest@tencent.com
// @TermsOfServiceUrl http://wetest.qq.com
package routers

import (
	"bfimpl/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",

		beego.NSNamespace("/user", beego.NSInclude(&controllers.UserController{})),

	)
	beego.AddNamespace(ns)
}
