package service

import (
	"context"
	"time"

	"github.com/Woodfyn/file-api/internal/core"
	"github.com/Woodfyn/file-api/pkg/auth"
	"github.com/Woodfyn/file-api/pkg/hash"
)

type authMongo interface {
	CreateUser(ctx context.Context, user *core.User) error
	GetUserByPassword(ctx context.Context, password string) (*core.User, error)
}

type Auth struct {
	authMongo authMongo

	hasher hash.PasswordHasher

	tokenManager   auth.TokenManager
	accessTokenTTL time.Duration
	refreshTTL     time.Duration
}

func NewAuth(authMongo authMongo, hasher hash.PasswordHasher, tokenManager auth.TokenManager, acssTokenTTL, refreshTTL time.Duration) *Auth {
	return &Auth{
		authMongo: authMongo,

		hasher: hasher,

		tokenManager:   tokenManager,
		accessTokenTTL: acssTokenTTL,
		refreshTTL:     refreshTTL,
	}
}

func (a *Auth) SignUp(ctx context.Context, input core.SignUpRequest) error {
	passwordHash, err := a.hasher.Hash(input.Password)
	if err != nil {
		return err
	}

	if err := a.authMongo.CreateUser(ctx, &core.User{
		Email:        input.Email,
		Password:     passwordHash,
		RegisteredAt: time.Now().Format(time.DateTime),
	}); err != nil {
		return err
	}

	return nil
}

func (a *Auth) SignIn(ctx context.Context, input core.SignInRequest) (*core.TokenResp, error) {
	passwordHash, err := a.hasher.Hash(input.Password)
	if err != nil {
		return nil, err
	}

	user, err := a.authMongo.GetUserByPassword(ctx, passwordHash)
	if err != nil {
		return nil, err
	}

	accessToken, err := a.tokenManager.NewJWT(user.ID.Hex(), a.accessTokenTTL)
	if err != nil {
		return nil, err
	}

	refreshToken, err := a.tokenManager.NewJWT(user.ID.Hex(), a.refreshTTL)
	if err != nil {
		return nil, err
	}

	return &core.TokenResp{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (a *Auth) Refresh(ctx context.Context, refreshToken string) (*core.TokenResp, error) {
	userId, err := a.tokenManager.Parse(refreshToken)
	if err != nil {
		return nil, err
	}

	ok := a.tokenManager.IsTokenExpired(refreshToken)
	if !ok {
		return nil, core.ErrInvalidRefreshToken
	}

	accessToken, err := a.tokenManager.NewJWT(userId, a.accessTokenTTL)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := a.tokenManager.NewJWT(userId, a.refreshTTL)
	if err != nil {
		return nil, err
	}

	return &core.TokenResp{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (a *Auth) Parse(token string) (string, error) {
	return a.tokenManager.Parse(token)
}

func (a *Auth) IsTokenExpired(token string) bool {
	return a.tokenManager.IsTokenExpired(token)
}
