package login

import (
	"context"
	"errors"
	"fmt"
	"github.com/p1xray/pxr-sso/internal/oidc"
	"github.com/p1xray/pxr-sso/internal/oidc/application/generator"
	"github.com/p1xray/pxr-sso/internal/oidc/domain/dto"
	"github.com/p1xray/pxr-sso/internal/oidc/domain/entity"
	"github.com/p1xray/pxr-sso/internal/oidc/infrastructure"
	"github.com/p1xray/pxr-sso/internal/oidc/infrastructure/repository"
	"github.com/p1xray/pxr-sso/pkg/logger/sl"
	"log/slog"
)

// URIBuilder is the URI builder to redirect the user agent to a specified URI.
type URIBuilder interface {
	// BuildConsentRedirectURI return the URI to redirect to the confirming consent page.
	BuildConsentRedirectURI(requestURI string) (string, error)

	// BuildCallbackRedirectURI return the URI to redirect to the client's callback page.
	BuildCallbackRedirectURI(rawURL, authorizationCode, state string) (string, error)
}

// UserReader is the user data reader from storage.
type UserReader interface {
	UserByUsername(ctx context.Context, username string, opts ...repository.UserOption) (dto.User, error)
}

// ClientReader is the client data reader from storage.
type ClientReader interface {
	ClientByCode(ctx context.Context, code string, opts ...repository.ClientOption) (dto.Client, error)
}

// SessionSaver is the saver session data to the storage.
type SessionSaver interface {
	// SaveSession saves the session data to the storage.
	SaveSession(ctx context.Context, session entity.Session) (dto.Session, error)
}

// AuthorizationRequestReader is the authorization request data reader from storage.
type AuthorizationRequestReader interface {
	// AuthorizationRequest returns the authorization request data reader from storage.
	AuthorizationRequest(ctx context.Context, requestURI string) (dto.ValidatedAuthorizeRequest, error)
}

// SessionCookieEncoder is the session cookie data encoder.
type SessionCookieEncoder interface {
	// Encode encodes the session cookie data.
	Encode(session dto.AuthorizedSession) (string, error)
}

// Login is the handler for logging in user request.
type Login interface {
	// Execute processes a logging in user request.
	Execute(ctx context.Context, loginRequest dto.LoginRequest) (dto.LoginResponse, error)
}

// usecase is the use case which handles the processing of logging in user request.
type usecase struct {
	log                        *slog.Logger
	uriBuilder                 URIBuilder
	userReader                 UserReader
	clientReader               ClientReader
	sessionSaver               SessionSaver
	authorizationRequestReader AuthorizationRequestReader
	sessionCookieEncoder       SessionCookieEncoder
}

// New creates a new use case which handles the processing of logging in user request.
func New(
	log *slog.Logger,
	uriBuilder URIBuilder,
	userReader UserReader,
	clientReader ClientReader,
	sessionSaver SessionSaver,
	authorizationRequestReader AuthorizationRequestReader,
	sessionCookieEncoder SessionCookieEncoder,
) *usecase {
	return &usecase{
		log:                        log,
		uriBuilder:                 uriBuilder,
		userReader:                 userReader,
		clientReader:               clientReader,
		sessionSaver:               sessionSaver,
		authorizationRequestReader: authorizationRequestReader,
		sessionCookieEncoder:       sessionCookieEncoder,
	}
}

// Execute processes a logging in user request.
//
// This method validating request parameters and check user credentials.
// If successful, creates a new session and redirects the user agent to the consent confirmation page
// or the client callback page. Otherwise, returns an error.
func (u *usecase) Execute(ctx context.Context, loginRequest dto.LoginRequest) (dto.LoginResponse, error) {
	const op = "processing of log in a user"
	const logTag = "[pxr-sso-use-case-login]"

	log := u.log.With(
		slog.String("request_uri", loginRequest.RequestURI()),
		slog.String("username", loginRequest.Username()),
	)
	log.Info(logTag + " attempting to log in a user")

	if err := loginRequest.Validate(); err != nil {
		log.Warn(logTag+" validate login request", err.Error())

		return dto.LoginResponse{}, oidc.InvalidRequestError(err)
	}

	authorizationRequest, err := u.authorizationRequestReader.AuthorizationRequest(ctx, loginRequest.RequestURI())
	if err != nil && !errors.Is(err, infrastructure.ErrEntityNotFound) {
		log.Error(logTag+" get authorization request", sl.Err(err))

		return dto.LoginResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	user, err := u.userReader.UserByUsername(ctx, loginRequest.Username(), repository.WithRoles())
	if err != nil && !errors.Is(err, infrastructure.ErrEntityNotFound) {
		log.Error(logTag+" get user by username", sl.Err(err))

		return dto.LoginResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	client, err := u.clientReader.ClientByCode(ctx, authorizationRequest.ClientID())
	if err != nil && !errors.Is(err, infrastructure.ErrEntityNotFound) {
		log.Error(logTag+" get client by code", sl.Err(err))

		return dto.LoginResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	userEntity, err := entity.NewExistUser(user)
	if err != nil {
		log.Warn(logTag+" create user entity", sl.Err(err))

		return dto.LoginResponse{}, oidc.InvalidUserCredentialsError()
	}

	if err = userEntity.CheckCredentials(loginRequest.Username(), loginRequest.Password()); err != nil {
		log.Warn(logTag+" check user credentials", sl.Err(err))

		return dto.LoginResponse{}, oidc.InvalidUserCredentialsError()
	}

	session, err := entity.NewSession(client, user)
	if err != nil {
		log.Warn(logTag+" create session", sl.Err(err))

		return dto.LoginResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	savedSession, err := u.sessionSaver.SaveSession(ctx, session)
	if err != nil {
		log.Warn(logTag+" save session", sl.Err(err))

		return dto.LoginResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	// TODO: get last user session and check granted (from session) and pending (from request) scopes.
	// 	If scopes are equals, then set scopes to actual session and redirect to callback (execute grant use case),
	// 	otherwise redirect to consent page.

	consentRedirectURI, err := u.uriBuilder.BuildConsentRedirectURI(loginRequest.RequestURI())
	if err != nil {
		log.Error(logTag+" build consent redirect uri", sl.Err(err))

		return dto.LoginResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	sessionToEncode := dto.NewAuthorizedSession(
		savedSession.Code().String(),
		user.ID(),
		savedSession.AuthTime(),
		"pxr.soo",
	)
	encodedSession, err := u.sessionCookieEncoder.Encode(sessionToEncode)
	if err != nil {
		log.Error(logTag+" encode session", sl.Err(err))

		return dto.LoginResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Debug(logTag + " user logged in successfully")

	sessionCookieName := generator.SessionCookieName(client.Code())
	sessionCookie := dto.NewSessionCookie(sessionCookieName, encodedSession)
	return dto.NewLoginResponse(consentRedirectURI, sessionCookie), nil
}
