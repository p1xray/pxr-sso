package grant

import (
	"context"
	"fmt"
	"github.com/p1xray/pxr-sso/internal/oidc/application/generator"
	"github.com/p1xray/pxr-sso/internal/oidc/domain/dto"
	"github.com/p1xray/pxr-sso/internal/oidc/domain/entity"
	"time"
)

type CallbackURIBuilder interface {
	BuildCallbackRedirectURI(rawURL, authorizationCode, state string) (string, error)
}

type AuthorizedGrantSaver interface {
	SaveAuthorizedGrant(ctx context.Context, authorizationCode string, authorizedGrant dto.AuthorizedGrant) error
}

type SessionUpdater interface {
	UpdateSession(ctx context.Context, session entity.Session) (dto.Session, error)
}

type AuthorizedGrant interface {
	Execute(ctx context.Context, data dto.AuthorizeContext) (string, error)
}

type usecase struct {
	uriBuilder           CallbackURIBuilder
	sessionUpdater       SessionUpdater
	authorizedGrantSaver AuthorizedGrantSaver
}

// NewUseCase returns the new use case for processing authorized grant.
func NewUseCase(
	uriBuilder CallbackURIBuilder,
	sessionUpdater SessionUpdater,
	authorizedGrantSaver AuthorizedGrantSaver,
) *usecase {
	return &usecase{
		uriBuilder:           uriBuilder,
		sessionUpdater:       sessionUpdater,
		authorizedGrantSaver: authorizedGrantSaver,
	}
}

// Execute ...
func (u *usecase) Execute(ctx context.Context, data dto.AuthorizeContext) (string, error) {
	const op = "process authorized grant"

	session, err := data.SingleSession()
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	authTime := time.Now()
	scopes := data.RequestGrantedScopes()
	if data.RequestPrompt() == "none" {
		authTime = session.AuthTime()
		scopes = session.Scopes()
	}

	sessionEntity := entity.NewExistSession(session)
	sessionEntity.Update(authTime, scopes)

	updatedSession, err := u.sessionUpdater.UpdateSession(ctx, sessionEntity)
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
