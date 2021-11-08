package handler

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/gin-gonic/gin"
	"os"
	"strconv"
	"time"
)

func Producer(c *gin.Context) {
	fmt.Println("Producer start!")
	addr,err := primitive.NewNamesrvAddr("http://100.64.247.138:24009")
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
	fmt.Println("Producer end!")
	c.JSON(200, gin.H{
		"msg": "success",
	})
}

func PullConsumer(c *gin.Context) {
	fmt.Println("Consumer start!")
	topic := "Msg"

	// 消费者主动拉取消息
	// not
	c1,err := rocketmq.NewPullConsumer(
		consumer.WithGroupName("GID_Group"),
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{"http://100.64.247.138:24009"})))
	if err != nil {
		panic(err)
	}
	err = c1.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	queue := primitive.MessageQueue{
		Topic:      topic,
		BrokerName: "broker-a", // 使用broker的名称
		QueueId:    0,
	}

	err = c1.Shutdown()
	if err != nil {
		fmt.Println("Shutdown Pull Consumer error: ",err)
	}

	offset := int64(0)
	for  {
		resp,err := c1.PullFrom(context.Background(),queue,offset,10)
		if err != nil {
			if err == rocketmq.ErrRequestTimeout {
				fmt.Printf("timeout\n")
				time.Sleep(time.Second)
				continue
			}
			fmt.Printf("unexpected error: %v\n",err)
			return
		}
		if resp.Status == primitive.PullFound {
			fmt.Printf("pull message success. nextOffset: %d\n",resp.NextBeginOffset)
			for _, ext := range resp.GetMessageExts() {
				fmt.Printf("pull msg: %s\n",ext)
			}
		}
		offset = resp.NextBeginOffset
	}
}

func PushConsumer(c *gin.Context) {
	fmt.Println("Consumer start")
	topic := "Msg"

	// 消息主动推送给消费者
	c2,err := rocketmq.NewPushConsumer(
		consumer.WithGroupName("GID_Group"),
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{"http://100.64.247.138:24009"})),
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