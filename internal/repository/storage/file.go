package storage

import (
	"bytes"
	"context"
	"io"
	"log/slog"

	"cloud.google.com/go/storage"
	"github.com/Woodfyn/file-api/internal/core"
)

type FileStorage struct {
	bucket *storage.BucketHandle
}

func NewFileStorage(bucket *storage.BucketHandle) *FileStorage {
	return &FileStorage{
		bucket: bucket,
	}
}

func (f *FileStorage) Upload(ctx context.Context, file *core.File) error {
	object := f.bucket.Object(file.Name)
	writer := object.NewWriter(ctx)
	defer writer.Close()

	if _, err := io.Copy(writer, bytes.NewReader(file.Bytes)); err != nil {
		slog.Error("failed to upload file", "err", err)
		return err
	}

	if err := object.ACL().Set(context.Background(), storage.AllUsers, storage.RoleReader); err != nil {
		slog.Error("failed to set ACL", "err", err)
		return err
	}

	return nil
}

func (f *FileStorage) GetFiles(ctx context.Context) ([]*core.File, error) {

	return nil, nil
}
