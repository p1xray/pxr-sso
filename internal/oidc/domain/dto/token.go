package dto

type TokenRequest struct {
	grantType         string
	authorizationCode string
	redirectURI       string
	codeVerifier      string
	clientID          string
	clientSecret      string
}

func NewTokenRequest(
	grantType,
	authorizationCode,
	redirectURI,
	codeVerifier,
	clientID,
	clientSecret string,
) TokenRequest {
	return TokenRequest{
		grantType:         grantType,
		authorizationCode: authorizationCode,
		redirectURI:       redirectURI,
		codeVerifier:      codeVerifier,
		clientID:          clientID,
		clientSecret:      clientSecret,
	}
}

func (t *TokenRequest) GrantType() string {
	return t.grantType
}

func (t *TokenRequest) AuthorizationCode() string {
	return t.authorizationCode
}

func (t *TokenRequest) RedirectURI() string {
	return t.redirectURI
}

func (t *TokenRequest) CodeVerifier() string {
	return t.codeVerifier
}

func (t *TokenRequest) ClientID() string {
	return t.clientID
}

func (t *TokenRequest) ClientSecret() string {
	return t.clientSecret
}

type TokenResponse struct {
	accessToken  string
	tokenType    string
	expiresIn    int64
	refreshToken string
	idToken      string
}

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

func (t *TokenResponse) AccessToken() string {
	return t.accessToken
}

func (t *TokenResponse) TokenType() string {
	return t.tokenType
}

func (t *TokenResponse) ExpiresIn() int64 {
	return t.expiresIn
}

func (t *TokenResponse) RefreshToken() string {
	return t.refreshToken
}

func (t *TokenResponse) IDToken() string {
	return t.idToken
}

type Tokens struct {
	accessToken  Token
	refreshToken Token
	idToken      Token
}

func NewTokens(accessToken, refreshToken, idToken Token) Tokens {
	return Tokens{
		accessToken:  accessToken,
		refreshToken: refreshToken,
		idToken:      idToken,
	}
}

func (t *Tokens) AccessTokenString() string {
	return t.accessToken.String()
}

func (t *Tokens) AccessTokenType() string {
	return t.accessToken.TokenType()
}

func (t *Tokens) AccessTokenExpiresIn() int64 {
	return t.accessToken.ExpiresIn()
}

func (t *Tokens) RefreshTokenString() string {
	return t.refreshToken.String()
}

func (t *Tokens) IDTokenString() string {
	return t.idToken.String()
}

type Token struct {
	id        string
	tokenType string
	jwtString string
	expiresIn int64
}

func NewToken(id, tokenType, jwtString string, expiresIn int64) Token {
	return Token{
		id:        id,
		tokenType: tokenType,
		jwtString: jwtString,
		expiresIn: expiresIn,
	}
}

func (t *Token) ID() string {
	return t.id
}

func (t *Token) TokenType() string {
	return t.tokenType
}

func (t *Token) String() string {
	return t.jwtString
}

func (t *Token) ExpiresIn() int64 {
	return t.expiresIn
}
