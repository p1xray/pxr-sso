package entity

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/p1xray/pxr-sso/internal/oauth/domain"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/builder"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/dto"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/generator"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/validator/authorize"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/validator/consent"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/validator/login"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/validator/register"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/validator/token"
	jwtclaims "github.com/p1xray/pxr-sso/pkg/jwt/claims"
	jwtcreator "github.com/p1xray/pxr-sso/pkg/jwt/creator"
	"github.com/p1xray/pxr-sso/pkg/nullable"
	"golang.org/x/crypto/bcrypt"
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
	err := o.createFlow(validatedData)
	if err != nil {
		return err
	}

	// build redirect URI to login page
	if err = o.createLoginRedirectURI(); err != nil {
		return err
	}

	return nil
}

func (o *OAuth) Login(data dto.Login) *domain.DisplayableError {
	// validate request parameters
	validator := login.NewValidator(data, o.client, o.flow, o.user)
	if err := validator.Validate(); err != nil {
		return err
	}

	if err := o.setFlowUsername(); err != nil {
		return domain.InternalError(err)
	}

	// generate authorization code
	if err := o.setFlowAuthorizationCode(); err != nil {
		return domain.InternalError(err)
	}

	// build callback redirect URI
	if err := o.createConsentRedirectURI(); err != nil {
		return domain.InternalError(err)
	}

	return nil
}

func (o *OAuth) Register(data dto.Register) *domain.DisplayableError {
	// validate request parameters
	validator := register.NewValidator(data, o.client, o.flow, o.user)
	if err := validator.Validate(); err != nil {
		return err
	}

	// create new user
	if err := o.createNewUser(data); err != nil {
		return domain.InternalError(err)
	}

	if err := o.setFlowUsername(); err != nil {
		return domain.InternalError(err)
	}

	// generate authorization code
	if err := o.setFlowAuthorizationCode(); err != nil {
		return domain.InternalError(err)
	}

	// build callback redirect URI
	if err := o.createConsentRedirectURI(); err != nil {
		return domain.InternalError(err)
	}

	return nil
}

func (o *OAuth) Consent(data dto.Consent) error {
	// validate request parameters
	validator := consent.NewValidator(data, o.client, o.flow)
	if err := validator.Validate(); err != nil {
		return err
	}

	// generate authorization code
	if err := o.setFlowAuthorizationCode(); err != nil {
		return err
	}

	// build callback redirect URI
	if err := o.createCallbackRedirectURI(); err != nil {
		return err
	}

	return nil
}

func (o *OAuth) ExchangeToken(data dto.ExchangeToken) (dto.Token, error) {
	// validate request parameters
	validator := token.NewValidator(data, o.client, o.flow)
	if err := validator.Validate(); err != nil {
		return dto.Token{}, err
	}

	// generate tokens
	scope := validator.ValidatedScope()
	tokens, err := o.generateTokens(scope)
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

func (o *OAuth) generateFlowID() (uuid.UUID, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("%s: %w", "generate flow ID", err)
	}

	return id, nil
}

func (o *OAuth) generateTokens(scope []string) (dto.Token, error) {
	const op = "generate tokens"

	accessTokenClaims, accessToken, err := o.generateAccessToken(scope)
	if err != nil {
		return dto.Token{}, fmt.Errorf("%s: %w", op, err)
	}

	_, refreshToken, err := o.generateRefreshToken()
	if err != nil {
		return dto.Token{}, fmt.Errorf("%s: %w", op, err)
	}

	_, idToken, err := o.generateIDToken(scope)
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

func (o *OAuth) generateAccessToken(scope []string) (jwtclaims.AccessTokenClaims, string, error) {
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
		// TODO: add to scope user permissions
		Scopes: scope,
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

func (o *OAuth) generateIDToken(scope []string) (jwtclaims.RefreshTokenClaims, string, error) {
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

func (o *OAuth) createLoginRedirectURI() error {
	flow, err := o.Flow()
	if err != nil {
		oauthErr := domain.ServerErrorOAuthError(err)
		o.HandleError(oauthErr, "")

		return fmt.Errorf("%s: %w", "create login redirect URI", err)
	}

	loginRedirectURI := o.uriBuilder.BuildLoginRedirectURI(flow)
	o.setRedirectURI(loginRedirectURI)

	return nil
}

func (o *OAuth) createConsentRedirectURI() error {
	flow, err := o.Flow()
	if err != nil {
		return fmt.Errorf("%s: %w", "create consent redirect URI", err)
	}

	consentRedirectURI := o.uriBuilder.BuildConsentRedirectURI(flow)
	o.setRedirectURI(consentRedirectURI)

	return nil
}

func (o *OAuth) createCallbackRedirectURI() error {
	flow, err := o.Flow()
	if err != nil {
		return fmt.Errorf("%s: %w", "create callback redirect URI", err)
	}

	callbackRedirectURI := o.uriBuilder.BuildCallbackRedirectURI(flow)
	o.setRedirectURI(callbackRedirectURI)

	return nil
}

func (o *OAuth) setRedirectURI(redirectURI string) {
	o.redirectURI = redirectURI
}

func (o *OAuth) RedirectURI() string {
	return o.redirectURI
}

func (o *OAuth) createFlow(validatedData dto.ValidatedAuthorize) error {
	id, err := o.generateFlowID()
	if err != nil {
		oauthErr := domain.ServerErrorOAuthError(err)
		o.HandleError(oauthErr, validatedData.RedirectURI())

		return fmt.Errorf("%s: %w", "create flow", err)
	}

	flow := dto.NewFlow(
		id,
		validatedData.ClientID(),
		validatedData.ResponseType(),
		validatedData.RedirectURI(),
		validatedData.CodeChallenge(),
		validatedData.CodeChallengeMethod(),
		validatedData.State(),
		validatedData.Scope(),
	)
	o.setFlow(flow)

	return nil
}

func (o *OAuth) updateFlow(setters ...dto.FlowOption) error {
	flow, err := o.Flow()
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
		flow.Scope(),
		setters...,
	)
	o.setFlow(updatedFlow)

	return nil
}

func (o *OAuth) setFlowAuthorizationCode() error {
	code := generator.NewAuthorizationCode()
	if err := o.updateFlow(dto.WithAuthorizationCode(code)); err != nil {
		return fmt.Errorf("set flow authorization code: %w", err)
	}

	return nil
}

func (o *OAuth) setFlowUsername() error {
	user, err := o.User()
	if err != nil {
		return fmt.Errorf("set flow username: %w", err)
	}

	if err = o.updateFlow(dto.WithUsername(user.Username())); err != nil {
		return fmt.Errorf("set flow username: %w", err)
	}

	return nil
}

func (o *OAuth) setFlow(flow dto.Flow) {
	o.flow = nullable.Some(flow)
}

func (o *OAuth) Flow() (dto.Flow, error) {
	if o.flow.IsNone() {
		return dto.Flow{}, fmt.Errorf("get flow: flow not initialized")
	}

	return o.flow.Unwrap(), nil
}

func (o *OAuth) createNewUser(data dto.Register) error {
	passwordHash, err := o.generatePasswordHash(data.Password())
	if err != nil {
		return fmt.Errorf("%s: %w", "create new user", err)
	}

	user := dto.NewRegisteringUser(data.Username(), passwordHash, data.FullName())
	o.setUser(user)

	return nil
}

func (o *OAuth) generatePasswordHash(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("%s: %w", "generate password hash", err)
	}

	return string(passwordHash), nil
}

func (o *OAuth) setUser(user dto.User) {
	o.user = nullable.Some(user)
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
