package storage

import (
	"context"
	"errors"
	"pxr-sso/internal/dto"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

// SSOStorage provides access to data storage.
type SSOStorage interface {
	// UserByUsername returns a user from the storage by their username.
	UserByUsername(ctx context.Context, username string) (*dto.User, error)

	// UserClient returns the user's client by code.
	UserClient(ctx context.Context, userID int64, clientCode string) (*dto.Client, error)
}
