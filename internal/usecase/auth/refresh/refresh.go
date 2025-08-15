package refresh

import (
	"context"
	"errors"
	"fmt"
	"github.com/p1xray/pxr-sso/internal/config"
	"github.com/p1xray/pxr-sso/internal/dto"
	"github.com/p1xray/pxr-sso/internal/entity"
	"github.com/p1xray/pxr-sso/internal/infrastructure"
	"github.com/p1xray/pxr-sso/internal/lib/logger/sl"
	"github.com/p1xray/pxr-sso/internal/usecase"
	jwtparser "github.com/p1xray/pxr-sso/pkg/jwt/parser"
	"log/slog"
)

type AuthRepository interface {
	ClientByCode(ctx context.Context, code string) (dto.Client, error)
	DataForRefreshTokens(ctx context.Context, refreshTokenID string) (dto.DataForRefreshTokens, error)

	Save(ctx context.Context, auth entity.Auth) error
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
	const op = "usecase.auth.refresh"

	log := uc.log.With(
		slog.String("op", op),
		slog.String("refresh token", data.RefreshToken),
		slog.String("client code", data.ClientCode),
	)
	log.Info("attempting to refresh tokens")

	// Get client from storage.
	client, err := uc.repo.ClientByCode(ctx, data.ClientCode)
	if err != nil {
		if errors.Is(err, infrastructure.ErrEntityNotFound) {
			log.Warn("client not found", sl.Err(err))

			return entity.Tokens{}, fmt.Errorf("%s: %w", op, usecase.ErrClientNotFound)
		}

		log.Error("error getting client from storage", sl.Err(err))
		return entity.Tokens{}, err
	}

	// Parse refresh token by client secret key.
	refreshTokenClaims, err := jwtparser.ParseRefreshToken(data.RefreshToken, []byte(client.SecretKey))
	if err != nil {
		log.Error("error parsing refresh token", sl.Err(err))

		return entity.Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	// Get data for refresh tokens from storage.
	storageRefreshTokensData, err := uc.repo.DataForRefreshTokens(ctx, refreshTokenClaims.ID)
	if err != nil {
		if errors.Is(err, infrastructure.ErrEntityNotFound) {
			log.Warn("session not found", sl.Err(err))

			return entity.Tokens{}, fmt.Errorf("%s: %w", op, usecase.ErrSessionNotFound)
		}

		log.Error("error getting session from storage", sl.Err(err))
		return entity.Tokens{}, err
	}

	// Create auth entity.
	auth, err := entity.NewAuth(
		uc.cfg.AccessTokenTTL,
		uc.cfg.RefreshTokenTTL,
		entity.WithUser(storageRefreshTokensData.User),
		entity.WithClient(client),
		entity.WithSession(storageRefreshTokensData.Session),
	)
	if err != nil {
		return entity.Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	// Refresh tokens.
	entityRefreshTokensParams := entity.RefreshTokensParams{
		UserAgent:   data.UserAgent,
		Fingerprint: data.Fingerprint,
		Issuer:      data.Issuer,
	}
	tokens, err := auth.RefreshTokens(entityRefreshTokensParams)
	if err != nil {
		log.Error("failed to refresh tokens", sl.Err(err))

		return entity.Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	// Save data to storage.
	err = uc.repo.Save(ctx, auth)
	if err != nil {
		log.Error("error saving data to storage.", sl.Err(err))

		return entity.Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("tokens refreshed successfully")

	return tokens, nil
}
