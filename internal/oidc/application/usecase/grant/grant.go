package grant

import (
	"context"
	"fmt"
	"github.com/p1xray/pxr-sso/internal/oidc/application/generator"
	"github.com/p1xray/pxr-sso/internal/oidc/domain/dto"
	"github.com/p1xray/pxr-sso/internal/oidc/domain/entity"
	"time"
)

// CallbackURIBuilder is the URI builder to redirect the user agent to the client's callback page.
type CallbackURIBuilder interface {
	// BuildCallbackRedirectURI return the URI to redirect to the client's callback page.
	BuildCallbackRedirectURI(rawURL, authorizationCode, state string) (string, error)
}

// AuthorizedGrantSaver is the saver authorized grant to a storage.
type AuthorizedGrantSaver interface {
	// SaveAuthorizedGrant saves the authorized grant to a storage.
	SaveAuthorizedGrant(ctx context.Context, authorizationCode string, authorizedGrant dto.AuthorizedGrant) error
}

// SessionSaver is the saver session data to the storage.
type SessionSaver interface {
	// SaveSession saves the session data to the storage.
	SaveSession(ctx context.Context, session entity.Session) (dto.Session, error)
}

// ScopeReader is the scopes data reader from storage.
type ScopeReader interface {
	// ScopesByCode returns a scopes data by code.
	ScopesByCode(ctx context.Context, codes []string) ([]dto.Scope, error)
}

// AuthorizationGrantFlow is the processor for OAuth 2.0 authorization grant flow.
type AuthorizationGrantFlow interface {
	// Execute processes OAuth 2.0 authorization grant flow.
	Execute(ctx context.Context, data dto.AuthorizeContext) (string, error)
}

// usecase is the use case which processes OAuth 2.0 authorization grant flow.
type usecase struct {
	uriBuilder           CallbackURIBuilder
	sessionSaver         SessionSaver
	scopeReader          ScopeReader
	authorizedGrantSaver AuthorizedGrantSaver
}

// NewUseCase creates a new use case for processing OAuth 2.0 authorization grant flow.
func NewUseCase(
	uriBuilder CallbackURIBuilder,
	sessionSaver SessionSaver,
	scopeReader ScopeReader,
	authorizedGrantSaver AuthorizedGrantSaver,
) *usecase {
	return &usecase{
		uriBuilder:           uriBuilder,
		sessionSaver:         sessionSaver,
		scopeReader:          scopeReader,
		authorizedGrantSaver: authorizedGrantSaver,
	}
}

// Execute processes OAuth 2.0 authorization grant flow.
//
// This method updates existing session, generates and stores OAuth 2.0 authorization code
// and finally generates the uri to redirect to the client's callback page.
func (u *usecase) Execute(ctx context.Context, data dto.AuthorizeContext) (string, error) {
	const op = "process authorized grant"

	session, err := data.SingleSession()
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	authTime := time.Now()
	scopeCodes := data.RequestGrantedScopes()
	if data.RequestPrompt() == "none" {
		authTime = session.AuthTime()
		scopeCodes = session.ScopeCodes()
	}

	scopes, err := u.scopeReader.ScopesByCode(ctx, scopeCodes)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	sessionEntity := entity.NewExistSession(session)
	sessionEntity.Update(authTime, scopes)

	updatedSession, err := u.sessionSaver.SaveSession(ctx, sessionEntity)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	authorizationCode := generator.AuthorizationCode()
	authorizedSession := dto.NewAuthorizedSession(
		updatedSession.CodeString(),
		updatedSession.UserID(),
		updatedSession.AuthTime(),
		updatedSession.IdentityProvider(),
	)
	authorizedGrant := dto.NewAuthorizedGrant(authorizedSession, data.Request())
	if err = u.authorizedGrantSaver.SaveAuthorizedGrant(ctx, authorizationCode, authorizedGrant); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return u.uriBuilder.BuildCallbackRedirectURI(data.RequestRedirectURI(), authorizationCode, data.RequestState())
}
