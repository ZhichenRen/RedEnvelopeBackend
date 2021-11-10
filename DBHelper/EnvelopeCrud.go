package DBHelper

import (
	"go-web/allocate"
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
	var envelopes []*Envelope
	condition := map[string]interface{}{
		"uid": uid,
	}
	err := _db.Table(Envelope{}.TableName()).Where(condition).Find(&envelopes).Error
	if err != nil {
		return nil, err
	}
	sort.SliceStable(envelopes, func(i, j int) bool {
		return envelopes[i].SnatchTime < envelopes[j].SnatchTime
	})
	return envelopes, nil
}

func GetEnvelopeByEID(eid int64) (envelope Envelope) {
	_db.Where("id = ?", eid).First(&envelope)
	return envelope
}

func CreateEnvelope(uid int64) (envelope Envelope, user User, err error) {
	snatchTime := time.Now().Unix()
	value := allocate.MoneyAllocate()
	// TODO
	// maxCount
	err = _db.Where("cur_count < ?", 50).First(&user, User{ID: uid}).Error
	if err == nil {
		envelope = Envelope{UID: uid, Opened: false,
			Value: value, SnatchTime: snatchTime}
		// TODO
		// there should be a error check
		_db.Create(&envelope)
		user.CurCount++
		_db.Save(&user)
	}
	return envelope, user, err
}

func OpenEnvelope(uid int64, eid int64) {
	user, _ := GetUser(uid)
	envelope := GetEnvelopeByEID(eid)
	envelope.Opened = true
	user.Amount += envelope.Value
	_db.Model(&user).Update("amount", user.Amount)
	_db.Model(&envelope).Update("opened", true)
	return
}

func OpenCheck(uid int64, eid int64) (envelope Envelope, error int) {
	err := _db.First(&envelope, Envelope{ID: eid, UID: uid}).Error
	if err != nil {
		error = 1
		return
	} else if envelope.Opened {
		error = 2
		return
	} else {
		error = 0
		return
	}
}
