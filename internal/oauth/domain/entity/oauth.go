package entity

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/p1xray/pxr-sso/internal/oauth/domain"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/builder"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/dto"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/generator"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/validator/authorize"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/validator/login"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/validator/token"
	jwtclaims "github.com/p1xray/pxr-sso/pkg/jwt/claims"
	jwtcreator "github.com/p1xray/pxr-sso/pkg/jwt/creator"
	"github.com/p1xray/pxr-sso/pkg/nullable"
	"time"
)

type OAuth struct {
	uriBuilder *builder.URI

	client nullable.Nullable[dto.Client]
	flow   nullable.Nullable[dto.Flow]
	user   nullable.Nullable[dto.User]

	err         *domain.OAuthError
	redirectURI string
}

func NewOAuth(setters ...OAuthOption) *OAuth {
	oauth := &OAuth{
		client: nullable.None[dto.Client](),
		flow:   nullable.None[dto.Flow](),
		user:   nullable.None[dto.User](),
	}

	for _, setter := range setters {
		setter(oauth)
	}

	return oauth
}

func (o *OAuth) Authorize(data dto.Authorize) error {
	// validate request parameters
	validator := authorize.NewValidator(data, o.client)
	if err := o.validateAuthorizeRequestParams(validator); err != nil {
		return err
	}

	validatedData := validator.ValidatedData()

	// create flow data
	flow, err := o.createFlow(validatedData)
	if err != nil {
		return err
	}

	o.setFlow(flow)

	// build redirect URI to login page
	loginRedirectURI := o.uriBuilder.BuildLoginRedirectURI(flow)
	o.setRedirectURI(loginRedirectURI)

	return nil
}

func (o *OAuth) Login(data dto.Login) *domain.DisplayableError {
	// validate request parameters
	validator := login.NewValidator(data, o.client, o.flow, o.user)
	if err := validator.Validate(); err != nil {
		return err
	}

	// generate authorization code
	authorizationCode := generator.NewAuthorizationCode()

	// update flow data
	if err := o.updateFlow(authorizationCode); err != nil {
		return domain.InternalError(err)
	}

	// build callback redirect URI
	callbackRedirectURI := o.uriBuilder.BuildCallbackRedirectURI(data.RedirectURI(), authorizationCode, data.State())
	o.setRedirectURI(callbackRedirectURI)

	return nil
}

func (o *OAuth) ExchangeToken(data dto.ExchangeToken) (dto.Token, error) {
	// validate request parameters
	validator := token.NewValidator(data, o.client, o.flow)
	if err := validator.Validate(); err != nil {
		return dto.Token{}, err
	}

	// generate tokens
	tokens, err := o.generateTokens()
	if err != nil {
		return dto.Token{}, err
	}

	return tokens, nil
}

func (o *OAuth) validateAuthorizeRequestParams(validator *authorize.Validator) error {
	if err := validator.Validate(); err != nil {
		validatedData := validator.ValidatedData()
		o.HandleError(err, validatedData.RedirectURI())

		return fmt.Errorf("%s: %w", "validate request parameters", err.Unwrap())
	}

	return nil
}

func (o *OAuth) createFlow(validatedData dto.ValidatedAuthorize) (dto.Flow, error) {
	id, err := o.generateFlowID()
	if err != nil {
		oauthErr := domain.ServerErrorOAuthError(err)
		o.HandleError(oauthErr, validatedData.RedirectURI())

		return dto.Flow{}, fmt.Errorf("%s: %w", "create flow", err)
	}

	flow := dto.NewFlow(
		id,
		validatedData.ClientID(),
		validatedData.ResponseType(),
		validatedData.RedirectURI(),
		validatedData.CodeChallenge(),
		validatedData.CodeChallengeMethod(),
		validatedData.State(),
	)

	return flow, nil
}

func (o *OAuth) generateFlowID() (uuid.UUID, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("%s: %w", "generate flow ID", err)
	}

	return id, nil
}

func (o *OAuth) generateTokens() (dto.Token, error) {
	const op = "generate tokens"

	accessTokenClaims, accessToken, err := o.generateAccessToken()
	if err != nil {
		return dto.Token{}, fmt.Errorf("%s: %w", op, err)
	}

	_, refreshToken, err := o.generateRefreshToken()
	if err != nil {
		return dto.Token{}, fmt.Errorf("%s: %w", op, err)
	}

	_, idToken, err := o.generateIDToken()
	if err != nil {
		return dto.Token{}, fmt.Errorf("%s: %w", op, err)
	}

	tokens := dto.NewToken(
		accessToken,
		accessTokenClaims.TokenType,
		refreshToken,
		idToken,
		jwtclaims.NumericDateToInt64(accessTokenClaims.Expiry),
	)

	return tokens, nil
}

func (o *OAuth) generateAccessToken() (jwtclaims.AccessTokenClaims, string, error) {
	const op = "access token"

	user, err := o.User()
	if err != nil {
		return jwtclaims.AccessTokenClaims{}, "", fmt.Errorf("%s: %w", op, err)
	}

	client, err := o.Client()
	if err != nil {
		return jwtclaims.AccessTokenClaims{}, "", fmt.Errorf("%s: %w", op, err)
	}

	createAccessTokenData := jwtcreator.AccessTokenCreateData{
		Subject: user.IDString(),
		// TODO: get this from request parameters and validate with client audiences
		Audiences: []string{"http://localhost:3000"},
		// TODO: get this from request parameters and validate with client scopes
		Scopes: []string{},
		// TODO: get this from from proto request
		Issuer: "http://localhost:6003",
		// TODO: get this from config
		TTL: 1 * time.Hour,
		Key: []byte(client.SecretKey),
	}
	claims, accessToken, err := jwtcreator.NewAccessToken(createAccessTokenData)
	if err != nil {
		return jwtclaims.AccessTokenClaims{}, "", fmt.Errorf("%s: %w", op, err)
	}

	return claims, accessToken, nil
}

func (o *OAuth) generateRefreshToken() (jwtclaims.RefreshTokenClaims, string, error) {
	const op = "refresh token"

	client, err := o.Client()
	if err != nil {
		return jwtclaims.RefreshTokenClaims{}, "", fmt.Errorf("%s: %w", op, err)
	}

	// TODO: get refresh token TTL from config
	claims, refreshToken, err := jwtcreator.NewRefreshToken([]byte(client.SecretKey), 24*time.Hour)
	if err != nil {
		return jwtclaims.RefreshTokenClaims{}, "", fmt.Errorf("%s: %w", op, err)
	}

	return claims, refreshToken, nil
}

func (o *OAuth) generateIDToken() (jwtclaims.RefreshTokenClaims, string, error) {
	const op = "id token"

	// TODO: implement this
	return jwtclaims.RefreshTokenClaims{}, "", nil
}

func (o *OAuth) HandleError(err *domain.OAuthError, redirectURI string) {
	o.setError(err)

	errorRedirectURI := o.uriBuilder.BuildErrorRedirectURI(redirectURI, o.err)
	o.setRedirectURI(errorRedirectURI)
}

func (o *OAuth) setError(err *domain.OAuthError) {
	o.err = err
}

func (o *OAuth) setRedirectURI(redirectURI string) {
	o.redirectURI = redirectURI
}

func (o *OAuth) RedirectURI() string {
	return o.redirectURI
}

func (o *OAuth) setFlow(flow dto.Flow) {
	o.flow = nullable.Some(flow)
}

func (o *OAuth) updateFlow(code string) error {
	flow, err := o.Flow()
	if err != nil {
		return fmt.Errorf("update flow: %w", err)
	}

	user, err := o.User()
	if err != nil {
		return fmt.Errorf("update flow: %w", err)
	}

	updatedFlow := dto.NewFlow(
		flow.ID(),
		flow.ClientID(),
		flow.ResponseType(),
		flow.RedirectURI(),
		flow.CodeChallenge(),
		flow.CodeChallengeMethod(),
		flow.State(),
		dto.WithAuthorizationCode(code),
		dto.WithUserID(user.ID()),
	)
	o.setFlow(updatedFlow)

	return nil
}

func (o *OAuth) Flow() (dto.Flow, error) {
	if o.flow.IsNone() {
		return dto.Flow{}, fmt.Errorf("get flow: flow not initialized")
	}

	return o.flow.Unwrap(), nil
}

func (o *OAuth) User() (dto.User, error) {
	if o.user.IsNone() {
		return dto.User{}, fmt.Errorf("get user: user not initialized")
	}

	return o.user.Unwrap(), nil
}

func (o *OAuth) Client() (dto.Client, error) {
	if o.client.IsNone() {
		return dto.Client{}, fmt.Errorf("get client: client not initialized")
	}

	return o.client.Unwrap(), nil
}
