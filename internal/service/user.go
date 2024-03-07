package service

import (
	"context"
	"time"

	"github.com/Woodfyn/file-api/internal/core"
	"github.com/Woodfyn/file-api/internal/repository/mongo"
	"github.com/Woodfyn/file-api/internal/repository/redis"
	"github.com/Woodfyn/file-api/pkg/auth"
	"github.com/Woodfyn/file-api/pkg/hash"
)

type UserService struct {
	repo      mongo.Users
	redisRepo redis.Tokens
	hasher    hash.PasswordHasher

	tokenManager auth.TokenManager
	acssTokenTTL time.Duration
	refreshTTL   time.Duration
}

func NewUserService(repo mongo.Users, redisRepo redis.Tokens, hasher hash.PasswordHasher, tokenManager auth.TokenManager, acssTokenTTL, refreshTTL time.Duration) *UserService {
	return &UserService{
		repo:      repo,
		redisRepo: redisRepo,
		hasher:    hasher,

		tokenManager: tokenManager,
		acssTokenTTL: acssTokenTTL,
		refreshTTL:   refreshTTL,
	}
}

func (s *UserService) SignUp(ctx context.Context, inp core.SingUpRequest) error {
	hashPassword, err := s.hasher.Hash(inp.Password)
	if err != nil {
		return err
	}

	userCreate := core.User{
		Email:        inp.Email,
		Password:     hashPassword,
		RegisteredAt: time.Now(),
	}

	if err := s.repo.CreateUser(ctx, &userCreate); err != nil {
		return err
	}

	return err
}

func (s *UserService) SignIn(ctx context.Context, inp core.SingInRequest) (string, string, error) {
	return "", "", nil
}
