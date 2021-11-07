package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go-web/database"
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
	user, err := rdb.HGetAll("User:" + userId).Result()
	if err != nil {
		fmt.Println(err)
	}
	// 后期与mysql对接，不存在时读取用户信息
	if len(user) == 0 {
		userInfo := make(map[string]interface{})
		userInfo["CurCount"] = 0
		userInfo["Amount"] = 0
		rdb.HMSet("User:" + userId, userInfo)
		user, err = rdb.HGetAll("User:" + userId).Result()
		if err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println("User:", user)
	maxCount, err := rdb.Get("MaxCount").Result()
	//随机数判断用户是否抢到红包，后期需要替换
	if user["CurCount"] < maxCount{
		snatchTime := time.Now().Unix()
		envelopeId, err := rdb.Incr("EnvelopeId").Result()
		if err != nil {
			fmt.Println(err)
		}
		curCount, err := rdb.HIncrBy("User:" + userId, "CurCount", 1).Result()
		if err != nil {
			fmt.Println(err)
		}

		envelope := make(map[string]interface{})
		envelope["Value"] = 0
		envelope["Opened"] = false
		envelope["SnatchTime"] = snatchTime
		_, err = rdb.HMSet("Envelope:" + strconv.Itoa(int(envelopeId)), envelope).Result()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Envelope:", envelope)

		_, err = rdb.LPush("User:" + userId + ":Envelopes", envelopeId).Result()
		if err != nil {
			fmt.Println(err)
		}
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
	_, err := rdb.Get("Count").Result()
	if err == nil {
		result, err := rdb.Incr("count").Result()
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

func main() {
	rdb = database.InitRedisClient(rdb)
	r := gin.Default()
	r.GET("/ping", Ping)
	r.POST("/snatch", SnatchHandler)
	r.POST("/open", OpenHandler)
	r.POST("/get_wallet_list", WalletListHandler)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	//envelopeId , err := rdb.Get("EnvelopeId").Result()
	//if err == nil {
	//	rdb.Incr("EnvelopeId")
	//	envelope := &Envelope{
	//		EnvelopeId: envelopeId,
	//		Value: 0,
	//		Opened: false,
	//		SnatchTime: time.Now().Unix(),
	//	}
	//	data, _ := json.Marshal(envelope)
	//	err := rdb.Set(envelopeId, data, -1).Err()
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//}
	//
	//envelopeJson, err := rdb.Get("100000000000").Result()
	//if err == nil {
	//	var envelope Envelope
	//	json.Unmarshal([]byte(envelopeJson), &envelope)
	//	fmt.Println(envelope)
	//}

}

