package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-web/handler"
)

func main() {
	err := handler.InitClient()
	if err != nil {
		fmt.Println("Connection failed")
		return
	}

	fmt.Println("Connection succeeded")
	r := gin.Default()
	//r.GET("/ping", handler.Ping)
	r.POST("/snatch", handler.SnatchHandler)
	r.POST("/open", handler.OpenHandler)
	r.POST("/get_wallet_list", handler.WalletListHandler)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
