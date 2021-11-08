package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-web/handler"
	"go-web/utils"
)

func main() {
	err := handler.InitClient()
	if err != nil {
		fmt.Println("Connection failed")
		return
	}

	fmt.Println("Connection succeeded")
	_, _, err = utils.CreateEnvelope(int64(123))
	fmt.Println(err)
	r := gin.Default()
	//r.GET("/ping", handler.Ping)
	r.POST("/snatch", handler.SnatchHandler)
	//r.POST("/open", OpenHandler)
	//r.POST("/get_wallet_list", WalletListHandler)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
