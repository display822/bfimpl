package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["bfimpl/controllers:AmountController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:AmountController"],
		beego.ControllerComments{
			Method:           "AddAmount",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:AmountController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:AmountController"],
		beego.ControllerComments{
			Method:           "DelayInAmount",
			Router:           `/delay/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:AmountController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:AmountController"],
		beego.ControllerComments{
			Method:           "GetAmounts",
			Router:           `/list`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:AmountController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:AmountController"],
		beego.ControllerComments{
			Method:           "GetAmountLogs",
			Router:           `/log`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:AmountController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:AmountController"],
		beego.ControllerComments{
			Method:           "GetTaskAmountLogs",
			Router:           `/tasklog`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:AmountController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:AmountController"],
		beego.ControllerComments{
			Method:           "SwitchAmount",
			Router:           `/switch`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:ClientController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:ClientController"],
		beego.ControllerComments{
			Method:           "AddClient",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:ClientController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:ClientController"],
		beego.ControllerComments{
			Method:           "UpdateClient",
			Router:           `/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:ClientController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:ClientController"],
		beego.ControllerComments{
			Method:           "GetClient",
			Router:           `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:ClientController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:ClientController"],
		beego.ControllerComments{
			Method:           "GetClients",
			Router:           `/list`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:LoginController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:LoginController"],
		beego.ControllerComments{
			Method:           "Login",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:ServiceController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:ServiceController"],
		beego.ControllerComments{
			Method:           "AddService",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:ServiceController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:ServiceController"],
		beego.ControllerComments{
			Method:           "UpdateService",
			Router:           `/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:ServiceController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:ServiceController"],
		beego.ControllerComments{
			Method:           "GetServices",
			Router:           `/list`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:ServiceController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:ServiceController"],
		beego.ControllerComments{
			Method:           "SwitchService",
			Router:           `/switch`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:TaskController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:TaskController"],
		beego.ControllerComments{
			Method:           "NewTask",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:TaskController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:TaskController"],
		beego.ControllerComments{
			Method:           "Task",
			Router:           `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:TaskController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:TaskController"],
		beego.ControllerComments{
			Method:           "AssignTask",
			Router:           `/assign/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:TaskController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:TaskController"],
		beego.ControllerComments{
			Method:           "TaskBackAmount",
			Router:           `/backamount/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:TaskController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:TaskController"],
		beego.ControllerComments{
			Method:           "CancelTask",
			Router:           `/cancel/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:TaskController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:TaskController"],
		beego.ControllerComments{
			Method:           "CommentTask",
			Router:           `/comment/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:TaskController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:TaskController"],
		beego.ControllerComments{
			Method:           "TaskComments",
			Router:           `/comment/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:TaskController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:TaskController"],
		beego.ControllerComments{
			Method:           "ConfirmTask",
			Router:           `/confirm/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:TaskController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:TaskController"],
		beego.ControllerComments{
			Method:           "TaskDashboard",
			Router:           `/dashboard`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:TaskController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:TaskController"],
		beego.ControllerComments{
			Method:           "EndTask",
			Router:           `/end/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:TaskController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:TaskController"],
		beego.ControllerComments{
			Method:           "ExecuteTask",
			Router:           `/execute/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:TaskController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:TaskController"],
		beego.ControllerComments{
			Method:           "GetTaskExeInfo",
			Router:           `/exeinfo/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:TaskController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:TaskController"],
		beego.ControllerComments{
			Method:           "FinishTask",
			Router:           `/finish/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:TaskController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:TaskController"],
		beego.ControllerComments{
			Method:           "TaskToday",
			Router:           `/focus`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:TaskController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:TaskController"],
		beego.ControllerComments{
			Method:           "FrozenTask",
			Router:           `/frozen/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:TaskController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:TaskController"],
		beego.ControllerComments{
			Method:           "TaskImportant",
			Router:           `/high`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:TaskController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:TaskController"],
		beego.ControllerComments{
			Method:           "TaskList",
			Router:           `/list`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:TaskController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:TaskController"],
		beego.ControllerComments{
			Method:           "PauseTask",
			Router:           `/pause/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:TaskController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:TaskController"],
		beego.ControllerComments{
			Method:           "SaveTaskDetail",
			Router:           `/save/:id`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:TaskController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:TaskController"],
		beego.ControllerComments{
			Method:           "TaskTags",
			Router:           `/tags`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:UserController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:UserController"],
		beego.ControllerComments{
			Method:           "AddUser",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:UserController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:UserController"],
		beego.ControllerComments{
			Method:           "Implementers",
			Router:           `/impls`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:UserController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:UserController"],
		beego.ControllerComments{
			Method:           "GroupLeaders",
			Router:           `/leaders`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:UserController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:UserController"],
		beego.ControllerComments{
			Method:           "UserList",
			Router:           `/list`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
