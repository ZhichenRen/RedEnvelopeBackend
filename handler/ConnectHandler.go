package handler

import (
	"github.com/go-redis/redis"
)

var rdb *redis.Client
var number int64

func InitClient() (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "221.194.149.10:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	number = 0
	_, err = rdb.Ping().Result()
	return err
}
