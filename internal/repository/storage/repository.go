package storage

import (
	"context"

	"cloud.google.com/go/storage"
	"github.com/Woodfyn/file-api/internal/core"
)

type Files interface {
	Upload(ctx context.Context, file *core.File) error
	GetFiles(ctx context.Context) ([]*core.File, error)
}

type Repository struct {
	Files Files
}

func NewRepository(storage *storage.BucketHandle) *Repository {
	return &Repository{
		Files: NewFileStorage(storage),
	}
}
