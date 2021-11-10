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
		if userId != realUId {
			c.JSON(200, gin.H{
				"code":    1,
				"message": "no authorization",
			})
		} else if opened == "0" {
			rdb.HSet("Envelope:"+envelopeId, "opened", true)
			updateAmount(userId, value)
			// TODO
			// write to MySQL
			// OpenEnvelope should be deleted
			DBHelper.OpenEnvelope(uid, eid)
			c.JSON(200, gin.H{
				"code":    0,
				"message": "success",
				"data": gin.H{
					"value": value,
				},
			})
		} else {
			c.JSON(200, gin.H{
				"code":    2,
				"message": "had been opened",
			})
		}
	} else {
		envelope, err := DBHelper.OpenCheck(uid, eid)
		if err == 0 {
			updateAmountInt(userId, envelope.Value)
			updateOpened(eid)
			// TODO
			// write to MySQL
			// OpenEnvelope should be deleted
			// DBHelper.OpenEnvelope(uid, eid)
			DBHelper.OpenEnvelope(uid, eid)
			c.JSON(200, gin.H{
				"code":    0,
				"message": "success",
				"data": gin.H{
					"value": envelope.Value,
				},
			})
		} else {
			c.JSON(200, gin.H{
				"code":    err,
				"message": "fail",
			})
		}
	}
}
