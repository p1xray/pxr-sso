package consent

import (
	"context"
	"errors"
	"fmt"
	"github.com/p1xray/pxr-sso/internal/oidc"
	"github.com/p1xray/pxr-sso/internal/oidc/application/validator"
	"github.com/p1xray/pxr-sso/internal/oidc/domain/dto"
	"github.com/p1xray/pxr-sso/internal/oidc/infrastructure"
	"github.com/p1xray/pxr-sso/internal/oidc/infrastructure/repository"
	"github.com/p1xray/pxr-sso/pkg/logger/sl"
	"log/slog"
)

type UserReader interface {
	User(ctx context.Context, id int64, opts ...repository.UserOption) (dto.User, error)
}

type SessionReader interface {
	SessionByCode(ctx context.Context, code string, opts ...repository.SessionOption) (dto.Session, error)
}

type AuthorizationRequestReader interface {
	AuthorizationRequest(ctx context.Context, requestURI string) (dto.ValidatedAuthorizeRequest, error)
}

type SessionCookieDecoder interface {
	Decode(encoded string) (dto.AuthorizedSession, error)
}

type AuthorizedGrantProcessor interface {
	Execute(ctx context.Context, data dto.AuthorizeContext) (string, error)
}

type Consent interface {
	Execute(ctx context.Context, consentRequest dto.ConsentRequest) (string, error)
}

type usecase struct {
	log                        *slog.Logger
	userReader                 UserReader
	sessionReader              SessionReader
	authorizationRequestReader AuthorizationRequestReader
	sessionCookieDecoder       SessionCookieDecoder
	authorizedGrantProcessor   AuthorizedGrantProcessor
}

func New(
	log *slog.Logger,
	userReader UserReader,
	sessionReader SessionReader,
	authorizationRequestReader AuthorizationRequestReader,
	sessionCookieDecoder SessionCookieDecoder,
	authorizedGrantProcessor AuthorizedGrantProcessor,
) *usecase {
	return &usecase{
		log:                        log,
		userReader:                 userReader,
		sessionReader:              sessionReader,
		authorizationRequestReader: authorizationRequestReader,
		sessionCookieDecoder:       sessionCookieDecoder,
		authorizedGrantProcessor:   authorizedGrantProcessor,
	}
}

// Execute ...
func (u *usecase) Execute(ctx context.Context, consentRequest dto.ConsentRequest) (string, error) {
	const op = "processing of confirm consent"
	const logTag = "[pxr-sso-use-case-consent]"

	log := u.log.With(
		slog.String("request_uri", consentRequest.RequestURI()),
		sl.Strings("scopes", consentRequest.Scopes()),
	)
	log.Info(logTag + " attempting to confirm consent")

	if err := consentRequest.Validate(); err != nil {
		log.Warn(logTag+" validate consent request", err.Error())

		return "", oidc.InvalidRequestError(err)
	}

	authorizationRequest, err := u.authorizationRequestReader.AuthorizationRequest(ctx, consentRequest.RequestURI())
	if err != nil && !errors.Is(err, infrastructure.ErrEntityNotFound) {
		log.Error(logTag+" get authorization request", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	sessionCookie, err := consentRequest.SingleSessionCookie()
	if err != nil {
		log.Error(logTag+" get single session cookie", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	decodedSessionCookie, err := u.sessionCookieDecoder.Decode(sessionCookie.Value())
	if err != nil {
		log.Error(logTag+" decode session", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	session, err := u.sessionReader.SessionByCode(
		ctx,
		decodedSessionCookie.ID(),
		repository.WithClient(),
		repository.WithUser(),
	)
	if err != nil {
		log.Error(logTag+" get session by code", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	consentRequestValidator := validator.NewConsentRequestValidator(consentRequest, authorizationRequest)
	validatedRequest := consentRequestValidator.Validate()

	authorizeCtx := dto.NewAuthorizeContext(validatedRequest, session.Client(), []dto.Session{session})
	redirectURI, err := u.authorizedGrantProcessor.Execute(ctx, authorizeCtx)
	if err != nil {
		log.Error(logTag+" process authorize grant", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info(logTag + " confirmed consent successfully")

	return redirectURI, nil
}
