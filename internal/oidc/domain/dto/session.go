package dto

import (
	"github.com/google/uuid"
	"time"
)

type Session struct {
	id               int64
	code             uuid.UUID
	client           Client
	user             User
	authTime         time.Time
	identityProvider string
	scopes           []Scope
}

func NewSession(
	id int64,
	code uuid.UUID,
	client Client,
	user User,
	authTime time.Time,
	identityProvider string,
	scopes []Scope,
) Session {
	return Session{
		id:               id,
		code:             code,
		client:           client,
		user:             user,
		authTime:         authTime,
		identityProvider: identityProvider,
		scopes:           scopes,
	}
}

func (s *Session) ID() int64 {
	return s.id
}

func (s *Session) Code() uuid.UUID {
	return s.code
}

func (s *Session) CodeString() string {
	return s.code.String()
}

func (s *Session) Client() Client {
	return s.client
}

func (s *Session) User() User {
	return s.user
}

func (s *Session) UserID() int64 {
	return s.user.ID()
}

func (s *Session) AuthTime() time.Time {
	return s.authTime
}

func (s *Session) IdentityProvider() string {
	return s.identityProvider
}

func (s *Session) Scopes() []Scope {
	return s.scopes
}

func (s *Session) ScopeCodes() []string {
	scopeCodes := make([]string, len(s.Scopes()))
	for i, scope := range s.Scopes() {
		scopeCodes[i] = scope.Code()
	}

	return scopeCodes
}

type SessionCookie struct {
	name  string
	value string
}

func NewSessionCookie(name, value string) SessionCookie {
	return SessionCookie{
		name:  name,
		value: value,
	}
}

func (s *SessionCookie) Name() string {
	return s.name
}

func (s *SessionCookie) Value() string {
	return s.value
}

type AuthorizedSession struct {
	id               string
	subject          int64
	authTime         time.Time
	identityProvider string
}

func NewAuthorizedSession(
	id string,
	subject int64,
	authTime time.Time,
	identityProvider string,
) AuthorizedSession {
	return AuthorizedSession{
		id:               id,
		subject:          subject,
		authTime:         authTime,
		identityProvider: identityProvider,
	}
}

func (s *AuthorizedSession) ID() string {
	return s.id
}

func (s *AuthorizedSession) Subject() int64 {
	return s.subject
}

func (s *AuthorizedSession) AuthTime() time.Time {
	return s.authTime
}

func (s *AuthorizedSession) IdentityProvider() string {
	return s.identityProvider
}
