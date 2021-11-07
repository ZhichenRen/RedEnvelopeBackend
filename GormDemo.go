package main

// TODO
// or mey a question
// how to connect to different schema
import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Config : db config
type Config struct {
	// gorm.Model
	username string
	password string
	host     string
	port     uint
	dbName   string
	timeout  string
}

type Envelope struct {
	envelope   string
	value      int
	opened     int8
	snatchTime int64 // we may change the type into timestamp?
}

func main() {
	// set the config
	var config Config
	config = Config{
		"group9",
		"Group9@haha",
		"124.238.238.165",
		3306,
		"red_envelope",
		"10s",
	}
	// dsn
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%s",
		config.username, config.password, config.host, config.port, config.dbName, config.timeout)
	// connect to mysql
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed, error=" + err.Error())
	}

	var envelope Envelope
	envelope = Envelope{
		"red",
		1,
		0,
		100,
	}
	fmt.Println(db.Create(&envelope).Error)

}
