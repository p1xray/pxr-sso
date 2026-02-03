package dto

type Token struct {
	accessToken  string
	tokenType    string
	expiresIn    uint32
	refreshToken string
	idToken      string
}

func NewToken(
	accessToken,
	tokenType,
	refreshToken,
	idToken string,
	expiresIn uint32,
) Token {
	return Token{
		accessToken:  accessToken,
		tokenType:    tokenType,
		expiresIn:    expiresIn,
		refreshToken: refreshToken,
		idToken:      idToken,
	}
}

func (t *Token) AccessToken() string {
	return t.accessToken
}

func (t *Token) TokenType() string {
	return t.tokenType
}

func (t *Token) ExpiresIn() uint32 {
	return t.expiresIn
}

func (t *Token) RefreshToken() string {
	return t.refreshToken
}

func (t *Token) IDToken() string {
	return t.idToken
}
