package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-web/allocate"
	"go-web/handler"
	"go-web/tokenbucket"
	"time"
)

func setupRouter() *gin.Engine{
	allocate.Init()
	handler.InitClient()
	handler.InitProducer()
	fmt.Println(handler.GetProducer())
	r := gin.Default()
	r.Use(tokenbucket.NewLimiter(10000, 10000, 500*time.Millisecond))
	r.GET("/ping", handler.Ping)
	r.POST("/snatch", handler.SnatchHandler)
	r.POST("/open", handler.OpenHandler)
	r.POST("/get_wallet_list", handler.WalletListHandler)
	return r
}

func main() {
	r := setupRouter()
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	handler.CloseProducer()
}
