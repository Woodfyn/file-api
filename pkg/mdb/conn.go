package mdb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ConnInfo struct {
	URI      string
	Username string
	Password string
}

func NewMongoClient(ctx context.Context, info ConnInfo) (*mongo.Client, error) {
	opts := options.Client()
	opts.SetAuth(options.Credential{
		Username: info.Username,
		Password: info.Password,
	})
	opts.ApplyURI(info.URI)

	dbClient, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	if err := dbClient.Ping(context.Background(), nil); err != nil {
		return nil, err
	}

	return dbClient, nil
}
