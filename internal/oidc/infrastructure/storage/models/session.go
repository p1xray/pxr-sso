package models

import (
	"github.com/google/uuid"
	"time"
)

type Session struct {
	ID               int64
	Code             uuid.UUID
	ClientID         int64
	UserID           int64
	AuthTime         time.Time
	IdentityProvider string
	CreatedAt        time.Time
	UpdatedAt        time.Time

	Client            Client
	User              User
	GrantedScopeLinks []SessionGrantedScopeLink
}

type SessionOption func(*Session)

func SessionCreated() SessionOption {
	now := time.Now()
	return func(s *Session) {
		s.CreatedAt = now
		s.UpdatedAt = now
	}
}

func SessionUpdated() SessionOption {
	return func(s *Session) {
		s.UpdatedAt = time.Now()
	}
}
