package storage

import (
	"context"
	"mime/multipart"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type File struct {
	s3        *s3.Client
	presigner *s3.PresignClient

	bucketName string
}

func NewFile(s3 *s3.Client, presign *s3.PresignClient, bucketName string) *File {
	return &File{
		s3:        s3,
		presigner: presign,

		bucketName: bucketName,
	}
}

func (f *File) UploadFile(ctx context.Context, key string, file multipart.File) error {
	_, err := f.s3.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(f.bucketName),
		Key:    aws.String(key),
		Body:   file,
	})

	return err
}

func (f *File) GetFile(ctx context.Context, key string) (*v4.PresignedHTTPRequest, error) {
	return f.presigner.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(f.bucketName),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(60 * int64(time.Second))
	})
}
