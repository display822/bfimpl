package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["bfimpl/controllers:ServiceController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:ServiceController"],
		beego.ControllerComments{
			Method: "AddService",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Filters: nil,
			Params: nil})

	beego.GlobalControllerRouter["bfimpl/controllers:ServiceController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:ServiceController"],
		beego.ControllerComments{
			Method: "UpdateService",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Filters: nil,
			Params: nil})

	beego.GlobalControllerRouter["bfimpl/controllers:ServiceController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:ServiceController"],
		beego.ControllerComments{
			Method: "SwitchService",
			Router: `/switch`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Filters: nil,
			Params: nil})

	beego.GlobalControllerRouter["bfimpl/controllers:UserController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:UserController"],
		beego.ControllerComments{
			Method: "AddUser",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Filters: nil,
			Params: nil})

	beego.GlobalControllerRouter["bfimpl/controllers:UserController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:UserController"],
		beego.ControllerComments{
			Method: "GroupLeaders",
			Router: `/leaders`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Filters: nil,
			Params: nil})

}
