package register

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

type URIBuilder interface {
	BuildConsentRedirectURI(requestURI string) (string, error)
	BuildCallbackRedirectURI(rawURL, authorizationCode, state string) (string, error)
}

type UserReader interface {
	IsUserExistByUsername(ctx context.Context, username string) (bool, error)
}

type UserSaver interface {
	SaveUser(ctx context.Context, user entity.User) (dto.User, error)
}

type ClientReader interface {
	ClientByCode(ctx context.Context, code string, opts ...repository.ClientOption) (dto.Client, error)
}

type SessionSaver interface {
	SaveSession(ctx context.Context, session entity.Session) (dto.Session, error)
}

type AuthorizationRequestReader interface {
	AuthorizationRequest(ctx context.Context, requestURI string) (dto.ValidatedAuthorizeRequest, error)
}

type SessionCookieEncoder interface {
	Encode(session dto.AuthorizedSession) (string, error)
}

type Register interface {
	Execute(ctx context.Context, registerRequest dto.RegisterRequest) (dto.RegisterResponse, error)
}

type usecase struct {
	log                        *slog.Logger
	uriBuilder                 URIBuilder
	userReader                 UserReader
	userSaver                  UserSaver
	clientReader               ClientReader
	sessionSaver               SessionSaver
	authorizationRequestReader AuthorizationRequestReader
	sessionCookieEncoder       SessionCookieEncoder
}

func New(
	log *slog.Logger,
	uriBuilder URIBuilder,
	userReader UserReader,
	userSaver UserSaver,
	clientReader ClientReader,
	sessionSaver SessionSaver,
	authorizationRequestReader AuthorizationRequestReader,
	sessionCookieEncoder SessionCookieEncoder,
) *usecase {
	return &usecase{
		log:                        log,
		uriBuilder:                 uriBuilder,
		userReader:                 userReader,
		userSaver:                  userSaver,
		clientReader:               clientReader,
		sessionSaver:               sessionSaver,
		authorizationRequestReader: authorizationRequestReader,
		sessionCookieEncoder:       sessionCookieEncoder,
	}
}

// Execute executes the use-case for registering a new user.
func (u *usecase) Execute(ctx context.Context, registerRequest dto.RegisterRequest) (dto.RegisterResponse, error) {
	const op = "processing of register a user"
	const logTag = "[pxr-sso-use-case-register]"

	log := u.log.With(
		slog.String("request_uri", registerRequest.RequestURI()),
		slog.String("username", registerRequest.Username()),
	)
	log.Info(logTag + " attempting to register new user")

	if err := registerRequest.Validate(); err != nil {
		log.Warn(logTag+" validate register request", err.Error())

		return dto.RegisterResponse{}, oidc.InvalidRequestError(err)
	}

	authorizationRequest, err := u.authorizationRequestReader.AuthorizationRequest(ctx, registerRequest.RequestURI())
	if err != nil && !errors.Is(err, infrastructure.ErrEntityNotFound) {
		log.Error(logTag+" get authorization request", sl.Err(err))

		return dto.RegisterResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	isUserExist, err := u.userReader.IsUserExistByUsername(ctx, registerRequest.Username())
	if err != nil {
		log.Error(logTag+" get user by username", sl.Err(err))

		return dto.RegisterResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	if isUserExist {
		return dto.RegisterResponse{}, oidc.ErrOAuthUserAlreadyExists
	}

	client, err := u.clientReader.ClientByCode(ctx, authorizationRequest.ClientID(), repository.WithDefaultRoles())
	if err != nil && !errors.Is(err, infrastructure.ErrEntityNotFound) {
		log.Error(logTag+" get client by code", sl.Err(err))

		return dto.RegisterResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	user, err := entity.NewUser(
		registerRequest.Username(),
		registerRequest.Password(),
		registerRequest.FullName(),
		client.ID(),
		client.DefaultRoles(),
	)
	if err != nil {
		log.Error(logTag+" create user", sl.Err(err))

		return dto.RegisterResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	savedUser, err := u.userSaver.SaveUser(ctx, user)
	if err != nil {
		log.Error(logTag+" save user", sl.Err(err))

		return dto.RegisterResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	session, err := entity.NewSession(client, savedUser)
	if err != nil {
		log.Warn(logTag+" create session", sl.Err(err))

		return dto.RegisterResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	savedSession, err := u.sessionSaver.SaveSession(ctx, session)
	if err != nil {
		log.Warn(logTag+" save session", sl.Err(err))

		return dto.RegisterResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	// TODO: get last user session and check granted (from session) and pending (from request) scopes.
	// 	If scopes are equals, then set scopes to actual session and redirect to callback (execute grant use case),
	// 	otherwise redirect to consent page.

	consentRedirectURI, err := u.uriBuilder.BuildConsentRedirectURI(registerRequest.RequestURI())
	if err != nil {
		log.Error(logTag+" build consent redirect uri", sl.Err(err))

		return dto.RegisterResponse{}, fmt.Errorf("%s: %w", op, err)
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

		return dto.RegisterResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info(logTag + " user registered successfully")

	sessionCookieName := generator.SessionCookieName(client.Code())
	sessionCookie := dto.NewSessionCookie(sessionCookieName, encodedSession)
	return dto.NewRegisterResponse(consentRedirectURI, sessionCookie), nil
}
