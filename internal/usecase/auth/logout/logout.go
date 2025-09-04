package logout

import (
	"context"
	"errors"
	"fmt"
	"github.com/p1xray/pxr-sso/internal/config"
	"github.com/p1xray/pxr-sso/internal/dto"
	"github.com/p1xray/pxr-sso/internal/entity"
	"github.com/p1xray/pxr-sso/internal/infrastructure"
	"github.com/p1xray/pxr-sso/internal/usecase"
	jwtparser "github.com/p1xray/pxr-sso/pkg/jwt/parser"
	"github.com/p1xray/pxr-sso/pkg/logger/sl"
	"log/slog"
)

type AuthRepository interface {
	ClientByCode(ctx context.Context, code string) (dto.Client, error)
	DataForLogout(ctx context.Context, refreshTokenID string) (dto.DataForLogout, error)

	Save(ctx context.Context, auth *entity.Auth) error
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

func (uc *UseCase) Execute(ctx context.Context, data Params) error {
	const op = "usecase.auth.logout"

	log := uc.log.With(
		slog.String("op", op),
		slog.String("refresh token", data.RefreshToken),
		slog.String("client code", data.ClientCode),
	)
	log.Info("attempting to user logout")

	// Get client from storage.
	client, err := uc.repo.ClientByCode(ctx, data.ClientCode)
	if err != nil {
		if errors.Is(err, infrastructure.ErrEntityNotFound) {
			log.Warn("client not found", sl.Err(err))

			return fmt.Errorf("%s: %w", op, usecase.ErrClientNotFound)
		}

		log.Error("error getting client from storage", sl.Err(err))
		return err
	}

	// Parse refresh token by client secret key.
	refreshTokenClaims, err := jwtparser.ParseRefreshToken(data.RefreshToken, []byte(client.SecretKey))
	if err != nil {
		log.Error("error parsing refresh token", sl.Err(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	// Get data for logout from storage.
	storageLogoutData, err := uc.repo.DataForLogout(ctx, refreshTokenClaims.ID)
	if err != nil {
		if errors.Is(err, infrastructure.ErrEntityNotFound) {
			log.Warn("session not found", sl.Err(err))

			return fmt.Errorf("%s: %w", op, usecase.ErrSessionNotFound)
		}

		log.Error("error getting session from storage", sl.Err(err))
		return err
	}

	// Create auth entity.
	auth, err := entity.NewAuth(
		uc.cfg.AccessTokenTTL,
		uc.cfg.RefreshTokenTTL,
		entity.WithAuthSession(storageLogoutData.Session),
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// Logout.
	if err = auth.Logout(); err != nil {
		log.Error("failed to logout", sl.Err(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	// Save data to storage.
	err = uc.repo.Save(ctx, &auth)
	if err != nil {
		log.Error("error saving data to storage.", sl.Err(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user logout successfully")

	return nil
}
