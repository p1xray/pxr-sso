package repository

import (
	"context"
	"fmt"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/dto"
	"github.com/p1xray/pxr-sso/internal/oauth/infrastructure/converter"
	"github.com/p1xray/pxr-sso/internal/oauth/infrastructure/storage/models"
)

type Storage interface {
	ClientByCode(ctx context.Context, code string) (models.Client, error)
	ClientAudiences(ctx context.Context, clientID int64) ([]models.Audience, error)
	UserByUsername(ctx context.Context, username string) (models.User, error)
}

type OAuth struct {
	storage Storage
}

func NewOAuthRepository(storage Storage) *OAuth {
	return &OAuth{
		storage: storage,
	}
}

func (o *OAuth) ClientByCode(ctx context.Context, code string) (dto.Client, error) {
	const op = "infrastructure.repository.ClientByCode"

	client, err := o.storage.ClientByCode(ctx, code)
	if err != nil {
		return dto.Client{}, fmt.Errorf("%s: %w", op, err)
	}

	clientAudiences, err := o.storage.ClientAudiences(ctx, client.ID)
	if err != nil {
		return dto.Client{}, fmt.Errorf("%s: %w", op, err)
	}

	clientDTO := converter.ToClientDTO(client, clientAudiences)

	return clientDTO, nil
}

func (o *OAuth) UserByUsername(ctx context.Context, username string) (dto.User, error) {
	const op = "infrastructure.repository.UserByUsername"

	user, err := o.storage.UserByUsername(ctx, username)
	if err != nil {
		return dto.User{}, fmt.Errorf("%s: %w", op, err)
	}

	userDTO := converter.ToUserDTO(user)

	return userDTO, nil
}
