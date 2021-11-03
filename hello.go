package main

import (
	"fmt"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
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

func main() {
	initClient()
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
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
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
