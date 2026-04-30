package generator

import "encoding/base64"

const authorizationCodeLength = 64

// AuthorizationCode generates a unique, cryptographically secure authorization code.
func AuthorizationCode() string {
	randomSymbols := randomBytes(authorizationCodeLength)
	code := base64.RawURLEncoding.EncodeToString(randomSymbols)

	return code
}
