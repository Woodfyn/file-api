package mongo

import (
	"context"

	"github.com/Woodfyn/file-api/internal/core"
	"go.mongodb.org/mongo-driver/mongo"
)

type FilesMongo struct {
	db *mongo.Client
}

func NewFilesMongo(db *mongo.Client) *FilesMongo {
	return &FilesMongo{db: db}
}

func (f *FilesMongo) CreateFile(ctx context.Context, file *core.File) error {
	return nil
}

func (f *FilesMongo) GetFiles(ctx context.Context) (*[]core.File, error) {
	return nil, nil
}
