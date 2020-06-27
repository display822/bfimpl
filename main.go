package main

import (
	_ "bfimpl/routers"
	"bfimpl/services/log"

	"github.com/astaxie/beego"
)

func main() {
	if beego.BConfig.RunMode == beego.DEV {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	log.GLogger.Info("server start at %s:%s",
		beego.AppConfig.String("httpaddr"),
		beego.AppConfig.String("httpport"))
	beego.Run()
}
