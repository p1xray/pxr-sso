package login

import (
	"context"
	"errors"
	"fmt"
	"github.com/p1xray/pxr-sso/internal/config"
	"github.com/p1xray/pxr-sso/internal/dto"
	"github.com/p1xray/pxr-sso/internal/entity"
	"github.com/p1xray/pxr-sso/internal/infrastructure"
	"github.com/p1xray/pxr-sso/internal/usecase"
	"github.com/p1xray/pxr-sso/pkg/logger/sl"
	"log/slog"
)

// Repository is a repository for log in use-case.
type Repository interface {
	DataForLogin(ctx context.Context, username, clientCode string) (dto.DataForLogin, error)
	Save(ctx context.Context, auth *entity.Auth) error
}

// UseCase is a use-case for logging in a user.
type UseCase struct {
	log  *slog.Logger
	cfg  config.TokensConfig
	repo Repository
}

// New returns new log in use-case.
func New(log *slog.Logger, cfg config.TokensConfig, repo Repository) *UseCase {
	return &UseCase{
		log:  log,
		cfg:  cfg,
		repo: repo,
	}
}

// Execute executes the use-case for logging in a user. If successful, new tokens are returned.
func (uc *UseCase) Execute(ctx context.Context, data Params) (entity.Tokens, error) {
	const op = "usecase.auth.login"

	log := uc.log.With(
		slog.String("op", op),
		slog.String("username", data.Username),
		slog.String("client code", data.ClientCode),
	)
	log.Info("attempting to login user")

	// Get user data from storage.
	storageLoginData, err := uc.repo.DataForLogin(ctx, data.Username, data.ClientCode)
	if err != nil {
		if errors.Is(err, infrastructure.ErrEntityNotFound) {
			return entity.Tokens{}, fmt.Errorf("%s: %w", op, usecase.ErrInvalidCredentials)
		}

		return entity.Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	// Create auth entity.
	auth, err := entity.NewAuth(
		uc.cfg.AccessTokenTTL,
		uc.cfg.RefreshTokenTTL,
		entity.WithAuthUser(storageLoginData.User),
		entity.WithAuthClient(storageLoginData.Client),
		entity.WithAuthSession(storageLoginData.Sessions...),
	)
	if err != nil {
		return entity.Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	// Log in.
	entityLoginParams := entity.LoginParams{
		Password:    data.Password,
		UserAgent:   data.UserAgent,
		Fingerprint: data.Fingerprint,
		Issuer:      data.Issuer,
	}
	tokens, err := auth.Login(entityLoginParams)
	if err != nil {
		log.Error("failed to login", sl.Err(err))

		return entity.Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	// Save data in storage.
	err = uc.repo.Save(ctx, &auth)
	if err != nil {
		log.Error("error saving data to storage.", sl.Err(err))

		return entity.Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user logged in successfully")

	return tokens, nil
}
