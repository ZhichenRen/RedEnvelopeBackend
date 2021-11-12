package main

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"go-web/dao"
	"strconv"
	"time"
)

func main() {
	fmt.Println("Consumer start!")
	db := dao.GetDB()
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
	}

	err = client.Subscribe("Msg", consumer.MessageSelector{}, func(ctx context.Context,
		msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		fmt.Printf("subscribe callback: %v \n", msgs)
		for i := 0; i < len(msgs); i++ {
			switch string(msgs[i].Body) {
			case "create_envelope":
				params := msgs[i].GetProperties()
				uid, err := strconv.Atoi(params["UID"])
				id, err := strconv.Atoi(params["EID"])
				snatchTime, err := strconv.Atoi(params["SnatchTime"])
				value, err := strconv.Atoi(params["Value"])
				if err == nil {
					envelope := dao.Envelope{UID: int64(uid), ID: int64(id), Opened: false, SnatchTime: int64(snatchTime), Value: value}
					db.Create(&envelope)
					fmt.Println(envelope)
					err := dao.UpdateCurCount(int64(uid))
					for ;err != nil; {
						err = dao.UpdateCurCount(int64(uid))
					}
				} else {
					fmt.Println("An error happened when writing database.")
				}
			case "open_envelope":
				params := msgs[i].GetProperties()
				uid, err := strconv.Atoi(params["UID"])
				eid, err := strconv.Atoi(params["EID"])
				user, err := dao.GetUser(int64(uid))
				envelope, err := dao.GetEnvelopeByEID(int64(eid))
				if err == nil {
					user.Amount += envelope.Value
					envelope.Opened = true
					db.Save(&user)
					db.Save(&envelope)
				} else {
					fmt.Println("An error happened when writing database.")
				}
			default:
				fmt.Println("Unknown message body.")
			}
		}
		return consumer.ConsumeSuccess, nil
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	// Note: start after subscribe
	err = client.Start()
	if err != nil {
		fmt.Println(err.Error())
	}
	time.Sleep(time.Hour)
	err = client.Shutdown()
	if err != nil {
		fmt.Printf("Shutdown Consumer error: %s", err.Error())
	}
}
