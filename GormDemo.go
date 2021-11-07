package main

// TODO
// or mey a question
// how to connect to different schema
import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
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
	envelopeId string `gorm:"primaryKey"`
	value      int
	opened     int8
	snatchTime int64
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
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic("Failed, error=" + err.Error())
	}

	//db.Model(&Envelope{}).Where("value = ?", 12).Update("value", 32)
	//envelope := db.Model(&Envelope{}).Where("envelopeId = ?", "123")
	var envelope Envelope
	db.Model(Envelope{}).Where("envelopId = ?", "123").First(envelope.envelopeId)
	fmt.Println(envelope.snatchTime)
}
