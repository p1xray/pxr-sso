package register

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
	"log/slog"
)

type AuthRepository interface {
	DataForRegister(ctx context.Context, username, clientCode string) (dto.DataForRegister, error)
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
	const op = "usecase.auth.register"

	log := uc.log.With(
		slog.String("op", op),
		slog.String("username", data.Username),
		slog.String("client code", data.ClientCode),
	)
	log.Info("attempting to register new user")

	// Get data from storage.
	storageData, err := uc.repo.DataForRegister(ctx, data.Username, data.ClientCode)
	if err != nil && !errors.Is(err, infrastructure.ErrEntityNotFound) {
		log.Error("error getting data from storage", sl.Err(err))

		return entity.Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	// Create auth entity.
	auth := entity.NewAuth(
		uc.cfg.AccessTokenTTL,
		uc.cfg.RefreshTokenTTL,
		entity.WithUser(storageData.User),
		entity.WithClient(storageData.Client),
	)

	// Register.
	entityRegisterParams := entity.RegisterParams{
		Username:      data.Username,
		Password:      data.Password,
		FullName:      data.FIO,
		DateOfBirth:   data.DateOfBirth,
		Gender:        data.Gender,
		AvatarFileKey: data.AvatarFileKey,
		UserAgent:     data.UserAgent,
		Fingerprint:   data.Fingerprint,
		Issuer:        data.Issuer,
	}
	tokens, err := auth.Register(entityRegisterParams)
	if err != nil {
		if errors.Is(err, entity.ErrUserExists) {
			log.Warn("user already exists", sl.Err(err))

			return entity.Tokens{}, fmt.Errorf("%s: %w", op, usecase.ErrUserExists)
		}

		log.Error("failed to register", sl.Err(err))

		return entity.Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	// Save data to storage.
	err = uc.repo.Save(ctx, auth)
	if err != nil {
		log.Error("error saving data to storage.", sl.Err(err))

		return entity.Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user register successfully")

	return tokens, nil
}
