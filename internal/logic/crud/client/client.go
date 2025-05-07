package clientcrud

import (
	"context"
	"fmt"
	"log/slog"
	"pxr-sso/internal/lib/logger/sl"
	"pxr-sso/internal/logic/crud"
	"pxr-sso/internal/logic/dto"
)

// CRUD provides methods for managing client data.
type CRUD struct {
	log            *slog.Logger
	clientProvider crud.ClientProvider
}

// New creates a new instance of the client's CRUD.
func New(
	log *slog.Logger,
	clientProvider crud.ClientProvider,
) *CRUD {
	return &CRUD{
		log:            log,
		clientProvider: clientProvider,
	}
}

// ClientByCode returns client by code.
func (c *CRUD) ClientByCode(ctx context.Context, code string) (dto.ClientDTO, error) {
	const op = "clientcrud.ClientByCode"

	log := c.log.With(
		slog.String("op", op),
		slog.String("code", code),
	)

	client, err := c.clientProvider.ClientByCode(ctx, code)
	if err != nil {
		log.Error("failed to get client from storage", sl.Err(err))

		return dto.ClientDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	clientData := dto.ClientDTO{
		ID:        client.ID,
		Code:      client.Code,
		SecretKey: client.SecretKey,
	}

	return clientData, nil
}

// UserClient returns the user's client by user ID and client code.
func (c *CRUD) UserClient(ctx context.Context, userID int64, clientCode string) (dto.ClientDTO, error) {
	const op = "clientcrud.UserClient"

	log := c.log.With(
		slog.String("op", op),
		slog.Int64("user ID", userID),
		slog.String("client code", clientCode),
	)

	client, err := c.clientProvider.UserClient(ctx, userID, clientCode)
	if err != nil {
		log.Error("failed to get client from storage", sl.Err(err))

		return dto.ClientDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	clientData := dto.ClientDTO{
		ID:        client.ID,
		Code:      client.Code,
		SecretKey: client.SecretKey,
	}

	return clientData, nil
}
