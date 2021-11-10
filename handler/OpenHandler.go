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
)

func OpenHandler(c *gin.Context) {
	userId, _ := c.GetPostForm("uid")
	envelopeId, _ := c.GetPostForm("envelope_id")
	uid, _ := strconv.ParseInt(userId, 10, 64)
	//eid, _ := strconv.ParseInt(envelopeId, 10, 64)
	envelopes, err := rdb.HGetAll("Envelope:" + envelopeId).Result()
	if err != nil {
		c.JSON(500, gin.H{
			"code": 1,
			"msg": "An error occurred when reading from redis",
		})
		return
	}
	if len(envelopes) == 0 {
		// key in redis expired, read from mysql
		envelopeList, err := dao.GetEnvelopesByUID(uid)
		// error or envelope not found
		if err != nil {
			c.JSON(500, gin.H{
				"code": 1,
				"msg": "A database error occurred or the envelope didn't exist.",
			})
			return
		}
		// write envelope set and envelopes
		for _, e := range envelopeList {
			writeEnvelopeToRedis(*e)
			writeEnvelopesSet(*e, userId)
		}
		envelopes, err = rdb.HGetAll("Envelope:" + envelopeId).Result()
		if err != nil {
			c.JSON(500, gin.H{
				"code": 1,
				"msg": "An error occurred when reading from redis",
			})
			return
		}
	}
	opened := envelopes["opened"]
	value := envelopes["value"]
	realUId := envelopes["uid"]
	if err != nil {
		fmt.Println(err)
	}
	if userId != realUId {
		c.JSON(403, gin.H{
			"code": 2,
			"message": "no authorization",
		})
		return
	}
	if opened == "0" {
		// write to redis
		rdb.HSet("Envelope:"+envelopeId, "opened", true)
		updateAmount(userId, value)
		// TODO
		// write to MySQL
		// OpenEnvelope should be deleted
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
			c.JSON(500, gin.H{
				"code": 1,
				"msg": "An error occurred when sending message.",
			})
		}
		wg.Wait()
		err = p.Shutdown()
		if err != nil {
			c.JSON(500, gin.H{
				"code": 1,
				"msg": "An error occurred when closing producer.",
			})
		}
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
