package mongo

import (
	"context"

	"github.com/Woodfyn/file-api/internal/core"
	"go.mongodb.org/mongo-driver/mongo"
)

type UsersMongo struct {
	db *mongo.Client
}

func NewUsersMongo(db *mongo.Client) *UsersMongo {
	return &UsersMongo{db: db}
}

func (u *UsersMongo) CreateUser(ctx context.Context, user *core.User) error {
	return nil
}

func (u *UsersMongo) GetUser(ctx context.Context, email string) (*core.User, error) {
	return nil, nil
}
