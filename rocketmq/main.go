package main

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"os"
	"time"
)

func main() {
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
