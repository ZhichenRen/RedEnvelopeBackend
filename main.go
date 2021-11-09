package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-web/allocate"
	"go-web/handler"
)

func main() {
	test()
	allocate.Init()
	handler.InitClient()
	r := gin.Default()
	r.GET("/ping", handler.Ping)
	r.POST("/snatch", handler.SnatchHandler)
	r.POST("/open", handler.OpenHandler)
	r.POST("/get_wallet_list", handler.WalletListHandler)
	r.GET("/produce", handler.Producer)
	r.GET("/pull_consume", handler.PullConsumer)
	r.GET("/push_consume", handler.PushConsumer)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func test() {
	config, err := NewConfig("config.yaml")
	if err == nil {
		fmt.Println(config)
	} else {
		fmt.Println(err)
	}
}

