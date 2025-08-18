package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"github.com/p1xray/pxr-sso/internal/infrastructure"
	"github.com/p1xray/pxr-sso/internal/infrastructure/storage/models"
	"strings"
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

func (s *Storage) User(ctx context.Context, id int64) (models.User, error) {
	const op = "sqlite.User"

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
	const op = "sqlite.UserByUsername"

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

func (s *Storage) CreateUser(ctx context.Context, user models.User) (int64, error) {
	const op = "sqlite.CreateUser"

	stmt, err := s.db.PrepareContext(ctx,
		`insert into users (
		   username,
		   password_hash,
		   fio,
		   date_of_birth,
		   gender,
		   avatar_file_key,
		   deleted,
		   created_at,
		   updated_at)
		values(?, ?, ?, ?, ?, ?, ?, ?, ?);`)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.ExecContext(
		ctx,
		user.Username,
		user.PasswordHash,
		user.FullName,
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
			return 0, fmt.Errorf("%s: %w", op, infrastructure.ErrEntityExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) UpdateUser(ctx context.Context, user models.User) error {
	const op = "sqlite.UpdateUser"

	stmt, err := s.db.PrepareContext(ctx,
		`update users
		 set username = ?,
			 password_hash = ?,
			 fio = ?,
			 date_of_birth = ?,
			 gender = ?,
			 avatar_file_key = ?,
			 deleted = ?,
			 created_at = ?,
			 updated_at = ?
		 where id = ?;`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.ExecContext(
		ctx,
		user.Username,
		user.PasswordHash,
		user.FullName,
		user.DateOfBirth,
		user.Gender,
		user.AvatarFileKey,
		user.Deleted,
		user.CreatedAt,
		user.UpdatedAt,
		user.ID,
	)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) RemoveUser(ctx context.Context, user models.User) error {
	const op = "sqlite.RemoveUser"

	stmt, err := s.db.PrepareContext(ctx,
		`update users
		 set deleted = ?,
			 updated_at = ?
		 where id = ?;`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.ExecContext(
		ctx,
		user.Deleted,
		user.UpdatedAt,
		user.ID,
	)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) RolesByUserID(ctx context.Context, userID int64) ([]models.Role, error) {
	const op = "sqlite.RolesByUserID"

	stmt, err := s.db.PrepareContext(ctx,
		`select
			 r.id,
			 r.code,
			 r.name,
			 r.description,
			 r.active,
			 r.deleted,
			 r.created_at,
			 r.updated_at
		 from roles r
			 join user_roles ur on ur.role_id = r.id
		 where r.active is true and ur.user_id = ?;`)
	if err != nil {
		return []models.Role{}, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := stmt.QueryContext(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	roles := make([]models.Role, 0)
	for rows.Next() {
		role := models.Role{}
		err = rows.Scan(
			&role.ID,
			&role.Code,
			&role.Name,
			&role.Description,
			&role.Active,
			&role.Deleted,
			&role.CreatedAt,
			&role.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		roles = append(roles, role)
	}

	return roles, nil
}

func (s *Storage) RolesByClientID(ctx context.Context, clientID int64) ([]models.Role, error) {
	const op = "sqlite.RolesByClientID"

	stmt, err := s.db.PrepareContext(ctx,
		`select
			  r.id,
			  r.code,
			  r.name,
			  r.description,
			  r.active,
			  r.deleted,
			  r.created_at,
			  r.updated_at
		  from roles r
			  join client_default_roles cdr on cdr.role_id = r.id
		  where r.active is true and cdr.client_id = ?;`)
	if err != nil {
		return []models.Role{}, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := stmt.QueryContext(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	roles := make([]models.Role, 0)
	for rows.Next() {
		role := models.Role{}
		err = rows.Scan(
			&role.ID,
			&role.Code,
			&role.Name,
			&role.Description,
			&role.Active,
			&role.Deleted,
			&role.CreatedAt,
			&role.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		roles = append(roles, role)
	}

	return roles, nil
}

func (s *Storage) PermissionsByUserID(ctx context.Context, userID int64) ([]models.Permission, error) {
	const op = "sqlite.PermissionsByUserID"

	stmt, err := s.db.PrepareContext(ctx,
		`select
			 p.id,
			 p.code,
			 p.description,
			 p.active,
			 p.deleted,
			 p.created_at,
			 p.updated_at
		 from permissions p
			 join role_permissions rp on rp.permission_id = p.id
			 join user_roles ur on ur.role_id = rp.role_id
		 where p.active is true and ur.user_id = ?;`)
	if err != nil {
		return []models.Permission{}, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := stmt.QueryContext(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	permissions := make([]models.Permission, 0)
	for rows.Next() {
		permission := models.Permission{}
		err = rows.Scan(
			&permission.ID,
			&permission.Code,
			&permission.Description,
			&permission.Active,
			&permission.Deleted,
			&permission.CreatedAt,
			&permission.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		permissions = append(permissions, permission)
	}

	return permissions, nil
}

func (s *Storage) PermissionsByRoleCodes(ctx context.Context, roleCodes []string) ([]models.Permission, error) {
	const op = "sqlite.PermissionsByRoleCodes"

	// Generate the placeholders for the IN clause.
	placeholders := make([]string, len(roleCodes))
	for i := range roleCodes {
		placeholders[i] = "?"
	}
	inClause := strings.Join(placeholders, ",")

	query := fmt.Sprintf(`select
		 p.id,
		 p.code,
		 p.description,
		 p.active,
		 p.deleted,
		 p.created_at,
		 p.updated_at
	 from permissions p
		 join role_permissions rp on rp.permission_id = p.id
		 join roles r on rp.role_id = r.id
	 where p.active is true and r.code in (%s);`, inClause)

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return []models.Permission{}, fmt.Errorf("%s: %w", op, err)
	}

	args := make([]interface{}, len(roleCodes))
	for i, rc := range roleCodes {
		args[i] = rc
	}

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	permissions := make([]models.Permission, 0)
	for rows.Next() {
		permission := models.Permission{}
		err = rows.Scan(
			&permission.ID,
			&permission.Code,
			&permission.Description,
			&permission.Active,
			&permission.Deleted,
			&permission.CreatedAt,
			&permission.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		permissions = append(permissions, permission)
	}

	return permissions, nil
}

func (s *Storage) SessionsByUserID(ctx context.Context, userID int64) ([]models.Session, error) {
	const op = "sqlite.SessionsByUserID"

	stmt, err := s.db.PrepareContext(ctx,
		`select
			 s.id,
			 s.user_id,
			 s.refresh_token,
			 s.user_agent,
			 s.fingerprint,
			 s.expires_at,
			 s.created_at,
			 s.updated_at
		 from sessions s
		 where s.user_id = ?;`)
	if err != nil {
		return []models.Session{}, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := stmt.QueryContext(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	sessions := make([]models.Session, 0)
	for rows.Next() {
		session := models.Session{}
		err = rows.Scan(
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
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		sessions = append(sessions, session)
	}

	return sessions, nil
}

func (s *Storage) SessionByRefreshTokenID(ctx context.Context, refreshTokenID string) (models.Session, error) {
	const op = "sqlite.SessionByRefreshTokenID"

	stmt, err := s.db.PrepareContext(ctx,
		`select
			 s.id,
			 s.user_id,
			 s.refresh_token,
			 s.user_agent,
			 s.fingerprint,
			 s.expires_at,
			 s.created_at,
			 s.updated_at
		 from sessions s
		 where s.refresh_token = ?;`)
	if err != nil {
		return models.Session{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, refreshTokenID)

	var session models.Session
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
			return models.Session{}, fmt.Errorf("%s: %w", op, infrastructure.ErrEntityNotFound)
		}

		return models.Session{}, fmt.Errorf("%s: %w", op, err)
	}

	return session, nil
}

func (s *Storage) CreateSession(ctx context.Context, session models.Session) (int64, error) {
	const op = "sqlite.CreateSession"

	stmt, err := s.db.PrepareContext(ctx,
		`insert into sessions (
			 user_id,
			 refresh_token,
			 user_agent,
			 fingerprint,
			 expires_at,
			 created_at,
			 updated_at)
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
			return 0, fmt.Errorf("%s: %w", op, infrastructure.ErrEntityExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) UpdateSession(ctx context.Context, session models.Session) error {
	const op = "sqlite.UpdateSession"

	stmt, err := s.db.PrepareContext(ctx,
		`update sessions
		 set user_id = ?,
			 refresh_token = ?,
			 user_agent = ?,
			 fingerprint = ?,
			 expires_at = ?,
			 created_at = ?,
			 updated_at = ?
		 where id = ?;`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.ExecContext(
		ctx,
		session.UserID,
		session.RefreshToken,
		session.UserAgent,
		session.Fingerprint,
		session.ExpiresAt,
		session.CreatedAt,
		session.UpdatedAt,
		session.ID,
	)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) RemoveSession(ctx context.Context, id int64) error {
	const op = "sqlite.RemoveSession"

	stmt, err := s.db.PrepareContext(ctx, `delete from sessions where id = ?;`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) ClientByCodeAndUserID(ctx context.Context, code string, userID int64) (models.Client, error) {
	const op = "sqlite.ClientByCodeAndUserID"

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
			 join user_clients uc on uc.client_id = c.id
		 where c.code = ? and uc.user_id = ?;`)
	if err != nil {
		return models.Client{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, code, userID)

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

func (s *Storage) ClientByCode(ctx context.Context, code string) (models.Client, error) {
	const op = "sqlite.ClientByCode"

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

func (s *Storage) CreateUserClientLink(ctx context.Context, userClientLink models.UserClientLink) (int64, error) {
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
			return 0, fmt.Errorf("%s: %w", op, infrastructure.ErrEntityExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) CreateUserRoleLink(ctx context.Context, userRoleLink models.UserRoleLink) (int64, error) {
	const op = "sqlite.CreateUserRoleLink"

	stmt, err := s.db.PrepareContext(ctx,
		`insert into user_roles (user_id, role_id, created_at, updated_at)
		 values (?, ?, ?, ?);`)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.ExecContext(
		ctx,
		userRoleLink.UserID,
		userRoleLink.RoleID,
		userRoleLink.CreatedAt,
		userRoleLink.UpdatedAt,
	)

	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
			return 0, fmt.Errorf("%s: %w", op, infrastructure.ErrEntityExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
