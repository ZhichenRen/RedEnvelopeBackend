package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-web/utils"
	"strconv"
)

func WalletListHandler(c *gin.Context) {
	userId, _ := c.GetPostForm("uid")
	uid, _ := strconv.ParseInt(userId, 10, 64)
	envelopeList, err := rdb.SMembers("User:" + userId + ":Envelopes").Result()
	users, _ := rdb.HGetAll("User:" + userId).Result()
	if err != nil {
		fmt.Println(err)
	}
	var curCount int
	var amount int
	var data []gin.H
	if len(users) == 0 {
		userInfo := make(map[string]interface{})
		user, _ := utils.GetUser(uid)
		amount = user.Amount
		curCount = user.CurCount
		userInfo["cur_count"] = user.CurCount
		userInfo["amount"] = user.Amount
		rdb.HMSet("User:"+userId, userInfo)
	} else {
		curCount, _ = strconv.Atoi(users["cur_count"])
		amount, _ = strconv.Atoi(users["amount"])
	}
	if len(envelopeList) == curCount {
		for _, envelopeId := range envelopeList {
			envelope, err := rdb.HGetAll("Envelope:" + envelopeId).Result()
			if err != nil {
				fmt.Println(err)
			}

			if len(envelope) != 0 {
				tmp := gin.H{}
				tmp["envelope_id"] = envelopeId
				tmp["snatch_time"] = envelope["snatch_time"]
				if envelope["opened"] == "0" {
					tmp["opened"] = false
				} else {
					tmp["opened"] = true
					tmp["value"] = envelope["value"]
				}
				data = append(data, tmp)
			} else {
				eid, _ := strconv.ParseInt(envelopeId, 10, 64)
				envelopeFromSql := utils.GetEnvelopeByEID(eid)
				tmp := gin.H{}
				tmp["envelope_id"] = envelopeId
				tmp["snatch_time"] = envelopeFromSql.SnatchTime
				if envelopeFromSql.Opened == false {
					tmp["opened"] = false
				} else {
					tmp["opened"] = true
					tmp["value"] = envelopeFromSql.Value
				}
				envelopeInfo := make(map[string]interface{})
				envelopeInfo["uid"] = envelopeFromSql.UID
				envelopeInfo["opened"] = envelopeFromSql.Opened
				envelopeInfo["value"] = envelopeFromSql.Value
				rdb.HMSet("Envelope:"+envelopeId, envelopeInfo)
				data = append(data, tmp)
			}
		}
	} else {
		fmt.Println(666)
		envelopes, _ := utils.GetEnvelopesByUID(uid)
		envelopeInfo := make(map[string]interface{})
		for _, envelope := range envelopes {
			tmp := gin.H{}
			tmp["envelope_id"] = envelope.ID
			tmp["snatch_time"] = envelope.SnatchTime
			if envelope.Opened == false {
				tmp["opened"] = false
			} else {
				tmp["opened"] = true
				tmp["value"] = envelope.Value
			}
			envelopeInfo["uid"] = envelope.UID
			envelopeInfo["opened"] = envelope.Opened
			envelopeInfo["value"] = envelope.Value
			rdb.HMSet("Envelope:"+strconv.Itoa(int(envelope.ID)), envelopeInfo)
			_, err = rdb.SAdd("User:"+userId+":Envelopes", strconv.Itoa(int(envelope.ID))).Result()
			data = append(data, tmp)
		}
	}

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"amount":        amount,
			"envelope_list": data,
		},
	})
}
