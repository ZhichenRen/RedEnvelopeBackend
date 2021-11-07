package dao

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type Envelope struct {
	ID         string `gorm:"size:12"`
	UID        string `gorm:"size:10"`
	Value      int
	Opened     bool
	SnatchTime int64
}

type PublicOpenedEnvelope struct {
	*Envelope           // 匿名嵌套
	UID       *struct{} `json:"uid,omitempty"`
}

func (Envelope) TableName() string {
	return "envelopes"
}

func GetEnvelopesByUID(uid string) ([]*Envelope, error) {
	dsn := "group9:Group9@haha@tcp(124.238.238.165:3306)/red_envelope?charset=utf8&parseTime=True&loc=Local&timeout=10s"
	// connect to mysql
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	var envelopes []*Envelope
	condition := map[string]interface{}{
		"uid": uid,
	}
	err = db.Table(Envelope{}.TableName()).Where(condition).Find(&envelopes).Error
	if err != nil {
		return nil, err
	}
	return envelopes, nil
}

func GetEnvelopeByEID(EID string) (envelope Envelope) {
	dsn := "group9:Group9@haha@tcp(124.238.238.165:3306)/red_envelope?charset=utf8&parseTime=True&loc=Local&timeout=10s"
	// connect to mysql
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.Where("id = ?", EID).First(&envelope)
	return envelope
}

func CreateEnvelope(user User) (envelope Envelope) {
	dsn := "group9:Group9@haha@tcp(124.238.238.165:3306)/red_envelope?charset=utf8&parseTime=True&loc=Local&timeout=10s"
	// connect to mysql
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}
	snatchTime := time.Now().UnixNano()
	// TODO
	// value should be a random number
	value := 10
	ID := "456"
	envelope = Envelope{ID: ID, UID: user.ID, Opened: false, Value: value, SnatchTime: snatchTime}
	db.Create(&envelope)
	return envelope
}
