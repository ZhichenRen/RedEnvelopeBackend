package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
)

func OpenHandler(c *gin.Context) {
	userId, _ := c.GetPostForm("uid")
	envelopeId, _ := c.GetPostForm("envelope_id")
	result, err := rdb.SIsMember("User:"+userId+":Envelopes", envelopeId).Result()
	if err != nil {
		fmt.Println(err)
	}
	if result == true {
		opened, err := rdb.HGet("Envelope:"+envelopeId, "Opened").Result()
		if err != nil {
			fmt.Println(err)
		}
		if opened != "0" {
			c.JSON(200, gin.H{
				"code":    0,
				"message": "您已经打开了此红包",
			})
		}
		maxAmount, err := rdb.Get("MaxAmount").Int()
		value := rand.Intn(maxAmount)

		err = rdb.HSet("Envelope:"+envelopeId, "Opened", true).Err()
		err = rdb.HSet("Envelope:"+envelopeId, "Value", value).Err()
		err = rdb.HIncrBy("User:"+userId, "Amount", int64(value)).Err()

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
			"message": "您并不拥有此红包",
		})
	}
}
