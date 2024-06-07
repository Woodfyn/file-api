package core

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type File struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    string             `bson:"user_id"`
	Name      string             `bson:"name"`
	Size      int64              `bson:"size"`
	CreatedAt string             `bson:"created_at"`
}

type GetAllFilesResp struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Name      string `json:"name"`
	Size      int64  `json:"size"`
	Url       string `json:"url"`
	CreatedAt string `json:"created_at"`
}
