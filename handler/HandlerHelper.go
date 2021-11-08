package handler

import (
	"go-web/utils"
	"strconv"
)

func writeUserToRedis(user utils.User) {
	userInfo := make(map[string]interface{})
	userInfo["cur_count"] = user.CurCount
	userInfo["amount"] = user.Amount
	rdb.HMSet("User:"+strconv.FormatInt(user.ID, 10), userInfo)
}

func writeEnvelopesSet(envelope utils.Envelope, userId string) {
	EID := strconv.FormatInt(envelope.ID, 10)
	rdb.SAdd("User:"+userId+":Envelopes", EID)
}

func writeEnvelopeToRedis(envelope utils.Envelope) {
	envelopeInfo := make(map[string]interface{})
	envelopeInfo["value"] = envelope.Value
	envelopeInfo["opened"] = envelope.Opened
	envelopeInfo["uid"] = envelope.UID
	envelopeInfo["snatch_time"] = envelope.SnatchTime
	rdb.HMSet("Envelope:"+strconv.FormatInt(envelope.ID, 10), envelopeInfo)
}
