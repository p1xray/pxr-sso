package register

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

// Repository is a repository for register a new user use-case.
type Repository interface {
	DataForRegister(ctx context.Context, username, clientCode string) (dto.DataForRegister, error)
	Save(ctx context.Context, auth *entity.Auth) error
}

type Handler interface {
	SendToKafka(clientCode string, user entity.User) error
}

// UseCase is a use-case for registering a new user.
type UseCase struct {
	log     *slog.Logger
	cfg     config.TokensConfig
	repo    Repository
	handler Handler
}

// New returns new register a new user use-case.
func New(log *slog.Logger, cfg config.TokensConfig, repo Repository, handler Handler) *UseCase {
	return &UseCase{
		log:     log,
		cfg:     cfg,
		repo:    repo,
		handler: handler,
	}
}

// Execute executes the use-case for registering a new user. If successful, new tokens are returned.
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
	auth, err := entity.NewAuth(
		uc.cfg.AccessTokenTTL,
		uc.cfg.RefreshTokenTTL,
		entity.WithAuthUser(storageData.User),
		entity.WithAuthClient(storageData.Client),
		entity.WithAuthDefaultRoles(storageData.ClientDefaultRoles...),
		entity.WithAuthDefaultPermissionCodes(storageData.ClientDefaultPermissionCodes...),
	)
	if err != nil {
		return entity.Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

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
	err = auth.Register(entityRegisterParams)
	if err != nil {
		if errors.Is(err, entity.ErrUserExists) {
			log.Warn("user already exists", sl.Err(err))

			return entity.Tokens{}, fmt.Errorf("%s: %w", op, usecase.ErrUserExists)
		}

		log.Error("failed to register", sl.Err(err))

		return entity.Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	// Save user data to storage.
	err = uc.repo.Save(ctx, &auth)
	if err != nil {
		log.Error("error saving user data to storage.", sl.Err(err))

		return entity.Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	// Create new session for saved user.
	tokens, err := auth.CreateNewSession(data.Issuer, data.UserAgent, data.Fingerprint)
	if err != nil {
		log.Error("error creating new session for registered user.", sl.Err(err))

		return entity.Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	// Save session data to storage.
	err = uc.repo.Save(ctx, &auth)
	if err != nil {
		log.Error("error saving session data to storage.", sl.Err(err))

		return entity.Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user register successfully")

	if err = uc.handler.SendToKafka(auth.ClientCode(), auth.User); err != nil {
		log.Error("error sending registered new user data to kafka", sl.Err(err))

		return entity.Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	return tokens, nil
}
