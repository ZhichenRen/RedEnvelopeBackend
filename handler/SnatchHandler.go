package handler

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/gin-gonic/gin"
	"go-web/dao"
	"strconv"
	"sync"
	"time"
)

func SnatchHandler(c *gin.Context) {
	// bool???
	userId, flag := c.GetPostForm("uid")
	fmt.Println("SnatchHandler label -1, GetPostForm uid", flag)
	// string -> int64
	uid, err := strconv.ParseInt(userId, 10, 64)
	fmt.Println("SnatchHandler label -2, ParseInt", err)
	user, err := rdb.HGetAll("User:" + userId).Result()
	fmt.Println("SnatchHandler label 1, get user from redis", err)
	// TODO how to get maxCount
	// maxCount, err := rdb.Get("MaxCount").Result()
	// search in mysql
	if len(user) == 0 {
		users, err := dao.GetUser(uid)
		fmt.Println("SnatchHandler label 2, get user from mysql", err)
		writeUserToRedis(users)
		user, err = rdb.HGetAll("User:" + userId).Result()
		fmt.Println("SnatchHandler label 3, get user from redis", err)
		if err != nil {
			c.JSON(500, gin.H{
				"code": 1,
				"msg":  "A database error occurred.",
			})
			return
		}
	}

	// cheat detection
	snatchCount, err := rdb.Get("User:" + userId + ":Snatch").Int64()
	fmt.Println("SnatchHandler label 4, get user from redis", err)
	if snatchCount == 0 {
		err = rdb.Set("User:"+userId+":Snatch", 1, 10000000000).Err()
		fmt.Println("SnatchHandler label 5, set user in redis", err)
	} else {
		snatchCount, err = rdb.Incr("User:" + userId + ":Snatch").Result()
		fmt.Println("SnatchHandler label 6, increase userId", err)
		if snatchCount > 10 {
			c.JSON(403, gin.H{
				"code": 2,
				"msg":  "系统检测到你在作弊！",
			})
			return
		}
	}
	if err != nil {
		c.JSON(500, gin.H{
			"code": 1,
			"msg":  "A database error occurred.",
		})
		return
	}

	maxCount := 10
	curCount, err := strconv.Atoi(user["cur_count"])
	fmt.Println("SnatchHandler label -4, Atoi", err)
	if curCount < maxCount {
		// TODO
		// OUR CODE HERE
		// 随机数判断用户是否抢到红包，后期需要替换
		// ...
		curCount, err := rdb.HIncrBy("User:"+userId, "cur_count", 1).Result()
		fmt.Println("SnatchHandler label 7, increase cur_count", err)
		if err != nil {
			fmt.Println(err)
		}
		newEnvelope := createEnvelope(userId)
		writeEnvelopesSet(newEnvelope, userId)

		// message queue
		p := GetProducer()
		var wg sync.WaitGroup
		topic := "Msg"
		params := make(map[string]string)
		params["UID"] = strconv.FormatInt(newEnvelope.UID, 10)
		params["EID"] = strconv.FormatInt(newEnvelope.ID, 10)
		params["Value"] = strconv.Itoa(newEnvelope.Value)
		params["SnatchTime"] = strconv.Itoa(int(time.Now().Unix()))
		message := primitive.NewMessage(topic, []byte("create_envelope"))
		message.WithProperties(params)
		fmt.Println(params)
		wg.Add(1)
		err = p.SendAsync(context.Background(),
			func(ctx context.Context, result *primitive.SendResult, e error) {
				if e != nil {
					fmt.Printf("receive message error: %s\n", err)
				} else {
					fmt.Printf("send message success: result=%s\n", result.String())
				}
				wg.Done()
			}, message)
		if err != nil {
			fmt.Printf("SnatchHandler label 9, an error occurred when sending message:%s\n", err)
			fmt.Println(message)
			c.JSON(500, gin.H{
				"code": 1,
				"msg":  "An error occurred when sending message.",
			})
			return
		}
		wg.Wait()
		//err = p.Shutdown()
		//fmt.Println("SnatchHandler label 10, shutdown", err)
		//if err != nil {
		//	c.JSON(500, gin.H{
		//		"code": 1,
		//		"msg":  "An error occurred when closing producer.",
		//	})
		//	return
		//}
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
			"code": 2,
			"msg":  "很抱歉，您没有抢到红包，这可能是因为手气不佳或已达上限",
		})
	}
}
