package core

import (
	"time"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

type User struct {
	ID           int       `json:"id" bson:"_id"`
	Email        string    `json:"email" bson:"email"`
	Password     string    `json:"password" bson:"password"`
	RegisteredAt time.Time `json:"registered_at" bson:"registered_at"`
}

type SingUpRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (s *SingUpRequest) Validate() error {
	return validate.Struct(s)
}

type SingInRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (s *SingInRequest) Validate() error {
	return validate.Struct(s)
}
