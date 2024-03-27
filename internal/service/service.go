package service

import (
	"context"
	"time"

	"github.com/Woodfyn/file-api/internal/core"
	"github.com/Woodfyn/file-api/internal/repository/mongo"
	"github.com/Woodfyn/file-api/internal/repository/redis"
	"github.com/Woodfyn/file-api/internal/repository/storage"
	"github.com/Woodfyn/file-api/pkg/auth"
	"github.com/Woodfyn/file-api/pkg/hash"
)

type Files interface {
	Upload(ctx context.Context, file *core.CreateFileDTO) error
	GetFiles(ctx context.Context) ([]*core.File, error)
}

type Users interface {
	SignUp(ctx context.Context, inp core.SignUpRequest) error
	SignIn(ctx context.Context, inp core.SignInRequest) (string, string, error)
	Refresh(ctx context.Context, refreshToken string) (string, error)
}

type Service struct {
	Files Files
	Users Users
}

type Deps struct {
	MongoRepo   *mongo.Repository
	RedisRepos  *redis.Repository
	StorageRepo *storage.Repository
	Hasher      hash.PasswordHasher

	TokenManager auth.TokenManager
	AcssTokenTTL time.Duration
	RefreshTTL   time.Duration
}

func NewService(deps Deps) *Service {
	return &Service{
		Files: NewFileService(deps.MongoRepo.Files, deps.StorageRepo.Files),
		Users: NewUserService(deps.MongoRepo.Users, deps.RedisRepos.Tokens, deps.Hasher, deps.TokenManager, deps.AcssTokenTTL, deps.RefreshTTL),
	}
}
