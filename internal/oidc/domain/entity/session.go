package entity

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/p1xray/pxr-sso/internal/oidc/domain/dto"
	"github.com/p1xray/pxr-sso/internal/oidc/domain/enum"
	"time"
)

// Session represents a model of an authentication session for a logged-in user,
// capturing essential details about the user's authentication state and interactions within the system.
type Session struct {
	id               int64
	code             uuid.UUID
	client           dto.Client
	user             dto.User
	authTime         time.Time
	identityProvider string
	scopes           []dto.Scope

	dataStatus enum.DataStatus
}

func NewSession(client dto.Client, user dto.User) (Session, error) {
	code, err := uuid.NewV7()
	if err != nil {
		return Session{}, fmt.Errorf("generate session code: %w", err)
	}

	session := Session{
		id:               0,
		code:             code,
		client:           client,
		user:             user,
		authTime:         time.Now(),
		identityProvider: "pxr.sso",
		scopes:           make([]dto.Scope, 0),
		dataStatus:       enum.DataStatusToCreate,
	}

	return session, nil
}

func NewExistSession(data dto.Session) Session {
	return Session{
		id:               data.ID(),
		code:             data.Code(),
		client:           data.Client(),
		user:             data.User(),
		authTime:         data.AuthTime(),
		identityProvider: data.IdentityProvider(),
		scopes:           data.Scopes(),
	}
}

func (s *Session) Update(authTime time.Time, scopes []dto.Scope) {
	s.authTime = authTime
	s.scopes = scopes
	s.dataStatus = enum.DataStatusToUpdate
}

func (s *Session) UpdateScopes(scopes []dto.Scope) {
	s.scopes = scopes
	s.dataStatus = enum.DataStatusToUpdate
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

func (s *Session) Client() dto.Client {
	return s.client
}

func (s *Session) User() dto.User {
	return s.user
}

func (s *Session) AuthTime() time.Time {
	return s.authTime
}

func (s *Session) IdentityProvider() string {
	return s.identityProvider
}

func (s *Session) Scopes() []dto.Scope {
	return s.scopes
}

func (s *Session) DataStatus() enum.DataStatus {
	return s.dataStatus
}
