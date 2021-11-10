package handler

import (
	"fmt"
	"go-web/allocate"
	"go-web/dao"
	"strconv"
	"time"
)

func writeUserToRedis(user dao.User) {
	userInfo := make(map[string]interface{})
	userInfo["cur_count"] = user.CurCount
	userInfo["amount"] = user.Amount
	rdb.HMSet("User:"+strconv.FormatInt(user.ID, 10), userInfo)
}

func writeEnvelopesSet(envelope dao.Envelope, userId string) {
	EID := strconv.FormatInt(envelope.ID, 10)
	rdb.SAdd("User:"+userId+":Envelopes", EID)
}

func writeEnvelopeToRedis(envelope dao.Envelope) {
	envelopeInfo := make(map[string]interface{})
	envelopeInfo["value"] = envelope.Value
	envelopeInfo["opened"] = envelope.Opened
	envelopeInfo["uid"] = envelope.UID
	envelopeInfo["snatch_time"] = envelope.SnatchTime
	rdb.HMSet("Envelope:"+strconv.FormatInt(envelope.ID, 10), envelopeInfo)
}

func updateAmount(UserId string, value string) {
	users, _ := rdb.HGetAll("User:" + UserId).Result()
	uid, _ := strconv.ParseInt(UserId, 10, 64)
	valueInt, _ := strconv.Atoi(value)
	if len(users) == 0 {
		user, _ := dao.GetUser(uid)
		user.Amount += valueInt
		writeUserToRedis(user)
	} else {
		curAmount, _ := strconv.Atoi(users["amount"])
		users["amount"] = strconv.Itoa(curAmount + valueInt)
	}
}

func updateAmountInt(userId string, value int) {
	users, _ := rdb.HGetAll("User:" + userId).Result()
	uid, _ := strconv.ParseInt(userId, 10, 64)
	if len(users) == 0 {
		user, _ := dao.GetUser(uid)
		user.Amount += value
		writeUserToRedis(user)
	} else {
		curAmount, _ := strconv.Atoi(users["amount"])
		users["amount"] = strconv.Itoa(curAmount + value)
	}
}

func updateOpened(eid int64) {
	envelope, _ := dao.GetEnvelopeByEID(eid)
	envelopeInfo := make(map[string]interface{})
	envelopeInfo["value"] = envelope.Value
	envelopeInfo["opened"] = true
	envelopeInfo["uid"] = envelope.UID
	envelopeInfo["snatch_time"] = envelope.SnatchTime
	rdb.HMSet("Envelope:"+strconv.FormatInt(envelope.ID, 10), envelopeInfo)
}

func updateCurCount(userId string) (curCount int64) {
	users, _ := rdb.HGetAll("User:" + userId).Result()
	uid, _ := strconv.ParseInt(userId, 10, 64)
	if len(users) == 0 {
		user, _ := dao.GetUser(uid)
		user.CurCount++
		writeUserToRedis(user)
	} else {
		curCount, _ = rdb.HIncrBy("User:"+userId, "cur_count", 1).Result()
	}
	return
}

func createEnvelope(userId string) (envelope dao.Envelope) {
	money := allocate.MoneyAllocate()
	snatchTime := time.Now().Unix()
	uid, err := strconv.ParseInt(userId, 10, 64)
	eid, err := rdb.Incr("EnvelopeId").Result()
	if err != nil {
		fmt.Println(err)
	}
	envelope = dao.Envelope{
		ID:         eid,
		UID:        uid,
		Value:      money,
		SnatchTime: snatchTime,
	}
	writeEnvelopeToRedis(envelope)
	return
}
