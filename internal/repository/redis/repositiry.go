package redis

import (
	"context"

	"github.com/go-redis/redis"
)

type Tokens interface {
	SetTokenSession(ctx context.Context, userId int, token string) error
	CheckTokenSession(ctx context.Context, userId int, token string) (bool, error)
}

type Repository struct {
	Tokens Tokens
}

func NewRepository(db *redis.Client) *Repository {
	return &Repository{
		Tokens: NewTokenRedis(db),
	}
}
