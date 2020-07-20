package main

import (
	_ "bfimpl/routers"
	"bfimpl/services/log"

	"bfimpl/controllers"

	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/robfig/cron/v3"
)

func main() {
	if beego.BConfig.RunMode == beego.DEV {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	log.GLogger.Info("server start at %s:%s",
		beego.AppConfig.String("httpaddr"),
		beego.AppConfig.String("httpport"))

	go func() {
		time.Sleep(time.Second * 4)
		controllers.AmountDelayOut()
	}()
	c := cron.New()
	_, err := c.AddFunc("* 2 * * *", func() {
		controllers.AmountDelayOut()
	})
	if err != nil {
		logs.Error("start delay amount:%s", err.Error())
		return
	}
	c.Start()
	defer c.Stop()
	beego.Run()
}
