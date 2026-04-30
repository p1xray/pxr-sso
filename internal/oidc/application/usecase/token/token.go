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
	Generate(scope []string, audience string, user dto.User, client dto.Client) (dto.Tokens, error)
}

type Token interface {
	Execute(ctx context.Context, tokenRequest dto.TokenRequest) (dto.TokenResponse, error)
}

type usecase struct {
	log                    *slog.Logger
	authorizedGrantReader  AuthorizedGrantReader
	authorizedGrantRemover AuthorizedGrantRemover
	userReader             UserReader
	clientReader           ClientReader
	tokensGenerator        TokensGenerator
}

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

// Execute ...
func (u *usecase) Execute(ctx context.Context, tokenRequest dto.TokenRequest) (dto.TokenResponse, error) {
	const op = "processing of token issuance"
	const logTag = "[pxr-sso-use-case-token]"

	log := u.log.With(
		slog.String("grant_type", tokenRequest.GrantType()),
		slog.String("authorization_code", tokenRequest.AuthorizationCode()),
		slog.String("redirect_uri", tokenRequest.RedirectURI()),
		slog.String("code_verifier", tokenRequest.CodeVerifier()),
		slog.String("client_id", tokenRequest.ClientID()),
	)
	log.Info(logTag + " attempting to token issuance")

	authorizedGrant, err := u.authorizedGrantReader.AuthorizedGrant(ctx, tokenRequest.AuthorizationCode())
	if err != nil {
		log.Error(logTag+" get authorized grant", err.Error())

		return dto.TokenResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	clientID := authorizedGrant.RequestClientID()
	client, err := u.clientReader.ClientByCode(ctx, clientID)
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

	tokenRequestValidator := validator.NewTokenRequestValidator(tokenRequest, authorizedGrant.Request(), client)
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
