package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-web/utils"
	"strconv"
)

func SnatchHandler(c *gin.Context) {
	userId, _ := c.GetPostForm("uid")
	// string -> int64
	uid, err := strconv.ParseInt(userId, 10, 64)
	users, err := rdb.HGetAll("User:" + userId).Result()
	if err != nil {
		fmt.Println(err)
	}
	// TODO how to get maxCount
	// maxCount, err := rdb.Get("MaxCount").Result()
	maxCount := 10
	curCount, _ := strconv.Atoi(users["cur_count"])
	// search in mysql
	if len(users) == 0 {
		newEnvelope, user, err := utils.CreateEnvelope(uid)
		if err == nil {
			writeUserToRedis(user)
			users, err = rdb.HGetAll("User:" + userId).Result()
			if err != nil {
				fmt.Println(err)
			}
			writeEnvelopesSet(newEnvelope, userId)
			c.JSON(200, gin.H{
				"code": 0,
				"msg":  "success",
				"data": gin.H{
					"envelope_id": newEnvelope.ID,
					"max_count":   maxCount,
					"cur_count":   user.CurCount,
				},
			})
		} else {
			c.JSON(200, gin.H{
				"code": 1,
				"msg":  "User not existed",
			})
		}
	} else if curCount < maxCount {
		//fmt.Println("User:", users)
		// TODO
		// OUR CODE HERE
		// 随机数判断用户是否抢到红包，后期需要替换
		// ...
		curCount, err := rdb.HIncrBy("User:"+userId, "cur_count", 1).Result()
		if err != nil {
			fmt.Println(err)
		}
		// TODO
		// value should be random
		newEnvelope, _, err := utils.CreateEnvelope(uid)
		writeEnvelopeToRedis(newEnvelope)
		writeEnvelopesSet(newEnvelope, userId)
		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "success",
			"data": gin.H{
				"envelope_id": newEnvelope.ID,
				"max_count":   maxCount,
				"cur_count":   curCount,
			},
		})
	} else {
		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "fail",
		})
	}
}
