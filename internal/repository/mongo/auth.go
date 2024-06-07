package mongo

import (
	"context"

	"github.com/Woodfyn/file-api/internal/core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Auth struct {
	db *mongo.Collection
}

func NewAuth(db *mongo.Database) *Auth {
	return &Auth{
		db: db.Collection(USER_COLLECTION),
	}
}

func (a *Auth) CreateUser(ctx context.Context, user *core.User) error {
	_, err := a.db.InsertOne(ctx, user)

	return err
}

func (a *Auth) GetUserByPassword(ctx context.Context, password string) (*core.User, error) {
	response := new(core.User)
	result := a.db.FindOne(ctx, bson.M{"password": password})
	if err := result.Err(); err != nil {
		return nil, err
	}

	if err := result.Decode(response); err != nil {
		return nil, err
	}

	return response, nil
}
