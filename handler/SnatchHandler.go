package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-web/DBHelper"
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
		err := DBHelper.CreateCheck(uid)
		if err == 0 {
			newEnvelope := createEnvelope(userId)
			updateCount := updateCurCount(userId)
			writeEnvelopesSet(newEnvelope, userId)
			// TODO
			// write to sql
			// CreateEnvelope should be deleted
			DBHelper.CreateEnvelope(newEnvelope)

			c.JSON(200, gin.H{
				"code": 0,
				"msg":  "success",
				"data": gin.H{
					"envelope_id": newEnvelope.ID,
					"max_count":   maxCount,
					"cur_count":   updateCount,
				},
			})
		} else {
			c.JSON(200, gin.H{
				"code": err,
				"msg":  "failed",
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
		newEnvelope := createEnvelope(userId)
		writeEnvelopesSet(newEnvelope, userId)
		// TODO
		// write to sql
		// CreateEnvelope should be deleted
		DBHelper.CreateEnvelope(newEnvelope)

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
			"code": 1,
			"msg":  "exceed",
		})
	}
}
