package crud

import (
	"context"
	"pxr-sso/internal/storage/domain"
)

// UserProvider represents a user provider from storage.
type UserProvider interface {
	// User returns a user from the storage by ID.
	User(ctx context.Context, id int64) (domain.User, error)

	// UserByUsername returns a user from the storage by their username.
	UserByUsername(ctx context.Context, username string) (domain.User, error)
}

// UserSaver represents a user saver in storage.
type UserSaver interface {
	// CreateUser creates a new user in the storage and returns new user ID.
	CreateUser(ctx context.Context, user domain.User) (int64, error)

	// CreateUserClientLink creates a user's client link and returns new link ID.
	CreateUserClientLink(ctx context.Context, userClientLink domain.UserClientLink) (int64, error)
}

// ClientProvider represents a client provider from storage.
type ClientProvider interface {
	// ClientByCode returns client by their code from storage.
	ClientByCode(ctx context.Context, code string) (domain.Client, error)

	// UserClient returns the user's client from the storage by code.
	UserClient(ctx context.Context, userID int64, clientCode string) (domain.Client, error)
}

// PermissionProvider represents a permission provider from storage.
type PermissionProvider interface {
	// UserPermissions returns the user's permissions from the storage.
	UserPermissions(ctx context.Context, userID int64) ([]domain.Permission, error)
}

// SessionProvider represents a session provider from storage.
type SessionProvider interface {
	// SessionByRefreshToken returns a session by its refresh token.
	SessionByRefreshToken(ctx context.Context, refreshToken string) (domain.Session, error)
}

// SessionSaver represents a session saver in storage.
type SessionSaver interface {
	// CreateSession creates a new session in the storage and returns new session ID.
	CreateSession(ctx context.Context, session domain.Session) (int64, error)

	// RemoveSession removes a session by ID.
	RemoveSession(ctx context.Context, id int64) error
}
