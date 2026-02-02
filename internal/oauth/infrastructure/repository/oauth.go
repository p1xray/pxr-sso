package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/dto"
	"github.com/p1xray/pxr-sso/internal/oauth/infrastructure"
	"github.com/p1xray/pxr-sso/internal/oauth/infrastructure/converter"
	"github.com/p1xray/pxr-sso/internal/oauth/infrastructure/storage/models"
	"github.com/p1xray/pxr-sso/pkg/logger/sl"
	"log/slog"
)

type Storage interface {
	ClientByCode(ctx context.Context, code string) (models.Client, error)
	ClientAudiences(ctx context.Context, clientID int64) ([]models.Audience, error)
}

type Auth struct {
	log     *slog.Logger
	storage Storage
}

func NewOAuthRepository(log *slog.Logger, storage Storage) *Auth {
	return &Auth{
		log:     log,
		storage: storage,
	}
}

func (a *Auth) ClientByCode(ctx context.Context, code string) (dto.Client, error) {
	const op = "repository.auth.ClientByCode"

	log := a.log.With(
		slog.String("op", op),
		slog.String("code", code),
	)

	client, err := a.storage.ClientByCode(ctx, code)
	if err != nil {
		if errors.Is(err, infrastructure.ErrEntityNotFound) {
			log.Warn("client not found", sl.Err(err))
		} else {
			log.Error("error getting user client", sl.Err(err))
		}

		return dto.Client{}, fmt.Errorf("%s: %w", op, err)
	}

	clientAudiences, err := a.storage.ClientAudiences(ctx, client.ID)
	if err != nil {
		log.Error("error getting client audiences", sl.Err(err))

		return dto.Client{}, fmt.Errorf("%s: %w", op, err)
	}

	clientDTO := converter.ToClientDTO(client, clientAudiences)

	return clientDTO, nil
}
