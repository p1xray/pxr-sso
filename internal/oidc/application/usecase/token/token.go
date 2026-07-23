package token

import (
	"context"
	"fmt"
	"github.com/p1xray/pxr-sso/internal/oidc/application/validator"
	"github.com/p1xray/pxr-sso/internal/oidc/domain/dto"
	"github.com/p1xray/pxr-sso/internal/oidc/infrastructure/repository"
	"log/slog"
)

type AuthorizedGrantReader interface {
	AuthorizedGrant(ctx context.Context, authorizationCode string) (dto.AuthorizedGrant, error)
}

type AuthorizedGrantRemover interface {
	RemoveAuthorizedGrant(ctx context.Context, authorizationCode string) error
}

type UserReader interface {
	User(ctx context.Context, id int64, opts ...repository.UserOption) (dto.User, error)
}

type ClientReader interface {
	ClientByCode(ctx context.Context, code string, opts ...repository.ClientOption) (dto.Client, error)
}

type TokensGenerator interface {
	Generate(
		scope []string,
		audience string,
		user dto.User,
		client dto.Client,
		session dto.AuthorizedSession,
	) (dto.Tokens, error)
}

// RequestProcessor processes incoming token requests from clients, ensuring they
// are valid and authorized before issuing the appropriate token response.
//
// Depending on the request type and granted permissions, the response can
// include various types of tokens such as Access Tokens, Refresh Tokens and ID
// Tokens.
//
// This interface abstracts the core logic behind token issuance in compliance
// with OAuth 2.0 and OpenID Connect standards. Implementations are responsible
// for validating the token request details, determining the types of tokens to
// issue based on the request's scope and authorization, and generating a token
// response that conforms to the protocol specifications. While the typical
// response includes an Access Token and, in the case of OpenID Connect, an ID
// Token, the exact contents of the response may vary based on the request
// parameters and server policies.
type RequestProcessor interface {
	Execute(ctx context.Context, request dto.TokenRequest) (dto.TokenResponse, error)
}

// Processes token requests in compliance with OAuth 2.0 and OpenID Connect
// standards, handling various types of token requests such as authorization code
// and refresh token.
//
// Generates the appropriate token responses including access tokens, refresh
// tokens, and ID tokens.
type usecase struct {
	log                    *slog.Logger
	authorizedGrantReader  AuthorizedGrantReader
	authorizedGrantRemover AuthorizedGrantRemover
	userReader             UserReader
	clientReader           ClientReader
	tokensGenerator        TokensGenerator
}

// New returns a new instance of the token requests processing usecase with required dependencies.
func New(
	log *slog.Logger,
	authorizedGrantReader AuthorizedGrantReader,
	authorizedGrantRemover AuthorizedGrantRemover,
	userReader UserReader,
	clientReader ClientReader,
	tokensGenerator TokensGenerator,
) *usecase {
	return &usecase{
		log:                    log,
		authorizedGrantReader:  authorizedGrantReader,
		authorizedGrantRemover: authorizedGrantRemover,
		userReader:             userReader,
		clientReader:           clientReader,
		tokensGenerator:        tokensGenerator,
	}
}

// Execute processes a token request, determining the necessary tokens to
// generate based on the request's scope and grant type. It generates an access
// token for every request and, depending on the scope, may also generate a
// refresh token and an ID token for OpenID Connect authentication.
func (u *usecase) Execute(ctx context.Context, request dto.TokenRequest) (dto.TokenResponse, error) {
	const op = "processing of token issuance"
	const logTag = "[pxr-sso-use-case-token]"

	log := u.log.With(
		slog.String("grant_type", request.GrantType()),
		slog.String("authorization_code", request.AuthorizationCode()),
		slog.String("redirect_uri", request.RedirectURI()),
		slog.String("code_verifier", request.CodeVerifier()),
		slog.String("client_id", request.ClientID()),
	)
	log.Info(logTag + " attempting to token issuance")

	authorizedGrant, err := u.authorizedGrantReader.AuthorizedGrant(ctx, request.AuthorizationCode())
	if err != nil {
		log.Error(logTag+" get authorized grant", err.Error())

		return dto.TokenResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	clientID := authorizedGrant.RequestClientID()
	client, err := u.clientReader.ClientByCode(ctx, clientID, repository.WithRedirectURIs())
	if err != nil {
		log.Error(logTag+" get client by code", err.Error())

		return dto.TokenResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	userID := authorizedGrant.SessionUserID()
	user, err := u.userReader.User(ctx, userID)
	if err != nil {
		log.Error(logTag+" get user by id", err.Error())

		return dto.TokenResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	tokenRequestValidator := validator.NewTokenRequestValidator(request, authorizedGrant.Request(), client)
	err = tokenRequestValidator.Validate()
	if err != nil {
		log.Warn(logTag+" validate token request", err.Error())

		return dto.TokenResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	tokens, err := u.tokensGenerator.Generate(
		authorizedGrant.RequestGrantedScopes(),
		authorizedGrant.RequestAudience(),
		user,
		client,
		authorizedGrant.Session(),
	)
	if err != nil {
		log.Error(logTag+" generate tokens", err.Error())

		return dto.TokenResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info(logTag + " issuance token successfully")

	tokenResponse := dto.NewTokenResponse(
		tokens.AccessTokenString(),
		tokens.AccessTokenType(),
		tokens.RefreshTokenString(),
		tokens.IDTokenString(),
		tokens.AccessTokenExpiresIn(),
	)

	return tokenResponse, nil
}
