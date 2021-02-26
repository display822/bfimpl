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
			Method:           "SwitchAmount",
			Router:           `/switch`,
			AllowHTTPMethods: []string{"put"},
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
			Method:           "GetAllAmounts",
			Router:           `/history`,
			AllowHTTPMethods: []string{"get"},
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
			Method:           "ChangeFinish",
			Router:           `/change/:id`,
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
			Method:           "TaskHistory",
			Router:           `/history/:id`,
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
			Method:           "PutUser",
			Router:           `/:id`,
			AllowHTTPMethods: []string{"put"},
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

	//oa=============
	beego.GlobalControllerRouter["bfimpl/controllers:DepartmentController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:DepartmentController"],
		beego.ControllerComments{
			Method:           "GetDepartments",
			Router:           `/list`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:DepartmentController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:DepartmentController"],
		beego.ControllerComments{
			Method:           "GetLevels",
			Router:           `/level/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:DepartmentController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:DepartmentController"],
		beego.ControllerComments{
			Method:           "GetServiceLine",
			Router:           `/service/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:EmployeeController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:EmployeeController"],
		beego.ControllerComments{
			Method:           "NewEmpEntry",
			Router:           `/new`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:EmployeeController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:EmployeeController"],
		beego.ControllerComments{
			Method:           "DeleteEmp",
			Router:           `/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:EmployeeController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:EmployeeController"],
		beego.ControllerComments{
			Method:           "GetEmpEntryInfo",
			Router:           `/entry/detail/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:EmployeeController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:EmployeeController"],
		beego.ControllerComments{
			Method:           "GetEmpEntryList",
			Router:           `/entry/list`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:EmployeeController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:EmployeeController"],
		beego.ControllerComments{
			Method:           "GetWorkflowNode",
			Router:           `/workflow/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:EmployeeController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:EmployeeController"],
		beego.ControllerComments{
			Method:           "CommitWorkflowNode",
			Router:           `/workflow/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:EmployeeController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:EmployeeController"],
		beego.ControllerComments{
			Method:           "NewEmpLeave",
			Router:           `/leave/:id`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:EmployeeController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:EmployeeController"],
		beego.ControllerComments{
			Method:           "CommitLeaveInfoNode",
			Router:           `/leave/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:EmployeeController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:EmployeeController"],
		beego.ControllerComments{
			Method:           "GetLeaveInfo",
			Router:           `/leave/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:EmployeeController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:EmployeeController"],
		beego.ControllerComments{
			Method:           "SaveEmpInfo",
			Router:           `/save`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:EmployeeController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:EmployeeController"],
		beego.ControllerComments{
			Method:           "GetEmpInfo",
			Router:           `/detail/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:EmployeeController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:EmployeeController"],
		beego.ControllerComments{
			Method:           "CreateEmpContract",
			Router:           `/contract/:id`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:EmployeeController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:EmployeeController"],
		beego.ControllerComments{
			Method:           "GetContracts",
			Router:           `/contracts`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:EmployeeController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:EmployeeController"],
		beego.ControllerComments{
			Method:           "GetEmpContracts",
			Router:           `/contracts/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:EmployeeController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:EmployeeController"],
		beego.ControllerComments{
			Method:           "UpdateContract",
			Router:           `/contract`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:EmployeeController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:EmployeeController"],
		beego.ControllerComments{
			Method:           "GetContract",
			Router:           `/contract/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:EmployeeController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:EmployeeController"],
		beego.ControllerComments{
			Method:           "DelContract",
			Router:           `/contract/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:EmployeeController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:EmployeeController"],
		beego.ControllerComments{
			Method:           "UpdateEmpStatus",
			Router:           `/status/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:EmployeeController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:EmployeeController"],
		beego.ControllerComments{
			Method:           "SearchEmp",
			Router:           `/search`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:EmployeeController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:EmployeeController"],
		beego.ControllerComments{
			Method:           "MoreContract",
			Router:           `/contract/continue`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:FileController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:FileController"],
		beego.ControllerComments{
			Method:           "Upload",
			Router:           `/upload`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:FileController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:FileController"],
		beego.ControllerComments{
			Method:           "SocialSecurityList",
			Router:           `/sslist`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:WorkController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:WorkController"],
		beego.ControllerComments{
			Method:           "ReqOvertime",
			Router:           `/overtime`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:WorkController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:WorkController"],
		beego.ControllerComments{
			Method:           "ApprovalUsers",
			Router:           `/approvals`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:WorkController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:WorkController"],
		beego.ControllerComments{
			Method:           "OvertimeById",
			Router:           `/overtime/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:WorkController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:WorkController"],
		beego.ControllerComments{
			Method:           "OvertimeList",
			Router:           `/overtime`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:WorkController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:WorkController"],
		beego.ControllerComments{
			Method:           "ApprovalOvertime",
			Router:           `/overtime`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:WorkController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:WorkController"],
		beego.ControllerComments{
			Method:           "DeleteOvertime",
			Router:           `/overtime/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:WorkController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:WorkController"],
		beego.ControllerComments{
			Method:           "ValidOvertime",
			Router:           `/overtime/:id/check`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:WorkController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:WorkController"],
		beego.ControllerComments{
			Method:           "ReqLeave",
			Router:           `/leave`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:WorkController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:WorkController"],
		beego.ControllerComments{
			Method:           "LeaveApprovalUsers",
			Router:           `/leave/approvals`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:WorkController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:WorkController"],
		beego.ControllerComments{
			Method:           "LeaveById",
			Router:           `/leave/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:WorkController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:WorkController"],
		beego.ControllerComments{
			Method:           "LeaveList",
			Router:           `/leave`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:WorkController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:WorkController"],
		beego.ControllerComments{
			Method:           "ApprovalLeave",
			Router:           `/leave`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:WorkController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:WorkController"],
		beego.ControllerComments{
			Method:           "DeleteLeave",
			Router:           `/leave/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:WorkController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:WorkController"],
		beego.ControllerComments{
			Method:           "ValidLeave",
			Router:           `/leave/:id/check`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:WorkController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:WorkController"],
		beego.ControllerComments{
			Method:           "GetProjects",
			Router:           `/overtime/projects`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:WorkController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:WorkController"],
		beego.ControllerComments{
			Method:           "LeaveListByDate",
			Router:           `/leavebydate`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:WorkController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:WorkController"],
		beego.ControllerComments{
			Method:           "RemainHoliday",
			Router:           `/remain/holiday`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:AttendanceController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:AttendanceController"],
		beego.ControllerComments{
			Method:           "ConfirmUserAttendance",
			Router:           `/attendance`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:AttendanceController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:AttendanceController"],
		beego.ControllerComments{
			Method:           "UploadAttendanceTmp",
			Router:           `/attendance/tmp`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:AttendanceController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:AttendanceController"],
		beego.ControllerComments{
			Method:           "GetAttendanceUserByDept",
			Router:           `/attendance/deptuser`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:AttendanceController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:AttendanceController"],
		beego.ControllerComments{
			Method:           "DeleteAttendanceTmp",
			Router:           `/attendance/tmp/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:AttendanceController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:AttendanceController"],
		beego.ControllerComments{
			Method:           "UpdateAttendanceTmp",
			Router:           `/attendance/tmp`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:AttendanceController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:AttendanceController"],
		beego.ControllerComments{
			Method:           "GetUserAttendanceTmps",
			Router:           `/attendance/tmp`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:AttendanceController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:AttendanceController"],
		beego.ControllerComments{
			Method:           "GetUserAttendanceTmpByDay",
			Router:           `/attendance/tmp/day`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:AttendanceController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:AttendanceController"],
		beego.ControllerComments{
			Method:           "GetUserAttendanceByMonth",
			Router:           `/attendance`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:AttendanceController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:AttendanceController"],
		beego.ControllerComments{
			Method:           "UpdateAttendance",
			Router:           `/attendance`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:AttendanceController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:AttendanceController"],
		beego.ControllerComments{
			Method:           "ExportPos",
			Router:           `/attendance/pos`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:AttendanceController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:AttendanceController"],
		beego.ControllerComments{
			Method:           "ExportData",
			Router:           `/attendance/data`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:AttendanceController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:AttendanceController"],
		beego.ControllerComments{
			Method:           "CreateAttendanceTmp",
			Router:           `/attendance/create/tmp`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:BenchController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:BenchController"],
		beego.ControllerComments{
			Method:           "GetMyApprove",
			Router:           `/myapprove`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:BenchController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:BenchController"],
		beego.ControllerComments{
			Method:           "GetMyRequest",
			Router:           `/myreq`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:ExpenseController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:ExpenseController"],
		beego.ControllerComments{
			Method:           "ParseDetailFile",
			Router:           `/details`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:ExpenseController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:ExpenseController"],
		beego.ControllerComments{
			Method:           "ApprovalUsers",
			Router:           `/approvals`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:ExpenseController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:ExpenseController"],
		beego.ControllerComments{
			Method:           "List",
			Router:           `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:ExpenseController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:ExpenseController"],
		beego.ControllerComments{
			Method:           "ReqExpense",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:ProjectController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:ProjectController"],
		beego.ControllerComments{
			Method:           "ParseProjectDetailFile",
			Router:           `/details`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:ExpenseController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:ExpenseController"],
		beego.ControllerComments{
			Method:           "GetProjects",
			Router:           `/projects`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:ProjectController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:ProjectController"],
		beego.ControllerComments{
			Method:           "FilterField",
			Router:           `/filter`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:ExpenseController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:ExpenseController"],
		beego.ControllerComments{
			Method:           "ApprovalExpense",
			Router:           `/`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:ExpenseController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:ExpenseController"],
		beego.ControllerComments{
			Method:           "PaidExpense",
			Router:           `/paid`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:ExpenseController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:ExpenseController"],
		beego.ControllerComments{
			Method:           "BatchPaidExpense",
			Router:           `/paid/batch`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:ProjectController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:ProjectController"],
		beego.ControllerComments{
			Method:           "List",
			Router:           `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:ExpenseController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:ExpenseController"],
		beego.ControllerComments{
			Method:           "ExpenseById",
			Router:           `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:ProjectController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:ProjectController"],
		beego.ControllerComments{
			Method:           "Create",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:ExpenseController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:ExpenseController"],
		beego.ControllerComments{
			Method:           "PaidInfo",
			Router:           `/paid/info`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:ExpenseController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:ExpenseController"],
		beego.ControllerComments{
			Method:           "DebitCard",
			Router:           `/:id/debit_card`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:ExpenseController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:ExpenseController"],
		beego.ControllerComments{
			Method:           "ExportUnpaid",
			Router:           `/export/unpaid`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:EngagementController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:EngagementController"],
		beego.ControllerComments{
			Method:           "ParseEngagementDetailFile",
			Router:           `/details`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:EngagementController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:EngagementController"],
		beego.ControllerComments{
			Method:           "Create",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:EngagementController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:EngagementController"],
		beego.ControllerComments{
			Method:           "Valid",
			Router:           `/valid`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:EngagementController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:EngagementController"],
		beego.ControllerComments{
			Method:           "List",
			Router:           `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:EngagementController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:EngagementController"],
		beego.ControllerComments{
			Method:           "PeriodList",
			Router:           `/period`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:EngagementController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:EngagementController"],
		beego.ControllerComments{
			Method:           "PeriodDetail",
			Router:           `/period/detail`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:EngagementController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:EngagementController"],
		beego.ControllerComments{
			Method:           "Project",
			Router:           `/project`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:LowPriceArticleController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:LowPriceArticleController"],
		beego.ControllerComments{
			Method:           "Create",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:DeviceController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:DeviceController"],
		beego.ControllerComments{
			Method:           "Create",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:LowPriceArticleController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:LowPriceArticleController"],
		beego.ControllerComments{
			Method:           "List",
			Router:           `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:DeviceController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:DeviceController"],
		beego.ControllerComments{
			Method:           "List",
			Router:           `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:DeviceController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:DeviceController"],
		beego.ControllerComments{
			Method:           "ListOutgoingByEmployee",
			Router:           `/employee/outgoing`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:DeviceController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:DeviceController"],
		beego.ControllerComments{
			Method:           "ListReturnByEmployee",
			Router:           `/employee/return`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:LowPriceArticleController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:LowPriceArticleController"],
		beego.ControllerComments{
			Method:           "Get",
			Router:           `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:LowPriceArticleController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:LowPriceArticleController"],
		beego.ControllerComments{
			Method:           "ListOutgoingByEmployee",
			Router:           `/employee/outgoing`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:LowPriceArticleRequisitionController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:LowPriceArticleRequisitionController"],
		beego.ControllerComments{
			Method:           "List",
			Router:           `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:LowPriceArticleRequisitionController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:LowPriceArticleRequisitionController"],
		beego.ControllerComments{
			Method:           "Outgoing",
			Router:           `/outgoing`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:LowPriceArticleRequisitionController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:LowPriceArticleRequisitionController"],
		beego.ControllerComments{
			Method:           "Return",
			Router:           `/:id/return`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:DeviceController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:DeviceController"],
		beego.ControllerComments{
			Method:           "Put",
			Router:           `/`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:LowPriceArticleRequisitionController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:LowPriceArticleRequisitionController"],
		beego.ControllerComments{
			Method:           "BatchReturn",
			Router:           `/return/batch`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:LowPriceArticleRequisitionController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:LowPriceArticleRequisitionController"],
		beego.ControllerComments{
			Method:           "Scrap",
			Router:           `/scrap`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:DeviceController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:DeviceController"],
		beego.ControllerComments{
			Method:           "Get",
			Router:           `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:DeviceController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:DeviceController"],
		beego.ControllerComments{
			Method:           "ReqDevice",
			Router:           `/apply`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:DeviceController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:DeviceController"],
		beego.ControllerComments{
			Method:           "ListApply",
			Router:           `/apply`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:DeviceController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:DeviceController"],
		beego.ControllerComments{
			Method:           "ApprovalDevice",
			Router:           `/apply/approval`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:DeviceController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:DeviceController"],
		beego.ControllerComments{
			Method:           "OutgoingDevice",
			Router:           `/:id/outgoing`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:DeviceController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:DeviceController"],
		beego.ControllerComments{
			Method:           "RevokeDevice",
			Router:           `/apply/:id/revoke`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:DeviceController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:DeviceController"],
		beego.ControllerComments{
			Method:           "ListDeviceApply",
			Router:           `/:id/apply`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:DeviceController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:DeviceController"],
		beego.ControllerComments{
			Method:           "ReceiveDevice",
			Router:           `/:id/receive`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:DeviceController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:DeviceController"],
		beego.ControllerComments{
			Method:           "DistributionDevice",
			Router:           `/:id/distribution`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:DeviceController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:DeviceController"],
		beego.ControllerComments{
			Method:           "ReturnDevice",
			Router:           `/:id/return`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:DeviceController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:DeviceController"],
		beego.ControllerComments{
			Method:           "GetDistributionProjects",
			Router:           `/distribution/projects`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:DeviceController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:DeviceController"],
		beego.ControllerComments{
			Method:           "GetApply",
			Router:           `/apply/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:DeviceController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:DeviceController"],
		beego.ControllerComments{
			Method:           "ApplyInfo",
			Router:           `/apply/info`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:DeviceController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:DeviceController"],
		beego.ControllerComments{
			Method:           "GetProjects",
			Router:           `/projects`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bfimpl/controllers:DeviceController"] = append(beego.GlobalControllerRouter["bfimpl/controllers:DeviceController"],
		beego.ControllerComments{
			Method:           "GetDistributionProjects",
			Router:           `/distribution/projects`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})
}
