package models

import "time"

type SessionGrantedScopeLink struct {
	ID        int64
	SessionID int64
	ScopeID   int64
	CreatedAt time.Time
	UpdatedAt time.Time

	Scope Scope
}

type SessionGrantedScopeLinkOption func(*SessionGrantedScopeLink)

func SessionGrantedScopeLinkCreated() SessionGrantedScopeLinkOption {
	now := time.Now()
	return func(l *SessionGrantedScopeLink) {
		l.CreatedAt = now
		l.UpdatedAt = now
	}
}
