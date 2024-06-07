package mongo

import (
	"context"

	"github.com/Woodfyn/file-api/internal/core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UsersMongo struct {
	db *mongo.Collection
}

func NewUsersMongo(db *mongo.Database) *UsersMongo {
	return &UsersMongo{
		db: db.Collection(USER_COLLECTION),
	}
}

func (u *UsersMongo) CreateUser(ctx context.Context, user *core.User) error {
	if _, err := u.db.InsertOne(ctx, user); err != nil {
		return err
	}

	return nil
}

func (u *UsersMongo) GetUser(ctx context.Context, password string) (*core.User, error) {
	response := new(core.User)

	result := u.db.FindOne(ctx, bson.M{"password": password})
	if err := result.Decode(response); err != nil {
		return nil, err
	}

	return response, nil
}
