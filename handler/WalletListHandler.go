package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func WalletListHandler(c *gin.Context) {
	userId, _ := c.GetPostForm("uid")
	fmt.Println(userId)

	envelopeList, err := rdb.SMembers("User:" + userId + ":Envelopes").Result()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(envelopeList)
	var data []gin.H
	for _, envelopeId := range envelopeList {
		envelope, err := rdb.HGetAll("Envelope:" + envelopeId).Result()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(envelope)
		tmp := gin.H{}
		tmp["envelope_id"] = envelopeId
		tmp["snatch_time"] = envelope["SnatchTime"]
		if envelope["Opened"] == "0" {
			tmp["opened"] = false
		} else {
			tmp["opened"] = true
			tmp["value"] = envelope["Value"]
		}
		data = append(data, tmp)
	}

	amount, err := rdb.HGet("User:"+userId, "Amount").Result()

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"amount":        amount,
			"envelope_list": data,
		},
	})
}
