package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["bfimpl/controllers:DemoController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:DemoController"],
        beego.ControllerComments{
            Method: "HttpRequestDemo",
            Router: `/http/get`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["bfimpl/controllers:DemoController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:DemoController"],
        beego.ControllerComments{
            Method: "GetMySQLInfo",
            Router: `/mysql/tables`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["bfimpl/controllers:DemoController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:DemoController"],
        beego.ControllerComments{
            Method: "SetRedis",
            Router: `/redis/string`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["bfimpl/controllers:DemoController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:DemoController"],
        beego.ControllerComments{
            Method: "GetRedis",
            Router: `/redis/string/:key`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["bfimpl/controllers:DemoController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:DemoController"],
        beego.ControllerComments{
            Method: "GetUser",
            Router: `/user/:uid`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["bfimpl/controllers:DemoController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:DemoController"],
        beego.ControllerComments{
            Method: "GetAllUsers",
            Router: `/users`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
