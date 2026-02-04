package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/p1xray/pxr-sso/internal/oauth/infrastructure"
	"github.com/p1xray/pxr-sso/internal/oauth/infrastructure/storage/models"
)

// Storage provides access to sqlite storage.
type Storage struct {
	db *sql.DB
}

// New creates a new instance of the SQLite store.
func New(storagePath string) (*Storage, error) {
	const op = "infrastructure.storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) ClientByCode(ctx context.Context, code string) (models.Client, error) {
	const op = "infrastructure.storage.sqlite.ClientByCode"

	stmt, err := s.db.PrepareContext(ctx,
		`select
			 c.id,
			 c.name,
			 c.code,
			 c.secret_key,
			 c.deleted,
			 c.created_at,
			 c.updated_at
		 from clients c
		 where c.code = ?;`)
	if err != nil {
		return models.Client{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, code)

	var client models.Client
	err = row.Scan(
		&client.ID,
		&client.Name,
		&client.Code,
		&client.SecretKey,
		&client.Deleted,
		&client.CreatedAt,
		&client.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Client{}, fmt.Errorf("%s: %w", op, infrastructure.ErrEntityNotFound)
		}

		return models.Client{}, fmt.Errorf("%s: %w", op, err)
	}

	return client, nil
}

func (s *Storage) ClientAudiences(ctx context.Context, clientID int64) ([]models.Audience, error) {
	const op = "infrastructure.storage.sqlite.ClientAudiences"

	stmt, err := s.db.PrepareContext(ctx,
		`select
			 ca.id,
			 ca.client_id,
			 ca.url,
			 ca.created_at,
			 ca.updated_at
		 from client_audiences ca
		 where ca.client_id = ?;`)
	if err != nil {
		return []models.Audience{}, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := stmt.QueryContext(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	audiences := make([]models.Audience, 0)
	for rows.Next() {
		audience := models.Audience{}
		err = rows.Scan(
			&audience.ID,
			&audience.ClientID,
			&audience.URL,
			&audience.CreatedAt,
			&audience.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		audiences = append(audiences, audience)
	}

	return audiences, nil
}

func (s *Storage) User(ctx context.Context, id int64) (models.User, error) {
	const op = "infrastructure.storage.sqlite.User"

	stmt, err := s.db.PrepareContext(ctx,
		`select
    		u.id,
    		u.username,
    		u.password_hash,
    		u.fio,
    		u.date_of_birth,
    		u.gender,
    		u.avatar_file_key,
    		u.deleted,
    		u.created_at,
    		u.updated_at
		from users u
		where u.id = ?;`)
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, id)

	var user models.User
	err = row.Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.FullName,
		&user.DateOfBirth,
		&user.Gender,
		&user.AvatarFileKey,
		&user.Deleted,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, infrastructure.ErrEntityNotFound)
		}

		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Storage) UserByUsername(ctx context.Context, username string) (models.User, error) {
	const op = "infrastructure.storage.sqlite.UserByUsername"

	stmt, err := s.db.PrepareContext(ctx,
		`select
    		u.id,
    		u.username,
    		u.password_hash,
    		u.fio,
    		u.date_of_birth,
    		u.gender,
    		u.avatar_file_key,
    		u.deleted,
    		u.created_at,
    		u.updated_at
		from users u
		where u.username = ?;`)
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, username)

	var user models.User
	err = row.Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.FullName,
		&user.DateOfBirth,
		&user.Gender,
		&user.AvatarFileKey,
		&user.Deleted,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, infrastructure.ErrEntityNotFound)
		}

		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}
