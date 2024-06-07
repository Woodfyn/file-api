package mongo

import (
	"context"

	"github.com/Woodfyn/file-api/internal/core"
	"go.mongodb.org/mongo-driver/mongo"
)

type Users interface {
	CreateUser(ctx context.Context, user *core.User) error
	GetUser(ctx context.Context, password string) (*core.User, error)
}

type Files interface {
	CreateFile(ctx context.Context, file *core.File) error
	GetFileByName(ctx context.Context, name string) (*core.File, error)
}

type Repository struct {
	Users Users
	Files Files
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		Users: NewUsersMongo(db),
		Files: NewFilesMongo(db),
	}
}
