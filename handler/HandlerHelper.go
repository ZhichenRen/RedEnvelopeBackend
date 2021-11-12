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
	err := rdb.HMSet("User:"+strconv.FormatInt(user.ID, 10), userInfo).Err()
	fmt.Println("HandlerHelper, writeUserToRedis label 1", err)
}

func writeEnvelopesSet(envelope dao.Envelope, userId string) {
	EID := strconv.FormatInt(envelope.ID, 10)
	err := rdb.SAdd("User:"+userId+":Envelopes", EID).Err()
	fmt.Println("HandlerHelper, writeEnvelopesSet label 1", err)
}

func writeEnvelopeToRedis(envelope dao.Envelope) {
	envelopeInfo := make(map[string]interface{})
	envelopeInfo["value"] = envelope.Value
	envelopeInfo["opened"] = envelope.Opened
	envelopeInfo["uid"] = envelope.UID
	envelopeInfo["snatch_time"] = envelope.SnatchTime
	err := rdb.HMSet("Envelope:"+strconv.FormatInt(envelope.ID, 10), envelopeInfo).Err()
	fmt.Println("HandlerHelper, writeEnvelopeToRedis label 1", err)
}

func updateAmount(UserId string, value string) {
	users, err := rdb.HGetAll("User:" + UserId).Result()
	fmt.Println("HandlerHelper, updateAmount label 1", err)
	uid, err := strconv.ParseInt(UserId, 10, 64)
	fmt.Println("HandlerHelper, updateAmount label 2", err)
	valueInt, err := strconv.Atoi(value)
	fmt.Println("HandlerHelper, updateAmount label 3", err)
	if len(users) == 0 {
		user, err := dao.GetUser(uid)
		fmt.Println("HandlerHelper, updateAmount label 4", err)
		user.Amount += valueInt
		writeUserToRedis(user)
	} else {
		curAmount, err := strconv.Atoi(users["amount"])
		fmt.Println("HandlerHelper, updateAmount label 5", err)
		//users["amount"] = strconv.Itoa(curAmount + valueInt)
		err = rdb.HSet("User:"+UserId, "amount", curAmount+valueInt).Err()
		fmt.Println("HandlerHelper, updateAmount label 6", err)
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
	moneyLeft, err := rdb.Get("TotalMoney").Int()
	fmt.Println("HandlerHelper, createEnvelope label 1", err)
	envelopeLeft, err := rdb.Get("EnvelopeNum").Int()
	fmt.Println("HandlerHelper, createEnvelope label 2", err)
	money := allocate.MoneyAllocate(int64(moneyLeft), int64(envelopeLeft))
	snatchTime := time.Now().Unix()
	uid, err := strconv.ParseInt(userId, 10, 64)
	fmt.Println("HandlerHelper, createEnvelope label 3", err)
	eid, err := rdb.Incr("EnvelopeId").Result()
	fmt.Println("HandlerHelper, createEnvelope label 4", err)
	envelope = dao.Envelope{
		ID:         eid,
		UID:        uid,
		Value:      money,
		SnatchTime: snatchTime,
	}
	err = rdb.IncrBy("TotalMoney", int64(-money)).Err()
	fmt.Println("HandlerHelper, createEnvelope label 5", err)
	err = rdb.IncrBy("EnvelopeNum", -1).Err()
	fmt.Println("HandlerHelper, createEnvelope label 6", err)
	writeEnvelopeToRedis(envelope)
	return
}
