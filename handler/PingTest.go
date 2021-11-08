package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

func Ping(c *gin.Context) {
	_, err := rdb.Get("Count").Result()
	if err == nil {
		result, err := rdb.Incr("Count").Result()
		if err != nil {
			fmt.Println(err)
		}
		c.JSON(200, gin.H{
			"message": "这个网页已经被访问了" + strconv.Itoa(int(result)) + "次。",
		})
	} else {
		rdb.Set("Count", 1, 0)
		c.JSON(200, gin.H{
			"message": "这个网页已经被访问了" + strconv.Itoa(1) + "次。",
		})
	}
}
