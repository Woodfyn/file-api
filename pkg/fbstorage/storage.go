package fbstorage

import (
	"context"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type FirebaseConfig struct {
	APIKey            string `json:"apiKey"`
	AuthDomain        string `json:"authDomain"`
	ProjectID         string `json:"projectId"`
	StorageBucket     string `json:"storageBucket"`
	MessagingSenderID string `json:"messagingSenderId"`
	AppID             string `json:"appId"`
	MeasurementID     string `json:"measurementId"`
}

func NewFBStorageClient(ctx context.Context, filename string) (*storage.Client, error) {
	opt := option.WithCredentialsFile(filename)
	client, err := storage.NewClient(ctx, opt)
	if err != nil {
		return nil, err
	}

	return client, nil
}
