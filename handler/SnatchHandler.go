package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-web/utils"
	"strconv"
)

func SnatchHandler(c *gin.Context) {
	userId, _ := c.GetPostForm("uid")
	uid, err := strconv.ParseInt(userId, 10, 64)
	users, err := rdb.HGetAll("User:" + userId).Result()
	if err != nil {
		fmt.Println(err)
	}
	// TODO how to get maxCount
	// maxCount, err := rdb.Get("MaxCount").Result()
	maxCount := 100
	curCount, _ := strconv.Atoi(users["cur_count"])
	// 开始对接
	if len(users) == 0 {
		userInfo := make(map[string]interface{})
		newEnvelope, user, err := utils.CreateEnvelope(uid)
		if err == nil {
			userInfo["cur_count"] = user.CurCount
			userInfo["amount"] = user.Amount
			rdb.HMSet("User:"+userId, userInfo)
			users, err = rdb.HGetAll("User:" + userId).Result()
			if err != nil {
				fmt.Println(err)
			}
			_, err = rdb.SAdd("User:"+userId+":Envelopes", strconv.Itoa(int(newEnvelope.ID))).Result()
			c.JSON(200, gin.H{
				"code": 0,
				"msg":  "success",
				"data": gin.H{
					"envelop_id": newEnvelope.ID,
					"max_count":  maxCount,
					"cur_count":  user.CurCount,
				},
			})
		} else {
			c.JSON(200, gin.H{
				"code": 1,
				"msg":  "failed",
			})
		}
	} else if curCount < maxCount {
		//fmt.Println("User:", users)
		// TODO
		// OUR CODE HERE
		// 随机数判断用户是否抢到红包，后期需要替换
		// ...
		if err != nil {
			fmt.Println(err)
		}
		curCount, err := rdb.HIncrBy("User:"+userId, "cur_count", 1).Result()
		if err != nil {
			fmt.Println(err)
		}
		envelopeInfo := make(map[string]interface{})
		// TODO
		// value should be random
		newEnvelope, _, err := utils.CreateEnvelope(uid)
		envelopeInfo["value"] = newEnvelope.Value
		envelopeInfo["opened"] = newEnvelope.Opened
		envelopeInfo["snatch_time"] = newEnvelope.SnatchTime
		_, err = rdb.HMSet("Envelope:"+strconv.Itoa(int(newEnvelope.ID)), envelopeInfo).Result()
		_, err = rdb.SAdd("User:"+userId+":Envelopes", strconv.Itoa(int(newEnvelope.ID))).Result()
		if err != nil {
			fmt.Println(err)
		}

		//_, err = rdb.SAdd("User:"+userId+":Envelopes", newEnvelope.ID).Result()
		if err != nil {
			fmt.Println(err)
		}
		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "success",
			"data": gin.H{
				"envelop_id": newEnvelope.ID,
				"max_count":  maxCount,
				"cur_count":  curCount,
			},
		})
	} else {
		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "fail",
		})
	}
}
