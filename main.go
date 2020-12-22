package main

import (
	_ "bfimpl/routers"
	"bfimpl/services/log"

	"bfimpl/controllers"

	"math/rand"
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

	c := cron.New()
	//额度过期 每天两点
	_, err := c.AddFunc("0 2 * * *", func() {
		controllers.AmountDelayOut()
	})
	if err != nil {
		logs.Error("start delay amount:%s", err.Error())
		return
	}
	//社保信息 每月16号两点
	_, err = c.AddFunc("0 2 16 * *", func() {
		controllers.GeneraSheBao()
	})
	if err != nil {
		logs.Error("genera shebao info:%s", err.Error())
		return
	}
	//年假增加 每月28号4点
	_, err = c.AddFunc("0 4 28 * *", func() {
		controllers.AddAnnual()
	})
	if err != nil {
		logs.Error("add annual:%s", err.Error())
		return
	}
	//清空去年年假 4月1号6点
	_, err = c.AddFunc("0 6 1 4 *", func() {
		controllers.DeleteAnnual()
	})
	if err != nil {
		logs.Error("add annual:%s", err.Error())
		return
	}
	c.Start()
	defer c.Stop()
	rand.Seed(time.Now().Unix())
	beego.Run()
}
