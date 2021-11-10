package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-web/dao"
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
		user, _ := dao.GetUser(uid)
		amount = user.Amount
		curCount = user.CurCount
		writeUserToRedis(user)
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
				envelopeFromSql, _ := dao.GetEnvelopeByEID(eid)
				tmp := gin.H{}
				tmp["envelope_id"] = envelopeId
				tmp["snatch_time"] = envelopeFromSql.SnatchTime
				if envelopeFromSql.Opened == false {
					tmp["opened"] = false
				} else {
					tmp["opened"] = true
					tmp["value"] = envelopeFromSql.Value
				}
				writeEnvelopeToRedis(envelopeFromSql)
				data = append(data, tmp)
			}
		}
	} else {
		envelopes, _ := dao.GetEnvelopesByUID(uid)
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
			writeEnvelopeToRedis(*envelope)
			writeEnvelopesSet(*envelope, userId)
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
