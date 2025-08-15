package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/p1xray/pxr-sso/internal/dto"
	"github.com/p1xray/pxr-sso/internal/entity"
	"github.com/p1xray/pxr-sso/internal/infrastructure"
	"github.com/p1xray/pxr-sso/internal/infrastructure/converter"
	"github.com/p1xray/pxr-sso/internal/infrastructure/storage/models"
	"github.com/p1xray/pxr-sso/internal/lib/logger/sl"
	"log/slog"
)

const emptyID = 0

type AuthStorage interface {
	User(ctx context.Context, id int64) (models.User, error)
	UserByUsername(ctx context.Context, username string) (models.User, error)
	CreateUser(ctx context.Context, user models.User) (int64, error)
	UpdateUser(ctx context.Context, user models.User) error
	RemoveUser(ctx context.Context, user models.User) error

	RolesByUserID(ctx context.Context, userID int64) ([]models.Role, error)

	PermissionsByUserID(ctx context.Context, userID int64) ([]models.Permission, error)

	SessionsByUserID(ctx context.Context, userID int64) ([]models.Session, error)
	SessionByRefreshTokenID(ctx context.Context, refreshTokenID string) (models.Session, error)
	CreateSession(ctx context.Context, session models.Session) (int64, error)
	UpdateSession(ctx context.Context, session models.Session) error
	RemoveSession(ctx context.Context, id int64) error

	ClientByCodeAndUserID(ctx context.Context, code string, userID int64) (models.Client, error)
	ClientByCode(ctx context.Context, code string) (models.Client, error)
}

type Auth struct {
	log     *slog.Logger
	storage AuthStorage
}

func NewAuthRepository(log *slog.Logger, storage AuthStorage) *Auth {
	return &Auth{
		log:     log,
		storage: storage,
	}
}

func (a *Auth) ClientByCode(ctx context.Context, code string) (dto.Client, error) {
	const op = "repository.auth.ClientByCode"

	log := a.log.With(
		slog.String("op", op),
		slog.String("code", code),
	)

	client, err := a.storage.ClientByCode(ctx, code)
	if err != nil {
		if errors.Is(err, infrastructure.ErrEntityNotFound) {
			log.Warn("client not found", sl.Err(err))
		} else {
			log.Error("error getting user client", sl.Err(err))
		}

		return dto.Client{}, fmt.Errorf("%s: %w", op, err)
	}

	clientDTO := converter.ToClientDTO(client)

	return clientDTO, nil
}

func (a *Auth) DataForLogin(ctx context.Context, username, clientCode string) (dto.DataForLogin, error) {
	const op = "repository.auth.DataForLogin"

	log := a.log.With(
		slog.String("op", op),
		slog.String("username", username),
		slog.String("client code", clientCode),
	)

	userDTO, err := a.userByUsername(ctx, log, username)
	if err != nil {
		return dto.DataForLogin{}, fmt.Errorf("%s: %w", op, err)
	}

	client, err := a.storage.ClientByCodeAndUserID(ctx, clientCode, userDTO.ID)
	if err != nil {
		if errors.Is(err, infrastructure.ErrEntityNotFound) {
			log.Warn("client not found", sl.Err(err))
		} else {
			log.Error("error getting user client", sl.Err(err))

			return dto.DataForLogin{}, fmt.Errorf("%s: %w", op, err)
		}
	}

	clientDTO := converter.ToClientDTO(client)

	userSessions, err := a.storage.SessionsByUserID(ctx, userDTO.ID)
	if err != nil {
		log.Error("error getting user sessions", sl.Err(err))

		return dto.DataForLogin{}, fmt.Errorf("%s: %w", op, err)
	}

	sessionsDTO := make([]dto.Session, len(userSessions))
	for i, userSession := range userSessions {
		sessionsDTO[i] = converter.ToSessionDTO(userSession)
	}

	return dto.DataForLogin{
		User:     userDTO,
		Client:   clientDTO,
		Sessions: sessionsDTO,
	}, nil
}

func (a *Auth) DataForRegister(ctx context.Context, username, clientCode string) (dto.DataForRegister, error) {
	const op = "repository.auth.DataForRegister"

	log := a.log.With(
		slog.String("op", op),
		slog.String("username", username),
		slog.String("client code", clientCode),
	)

	userDTO, err := a.userByUsername(ctx, log, username)
	if err != nil && !errors.Is(err, infrastructure.ErrEntityNotFound) {
		return dto.DataForRegister{}, fmt.Errorf("%s: %w", op, err)
	}

	clientDTO, err := a.ClientByCode(ctx, clientCode)
	if err != nil {
		return dto.DataForRegister{}, fmt.Errorf("%s: %w", op, err)
	}

	return dto.DataForRegister{
		User:   userDTO,
		Client: clientDTO,
	}, nil
}

func (a *Auth) DataForRefreshTokens(ctx context.Context, refreshTokenID string) (dto.DataForRefreshTokens, error) {
	const op = "repository.auth.DataForRefreshTokens"

	log := a.log.With(
		slog.String("op", op),
		slog.String("refresh token ID", refreshTokenID),
	)

	sessionDTO, err := a.sessionByRefreshTokenID(ctx, log, refreshTokenID)
	if err != nil {
		return dto.DataForRefreshTokens{}, fmt.Errorf("%s: %w", op, err)
	}

	userDTO, err := a.user(ctx, log, sessionDTO.UserID)

	return dto.DataForRefreshTokens{
		Session: sessionDTO,
		User:    userDTO,
	}, nil
}

func (a *Auth) DataForLogout(ctx context.Context, refreshTokenID string) (dto.DataForLogout, error) {
	const op = "repository.auth.DataForLogout"

	log := a.log.With(
		slog.String("op", op),
		slog.String("refresh token ID", refreshTokenID),
	)

	sessionDTO, err := a.sessionByRefreshTokenID(ctx, log, refreshTokenID)
	if err != nil {
		return dto.DataForLogout{}, fmt.Errorf("%s: %w", op, err)
	}

	return dto.DataForLogout{
		Session: sessionDTO,
	}, nil
}

func (a *Auth) Save(ctx context.Context, auth entity.Auth) error {
	const op = "repository.auth.Save"

	log := a.log.With(
		slog.String("op", op),
	)

	if err := a.SaveUser(ctx, auth.User); err != nil {
		log.Error("error saving user", sl.Err(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	for _, session := range auth.Sessions {
		if err := a.SaveSession(ctx, session); err != nil {
			log.Error("error saving session", sl.Err(err))

			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}

func (a *Auth) SaveUser(ctx context.Context, user entity.User) error {
	const op = "repository.auth.SaveUser"

	log := a.log.With(
		slog.String("op", op),
	)

	if user.IsToCreate() {
		if err := a.createUser(ctx, user); err != nil {
			log.Error("error creating user", sl.Err(err))

			return fmt.Errorf("%s: %w", op, err)
		}
	}

	if user.IsToUpdate() {
		if err := a.updateUser(ctx, user); err != nil {
			log.Error("error updating user", sl.Err(err))

			return fmt.Errorf("%s: %w", op, err)
		}
	}

	if user.IsToRemove() {
		if err := a.removeUser(ctx, user); err != nil {
			log.Error("error removing user", sl.Err(err))

			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}

func (a *Auth) createUser(ctx context.Context, user entity.User) error {
	userStorageModel := converter.ToUserStorage(user, models.UserCreated())

	id, err := a.storage.CreateUser(ctx, userStorageModel)
	if err != nil {
		return err
	}

	user.ID = id

	return nil
}

func (a *Auth) updateUser(ctx context.Context, user entity.User) error {
	if user.ID == emptyID {
		return infrastructure.ErrRequireIDToUpdate
	}

	userStorageModel := converter.ToUserStorage(user, models.UserUpdated())

	return a.storage.UpdateUser(ctx, userStorageModel)
}

func (a *Auth) removeUser(ctx context.Context, user entity.User) error {
	if user.ID == emptyID {
		return infrastructure.ErrRequireIDToRemove
	}

	userStorageModel := converter.ToUserStorage(user, models.UserRemoved())

	return a.storage.RemoveUser(ctx, userStorageModel)
}

func (a *Auth) SaveSession(ctx context.Context, session entity.Session) error {
	const op = "repository.auth.SaveSession"

	log := a.log.With(
		slog.String("op", op),
	)

	if session.IsToCreate() {
		if err := a.createSession(ctx, session); err != nil {
			log.Error("error creating session", sl.Err(err))

			return fmt.Errorf("%s: %w", op, err)
		}
	}

	if session.IsToUpdate() {
		if err := a.updateSession(ctx, session); err != nil {
			log.Error("error updating session", sl.Err(err))

			return fmt.Errorf("%s: %w", op, err)
		}
	}

	if session.IsToRemove() {
		if err := a.removeSession(ctx, session); err != nil {
			log.Error("error removing session", sl.Err(err))

			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}

func (a *Auth) createSession(ctx context.Context, session entity.Session) error {
	sessionStorageModel := converter.ToSessionStorage(session, models.SessionCreated())

	id, err := a.storage.CreateSession(ctx, sessionStorageModel)
	if err != nil {
		return err
	}

	session.ID = id

	return nil
}

func (a *Auth) updateSession(ctx context.Context, session entity.Session) error {
	if session.ID == emptyID {
		return infrastructure.ErrRequireIDToUpdate
	}

	sessionStorageModel := converter.ToSessionStorage(session, models.SessionUpdated())

	return a.storage.UpdateSession(ctx, sessionStorageModel)
}

func (a *Auth) removeSession(ctx context.Context, session entity.Session) error {
	if session.ID == emptyID {
		return infrastructure.ErrRequireIDToRemove
	}

	return a.storage.RemoveSession(ctx, session.ID)
}

func (a *Auth) user(ctx context.Context, log *slog.Logger, id int64) (dto.User, error) {
	user, err := a.storage.User(ctx, id)
	if err != nil {
		if errors.Is(err, infrastructure.ErrEntityNotFound) {
			log.Warn("user not found", sl.Err(err))
		} else {
			log.Error("error getting user", sl.Err(err))
		}

		return dto.User{}, err
	}

	userDTO, err := a.userWithRolesPermissions(ctx, log, user)
	if err != nil {
		return dto.User{}, err
	}

	return userDTO, nil
}

func (a *Auth) userByUsername(ctx context.Context, log *slog.Logger, username string) (dto.User, error) {
	user, err := a.storage.UserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, infrastructure.ErrEntityNotFound) {
			log.Warn("user not found", sl.Err(err))
		} else {
			log.Error("error getting user", sl.Err(err))

			return dto.User{}, err
		}
	}

	userDTO, err := a.userWithRolesPermissions(ctx, log, user)
	if err != nil {
		return dto.User{}, err
	}

	return userDTO, nil
}

func (a *Auth) userWithRolesPermissions(ctx context.Context, log *slog.Logger, user models.User) (dto.User, error) {
	userRoles, err := a.storage.RolesByUserID(ctx, user.ID)
	if err != nil {
		log.Error("error getting user roles", sl.Err(err))

		return dto.User{}, err
	}

	userPermissions, err := a.storage.PermissionsByUserID(ctx, user.ID)
	if err != nil {
		log.Error("error getting user permissions", sl.Err(err))

		return dto.User{}, err
	}

	userDTO := converter.ToUserDTO(user, userRoles, userPermissions)

	return userDTO, nil
}

func (a *Auth) sessionByRefreshTokenID(ctx context.Context, log *slog.Logger, refreshTokenID string) (dto.Session, error) {
	session, err := a.storage.SessionByRefreshTokenID(ctx, refreshTokenID)
	if err != nil {
		if errors.Is(err, infrastructure.ErrEntityNotFound) {
			log.Warn("session not found", sl.Err(err))
		} else {
			log.Error("error getting session", sl.Err(err))
		}

		return dto.Session{}, err
	}

	sessionDTO := converter.ToSessionDTO(session)

	return sessionDTO, nil
}
