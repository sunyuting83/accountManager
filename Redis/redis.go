package redis

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	MyRedis *redis.Client
	ctx     = context.Background()
)

func InitRedis(Host, Password string, DB int) {
	MyRedis = redis.NewClient(&redis.Options{
		Addr:     Host,
		Password: Password, // no password set
		DB:       DB,       // use default DB
	})
	_, err := MyRedis.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Redis connect ping failed, err:", err)
		os.Exit(1)
		return
	}
	fmt.Println("Redis connect succeeded")
}
func Set(key string, value string, t int64) bool {
	expire := time.Duration(t) * time.Second
	if err := MyRedis.Set(ctx, key, value, expire).Err(); err != nil {
		return false
	}
	return true
}

func Get(key string) string {
	result, err := MyRedis.Get(ctx, key).Result()
	if err != nil {
		return ""
	}
	return result
}

func Delete(key string) bool {
	_, err := MyRedis.Del(ctx, key).Result()
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func ExpireRedis(key string, t int64) bool {
	expire := time.Duration(t) * time.Second
	if err := MyRedis.Expire(ctx, key, expire).Err(); err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
