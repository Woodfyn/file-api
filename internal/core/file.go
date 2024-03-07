package core

import (
	"io"
)

type File struct {
	ID    string `json:"id" bson:"_id"`
	Name  string `json:"name" bson:"name"`
	Size  int64  `json:"size" bson:"size"`
	Bytes []byte `json:"file" bson:"bytes"`
}

type CreateFileDTO struct {
	Name   string `json:"name" validate:"required"`
	Size   int64  `json:"size" validate:"required"`
	Reader io.Reader
}

func (f *CreateFileDTO) Validate() error {
	return validate.Struct(f)
}
