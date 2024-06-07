package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis"
)

type Tokens interface {
	SetTokenSession(ctx context.Context, userId string, token string, ttl time.Duration) error
	GetTokenSession(ctx context.Context, userId string) (string, error)
}

type Repository struct {
	Tokens Tokens
}

func NewRepository(db *redis.Client) *Repository {
	return &Repository{
		Tokens: NewTokenRedis(db),
	}
}
