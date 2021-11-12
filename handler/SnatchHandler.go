package handler

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/gin-gonic/gin"
	"go-web/dao"
	"strconv"
	"sync"
	"time"
)

func SnatchHandler(c *gin.Context) {
	userId, _ := c.GetPostForm("uid")
	// string -> int64
	uid, err := strconv.ParseInt(userId, 10, 64)
	user, err := rdb.HGetAll("User:" + userId).Result()
	if err != nil {
		fmt.Println(err)
	}
	// TODO how to get maxCount
	// maxCount, err := rdb.Get("MaxCount").Result()
	// search in mysql
	if len(user) == 0 {
		users, err := dao.GetUser(uid)
		writeUserToRedis(users)
		user, err = rdb.HGetAll("User:" + userId).Result()
		if err != nil {
			c.JSON(500, gin.H{
				"code": 1,
				"msg": "A database error occurred.",
			})
			return
		}
	}

	// cheat detection
	snatchCount, err := rdb.Get("User:" + userId + ":Snatch").Int64()
	if snatchCount == 0 {
		err = rdb.Set("User:" + userId + ":Snatch", 1, 60000000000).Err()
	} else {
		snatchCount, err = rdb.Incr("User:" + userId + ":Snatch").Result()
		if snatchCount > 60 {
			c.JSON(403, gin.H{
				"code": 2,
				"msg": "系统检测到你在作弊！",
			})
			return
		}
	}
	if err != nil {
		c.JSON(500, gin.H{
			"code": 1,
			"msg": "A database error occurred.",
		})
		return
	}

	maxCount := 10
	curCount, _ := strconv.Atoi(user["cur_count"])
	if curCount < maxCount {
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

		// message queue
		p, err := rocketmq.NewProducer(
			producer.WithNsResolver(primitive.NewPassthroughResolver([]string{"http://100.64.247.138:24009"})),
			producer.WithRetry(2),
			producer.WithNamespace("MQ_INST_8149062485579066312_2586445845"),
			producer.WithCredentials(primitive.Credentials{
				AccessKey: "s7lec7baJkQeOBWS6Mb26vmV",
				SecretKey: "TiJYTqrIC7iLBK4UbpkgGJqM",
			}),
			producer.WithGroupName("GID_Group"),
		)
		if err != nil {
			fmt.Println("init producer error: " + err.Error())
			c.JSON(500, gin.H{
				"code": 1,
				"msg": "An error occurred when creating producer.",
			})
			return
		}
		err = p.Start()
		if err != nil {
			fmt.Printf("start producer error: %s", err.Error())
			c.JSON(500, gin.H{
				"code": 1,
				"msg": "An error occurred when starting producer.",
			})
			return
		}
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
					fmt.Printf("receive message error: %s\n", e)
				} else {
					fmt.Printf("send message success: result=%s\n", result.String())
				}
				wg.Done()
			}, message)
		if err != nil {
			fmt.Printf("An error occurred when sending message: %s\n", err)
			fmt.Println(message)
			c.JSON(500, gin.H{
				"code": 1,
				"msg": "An error occurred when sending message.",
			})
			return
		}
		wg.Wait()
		err = p.Shutdown()
		if err != nil {
			c.JSON(500, gin.H{
				"code": 1,
				"msg": "An error occurred when closing producer.",
			})
			return
		}
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
