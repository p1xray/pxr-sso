package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
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
	const op = "sqlite.UserByUsername"

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
			return domain.User{}, fmt.Errorf("%s: %w", op, storage.ErrEntityNotFound)
		}

		return domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, err
}

// UserPermissions returns the user's permissions from the storage.
func (s *Storage) UserPermissions(ctx context.Context, userID int64) ([]domain.Permission, error) {
	const op = "sqlite.UserPermissions"

	stmt, err := s.db.PrepareContext(ctx,
		`select p.id, p.code, p.description, p.active, p.deleted, p.created_at, p.updated_at from permissions p
			join role_permissions rp on rp.permission_id = p.id
			join user_roles ur on ur.role_id = rp.role_id
		where ur.user_id = ?;`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := stmt.QueryContext(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	permissions := make([]domain.Permission, 0)
	for rows.Next() {
		p := domain.Permission{}
		err := rows.Scan(&p.ID, &p.Code, &p.Description, &p.Active, &p.Deleted, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		permissions = append(permissions, p)
	}

	return permissions, nil
}

// UserClient returns the user's client from the storage by code.
func (s *Storage) UserClient(ctx context.Context, userID int64, clientCode string) (domain.Client, error) {
	const op = "sqlite.UserClient"

	stmt, err := s.db.PrepareContext(ctx,
		`select c.id, c.name, c.code, c.secret_key, c.deleted, c.created_at, c.updated_at from clients c
			join user_clients uc on uc.client_id = c.id
		where uc.user_id = ? and c.code = ?;`)
	if err != nil {
		return domain.Client{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, userID, clientCode)

	var client domain.Client
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
			return domain.Client{}, fmt.Errorf("%s: %w", op, storage.ErrEntityNotFound)
		}

		return domain.Client{}, fmt.Errorf("%s: %w", op, err)
	}

	return client, nil
}

// CreateSession creates a new session in the storage.
func (s *Storage) CreateSession(ctx context.Context, session domain.Session) (int64, error) {
	const op = "sqlite.CreateSession"

	stmt, err := s.db.PrepareContext(ctx,
		`insert into sessions (user_id, refresh_token, user_agent, fingerprint, expires_at, created_at, updated_at)
		values(?, ?, ?, ?, ?, ?, ?);`)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.ExecContext(
		ctx,
		session.UserID,
		session.RefreshToken,
		session.UserAgent,
		session.Fingerprint,
		session.ExpiresAt,
		session.CreatedAt,
		session.UpdatedAt,
	)

	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrEntityExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	newSessionID, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return newSessionID, nil
}

// ClientByCode returns client by their code from storage.
func (s *Storage) ClientByCode(ctx context.Context, code string) (domain.Client, error) {
	const op = "sqlite.ClientByCode"

	stmt, err := s.db.PrepareContext(ctx,
		`select c.id, c.name, c.code, c.secret_key, c.deleted, c.created_at, c.updated_at
		from clients c
		where c.code = ?;`)
	if err != nil {
		return domain.Client{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, code)

	var client domain.Client
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
			return domain.Client{}, fmt.Errorf("%s: %w", op, storage.ErrEntityNotFound)
		}

		return domain.Client{}, fmt.Errorf("%s: %w", op, err)
	}

	return client, nil
}

// CreateUser creates a new user in the storage and returns new user ID.
func (s *Storage) CreateUser(ctx context.Context, user domain.User) (int64, error) {
	const op = "sqlite.CreateUser"

	stmt, err := s.db.PrepareContext(ctx,
		`insert into users (username, password_hash, fio, date_of_birth, gender, avatar_file__key, deleted, created_at, updated_at)
		values(?, ?, ?, ?, ?, ?, ?, ?, ?);`)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.ExecContext(
		ctx,
		user.Username,
		user.PasswordHash,
		user.FIO,
		user.DateOfBirth,
		user.Gender,
		user.AvatarFileKey,
		user.Deleted,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrEntityExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	newUserID, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return newUserID, nil
}

// CreateUserClientLink creates a user's client link and returns new link ID.
func (s *Storage) CreateUserClientLink(ctx context.Context, userClientLink domain.UserClientLink) (int64, error) {
	const op = "sqlite.CreateUserClientLink"

	stmt, err := s.db.PrepareContext(ctx,
		`insert into user_clients (user_id, client_id, created_at, updated_at)
		values (?, ?, ?, ?);`)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.ExecContext(
		ctx,
		userClientLink.UserID,
		userClientLink.ClientID,
		userClientLink.CreatedAt,
		userClientLink.UpdatedAt,
	)

	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrEntityExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	newLinkID, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return newLinkID, nil
}

// RemoveSession removes a session by ID.
func (s *Storage) RemoveSession(ctx context.Context, id int64) error {
	const op = "sqlite.RemoveSession"

	stmt, err := s.db.PrepareContext(ctx, `delete from sessions where id=?;`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// User returns a user from the storage by ID.
func (s *Storage) User(ctx context.Context, id int64) (domain.User, error) {
	const op = "sqlite.User"

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
		where u.id = ?;`)
	if err != nil {
		return domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, id)

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
			return domain.User{}, fmt.Errorf("%s: %w", op, storage.ErrEntityNotFound)
		}

		return domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, err
}

// SessionByRefreshToken returns a session by its refresh token.
func (s *Storage) SessionByRefreshToken(ctx context.Context, refreshToken string) (domain.Session, error) {
	const op = "sqlite.SessionByRefreshToken"

	stmt, err := s.db.PrepareContext(ctx,
		`select id, user_id, refresh_token, user_agent, fingerprint, expires_at, created_at, updated_at
		from sessions s
		where s.refresh_token = ?;`)
	if err != nil {
		return domain.Session{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, refreshToken)

	var session domain.Session
	err = row.Scan(
		&session.ID,
		&session.UserID,
		&session.RefreshToken,
		&session.UserAgent,
		&session.Fingerprint,
		&session.ExpiresAt,
		&session.CreatedAt,
		&session.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Session{}, fmt.Errorf("%s: %w", op, storage.ErrEntityNotFound)
		}

		return domain.Session{}, fmt.Errorf("%s: %w", op, err)
	}

	return session, err
}
