package service

import (
	"context"
	"log/slog"
	"mime/multipart"
	"time"

	"github.com/Woodfyn/file-api/internal/core"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
)

type fileMongo interface {
	CreateFile(ctx context.Context, file *core.File) error
	GetFiles(ctx context.Context, userId string) ([]*core.File, error)
}

type fileS3 interface {
	UploadFile(ctx context.Context, key string, file multipart.File) error
	GetFile(ctx context.Context, key string) (*v4.PresignedHTTPRequest, error)
}

type File struct {
	mongo fileMongo
	s3    fileS3
}

func NewFile(mongo fileMongo, s3 fileS3) *File {
	return &File{
		mongo: mongo,
		s3:    s3,
	}
}

func (f *File) Upload(ctx context.Context, fileHeader *multipart.FileHeader, file multipart.File, userId string) error {
	if err := f.mongo.CreateFile(ctx, &core.File{
		UserID:    userId,
		Name:      fileHeader.Filename,
		Size:      fileHeader.Size,
		CreatedAt: time.Now().Format(time.DateTime),
	}); err != nil {
		return err
	}

	if err := f.s3.UploadFile(ctx, fileHeader.Filename, file); err != nil {
		return err
	}

	return nil
}

func (f *File) GetFiles(ctx context.Context, userId string) ([]*core.GetAllFilesResp, error) {
	slog.Info("GetFiles", "userId", userId)

	filesMongo, err := f.mongo.GetFiles(ctx, userId)
	if err != nil {
		return nil, err
	}

	slog.Info("GetFiles", "response", filesMongo)

	if len(filesMongo) == 0 {
		return nil, nil
	}

	var response []*core.GetAllFilesResp
	for _, fileMongo := range filesMongo {
		output, err := f.s3.GetFile(ctx, fileMongo.Name)
		if err != nil {
			return nil, err
		}

		response = append(response, &core.GetAllFilesResp{
			ID:        fileMongo.ID.Hex(),
			UserID:    fileMongo.UserID,
			Name:      fileMongo.Name,
			Size:      fileMongo.Size,
			Url:       output.URL,
			CreatedAt: fileMongo.CreatedAt,
		})
	}

	return response, nil
}
