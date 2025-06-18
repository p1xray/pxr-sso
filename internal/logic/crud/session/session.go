package sessioncrud

import (
	"context"
	"github.com/p1xray/pxr-sso/internal/logic/crud"
	"github.com/p1xray/pxr-sso/internal/logic/dto"
	"github.com/p1xray/pxr-sso/internal/storage/domain"
	"time"
)

// CRUD provides methods for managing session data.
type CRUD struct {
	sessionProvider crud.SessionProvider
	sessionSaver    crud.SessionSaver
}

// New creates a new instance of the session's CRUD.
func New(
	sessionProvider crud.SessionProvider,
	sessionSaver crud.SessionSaver,
) *CRUD {
	return &CRUD{
		sessionProvider: sessionProvider,
		sessionSaver:    sessionSaver,
	}
}

// SessionByRefreshToken returns a session by its refresh token.
func (c *CRUD) SessionByRefreshToken(ctx context.Context, refreshToken string) (domain.Session, error) {
	session, err := c.sessionProvider.SessionByRefreshToken(ctx, refreshToken)
	if err != nil {
		return domain.Session{}, err
	}

	return session, nil
}

// CreateSession creates a new session in the storage.
func (c *CRUD) CreateSession(ctx context.Context, session dto.CreateSessionDTO) error {
	now := time.Now()
	sessionToCreate := domain.Session{
		UserID:       session.UserID,
		RefreshToken: session.RefreshToken,
		UserAgent:    session.UserAgent,
		Fingerprint:  session.Fingerprint,
		ExpiresAt:    session.ExpiresAt,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if _, err := c.sessionSaver.CreateSession(ctx, sessionToCreate); err != nil {
		return err
	}

	return nil
}

// RemoveSession removes a session by ID.
func (c *CRUD) RemoveSession(ctx context.Context, id int64) error {
	if err := c.sessionSaver.RemoveSession(ctx, id); err != nil {
		return err
	}

	return nil
}
