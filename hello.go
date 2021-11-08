package main

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go-web/database"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var rdb *redis.Client

func SnatchHandler(c *gin.Context){
	userId, _ := c.GetPostForm("uid")
	user, err := rdb.HGetAll("User:" + userId).Result()
	if err != nil {
		fmt.Println(err)
	}
	// 后期与mysql对接，不存在时读取用户信息
	if len(user) == 0 {
		userInfo := make(map[string]interface{})
		userInfo["CurCount"] = 0
		userInfo["Amount"] = 0
		rdb.HMSet("User:" + userId, userInfo)
		user, err = rdb.HGetAll("User:" + userId).Result()
		if err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println("User:", user)
	maxCount, err := rdb.Get("MaxCount").Result()
	//随机数判断用户是否抢到红包，后期需要替换
	if user["CurCount"] < maxCount{
		snatchTime := time.Now().Unix()
		envelopeId, err := rdb.Incr("EnvelopeId").Result()
		if err != nil {
			fmt.Println(err)
		}
		curCount, err := rdb.HIncrBy("User:" + userId, "CurCount", 1).Result()
		if err != nil {
			fmt.Println(err)
		}

		envelope := make(map[string]interface{})
		envelope["Value"] = 0
		envelope["Opened"] = false
		envelope["SnatchTime"] = snatchTime
		_, err = rdb.HMSet("Envelope:" + strconv.Itoa(int(envelopeId)), envelope).Result()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Envelope:", envelope)

		_, err = rdb.SAdd("User:" + userId + ":Envelopes", envelopeId).Result()
		if err != nil {
			fmt.Println(err)
		}
		c.JSON(200, gin.H{
			"code": 0,
			"msg": "success",
			"data": gin.H{
				"envelop_id": envelopeId,
				"max_count": maxCount,
				"cur_count" : curCount,
			},
		})
	} else {
		c.JSON(200, gin.H{
			"code": 0,
			"msg": "fail",
		})
	}
}

func OpenHandler(c *gin.Context){
	userId, _ := c.GetPostForm("uid")
	envelopeId, _ := c.GetPostForm("envelope_id")

	result, err := rdb.SIsMember("User:" + userId + ":Envelopes", envelopeId).Result()
	if err != nil {
		fmt.Println(err)
	}
	if result == true {
		opened, err := rdb.HGet("Envelope:" + envelopeId, "Opened").Result()
		if err != nil {
			fmt.Println(err)
		}
		if opened != "0" {
			c.JSON(200, gin.H{
				"code": 0,
				"msg": "您已经打开了此红包",
			})
		}
		maxAmount, err := rdb.Get("MaxAmount").Int()
		value := rand.Intn(maxAmount)

		err = rdb.HSet("Envelope:" + envelopeId, "Opened", true).Err()
		err = rdb.HSet("Envelope:" + envelopeId, "Value", value).Err()
		err = rdb.HIncrBy("User:" + userId, "Amount", int64(value)).Err()

		c.JSON(200, gin.H{
			"code": 0,
			"msg": "success",
			"data": gin.H{
				"value": value,
			},
		})

	} else {
		c.JSON(200, gin.H{
			"code": 0,
			"msg": "您并不拥有此红包",
		})
	}
}

func WalletListHandler(c *gin.Context){
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

	amount, err := rdb.HGet("User:" + userId, "Amount").Result()

	c.JSON(200, gin.H{
		"code": 0,
		"msg": "success",
		"data": gin.H{
			"amount": amount,
			"envelope_list": data,
		},
	})
}

func Ping(c* gin.Context){
	_, err := rdb.Get("Count").Result()
	if err == nil {
		result, err := rdb.Incr("count").Result()
		if err != nil {
			fmt.Println(err)
		}
		c.JSON(200, gin.H{
			"message": "这个网页已经被访问了" + strconv.Itoa(int(result)) + "次。",
		})
	} else {
		rdb.Set("Count", 1, 0)
		c.JSON(200, gin.H{
			"message": "这个网页已经被访问了" + strconv.Itoa(1) + "次。",
		})
	}
}

func Producer(c *gin.Context) {
	addr,err := primitive.NewNamesrvAddr("http://MQ_INST_8149062485579066312_2586445845.cn-beijing.rocketmq-internal.ivolces.com:24009")
	if err != nil {
		panic(err)
	}
	topic := "Msg"
	p,err := rocketmq.NewProducer(
		producer.WithGroupName("GID_Group"),
		producer.WithNameServer(addr),
		producer.WithCreateTopicKey(topic),
		producer.WithRetry(1))
	if err != nil {
		panic(err)
	}

	err = p.Start()
	if err != nil {
		panic(err)
	}

	// 发送异步消息
	res,err := p.SendSync(context.Background(),primitive.NewMessage(topic,[]byte("send sync message")))
	if err != nil {
		fmt.Printf("send sync message error:%s\n",err)
	} else {
		fmt.Printf("send sync message success. result=%s\n",res.String())
	}

	// 发送消息后回调
	err = p.SendAsync(context.Background(), func(ctx context.Context, result *primitive.SendResult, err error) {
		if err != nil {
			fmt.Printf("receive message error:%v\n",err)
		} else {
			fmt.Printf("send message success. result=%s\n",result.String())
		}
	},primitive.NewMessage(topic,[]byte("send async message")))
	if err != nil {
		fmt.Printf("send async message error:%s\n",err)
	}

	// 批量发送消息
	var msgs []*primitive.Message
	for i := 0; i < 5; i++ {
		msgs = append(msgs, primitive.NewMessage(topic,[]byte("batch send message. num:"+strconv.Itoa(i))))
	}
	res,err = p.SendSync(context.Background(),msgs...)
	if err != nil {
		fmt.Printf("batch send sync message error:%s\n",err)
	} else {
		fmt.Printf("batch send sync message success. result=%s\n",res.String())
	}

	// 发送延迟消息
	msg := primitive.NewMessage(topic,[]byte("delay send message"))
	msg.WithDelayTimeLevel(3)
	res,err = p.SendSync(context.Background(),msg)
	if err != nil {
		fmt.Printf("delay send sync message error:%s\n",err)
	} else {
		fmt.Printf("delay send sync message success. result=%s\n",res.String())
	}

	// 发送带有tag的消息
	msg1 := primitive.NewMessage(topic,[]byte("send tag message"))
	msg1.WithTag("tagA")
	res,err = p.SendSync(context.Background(),msg1)
	if err != nil {
		fmt.Printf("send tag sync message error:%s\n",err)
	} else {
		fmt.Printf("send tag sync message success. result=%s\n",res.String())
	}

	err = p.Shutdown()
	if err != nil {
		panic(err)
	}
}

func Consumer(c *gin.Context) {
	topic := "Msg"

	// 消息主动推送给消费者
	c2,err := rocketmq.NewPushConsumer(
		consumer.WithGroupName("GID_Group"),
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{"http://MQ_INST_8149062485579066312_2586445845.cn-beijing.rocketmq-internal.ivolces.com:24009"})),
		consumer.WithConsumeFromWhere(consumer.ConsumeFromFirstOffset), // 选择消费时间(首次/当前/根据时间)
		consumer.WithConsumerModel(consumer.BroadCasting)) // 消费模式(集群消费:消费完其他人不能再读取/广播消费：所有人都能读)
	if err != nil {
		panic(err)
	}

	err = c2.Subscribe(
		topic,consumer.MessageSelector{
			Type: consumer.TAG,
			Expression: "*", // 可以 TagA || TagB
		},
		func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
			orderlyCtx,_ := primitive.GetOrderlyCtx(ctx)
			fmt.Printf("orderly context: %v\n",orderlyCtx)
			for i := range msgs {
				fmt.Printf("Subscribe callback: %v\n",msgs[i])
			}
			return consumer.ConsumeSuccess,nil
		})
	if err != nil {
		fmt.Printf("Subscribe error:%s\n",err)
	}

	err = c2.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	time.Sleep(time.Minute)
	err = c2.Shutdown()
	if err != nil {
		fmt.Println("Shutdown Consumer error: ",err)
	}
}

func main() {
	rdb = database.InitRedisClient(rdb)
	r := gin.Default()
	r.GET("/ping", Ping)
	r.POST("/snatch", SnatchHandler)
	r.POST("/open", OpenHandler)
	r.POST("/get_wallet_list", WalletListHandler)
	r.GET("/produce", Producer)
	r.GET("/consume", Consumer)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

