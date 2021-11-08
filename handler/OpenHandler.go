package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-web/utils"
	"strconv"
)

func OpenHandler(c *gin.Context) {
	userId, _ := c.GetPostForm("uid")
	envelopeId, _ := c.GetPostForm("envelope_id")
	uid, _ := strconv.ParseInt(userId, 10, 64)
	eid, _ := strconv.ParseInt(envelopeId, 10, 64)
	//result, err := rdb.SIsMember("User:"+userId+":Envelopes", envelopeId).Result()
	envelopes, err := rdb.HGetAll("Envelope:" + envelopeId).Result()
	if err != nil {
		fmt.Println(err)
	}
	if len(envelopes) != 0 {
		opened := envelopes["opened"]
		value := envelopes["value"]
		realUId := envelopes["uid"]
		if err != nil {
			fmt.Println(err)
		}
		if opened != "0" {
			c.JSON(200, gin.H{
				"code":    0,
				"message": "您已经打开了此红包",
			})
		} else if opened == "0" && userId == realUId {
			_, user, _ := utils.OpenEnvelope(uid, eid)
			err = rdb.HSet("Envelope:"+envelopeId, "opened", true).Err()
			userInfo := make(map[string]interface{})
			userInfo["cur_count"] = user.CurCount
			userInfo["amount"] = user.Amount
			rdb.HMSet("User:"+userId, userInfo)
			c.JSON(200, gin.H{
				"code":    0,
				"message": "success",
				"data": gin.H{
					"value": value,
				},
			})
		} else {
			c.JSON(200, gin.H{
				"code":    0,
				"message": "no authorization",
			})
		}

	} else {
		envelope, user, err := utils.OpenEnvelope(uid, eid)
		if err == nil {
			userInfo := make(map[string]interface{})
			userInfo["cur_count"] = user.CurCount
			userInfo["amount"] = user.Amount
			envelopeInfo := make(map[string]interface{})
			envelopeInfo["id"] = envelope.ID
			envelopeInfo["uid"] = envelope.UID
			envelopeInfo["opened"] = envelope.Opened
			envelopeInfo["value"] = envelope.Value
			rdb.HMSet("Envelope:"+envelopeId, envelopeInfo)
			rdb.HMSet("User:"+userId, userInfo)
			c.JSON(200, gin.H{
				"code":    0,
				"message": "success",
				"data": gin.H{
					"value": envelope.Value,
				},
			})
		} else {
			c.JSON(200, gin.H{
				"code":    1,
				"message": "fail",
			})
		}
	}
}
