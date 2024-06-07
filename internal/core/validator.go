package core

import "github.com/go-playground/validator/v10"

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func (s *SignUpRequest) Validate() error {
	return validate.Struct(s)
}

func (s *SignInRequest) Validate() error {
	return validate.Struct(s)
}

func (s *RefreshTokenReq) Validate() error {
	return validate.Struct(s)
}
