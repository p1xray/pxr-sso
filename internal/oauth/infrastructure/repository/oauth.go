package repository

import (
	"context"
	"fmt"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/dto"
	"github.com/p1xray/pxr-sso/internal/oauth/infrastructure/converter"
	"github.com/p1xray/pxr-sso/internal/oauth/infrastructure/storage/models"
)

const emptyID = 0

type Storage interface {
	ClientByCode(ctx context.Context, code string) (models.Client, error)
	ClientAudiences(ctx context.Context, clientID int64) ([]models.Audience, error)

	User(ctx context.Context, id int64) (models.User, error)
	UserByUsername(ctx context.Context, username string) (models.User, error)
	CreateUser(ctx context.Context, user models.User) error

	CreateUserClientLink(ctx context.Context, link models.UserClientLink) error
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

func (o *OAuth) User(ctx context.Context, id int64) (dto.User, error) {
	const op = "infrastructure.repository.User"

	user, err := o.storage.User(ctx, id)
	if err != nil {
		return dto.User{}, fmt.Errorf("%s: %w", op, err)
	}

	userDTO := converter.ToUserDTO(user)

	return userDTO, nil
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

func (o *OAuth) CreateUser(ctx context.Context, user dto.User, clientID int64) error {
	const op = "infrastructure.repository.CreateUser"

	if err := o.createUser(ctx, user); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := o.createUserClientLink(ctx, user.ID(), clientID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (o *OAuth) createUser(ctx context.Context, user dto.User) error {
	userStorageModel := converter.ToUserStorage(models.User{}, user, models.UserCreated())
	if err := o.storage.CreateUser(ctx, userStorageModel); err != nil {
		return err
	}

	return nil
}

func (o *OAuth) createUserClientLink(ctx context.Context, userID int64, clientID int64) error {
	if userID == emptyID || clientID == emptyID {
		return fmt.Errorf("a non-null identifiers is required to create an user client link in storage")
	}

	userClientLinkStorageModel := converter.ToUserClientLinkStorage(userID, clientID, models.UserClientLinkCreated())
	if err := o.storage.CreateUserClientLink(ctx, userClientLinkStorageModel); err != nil {
		return err
	}

	return nil
}
