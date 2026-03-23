package config

import (
	"fmt"
	"sync"

	"github.com/go-redis/redis"
)

var initRedisOnce sync.Once
var isRedisInit = false
var RedisClient *redis.Client

// SetupRedis 初始化连接
func SetupRedis(cfg *Redis) (err error) {
	if isRedisInit != false {
		return nil
	}
	//fmt.Println("redis init")
	//fmt.Println(cfg)
	initRedisOnce.Do(func() {
		RedisClient = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
			Password: cfg.Password,
			DB:       cfg.DataBase,
		})
		isRedisInit = true
	})
	_, err = RedisClient.Ping().Result()
	//fmt.Println("redis init 结束...")
	return err
}

func CloseRedis(client *redis.Client) {
	_ = client.Close()
}

func GetRedis() *redis.Client {
	return RedisClient
}
