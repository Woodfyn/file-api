package rdb

import "github.com/go-redis/redis"

type ConnInfo struct {
	Addr     string
	Password string
}

func NewRdbClient(info ConnInfo) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     info.Addr,
		Password: info.Password,
	})
}
