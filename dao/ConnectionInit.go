package dao

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var _db *gorm.DB

func init() {
	//dsn := "group9:Group9@haha@tcp(rdsmysqlh1a4d645c087a17d2.rds.ivolces.com:3306)/red_envelope?charset=utf8&parseTime=True&loc=Local&timeout=10s"
	dsn := "group9:Group9@haha@tcp(124.238.238.165:3306)/red_envelope?charset=utf8&parseTime=True&loc=Local&timeout=10s"
	var err error
	_db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database, error=" + err.Error())
	}
	sqlDB, _ := _db.DB()
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(20)
}

func GetDB() *gorm.DB {
	fmt.Println(_db)
	return _db
}
