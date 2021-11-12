package handler

import (
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

var p rocketmq.Producer

func InitProducer() {
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
		return
	}
	err = p.Start()
	if err != nil {
		fmt.Printf("start producer error: %s", err.Error())
		return
	}
}

func CloseProducer() {
	err := p.Shutdown()
	if err != nil {
		fmt.Printf("An error occurred when closing producer: %s\n", err)
		return
	}
}