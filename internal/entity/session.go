package entity

import (
	"fmt"
	"github.com/p1xray/pxr-sso/internal/enum"
	"time"
)

type Session struct {
	ID             int64
	UserID         int64
	RefreshTokenID string
	UserAgent      string
	Fingerprint    string
	ExpiresAt      time.Time

	Tokens Tokens

	dataStatus enum.DataStatusEnum
}

func NewSession(userID int64, userAgent, fingerprint string, setters ...SessionOption) (Session, error) {
	session := Session{
		UserID:      userID,
		UserAgent:   userAgent,
		Fingerprint: fingerprint,
	}

	for _, setter := range setters {
		if err := setter(&session); err != nil {
			return Session{}, err
		}
	}

	return session, nil
}

func (s *Session) Validate(userAgent, fingerprint string) error {
	const op = "entity.Session.Validate"

	// Check session expiration time.
	now := time.Now()
	if s.ExpiresAt.Before(now) {
		return fmt.Errorf("%s: %w", op, ErrRefreshTokenExpired)
	}

	// Check session user agent and fingerprint.
	if s.UserAgent != userAgent || s.Fingerprint != fingerprint {
		return fmt.Errorf("%s: %w", op, ErrInvalidSession)
	}

	return nil
}

func (s *Session) SetToCreate() {
	s.dataStatus = enum.ToCreate
}

func (s *Session) SetToUpdate() {
	s.dataStatus = enum.ToUpdate
}

func (s *Session) SetToRemove() {
	s.dataStatus = enum.ToRemove
}

func (s *Session) IsToCreate() bool {
	return s.dataStatus == enum.ToCreate
}

func (s *Session) IsToUpdate() bool {
	return s.dataStatus == enum.ToUpdate
}

func (s *Session) IsToRemove() bool {
	return s.dataStatus == enum.ToRemove
}

func (s *Session) ResetDataStatus() {
	s.dataStatus = enum.None
}
