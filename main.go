package main

import (
	"fmt"
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
	limit1 := tokenbucket.NewLimiter(5000, 5000, 500*time.Millisecond)
	limit2 := tokenbucket.NewLimiter(5000, 5000, 500*time.Millisecond)
	limit3 := tokenbucket.NewLimiter(2000, 2000, 500*time.Millisecond)
	//pprof.Register(r)
	r.GET("/ping", handler.Ping)
	r.POST("/snatch", limit1, handler.SnatchHandler)
	r.POST("/open", limit2, handler.OpenHandler)
	r.POST("/get_wallet_list", limit3, handler.WalletListHandler)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	handler.CloseProducer()
}
