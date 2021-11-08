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
	curCount, _ := strconv.Atoi(users["curCount"])
	// 开始对接
	if len(users) == 0 {
		userInfo := make(map[string]interface{})
		newEnvelope, user, err := utils.CreateEnvelope(uid)
		if err == nil {
			userInfo["CurCount"] = user.CurCount
			userInfo["Amount"] = user.Amount
			rdb.HMSet("User:"+userId, userInfo)
			users, err = rdb.HGetAll("User:" + userId).Result()
			if err != nil {
				fmt.Println(err)
			}
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
				"data": gin.H{
					"envelop_id": 0,
					"max_count":  0,
					"cur_count":  0,
				},
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
		curCount, err := rdb.HIncrBy("User:"+userId, "CurCount", 1).Result()
		if err != nil {
			fmt.Println(err)
		}
		envelope := make(map[string]interface{})
		// TODO
		// value should be random
		newEnvelope, _, err := utils.CreateEnvelope(uid)
		envelope["Value"] = newEnvelope.Value
		envelope["Opened"] = newEnvelope.Opened
		envelope["SnatchTime"] = newEnvelope.SnatchTime
		envelope["EnvelopeId"] = newEnvelope.ID
		_, err = rdb.HMSet("Envelope:"+strconv.Itoa(int(newEnvelope.ID)), envelope).Result()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Envelope:", envelope)

		_, err = rdb.SAdd("User:"+userId+":Envelopes", newEnvelope.ID).Result()
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
