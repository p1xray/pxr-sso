package claims

import (
	"github.com/go-jose/go-jose/v4/jwt"
	"time"
)

// AuthenticationClaims represents how and when the user signed in.
//
// Defined in OpenID Connect Core 1.0, Section 2, Section 3.1.3.6, Section 3.3.2.11.
type AuthenticationClaims struct {
	// AuthTime represents the time when the authentication occurred, facilitating
	// checks against token freshness and replay attacks.
	AuthTime *jwt.NumericDate `json:"auth_time,omitempty"`

	// Nonce represent the string value used to associate a client session
	// with an ID Token.
	//
	// Used in OpenID Connect to mitigate replay attacks by binding a session
	// to a token.
	Nonce string `json:"nonce,omitempty"`

	// ACR represents the Authentication Context Class Reference.
	//
	// Used in OpenID Connect. This claim specifies the authentication context class
	// that the authentication performed satisfied. It allows clients to request and
	// services to assert the strength of an authentication process.
	ACR string `json:"acr,omitempty"`

	// AMR represents the Authentication Methods References.
	//
	// Used in OpenID Connect. This claim is used to specify the authentication
	// methods used in the authentication process. It provides transparency about how
	// the authentication was performed, such as 'pwd' for password-based or 'mfa'
	// for multifactor authentication.
	AMR []string `json:"amr,omitempty"`

	// AZP represents the authorized party - the party to which the ID Token was issued.
	//
	// Used in OpenID Connect. This claim is particularly useful in delegated
	// scenarios to identify the party using the ID Token.
	AZP string `json:"azp,omitempty"`

	// AtHash represents the hash of the access token issued.
	//
	// Used in OpenID Connect. This claim is included in the ID Token and is a hash
	// of the access token, allowing the recipient to validate the integrity of the
	// access token. It is particularly useful in implicit and authorization code
	// flows.
	AtHash string `json:"at_hash,omitempty"`

	// CHash represents the Code Hash Value.
	//
	// Used in OpenID Connect. This claim is used in the ID Token to provide a hash
	// of the authorization code. It ensures that the authorization code is bound to
	// the ID Token, enhancing the security of the code exchange process.
	CHash string `json:"c_hash,omitempty"`
}

// NewAuthenticationClaims returns a new AuthenticationClaims struct using the
// required authentication timestamp and applies optional configurations.
func NewAuthenticationClaims(authTime time.Time, opts ...AuthenticationClaimsOption) AuthenticationClaims {
	authenticationClaims := AuthenticationClaims{
		AuthTime: jwt.NewNumericDate(authTime),
	}

	for _, opt := range opts {
		opt(&authenticationClaims)
	}

	return authenticationClaims
}

// AuthenticationClaimsOption is how options for the AuthenticationClaims are set up.
type AuthenticationClaimsOption func(*AuthenticationClaims)

// WithNonce applies a unique client session identifier (nonce) to the token to
// mitigate replay attacks by ensuring token-to-session association.
func WithNonce(nonce string) AuthenticationClaimsOption {
	return func(c *AuthenticationClaims) {
		c.Nonce = nonce
	}
}

// WithReferences applies the Authentication Context Class Reference (ACR) and the
// Authentication Methods References (AMR) slice to specify the security strength
// and the explicit list of authentication methods (e.g., "pwd", "mfa") used
// during login.
func WithReferences(acr string, amr []string) AuthenticationClaimsOption {
	return func(c *AuthenticationClaims) {
		c.ACR = acr
		c.AMR = amr
	}
}

// WithAuthorizedParty applies the specific party (OAuth 2.0 client ID) to
// which the ID Token was issued, useful in delegated authorization scenarios.
func WithAuthorizedParty(azp string) AuthenticationClaimsOption {
	return func(c *AuthenticationClaims) {
		c.AZP = azp
	}
}

// WithHash applies validation hashes for the Access Token (at_hash) and
// Authorization Code (c_hash) to ensure client-side token integrity.
func WithHash(atHash, cHash string) AuthenticationClaimsOption {
	return func(c *AuthenticationClaims) {
		c.AtHash = atHash
		c.CHash = cHash
	}
}
