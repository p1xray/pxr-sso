package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"pxr-sso/internal/domain"
)

// Storage provides access to sqlite storage.
type Storage struct {
	db *sql.DB
}

// New creates a new instance of the SQLite store.
func New(storagePath string) (*Storage, error) {
	const op = "sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

// UserByUsername returns a user from the storage by their username.
func (s *Storage) UserByUsername(ctx context.Context, username string) (domain.User, error) {

}

// UserPermissions returns the user's permissions from the storage.
func (s *Storage) UserPermissions(ctx context.Context, userID int64) ([]domain.Permission, error) {
	
}

// UserClient returns the user's client from the storage by code.
func (s *Storage) UserClient(ctx context.Context, userID int64, clientCode string) (domain.Client, error) {

}

// CreateSession creates a new session in the storage.
func (s *Storage) CreateSession(ctx context.Context, session domain.Session) error {

}
