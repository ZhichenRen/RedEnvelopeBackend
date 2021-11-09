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
	"sync"
	"time"
)

func Producer(c *gin.Context) {
	fmt.Println("Producer start!")
	topic := "Msg"
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
		os.Exit(0)
	}
	err = p.Start()
	if err != nil {
		fmt.Printf("start producer error: %s", err.Error())
		os.Exit(1)
	}
	var wg sync.WaitGroup
	// try to create 10 envelopes for a user
	for i := 0; i < 100; i++ {
		var params map[string]string
		params["UID"] = "1234"
		params["Value"] = "100"
		params["SnatchTime"] = strconv.Itoa(int(time.Now().Unix()))
		message := primitive.NewMessage(topic, []byte("create_envelope"))
		message.WithProperties(params)
		wg.Add(1)
		err := p.SendAsync(context.Background(),
			func(ctx context.Context, result *primitive.SendResult, e error) {
				if e != nil {
					fmt.Printf("receive message error: %s\n", err)
				} else {
					fmt.Printf("send message success: result=%s\n", result.String())
				}
				wg.Done()
			}, message)

		if err != nil {
			fmt.Printf("send message error: %s\n", err)
			c.JSON(400, gin.H{
				"msg": "An error happened when sending messages.",
			})
			return
		}
	}
	wg.Wait()
	err = p.Shutdown()
	if err != nil {
		fmt.Printf("shutdown producer error: %s", err.Error())
		c.JSON(400, gin.H{
			"msg": "An error happened when shutting down producer.",
		})
		return
	}
	c.JSON(200, gin.H{
		"msg": "success",
	})
}

func Consumer(c *gin.Context) {
	fmt.Println("Consumer start!")
	client, err := rocketmq.NewPushConsumer(
		consumer.WithGroupName("GID_Group"),
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{"http://100.64.247.138:24009"})),
		consumer.WithNamespace("MQ_INST_8149062485579066312_2586445845"),
		consumer.WithCredentials(primitive.Credentials{
			AccessKey: "s7lec7baJkQeOBWS6Mb26vmV",
			SecretKey: "TiJYTqrIC7iLBK4UbpkgGJqM",
		}),
	)
	if err != nil {
		fmt.Println("init consumer error: " + err.Error())
		os.Exit(0)
	}

	err = client.Subscribe("Msg", consumer.MessageSelector{}, func(ctx context.Context,
		msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		fmt.Printf("subscribe callback: %v \n", msgs)
		return consumer.ConsumeSuccess, nil
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	// Note: start after subscribe
	err = client.Start()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	time.Sleep(time.Hour)
	err = client.Shutdown()
	if err != nil {
		fmt.Printf("Shutdown Consumer error: %s", err.Error())
	}
}