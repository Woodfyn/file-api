package core

import (
	"errors"
)

var (
	ErrFileNotFound = errors.New("file not found")
	ErrTokenExpired = errors.New("token expired")
)
