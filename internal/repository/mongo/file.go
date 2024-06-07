package mongo

import (
	"context"

	"github.com/Woodfyn/file-api/internal/core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type File struct {
	db *mongo.Collection
}

func NewFile(db *mongo.Database) *File {
	return &File{
		db: db.Collection(FILE_COLLECTION),
	}
}

func (f *File) CreateFile(ctx context.Context, file *core.File) error {
	_, err := f.db.InsertOne(ctx, file)

	return err
}

func (f *File) GetFiles(ctx context.Context, userId string) ([]*core.File, error) {
	var response []*core.File
	cursor, err := f.db.Find(ctx, bson.M{"user_id": userId})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &response); err != nil {
		return nil, err
	}

	return response, nil
}
