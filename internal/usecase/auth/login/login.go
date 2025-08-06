package login

import (
	"context"
	"errors"
	"fmt"
	"github.com/p1xray/pxr-sso/internal/config"
	"github.com/p1xray/pxr-sso/internal/dto"
	"github.com/p1xray/pxr-sso/internal/entity"
	"github.com/p1xray/pxr-sso/internal/lib/logger/sl"
	"github.com/p1xray/pxr-sso/internal/logic/service"
	"github.com/p1xray/pxr-sso/internal/storage"
	"log/slog"
)

type AuthRepository interface {
	UserData(ctx context.Context, username string, clientCode string) (dto.User, error)

	UserCredentialsByUsername(ctx context.Context, username string) (dto.UserCredentials, error)
	UserPermissions(ctx context.Context, id int64) ([]dto.Permission, error)
	UserClientByCode(ctx context.Context, userID int64, clientCode string) (dto.Client, error)
	CreateSession(ctx context.Context, session entity.Session) error
}

type UseCase struct {
	log  *slog.Logger
	cfg  config.TokensConfig
	repo AuthRepository
}

func New(log *slog.Logger, cfg config.TokensConfig, repo AuthRepository) *UseCase {
	return &UseCase{
		log:  log,
		cfg:  cfg,
		repo: repo,
	}
}

func (uc *UseCase) Execute(ctx context.Context, data Params) (entity.Tokens, error) {
	const op = "auth.Login"

	log := uc.log.With(
		slog.String("op", op),
		slog.String("username", data.Username),
		slog.String("client code", data.ClientCode),
	)
	log.Info("attempting to login user")

	// Get user data from storage.
	user, err := uc.repo.UserData(ctx, data.Username, data.ClientCode)
	if err != nil {
		if errors.Is(err, storage.ErrEntityNotFound) {
			return entity.Tokens{}, fmt.Errorf("%s: %w", op, service.ErrInvalidCredentials)
		}

		return entity.Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	// Create auth entity.
	auth := entity.NewAuth(user, uc.cfg.AccessTokenTTL, uc.cfg.RefreshTokenTTL)

	// Log in.
	entityLoginParams := entity.LoginParams{
		Password:    data.Password,
		UserAgent:   data.UserAgent,
		Fingerprint: data.Fingerprint,
		Issuer:      data.Issuer,
	}
	if err = auth.Login(entityLoginParams); err != nil {
		log.Error("failed to login", sl.Err(err))

		return entity.Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	// Create new session in storage.
	err = uc.repo.CreateSession(ctx, auth.Session)
	if err != nil {
		log.Error("failed to create session", sl.Err(err))

		return entity.Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user logged in successfully")

	return auth.Tokens(), nil
}
