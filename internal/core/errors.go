package core

import (
	"errors"
)

var (
	ErrFileNotFound = errors.New("file not found")
	ErrTokenExpired = errors.New("token expired")

	ErrEmptyAuthHeader   = errors.New("empty auth header")
	ErrInvalidAuthHeader = errors.New("invalid auth header")

	ErrAccessTokenEmpty      = errors.New("access token is empty")
	ErrAccessTokenIsExpired  = errors.New("access token is expired")
	ErrInvalidAccessToken    = errors.New("invalid access token")
	ErrInvalidRefreshToken   = errors.New("invalid refresh token")
	ErrRefreshTokenIsExpired = errors.New("refresh token is expired")

	ErrInvalidPassword = errors.New("invalid password")
)
