package services

import "github.com/jinzhu/gorm"

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
	return dbSlave
}
