package authorize

import (
	"context"
	"errors"
	"fmt"
	"github.com/p1xray/pxr-sso/internal/oidc/application/validator"
	"github.com/p1xray/pxr-sso/internal/oidc/domain/dto"
	"github.com/p1xray/pxr-sso/internal/oidc/infrastructure"
	"github.com/p1xray/pxr-sso/internal/oidc/infrastructure/repository"
	"github.com/p1xray/pxr-sso/pkg/logger/sl"
	"github.com/p1xray/pxr-sso/pkg/nullable"
	"log/slog"
)

// ErrorURIBuilder is the URI builder to redirect the user agent to the error page.
type ErrorURIBuilder interface {
	// BuildErrorRedirectURI return the URI to redirect to the error page.
	BuildErrorRedirectURI(rawURL string, err error) (string, error)
}

// FlowSwitcher is the handler for defining the interaction flow.
type FlowSwitcher interface {
	// Switch processes the authorize context to decide which interaction flow to execute.
	Switch(ctx context.Context, data dto.AuthorizeContext) (string, error)
}

// ClientReader is the client data reader from storage.
type ClientReader interface {
	// ClientByCode returns a client data by code.
	ClientByCode(ctx context.Context, code string, opts ...repository.ClientOption) (dto.Client, error)
}

// SessionReader is the session data reader from storage.
type SessionReader interface {
	// SessionsByCode returns a sessions by code.
	SessionsByCode(ctx context.Context, codes []string, opts ...repository.SessionOption) ([]dto.Session, error)
}

type SessionCookieDecoder interface {
	Decode(encoded string) (dto.AuthorizedSession, error)
}

// Authorize is the handler for authorization requests, ensuring they are processed according to OAuth 2.0
// and OpenID Connect protocol specifications.
type Authorize interface {
	// Execute processes an authorization request.
	Execute(ctx context.Context, req dto.AuthorizeRequest) (string, error)
}

// usecase is the use case which handles the processing of authorization requests by validating and
// then processing these requests.
type usecase struct {
	log                  *slog.Logger
	clientReader         ClientReader
	sessionReader        SessionReader
	sessionCookieDecoder SessionCookieDecoder
	uriBuilder           ErrorURIBuilder
	switcher             FlowSwitcher
}

// NewUseCase creates a new use case for processing of authorization request.
func NewUseCase(
	log *slog.Logger,
	uriBuilder ErrorURIBuilder,
	clientReader ClientReader,
	sessionReader SessionReader,
	sessionCookieDecoder SessionCookieDecoder,
	switcher FlowSwitcher,
) *usecase {
	return &usecase{
		log:                  log,
		clientReader:         clientReader,
		sessionReader:        sessionReader,
		sessionCookieDecoder: sessionCookieDecoder,
		uriBuilder:           uriBuilder,
		switcher:             switcher,
	}
}

// Execute processes the authorization request, validating its parameters and generating an appropriate
// response that either grants or denies the authorization.
//
// This method ensures that only requests meeting the necessary validation criteria are processed,
// maintaining the integrity and security of the authorization flow.
func (u *usecase) Execute(ctx context.Context, req dto.AuthorizeRequest) (string, error) {
	const op = "processing of authorization request"
	const logTag = "[pxr-sso-use-case-authorize]"

	log := u.log.With(
		sl.Strings("response_type", req.ResponseType()),
		sl.Strings("client_id", req.ClientID()),
		sl.Strings("redirect_uri", req.RedirectURI()),
		sl.Strings("code_challenge", req.CodeChallenge()),
		sl.Strings("code_challenge_method", req.CodeChallengeMethod()),
		sl.Strings("state", req.State()),
		sl.Strings("audience", req.Audience()),
		sl.Strings("scope", req.Scope()),
		slog.Any("session", req.Sessions()),
	)
	log.Info(logTag + " attempting to initiate user authorization")

	nullableClient := nullable.None[dto.Client]()
	if len(req.ClientID()) == 1 {
		client, err := u.clientReader.ClientByCode(
			ctx,
			req.ClientID()[0],
			repository.WithRedirectURIs(),
			repository.WithAudiences(),
			repository.WithScopes())

		if err != nil && !errors.Is(err, infrastructure.ErrEntityNotFound) {
			log.Error(logTag+" get client by code", sl.Err(err))

			errorRedirectURI, buildErr := u.uriBuilder.BuildErrorRedirectURI("", err)
			if buildErr != nil {
				log.Error(logTag+" build error redirect URI", sl.Err(err))

				return "", fmt.Errorf("%s: %w", op, buildErr)
			}

			return errorRedirectURI, nil
		}

		nullableClient = nullable.Some(client)
	}

	sessionCodes := make([]string, len(req.Sessions()))
	for i, session := range req.Sessions() {
		decodedSessionCookie, err := u.sessionCookieDecoder.Decode(session.Value())
		if err != nil {
			log.Error(logTag+" decode session", sl.Err(err))

			return "", fmt.Errorf("%s: %w", op, err)
		}

		sessionCodes[i] = decodedSessionCookie.ID()
	}

	sessions, err := u.sessionReader.SessionsByCode(ctx, sessionCodes, repository.WithClient(), repository.WithUser())
	if err != nil && !errors.Is(err, infrastructure.ErrEntityNotFound) {
		log.Error(logTag+" get session by code", sl.Err(err))

		errorRedirectURI, buildErr := u.uriBuilder.BuildErrorRedirectURI("", err)
		if buildErr != nil {
			log.Error(logTag+" build error redirect URI", sl.Err(err))

			return "", fmt.Errorf("%s: %w", op, buildErr)
		}

		return errorRedirectURI, nil
	}

	authorizationRequestValidator := validator.NewAuthorizationRequestValidator(req, nullableClient, sessions)
	validatedAuthorizationRequest, err := authorizationRequestValidator.Validate()
	if err != nil {
		log.Warn(logTag+" validate authorization request", sl.Err(err))

		clientRedirectURI := validatedAuthorizationRequest.RedirectURI()
		errorRedirectURI, buildErr := u.uriBuilder.BuildErrorRedirectURI(clientRedirectURI, err)
		if buildErr != nil {
			log.Error(logTag+" build error redirect URI", sl.Err(err))

			return "", fmt.Errorf("%s: %w", op, buildErr)
		}

		return errorRedirectURI, nil
	}

	authorizeCtx := dto.NewAuthorizeContext(validatedAuthorizationRequest, nullableClient.Unwrap(), sessions)
	redirectURI, err := u.switcher.Switch(ctx, authorizeCtx)
	if err != nil {
		log.Error(logTag+" process authorize flow", sl.Err(err))

		errorRedirectURI, buildErr := u.uriBuilder.BuildErrorRedirectURI("", err)
		if buildErr != nil {
			log.Error(logTag+" build error redirect URI", sl.Err(err))

			return "", fmt.Errorf("%s: %w", op, buildErr)
		}

		return errorRedirectURI, nil
	}

	log.Info(logTag + " initiate user authorization successfully")

	return redirectURI, nil
}
