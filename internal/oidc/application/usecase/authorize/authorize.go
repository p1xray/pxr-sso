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

type ErrorURIBuilder interface {
	BuildErrorRedirectURI(rawURL string, err error) (string, error)
}

type FlowSwitcher interface {
	Switch(ctx context.Context, data dto.AuthorizeContext) (string, error)
}

type Repository interface {
	ClientByCode(ctx context.Context, code string, opts ...repository.ClientOption) (dto.Client, error)
	SessionsByCode(ctx context.Context, codes []string) ([]dto.Session, error)
}

type Authorize interface {
	Execute(ctx context.Context, req dto.AuthorizeRequest) (string, error)
}

// authorize is the use case which handles the processing of authorization requests by validating and
// then processing these requests based on defined business logic. It also includes
// the fetching of authorization requests when necessary.
type authorize struct {
	log        *slog.Logger
	repo       Repository
	uriBuilder ErrorURIBuilder
	switcher   FlowSwitcher
}

// NewUseCase returns new use case for processing of authorization request.
func NewUseCase(
	log *slog.Logger,
	repo Repository,
	uriBuilder ErrorURIBuilder,
	switcher FlowSwitcher,
) *authorize {
	return &authorize{
		log:        log,
		repo:       repo,
		uriBuilder: uriBuilder,
		switcher:   switcher,
	}
}

// Execute executes the processing of authorization request.
//
// This method ensures that only requests meeting the necessary validation criteria are processed,
// maintaining the integrity and security of the authorization flow.
func (uc *authorize) Execute(ctx context.Context, req dto.AuthorizeRequest) (string, error) {
	const op = "processing of authorization request"
	const logTag = "[pxr-sso-use-case-authorize]"

	log := uc.log.With(
		sl.Strings("response_type", req.ResponseType()),
		sl.Strings("client_id", req.ClientID()),
		sl.Strings("redirect_uri", req.RedirectURI()),
		sl.Strings("code_challenge", req.CodeChallenge()),
		sl.Strings("code_challenge_method", req.CodeChallengeMethod()),
		sl.Strings("state", req.State()),
		sl.Strings("audience", req.Audience()),
		sl.Strings("scope", req.Scope()),
		slog.Any("session", req.Session()),
	)
	log.Info(logTag + " attempting to initiate user authorization")

	nullableClient := nullable.None[dto.Client]()
	if len(req.ClientID()) == 1 {
		client, err := uc.repo.ClientByCode(
			ctx,
			req.ClientID()[0],
			repository.WithRedirectURIs(),
			repository.WithAudiences(),
			repository.WithScopes())

		if err != nil && !errors.Is(err, infrastructure.ErrEntityNotFound) {
			log.Error(logTag+" get client by code", sl.Err(err))

			errorRedirectURI, buildErr := uc.uriBuilder.BuildErrorRedirectURI("", err)
			if buildErr != nil {
				log.Error(logTag+" build error redirect URI", sl.Err(err))

				return "", fmt.Errorf("%s: %w", op, buildErr)
			}

			return errorRedirectURI, nil
		}

		nullableClient = nullable.Some(client)
	}

	sessions, err := uc.repo.SessionsByCode(ctx, req.Session())
	if err != nil && !errors.Is(err, infrastructure.ErrEntityNotFound) {
		log.Error(logTag+" get session by code", sl.Err(err))

		errorRedirectURI, buildErr := uc.uriBuilder.BuildErrorRedirectURI("", err)
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

		errorRedirectURI, buildErr := uc.uriBuilder.BuildErrorRedirectURI("", err)
		if buildErr != nil {
			log.Error(logTag+" build error redirect URI", sl.Err(err))

			return "", fmt.Errorf("%s: %w", op, buildErr)
		}

		return errorRedirectURI, nil
	}

	authorizeCtx := dto.NewAuthorizeContext(validatedAuthorizationRequest, nullableClient.Unwrap(), sessions)
	redirectURI, err := uc.switcher.Switch(ctx, authorizeCtx)
	if err != nil {
		log.Error(logTag+" process authorize flow", sl.Err(err))

		errorRedirectURI, buildErr := uc.uriBuilder.BuildErrorRedirectURI("", err)
		if buildErr != nil {
			log.Error(logTag+" build error redirect URI", sl.Err(err))

			return "", fmt.Errorf("%s: %w", op, buildErr)
		}

		return errorRedirectURI, nil
	}

	log.Info(logTag + " initiate user authorization successfully")

	return redirectURI, nil
}
