package main

import (
	"fmt"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"go-web/allocate"
	"go-web/handler"
	"go-web/tokenbucket"
	"time"
)

func main() {
	allocate.Init()
	handler.InitClient()
	handler.InitProducer()
	fmt.Println(handler.GetProducer())
	r := gin.Default()
	r.Use(tokenbucket.NewLimiter(2000, 2000, 500*time.Millisecond))
	pprof.Register(r)
	r.GET("/ping", handler.Ping)
	r.POST("/snatch", handler.SnatchHandler)
	r.POST("/open", handler.OpenHandler)
	r.POST("/get_wallet_list", handler.WalletListHandler)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	handler.CloseProducer()
}
