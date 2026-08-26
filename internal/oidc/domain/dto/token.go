package dto

// TokenRequest represents a request to get various types of tokens (e.g., access
// token, refresh token) from the authorization server.
//
// This is part of the OAuth 2.0 and OpenID Connect token exchange flow, where
// clients can request tokens based on different grant types like
// 'authorization_code', 'refresh_token' and others.
type TokenRequest struct {
	grantType         string
	authorizationCode string
	redirectURI       string
	codeVerifier      string
	clientID          string
}

// NewTokenRequest returns a new TokenRequest with the specified parameters.
func NewTokenRequest(
	grantType,
	authorizationCode,
	redirectURI,
	codeVerifier,
	clientID string,
) TokenRequest {
	return TokenRequest{
		grantType:         grantType,
		authorizationCode: authorizationCode,
		redirectURI:       redirectURI,
		codeVerifier:      codeVerifier,
		clientID:          clientID,
	}
}

// GrantType returns the grant type of the token request, indicating the method
// being used to get the token. Common values include 'authorization_code',
// 'refresh_token', etc.
func (t *TokenRequest) GrantType() string {
	return t.grantType
}

// AuthorizationCode returns the authorization code received from the
// authorization server. This is used in the authorization code grant type to
// exchange for an access token.
func (t *TokenRequest) AuthorizationCode() string {
	return t.authorizationCode
}

// RedirectURI returns the redirect URI where the response will be sent. This
// must match the redirect URI registered with the authorization server during
// the initial request.
func (t *TokenRequest) RedirectURI() string {
	return t.redirectURI
}

// CodeVerifier returns the code verifier used in the PKCE (Proof Key for Code
// Exchange) process. Required for public clients using the authorization code
// grant type to enhance security.
func (t *TokenRequest) CodeVerifier() string {
	return t.codeVerifier
}

// ClientID returns the client ID of the requesting client.
func (t *TokenRequest) ClientID() string {
	return t.clientID
}

// TokenResponse represents a successful response from the token request
// processor, containing the issued access token and related information.
type TokenResponse struct {
	accessToken  string
	tokenType    string
	expiresIn    int64
	refreshToken string
	idToken      string
}

// NewTokenResponse returns a new TokenResponse instance with the provided token
// details.
func NewTokenResponse(
	accessToken,
	tokenType,
	refreshToken,
	idToken string,
	expiresIn int64,
) TokenResponse {
	return TokenResponse{
		accessToken:  accessToken,
		tokenType:    tokenType,
		expiresIn:    expiresIn,
		refreshToken: refreshToken,
		idToken:      idToken,
	}
}

// AccessToken returns the access token issued by the authorization server.
func (t *TokenResponse) AccessToken() string {
	return t.accessToken
}

// TokenType returns the type of the token issued.
func (t *TokenResponse) TokenType() string {
	return t.tokenType
}

// ExpiresIn returns the lifetime in seconds of the access token.
func (t *TokenResponse) ExpiresIn() int64 {
	return t.expiresIn
}

// RefreshToken returns the optional refresh token that can be used to obtain new
// access tokens.
func (t *TokenResponse) RefreshToken() string {
	return t.refreshToken
}

// IDToken returns the optional ID token that provides identity information about the user.
func (t *TokenResponse) IDToken() string {
	return t.idToken
}

// Tokens represents an issued tokens by the authorization server.
//
// Access tokens authorize clients for resource access; refresh tokens enable
// long-lived sessions by allowing new access tokens to be obtained without
// re-authentication; ID tokens provide identity information about the user,
// crucial for OpenID Connect authentication flows.
type Tokens struct {
	accessToken  Token
	refreshToken Token
	idToken      Token
}

// NewTokens returns a new Tokens with the given access, refresh, and ID tokens.
func NewTokens(accessToken, refreshToken, idToken Token) Tokens {
	return Tokens{
		accessToken:  accessToken,
		refreshToken: refreshToken,
		idToken:      idToken,
	}
}

// AccessTokenString returns the raw string value of the access token.
func (t *Tokens) AccessTokenString() string {
	return t.accessToken.String()
}

// AccessTokenType returns the type of the access token (e.g., "Bearer").
func (t *Tokens) AccessTokenType() string {
	return t.accessToken.TokenType()
}

// AccessTokenExpiresIn returns the lifetime of the access token in seconds.
func (t *Tokens) AccessTokenExpiresIn() int64 {
	return t.accessToken.ExpiresIn()
}

// RefreshTokenString returns the raw string value of the refresh token.
func (t *Tokens) RefreshTokenString() string {
	return t.refreshToken.String()
}

// IDTokenString returns the raw string value of the ID token.
func (t *Tokens) IDTokenString() string {
	return t.idToken.String()
}

// Token represents an issued token by the authorization server.
type Token struct {
	id        string
	tokenType string
	jwtString string
	expiresIn int64
}

// NewToken returns a new Token instance with the provided parameters.
func NewToken(id, tokenType, jwtString string, expiresIn int64) Token {
	return Token{
		id:        id,
		tokenType: tokenType,
		jwtString: jwtString,
		expiresIn: expiresIn,
	}
}

// ID returns the unique identifier of the token.
func (t *Token) ID() string {
	return t.id
}

// TokenType returns the type of the token (e.g., "Bearer").
func (t *Token) TokenType() string {
	return t.tokenType
}

// String returns the raw JWT string representation of the token.
func (t *Token) String() string {
	return t.jwtString
}

// ExpiresIn returns the lifetime of the token in seconds.
func (t *Token) ExpiresIn() int64 {
	return t.expiresIn
}
