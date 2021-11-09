package main

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"go-web/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"strconv"
	"time"
)

func main() {
	fmt.Println("Consumer start!")
	dsn := "group9:Group9@haha@tcp(rdsmysqlh1a4d645c087a17d2.rds.ivolces.com:3306)/red_envelope?charset=utf8&parseTime=True&loc=Local&timeout=10s"
	// connect to mysql
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
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
		for i := 0; i < len(msgs); i++ {
			switch string(msgs[i].Body) {
			case "create_envelope":
				params := msgs[i].GetProperties()
				var envelope utils.Envelope
				var user utils.User
				uid, err := strconv.Atoi(params["UID"])
				snatchTime, err := strconv.Atoi(params["SnatchTime"])
				value, err := strconv.Atoi(params["Value"])
				err = db.Where("cur_count < ?", 50).First(&user, utils.User{ID: int64(uid)}).Error
				if err == nil {
					envelope.UID = int64(uid)
					envelope.SnatchTime = int64(snatchTime)
					envelope.Value = value
					envelope.Opened = false
					db.Create(envelope)
					user.CurCount++
					db.Save(&user)
					fmt.Println(envelope)
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
		os.Exit(-1)
	}
	time.Sleep(time.Hour)
	err = client.Shutdown()
	if err != nil {
		fmt.Printf("Shutdown Consumer error: %s", err.Error())
	}
}
