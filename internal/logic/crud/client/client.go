package clientcrud

import (
	"context"
	"pxr-sso/internal/logic/crud"
	"pxr-sso/internal/logic/dto"
)

// CRUD provides methods for managing client data.
type CRUD struct {
	clientProvider crud.ClientProvider
}

// New creates a new instance of the client's CRUD.
func New(clientProvider crud.ClientProvider) *CRUD {
	return &CRUD{
		clientProvider: clientProvider,
	}
}

// ClientByCode returns client by code.
func (c *CRUD) ClientByCode(ctx context.Context, code string) (dto.ClientDTO, error) {
	client, err := c.clientProvider.ClientByCode(ctx, code)
	if err != nil {
		return dto.ClientDTO{}, err
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
	client, err := c.clientProvider.UserClient(ctx, userID, clientCode)
	if err != nil {
		return dto.ClientDTO{}, err
	}

	clientData := dto.ClientDTO{
		ID:        client.ID,
		Code:      client.Code,
		SecretKey: client.SecretKey,
	}

	return clientData, nil
}
