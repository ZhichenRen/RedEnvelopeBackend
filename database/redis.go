package database

import (
	"fmt"
	"github.com/go-redis/redis"
)

//var rdb *redis.Client

func InitRedisClient(rdb *redis.Client) *redis.Client {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "221.194.149.10:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := rdb.Ping().Result()
	if err != nil {
		return nil
	}
	fmt.Println("Connect successfully!")
	return rdb
}
