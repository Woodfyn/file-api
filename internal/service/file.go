package service

import (
	"context"
	"io"

	"github.com/Woodfyn/file-api/internal/core"
	"github.com/Woodfyn/file-api/internal/repository/mongo"
)

type FileService struct {
	mongoRepo mongo.Files
}

func NewFileService(mongoRepo mongo.Files) *FileService {
	return &FileService{
		mongoRepo: mongoRepo,
	}
}

func (s *FileService) Upload(ctx context.Context, file *core.CreateFileDTO) error {
	return nil
}

func (s *FileService) GetFiles(ctx context.Context, w io.Writer) ([]*core.File, error) {
	return nil, nil
}
