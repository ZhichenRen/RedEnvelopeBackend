package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-web/dao"
	"sort"
	"strconv"
)

func WalletListHandler(c *gin.Context) {
	userId, flag := c.GetPostForm("uid")
	if flag == false {
		fmt.Println("WalletListHandler label -1, GetPostForm", flag)
	}
	uid, err := strconv.ParseInt(userId, 10, 64)
	logError("WalletListHandler", -2, err)
	envelopeList, err := rdb.SMembers("User:" + userId + ":Envelopes").Result()
	logError("WalletListHandler", 1, err)
	users, err := rdb.HGetAll("User:" + userId).Result()
	logError("WalletListHandler", 2, err)
	var curCount int
	var amount int
	var data []gin.H
	if len(users) == 0 {
		user, err := dao.GetUser(uid)
		logError("WalletListHandler", 3, err)
		if err != nil {
			c.JSON(200, gin.H{
				"code": 2,
				"msg": "用户ID不存在",
			})
			return
		}
		amount = user.Amount
		curCount = user.CurCount
		writeUserToRedis(user)
	} else {
		curCount, err = strconv.Atoi(users["cur_count"])
		logError("WalletListHandler", -3, err)
		amount, err = strconv.Atoi(users["amount"])
		logError("WalletListHandler", -4, err)
	}
	if len(envelopeList) == curCount {
		for _, envelopeId := range envelopeList {
			envelope, err := rdb.HGetAll("Envelope:" + envelopeId).Result()
			logError("WalletListHandler", 4, err)
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
				eid, err := strconv.ParseInt(envelopeId, 10, 64)
				logError("WalletListHandler", -5, err)
				envelopeFromSql, err := dao.GetEnvelopeByEID(eid)
				logError("WalletListHandler", 5, err)
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
		envelopes, err := dao.GetEnvelopesByUID(uid)
		logError("WalletListHandler", 6, err)
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
	sort.SliceStable(data, func(i, j int) bool {
		snatchTimeI, ok := data[i]["snatch_time"].(int64)
		snatchTimeJ, ok := data[j]["snatch_time"].(int64)
		if ok == false {
			fmt.Println("Error happen when convert interface{} to int64!")
			return false
		}
		return snatchTimeI > snatchTimeJ
	})
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"amount":        amount,
			"envelope_list": data,
		},
	})
}
