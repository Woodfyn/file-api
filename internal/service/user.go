package service

import (
	"context"
	"time"

	"github.com/Woodfyn/file-api/internal/core"
	"github.com/Woodfyn/file-api/internal/repository/mongo"
	"github.com/Woodfyn/file-api/internal/repository/redis"
	"github.com/Woodfyn/file-api/pkg/auth"
	"github.com/Woodfyn/file-api/pkg/hash"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService struct {
	repo      mongo.Users
	redisRepo redis.Tokens
	hasher    hash.PasswordHasher

	tokenManager   auth.TokenManager
	accessTokenTTL time.Duration
	refreshTTL     time.Duration
}

func NewUserService(repo mongo.Users, redisRepo redis.Tokens, hasher hash.PasswordHasher, tokenManager auth.TokenManager, acssTokenTTL, refreshTTL time.Duration) *UserService {
	return &UserService{
		repo:      repo,
		redisRepo: redisRepo,
		hasher:    hasher,

		tokenManager:   tokenManager,
		accessTokenTTL: acssTokenTTL,
		refreshTTL:     refreshTTL,
	}
}

func (s *UserService) SignUp(ctx context.Context, inp core.SignUpRequest) error {
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

	return nil
}

func (s *UserService) SignIn(ctx context.Context, inp core.SignInRequest) (string, string, error) {
	hashPassword, err := s.hasher.Hash(inp.Password)
	if err != nil {
		return "", "", err
	}

	user, err := s.repo.GetUser(ctx, hashPassword)
	if err != nil {
		return "", "", err
	}

	refreshToken, accessToken, err := s.genereteTokens(ctx, user.ID)
	if err != nil {
		return "", "", err
	}

	return refreshToken, accessToken, nil
}

func (s *UserService) genereteTokens(ctx context.Context, userId primitive.ObjectID) (string, string, error) {
	accessToken, err := s.tokenManager.NewJWT(userId.String(), s.accessTokenTTL)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.tokenManager.NewRefreshToken()
	if err != nil {
		return "", "", err
	}

	if err := s.redisRepo.SetTokenSession(ctx, userId.Hex(), refreshToken, s.refreshTTL); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *UserService) Refresh(ctx context.Context, refreshToken string) (string, error) {
	userId, err := s.redisRepo.GetTokenSession(ctx, refreshToken)
	if err != nil {
		return "", core.ErrTokenExpired
	}

	accessToken, err := s.tokenManager.NewJWT(userId, s.accessTokenTTL)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
