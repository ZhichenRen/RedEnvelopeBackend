package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"strconv"
)

var rdb *redis.Client

func initClient() (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "221.194.149.10:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err = rdb.Ping().Result()
	if err != nil {
		return err
	}
	fmt.Println("Connect successfully!")
	return nil
}

func redisExample() {
	err := rdb.Set("score", 100, 0).Err()
	if err != nil {
		fmt.Printf("set score failed, err:%v\n", err)
		return
	}

	val, err := rdb.Get("score").Result()
	if err != nil {
		fmt.Printf("get score failed, err:%v\n", err)
		return
	}
	fmt.Println("score", val)

	val2, err := rdb.Get("name").Result()
	if err == redis.Nil {
		fmt.Println("name does not exist")
	} else if err != nil {
		fmt.Printf("get name failed, err:%v\n", err)
		return
	} else {
		fmt.Println("name", val2)
	}
}

func SnatchHandler(c *gin.Context){
	userId, _ := c.GetPostForm("uid")
	fmt.Println(userId)
	envelopeId := 123
	maxCount := 5
	curCount := 3

	c.JSON(200, gin.H{
		"code": 0,
		"msg": "success",
		"data": gin.H{
			"envelop_id": envelopeId,
			"max_count": maxCount,
			"cur_count" : curCount,
		},
	})
}

func OpenHandler(c *gin.Context){
	userId, _ := c.GetPostForm("uid")
	fmt.Println(userId)

	value := 50

	c.JSON(200, gin.H{
		"code": 0,
		"msg": "success",
		"data": gin.H{
			"value": value,
		},
	})
}

func WalletListHandler(c *gin.Context){
	userId, _ := c.GetPostForm("uid")
	fmt.Println(userId)

	envelops := []gin.H{
		{
			"envelop_id": 123,
			"value": 50,
			"opened": true,
			"snatch_time": 1634551711,
		},
	}
	amount := 50

	c.JSON(200, gin.H{
		"code": 0,
		"msg": "success",
		"data": gin.H{
			"amount": amount,
			"envelope_list": envelops,
		},
	})
}

func Ping(c* gin.Context){
	countString, getError := rdb.Get("count").Result()
	var value int
	if getError == nil {
		count, countError := strconv.Atoi(countString)
		if countError != nil {
			fmt.Println(countError)
		}
		value = count + 1
		setError := rdb.Set("count", strconv.Itoa(value), 0).Err()
		if setError != nil {
			fmt.Println(setError)
		}
	} else {
		fmt.Println(getError)
		setError := rdb.Set("count", strconv.Itoa(1), 0).Err()
		if setError != nil {
			fmt.Println(setError)
		}
		value = 1
	}
	c.JSON(200, gin.H{
		"message": "这个网页已经被访问了" + strconv.Itoa(value) + "次！",
	})
}

func main() {
	initClient()
	r := gin.Default()
	r.GET("/ping", Ping)
	r.POST("/snatch", SnatchHandler)
	r.POST("/open", OpenHandler)
	r.POST("/get_wallet_list", WalletListHandler)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

