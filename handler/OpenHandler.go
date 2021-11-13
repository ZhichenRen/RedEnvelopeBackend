package handler

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/gin-gonic/gin"
	"go-web/dao"
	"strconv"
	"sync"
)

func OpenHandler(c *gin.Context) {
	userId, flag := c.GetPostForm("uid")
	if flag == false {
		fmt.Println("OpenHandler label -1, GetPostForm", flag)
	}
	envelopeId, flag := c.GetPostForm("envelope_id")
	if flag == false {
		fmt.Println("OpenHandler label -2, GetPostForm", flag)
	}
	uid, err := strconv.ParseInt(userId, 10, 64)
	if flag == false {
		fmt.Println("OpenHandler label -3, ParseInt", err)
	}
	//eid, _ := strconv.ParseInt(envelopeId, 10, 64)
	envelopes, err := rdb.HGetAll("Envelope:" + envelopeId).Result()
	logError("OpenHandler", 1, err)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 1,
			"msg":  "An error occurred when reading from redis",
		})
		return
	}

	if len(envelopes) == 0 {
		// key in redis expired, read from mysql
		envelopeList, err := dao.GetEnvelopesByUID(uid)
		logError("OpenHandler", 2, err)
		// error or envelope not found
		if err != nil {
			c.JSON(500, gin.H{
				"code": 1,
				"msg":  "A database error occurred or the envelope didn't exist.",
			})
			return
		}
		// write envelope set and envelopes
		for _, e := range envelopeList {
			writeEnvelopeToRedis(*e)
			writeEnvelopesSet(*e, userId)
		}
		envelopes, err = rdb.HGetAll("Envelope:" + envelopeId).Result()
		logError("OpenHandler", 3, err)
		if err != nil {
			c.JSON(500, gin.H{
				"code": 1,
				"msg":  "An error occurred when reading from redis",
			})
			return
		}
	}
	opened := envelopes["opened"]
	value := envelopes["value"]
	realUId := envelopes["uid"]
	if userId != realUId {
		c.JSON(200, gin.H{
			"code":    2,
			"message": "这个红包不属于您，您无权打开！",
		})
		return
	}
	if opened == "0" {
		// write to redis
		rdb.HSet("Envelope:"+envelopeId, "opened", true)
		updateAmount(userId, value)

		p := GetProducer()
		var wg sync.WaitGroup
		topic := "Msg"
		params := make(map[string]string)
		params["UID"] = userId
		params["EID"] = envelopeId
		message := primitive.NewMessage(topic, []byte("open_envelope"))
		message.WithProperties(params)
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
			fmt.Printf("OpenHandler label 6, an error occurred when sending message:%s\n", err)
			fmt.Println(message)
		}
		if err != nil {
			c.JSON(500, gin.H{
				"code": 1,
				"msg":  "An error occurred when sending message.",
			})
		}
		wg.Wait()
		//err = p.Shutdown()
		//fmt.Println("OpenHandler label 7, shutdown", err)
		//if err != nil {
		//	c.JSON(500, gin.H{
		//		"code": 1,
		//		"msg":  "An error occurred when closing producer.",
		//	})
		//}
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
			"message": "这个红包已经被打开了！",
		})
	}
}
