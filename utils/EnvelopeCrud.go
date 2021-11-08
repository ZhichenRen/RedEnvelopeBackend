package utils

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sort"
	"time"
)

type Envelope struct {
	ID         int64
	UID        int64
	Value      int
	Opened     bool
	SnatchTime int64
}

func (Envelope) TableName() string {
	return "envelopes"
}

func GetEnvelopesByUID(uid int64) ([]*Envelope, error) {
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
	sort.SliceStable(envelopes, func(i, j int) bool {
		return envelopes[i].SnatchTime < envelopes[j].SnatchTime
	})
	return envelopes, nil
}

func GetEnvelopeByEID(eid int64) (envelope Envelope) {
	dsn := "group9:Group9@haha@tcp(124.238.238.165:3306)/red_envelope?charset=utf8&parseTime=True&loc=Local&timeout=10s"
	// connect to mysql
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.Where("id = ?", eid).First(&envelope)
	return envelope
}

func CreateEnvelope(uid int64) (envelope Envelope, user User, err error) {
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
	// TODO
	// maxCount
	err = db.Where("cur_count < ?", 50).First(&user, User{ID: uid}).Error
	if err == nil {
		envelope = Envelope{UID: uid, Opened: false,
			Value: value, SnatchTime: snatchTime}
		// TODO
		// there should be a error check
		db.Create(&envelope)
		user.CurCount++
		db.Save(&user)
	}
	return envelope, user, err
}

func OpenEnvelope(uid int64, eid int64) (envelope Envelope, user User, err error) {
	dsn := "group9:Group9@haha@tcp(124.238.238.165:3306)/red_envelope?charset=utf8&parseTime=True&loc=Local&timeout=10s"
	// connect to mysql
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}
	err = db.First(&envelope, Envelope{ID: eid, UID: uid, Opened: false}).Error
	db.First(&user, User{ID: uid})
	if err == nil {
		envelope.Opened = true
		user.Amount += envelope.Value
		db.Model(&user).Update("amount", user.Amount)
		db.Model(&envelope).Update("opened", envelope.Opened)
	}
	return
}
