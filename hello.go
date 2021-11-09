package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go-web/database"
	"math/rand"
	"strconv"
	"time"
)

var rdb *redis.Client

type Envelope struct {
	EnvelopeId string
	Value      int
	Opened     bool
	SnatchTime int64
}

type User struct {
	CurCount     int
	Amount       int
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

func SnatchHandler(c *gin.Context) {
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
		rdb.HMSet("User:"+userId, userInfo)
		user, err = rdb.HGetAll("User:" + userId).Result()
		if err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println("User:", user)
	maxCount, err := rdb.Get("MaxCount").Result()
	//随机数判断用户是否抢到红包，后期需要替换
	if user["CurCount"] < maxCount {
		snatchTime := time.Now().Unix()
		envelopeId, err := rdb.Incr("EnvelopeId").Result()
		if err != nil {
			fmt.Println(err)
		}
		curCount, err := rdb.HIncrBy("User:"+userId, "CurCount", 1).Result()
		if err != nil {
			fmt.Println(err)
		}

		envelope := make(map[string]interface{})
		envelope["Value"] = 0
		envelope["Opened"] = false
		envelope["SnatchTime"] = snatchTime
		_, err = rdb.HMSet("Envelope:"+strconv.Itoa(int(envelopeId)), envelope).Result()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Envelope:", envelope)

		_, err = rdb.SAdd("User:"+userId+":Envelopes", envelopeId).Result()
		if err != nil {
			fmt.Println(err)
		}
		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "success",
			"data": gin.H{
				"envelop_id": envelopeId,
				"max_count":  maxCount,
				"cur_count":  curCount,
			},
		})
	} else {
		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "fail",
		})
	}
}

func OpenHandler(c *gin.Context) {
	userId, _ := c.GetPostForm("uid")
	envelopeId, _ := c.GetPostForm("envelope_id")

	result, err := rdb.SIsMember("User:"+userId+":Envelopes", envelopeId).Result()
	if err != nil {
		fmt.Println(err)
	}
	if result == true {
		opened, err := rdb.HGet("Envelope:"+envelopeId, "Opened").Result()
		if err != nil {
			fmt.Println(err)
		}
		if opened != "0" {
			c.JSON(200, gin.H{
				"code": 0,
				"msg":  "您已经打开了此红包",
			})
		}
		maxAmount, err := rdb.Get("MaxAmount").Int()
		value := rand.Intn(maxAmount)

		err = rdb.HSet("Envelope:"+envelopeId, "Opened", true).Err()
		err = rdb.HSet("Envelope:"+envelopeId, "Value", value).Err()
		err = rdb.HIncrBy("User:"+userId, "Amount", int64(value)).Err()

		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "success",
			"data": gin.H{
				"value": value,
			},
		})

	} else {
		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "您并不拥有此红包",
		})
	}
}

func WalletListHandler(c *gin.Context) {
	userId, _ := c.GetPostForm("uid")
	fmt.Println(userId)

	envelopeList, err := rdb.SMembers("User:" + userId + ":Envelopes").Result()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(envelopeList)
	var data []gin.H
	for _, envelopeId := range envelopeList {
		envelope, err := rdb.HGetAll("Envelope:" + envelopeId).Result()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(envelope)
		tmp := gin.H{}
		tmp["envelope_id"] = envelopeId
		tmp["snatch_time"] = envelope["SnatchTime"]
		if envelope["Opened"] == "0" {
			tmp["opened"] = false
		} else {
			tmp["opened"] = true
			tmp["value"] = envelope["Value"]
		}
		data = append(data, tmp)
	}

	amount, err := rdb.HGet("User:"+userId, "Amount").Result()

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"amount":        amount,
			"envelope_list": data,
		},
	})
}

func Ping(c *gin.Context) {
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
}
