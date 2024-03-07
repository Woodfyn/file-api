package redis

import (
	"context"

	"github.com/go-redis/redis"
)

type TokenRedis struct {
	db *redis.Client
}

func NewTokenRedis(db *redis.Client) *TokenRedis {
	return &TokenRedis{db: db}
}

func (t *TokenRedis) SetTokenSession(ctx context.Context, userId int, token string) error {
	return nil
}

func (t *TokenRedis) CheckTokenSession(ctx context.Context, userId int, token string) (bool, error) {
	return false, nil
}
