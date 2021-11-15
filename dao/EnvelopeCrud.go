package dao

import (
	"fmt"
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
	var envelopes []*Envelope
	condition := map[string]interface{}{
		"uid": uid,
	}
	err := _db.Table(Envelope{}.TableName()).Where(condition).Find(&envelopes).Order("snatch_time DESC").Error
	if err != nil {
		return nil, err
	}
	//sort.SliceStable(envelopes, func(i, j int) bool {
	//	return envelopes[i].SnatchTime > envelopes[j].SnatchTime
	//})
	return envelopes, nil
}

func GetEnvelopeByEID(eid int64) (Envelope, error) {
	var envelope Envelope
	err := _db.Where("id = ?", eid).First(&envelope).Error
	return envelope, err
}

func CreateCheck(uid int64) (errorCode int) {
	// TODO
	// maxCount
	fmt.Println(uid)
	var user User
	err := _db.Where("id = ?", uid).First(&user).Error
	fmt.Println(err)
	if err != nil {
		errorCode = 2
	} else if user.CurCount > 10 {
		errorCode = 1
	} else {
		errorCode = 0
	}
	return
}

func CreateEnvelope(envelope Envelope) {
	user, _ := GetUser(envelope.UID)
	user.CurCount++
	_db.Model(&user).Update("cur_count", user.CurCount)
	_db.Create(&envelope)
}

func OpenEnvelope(uid int64, eid int64) {
	user, _ := GetUser(uid)
	envelope, _ := GetEnvelopeByEID(eid)
	envelope.Opened = true
	user.Amount += envelope.Value
	_db.Model(&user).Update("amount", user.Amount)
	_db.Model(&envelope).Update("opened", true)
	return
}

func OpenCheck(uid int64, eid int64) (envelope Envelope, errorCode int) {
	err := _db.First(&envelope, Envelope{ID: eid, UID: uid}).Error
	if err != nil {
		errorCode = 1
		return
	} else if envelope.Opened {
		errorCode = 2
		return
	} else {
		errorCode = 0
		return
	}
}
