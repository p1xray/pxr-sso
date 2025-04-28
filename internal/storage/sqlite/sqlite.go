package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"pxr-sso/internal/domain"
	"pxr-sso/internal/storage"
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
	const op = "sqlite.SaveUser"

	stmt, err := s.db.PrepareContext(ctx,
		`select
    		u.id,
    		u.username,
    		u.password_hash,
    		u.fio,
    		u.date_of_birth,
    		u.gender,
    		u.avatar_file__key,
    		u.deleted,
    		u.created_at,
    		u.updated_at
		from users u
		where u.username = ?;`)
	if err != nil {
		return domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, username)

	var user domain.User
	err = row.Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.FIO,
		&user.DateOfBirth,
		&user.Gender,
		&user.AvatarFileKey,
		&user.Deleted,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}

		return domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, err
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
