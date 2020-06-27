package services

import (
	"github.com/jinzhu/gorm"
	"github.com/astaxie/beego"
	"bfimpl/services/log"
)

var dbMaster *gorm.DB
var dbSlave *gorm.DB

func SetDbConnection(master, slave *gorm.DB) {
	if slave == nil {
		slave = master
	}
	dbMaster = master
	dbSlave = slave
}

func Master() *gorm.DB {
	return dbMaster
}

func Slave() *gorm.DB {
	if dbSlave == nil {
		DBInit()
	}
	return dbSlave
}

func DBInit()  {
	dsn := beego.AppConfig.String("dsn")
	if dsn == "" {
		return
	}
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		log.GLogger.Critical("error conect to db, err=%v", err)
		return
	}
	log.GLogger.Info("init db connection ok")
	db.SetLogger(log.GLogger)
	db.DB().SetMaxOpenConns(30)
	db.DB().SetMaxIdleConns(10)
	SetDbConnection(db, db)
}