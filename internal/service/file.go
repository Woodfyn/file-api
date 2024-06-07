package service

import (
	"context"

	"github.com/Woodfyn/file-api/internal/core"
	"github.com/Woodfyn/file-api/internal/repository/mongo"
	"github.com/Woodfyn/file-api/internal/repository/storage"
)

type FileService struct {
	mongoRepo   mongo.Files
	storageRepo storage.Files
}

func NewFileService(mongoRepo mongo.Files, storageRepo storage.Files) *FileService {
	return &FileService{
		mongoRepo:   mongoRepo,
		storageRepo: storageRepo,
	}
}

func (s *FileService) Upload(ctx context.Context, fileDTO *core.CreateFileDTO) error {
	file, err := core.NewFile(fileDTO)
	if err != nil {
		return err
	}

	if err := s.storageRepo.Upload(ctx, file); err != nil {
		return err
	}

	if err := s.mongoRepo.CreateFile(ctx, file); err != nil {
		return err
	}

	return nil
}

func (s *FileService) GetFiles(ctx context.Context) ([]*core.File, error) {
	files, err := s.storageRepo.GetFiles(ctx)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		filesMongo, err := s.mongoRepo.GetFileByName(ctx, file.Name)
		if err != nil {
			return nil, err
		}

		file.ID = filesMongo.ID
	}

	return nil, nil
}
