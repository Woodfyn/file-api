package core

import (
	"io"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type File struct {
	ID    primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name  string             `json:"name" bson:"name"`
	Size  int64              `json:"size" bson:"size"`
	Type  string             `json:"type" bson:"type"`
	Bytes []byte             `json:"file" bson:"file"`
}

type CreateFileDTO struct {
	Name   string `json:"name"`
	Size   int64  `json:"size"`
	Type   string `json:"type"`
	Reader io.Reader
}

func NewFile(dto *CreateFileDTO) (*File, error) {
	bytes, err := io.ReadAll(dto.Reader)
	if err != nil {
		return nil, err
	}

	return &File{
		Name:  dto.Name,
		Size:  dto.Size,
		Type:  dto.Type,
		Bytes: bytes,
	}, nil
}
