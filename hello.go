package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"math/rand"
	"strconv"
	"time"
)

var rdb *redis.Client

type Envelope struct {
	EnvelopeId string
	Value int
	Opened bool
	SnatchTime int64
}

type User struct {
	CurCount int
	Amount int
	EnvelopeList []string
}

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

func SnatchHandler(c *gin.Context){
	userId, _ := c.GetPostForm("uid")
	fmt.Println(userId)

	randNumber := rand.Intn(100)
	//随机数判断用户是否抢到红包，后期需要替换
	if randNumber < 5 {
		snatchTime := time.Now().Unix()
		fmt.Println(snatchTime)
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
	} else {
		c.JSON(200, gin.H{
			"code": 0,
			"msg": "fail",
		})
	}
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
	//r := gin.Default()
	//r.GET("/ping", Ping)
	//r.POST("/snatch", SnatchHandler)
	//r.POST("/open", OpenHandler)
	//r.POST("/get_wallet_list", WalletListHandler)
	//r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	envelopeId , err := rdb.Get("EnvelopeId").Result()
	if err == nil {
		rdb.Incr("EnvelopeId")
		envelope := &Envelope{
			EnvelopeId: envelopeId,
			Value: 0,
			Opened: false,
			SnatchTime: time.Now().Unix(),
		}
		data, _ := json.Marshal(envelope)
		err := rdb.Set(envelopeId, data, -1).Err()
		if err != nil {
			fmt.Println(err)
		}
	}

	envelopeJson, err := rdb.Get("100000000000").Result()
	if err == nil {
		var envelope Envelope
		json.Unmarshal([]byte(envelopeJson), &envelope)
		fmt.Println(envelope)
	}

}

