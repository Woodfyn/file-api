package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis"
)

type TokenRedis struct {
	db *redis.Client
}

func NewTokenRedis(db *redis.Client) *TokenRedis {
	return &TokenRedis{db: db}
}

func (t *TokenRedis) SetTokenSession(ctx context.Context, userId string, token string, ttl time.Duration) error {
	if err := t.db.Set(token, userId, ttl).Err(); err != nil {
		return err
	}

	return nil
}

func (t *TokenRedis) GetTokenSession(ctx context.Context, token string) (string, error) {
	id, err := t.db.Get(token).Result()
	if err != nil {
		return "", err
	}

	return id, err
}
