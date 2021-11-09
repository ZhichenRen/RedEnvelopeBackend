package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-web/DBHelper"
	"strconv"
)

func OpenHandler(c *gin.Context) {
	userId, _ := c.GetPostForm("uid")
	envelopeId, _ := c.GetPostForm("envelope_id")
	uid, _ := strconv.ParseInt(userId, 10, 64)
	eid, _ := strconv.ParseInt(envelopeId, 10, 64)
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
				"message": "had been opened",
			})
		} else if opened == "0" && userId == realUId {
			_, user, _ := DBHelper.OpenEnvelope(uid, eid)
			err = rdb.HSet("Envelope:"+envelopeId, "opened", true).Err()
			writeUserToRedis(user)
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
		envelope, user, err := DBHelper.OpenEnvelope(uid, eid)
		if err == nil {
			writeUserToRedis(user)
			writeEnvelopeToRedis(envelope)
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
