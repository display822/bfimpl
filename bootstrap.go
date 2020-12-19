package main

import (
	"bfimpl/controllers"
	"bfimpl/services"
	"bfimpl/services/log"
	"fmt"
	"net/http"
	"os"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/go-redis/redis"
)

func init() {
	beego.ErrorHandler("404", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte(`{"ret":404,"msg":"API not found"}`))
		writer.WriteHeader(404)
	})
}

// 初始化日志系统
func init() {
	l := log.NewLogger()
	if beego.BConfig.RunMode == beego.DEV {
		l.SetLevel(logs.LevelDebug)
		_ = l.SetLogger(logs.AdapterConsole)
	} else {
		l.SetLevel(logs.LevelInfo)
	}
	_ = os.MkdirAll("logs", 0755)
	_ = l.SetLogger(logs.AdapterFile,
		fmt.Sprintf(`{"filename":"logs%v%v.log", "daily":true, "maxdays":%d}`,
			string(os.PathSeparator),
			beego.BConfig.AppName,
			10,
		))
	tlogaddr := beego.AppConfig.DefaultString("tlogaddr", "")
	if tlogaddr != "" && tlogaddr != "127.0.0.1:6666" {
		sname := beego.AppConfig.DefaultString("servicename", "wetest")
		tname := beego.AppConfig.DefaultString("toolname", beego.BConfig.AppName)
		logs.Register("dlogs", log.NewDlogsWriter)
		dlogconfig := fmt.Sprintf(`{
			"servicename":"%s",
			"toolname":"%v",
			"nettype":"udp",
			"tlogipport":"%v"}`,
			tname,
			sname,
			tlogaddr,
		)
		err := l.SetLogger("dlogs", dlogconfig)
		if err != nil {
			l.Error("init dlogs error %v", err)
		} else {
			l.Info("init dlogs ok, tlogaddr=%v, servicename=%v toolname=%v", tlogaddr, sname, tname)
		}
	}
	log.SetLogger(l)
	// 是否自动打印错误返回日志
	controllers.AutoLogError = beego.AppConfig.DefaultBool("autologerror", true)
}

// 初始化redis连接
func init() {
	redisaddr := beego.AppConfig.String("redisaddr")
	if redisaddr == "" {
		return
	}
	c := redis.NewClient(&redis.Options{
		Addr: redisaddr,
	})
	services.SetRedisClient(c)
	log.GLogger.Info("init redis addr=%v", redisaddr)
}

// 初始化的数据库连接
func init() {
	//services.DBInit()
}

func init() {
	services.MailInit()
}
