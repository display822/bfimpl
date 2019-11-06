# api-project-go
api-project-go是WeTest用于快速开发后台接口的框架。

## 版本声明

* 重点:
    * 虽然使用go modules，但是beego的脚手架命令bee，要求项目必须在GOPATH下创建运行
    * 但是go module在GOPATH下默认不开启，需要导出变量
```bash
export GO111MODULE=on
# 或者
source env.sh
```

* go > 1.12（支持go modules）
* beego >= 1.11.1 
* Goland IDE请激活vgo模式 

## 初始安装环境
```shell
git clone [this_repo_url] api-project-go
mv api-project-go [your_app_name]
go run tool.go [your_app_name]
cd your_app_name && rm -rf .git
go build
```
ps: your app name代表你的项目名称，上面一定是相同的。否则vgo无法使用

## 根据项目名字修改配置
* conf/app.conf
* appname = apidemo ，修改为服务名称
* servicename = wetest，修改为业务名称，比如cloud/fsight等


## 定制的一些东西

* 增加了services层
    * services/storeapi，提供文件平台的访问客户端
    * services/log，beego.Info/Debug等日志模式官方已经废弃，请用log.Logger
* 框架初始化代码请放到bootstrap.go当中
* models定义
    * models下为数据模型定义
    * models/forms为http请求表单数据定义
* BaseContoller
    * 所有conntroller需要匿名继承BaseController
    * BaseController中，Correct/Error系列函数返回json数据
* 定制了404返回页面

## 程序启停

* 建议使用graceful更新方式conf/app.conf中Graceful配置
* start.sh支持 [start | stop | reload | status | rdoc]
* rdoc为运行在dev模式下，但是启动了swagger日志

## 样例，请学会样例后再开始开发接口
demo.go中展示了普通接口、redis访问接口、数据库访问接口、http访问口

* 普通接口：不同接口类型的注释方法
    * GetUser
    * GetAllUsers
* redis接口
    * SetRedis: 如何利用valid校验传入参数
    * GetRedis：路径参数[path]如何获取

* mysql接口
    * GetMySQLInfo: 如何读取数据

* HTTP转发接口
    * HttpRequestDemo: 如何利用beego的httplib模块发送请求    

## 组件文档
* ORM推荐组件：[gorm](https://gorm.io/docs/)
    * 连接池 [connection pool](https://gorm.io/docs/generic_interface.html)
    * 方法串行使用 [method chaining](https://gorm.io/docs/method_chaining.html)
    * 模型定义 [model](https://gorm.io/docs/models.html)
    * ORM查询方式 [CRUD] (https://gorm.io/docs/create.html)
    * 原生查询方式 [SQL builder](https://gorm.io/docs/sql_builder.html)
    
* ORM可选组件: [xorm](https://xorm.io/zh/)
    * 主从分离 [engine group](https://github.com/go-xorm/xorm)

* ORM可选组件: [beego orm](https://beego.me/docs/mvc/model/orm.md)

* Redis组件: [go-redis](https://github.com/go-redis/redis)

* HTTP库组件
    * [beego 原生](https://github.com/astaxie/beego/tree/develop/httplib)
    * [grequests](https://github.com/levigross/grequests) (推荐)
    * [heimdall](https://github.com/gojek/heimdall)

## 接口注释标准

##### 注释格式说明地址
* [https://beego.me/docs/advantage/docs.md](https://beego.me/docs/advantage/docs.md)

```swagger
// @Title 接口标题
// @Description 接口详细描述
// @Param	参数名	参数来源	参数类型	是否必选 "参数描述"
// @Success http响应码 {返回类型} 返回类型的model定义
// @Failure http响应码 错误信息
// @router beego的注释路由
```

##### GET接口注释示例

```swagger
// @Title 获取用户
// @Description 通过用户ID获取用户信息
// @Param	uid	path	string	true	"用户id"
// @Success 200 {string} models.User.Id
// @Failure 200 uid is empty
// @router /:uid [get]
```

##### POST结果注释示例

```swagger
// @Title 创建用户
// @Description 根据传递的信息创建新的用户
// @Param   body body  models.User true "包体中有用户信息"
// @Success 200 {string} models.User.Id
// @Failure 200 body is empty
// @router / [post]
```


## 热启动

默认开启了graceful配置，热更新服务。./start.sh reload可以发送HUP命令重启进程。
dev模式下关闭graceful模式，否则bee run调试不便。

## swagger接口测试

在dev模式下，会默认启动swagger。访问服务[http://yourhost:port/swagger/](http://yourhost:port/swagger/)即可

## 管理端口

dev模式下，会开启AdminPort，这样可以看到程序内部的统计。包括请求qps统计。路由统计等。
* EnableAdmin = true
* AdminAddr = "localhost"
* AdminPort = 8088