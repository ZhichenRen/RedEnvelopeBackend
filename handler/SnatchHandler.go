package handler

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/gin-gonic/gin"
	"go-web/dao"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

func SnatchHandler(c *gin.Context) {
	// bool???
	userId, flag := c.GetPostForm("uid")
	if flag == false {
		fmt.Println("SnatchHandler label -1, GetPostForm uid", flag)
	}
	isCheat, err := rdb.Get("User:" + userId + ":Cheat").Result()
	if err != nil{
		logError("SnatchHandler", 0, err)
		c.JSON(500, gin.H{
			"code": 1,
			"msg":  "A database error occurred.",
		})
		return
	}
	if isCheat == "1" {
		c.JSON(403, gin.H{
			"code": 2,
			"msg":  "您因为作弊被系统封禁！",
		})
		return
	}
	// string -> int64
	uid, err := strconv.ParseInt(userId, 10, 64)
	logError("SnatchHandler", -2, err)

	// cheat detection
	snatchCount, err := rdb.Get("User:" + userId + ":Snatch").Int64()
	logError("SnatchHandler", 4, err)
	if snatchCount == 0 {
		err = rdb.Set("User:"+userId+":Snatch", 1, 10000000000).Err()
		logError("SnatchHandler", 5, err)
	} else {
		snatchCount, err = rdb.Incr("User:" + userId + ":Snatch").Result()
		logError("SnatchHandler", 6, err)
		if snatchCount > 10 {
			err := rdb.Set("User:"+userId+":Cheat", 1, 43200000000000).Err()
			if err != nil {
				logError("SnatchHAndler", 0, err)
				c.JSON(500, gin.H{
					"code": 1,
					"msg":  "A database error occurred.",
				})
				return
			}
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

	user, err := rdb.HGetAll("User:" + userId).Result()
	logError("SnatchHandler", 1, err)
	if len(user) == 0 {
		users, err := dao.GetUser(uid)
		logError("SnatchHandler", 2, err)
		writeUserToRedis(users)
		user, err = rdb.HGetAll("User:" + userId).Result()
		logError("SnatchHandler", 3, err)
		if err != nil {
			c.JSON(500, gin.H{
				"code": 1,
				"msg":  "A database error occurred.",
			})
			return
		}
	}

	probability, err := rdb.Get("Probability").Int()
	if err != nil {
		logError("SnatchHandler", 0, err)
		c.JSON()
	}
	n := rand.Intn(100)
	if n >= probability {
		c.JSON(500, gin.H{
			"code": 1,
			"msg":  "A database error occurred.",
		})
		return
	}

	maxCount, err := rdb.Get("MaxCount").Int()
	curCount, err := strconv.Atoi(user["cur_count"])
	logError("SnatchHandler", -4, err)
	if curCount < maxCount {
		// TODO
		// OUR CODE HERE
		// 随机数判断用户是否抢到红包，后期需要替换
		// ...
		curCount, err := rdb.HIncrBy("User:"+userId, "cur_count", 1).Result()
		logError("SnatchHandler", 7, err)
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
