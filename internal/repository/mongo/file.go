package mongo

import (
	"context"

	"github.com/Woodfyn/file-api/internal/core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type FilesMongo struct {
	db *mongo.Collection
}

func NewFilesMongo(db *mongo.Database) *FilesMongo {
	return &FilesMongo{
		db: db.Collection(FILE_COLLECTION),
	}
}

func (f *FilesMongo) CreateFile(ctx context.Context, file *core.File) error {
	if _, err := f.db.InsertOne(ctx, file); err != nil {
		return err
	}

	return nil
}

func (f *FilesMongo) GetFileByName(ctx context.Context, name string) (*core.File, error) {
	response := new(core.File)

	result := f.db.FindOne(ctx, bson.M{"name": name})
	if err := result.Decode(response); err != nil {
		return nil, err
	}

	return response, nil
}
