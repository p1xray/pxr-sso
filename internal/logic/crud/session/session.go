package sessioncrud

import (
	"context"
	"fmt"
	"log/slog"
	"pxr-sso/internal/domain"
	"pxr-sso/internal/lib/logger/sl"
	"pxr-sso/internal/logic/crud"
	"pxr-sso/internal/logic/dto"
	"time"
)

// CRUD provides methods for managing session data.
type CRUD struct {
	log             *slog.Logger
	sessionProvider crud.SessionProvider
	sessionSaver    crud.SessionSaver
}

// New creates a new instance of the session's CRUD.
func New(
	log *slog.Logger,
	sessionProvider crud.SessionProvider,
	sessionSaver crud.SessionSaver,
) *CRUD {
	return &CRUD{
		log:             log,
		sessionProvider: sessionProvider,
		sessionSaver:    sessionSaver,
	}
}

// CreateSession creates a new session in the storage.
func (c *CRUD) CreateSession(ctx context.Context, session dto.CreateSessionDTO) error {
	const op = "sessioncrud.CreateSession"

	log := c.log.With(
		slog.String("op", op),
		slog.Int64("user ID", session.UserID),
		slog.String("refresh token", session.RefreshToken),
	)

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
	if err := c.sessionSaver.CreateSession(ctx, sessionToCreate); err != nil {
		log.Error("failed to create session in storage", sl.Err(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
