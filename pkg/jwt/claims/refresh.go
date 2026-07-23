package claims

import "time"

// RefreshTokenClaims represents a set of claims for the JWT implementation of
// the Refresh Token.
//
// This token is used exclusively by the Authorization Server to securely extend
// sessions and issue new Access/ID tokens.
type RefreshTokenClaims struct {
	RegisteredClaims     `json:",inline"`
	AuthenticationClaims `json:",inline"`
	AuthorizationClaims  `json:",inline"`
	ClientClaims         `json:",inline"`
	SessionClaims        `json:",inline"`
}

// NewRefreshTokenClaims returns a new RefreshTokenClaims struct designed for
// managing long-lived sessions. It constructs the underlying RegisteredClaims
// with core expiration parameters and evaluates the provided functional options
// to inject user session states, authorization scopes, authentication histories,
// or client identifiers.
func NewRefreshTokenClaims(
	issuer string,
	subject string,
	audience []string,
	ttl time.Duration,
	opts ...RefreshTokenClaimsOption,
) RefreshTokenClaims {
	registeredClaims := NewRegisteredClaims(issuer, subject, audience, ttl)

	refreshTokenClaims := RefreshTokenClaims{
		RegisteredClaims: registeredClaims,
	}

	for _, opt := range opts {
		opt(&refreshTokenClaims)
	}

	return refreshTokenClaims
}

// RefreshTokenClaimsOption is how options for the RefreshTokenClaims are set up.
type RefreshTokenClaimsOption func(*RefreshTokenClaims)

// WithRefreshTokenAuthentication applies the original authentication context and
// timestamp into the refresh token, allowing the server to track session
// freshness during token rotation.
func WithRefreshTokenAuthentication(authTime time.Time, opts ...AuthenticationClaimsOption) RefreshTokenClaimsOption {
	return func(c *RefreshTokenClaims) {
		c.AuthenticationClaims = NewAuthenticationClaims(authTime, opts...)
	}
}

// WithRefreshTokenAuthorization applies the granted authorization scopes and
// optional privileges to the refresh token, restricting downstream access tokens
// to this exact boundary or its subset.
func WithRefreshTokenAuthorization(scopes []string, opts ...AuthorizationClaimsOption) RefreshTokenClaimsOption {
	return func(c *RefreshTokenClaims) {
		c.AuthorizationClaims = NewAuthorizationClaims(scopes, opts...)
	}
}

// WithRefreshTokenClient applies the specific registered OAuth 2.0 client
// application that requested it, preventing token misuse by unauthorized
// clients.
func WithRefreshTokenClient(clientID string) RefreshTokenClaimsOption {
	return func(c *RefreshTokenClaims) {
		c.ClientClaims = NewClientClaims(clientID)
	}
}

// WithRefreshTokenSession applies the active browser session identifier (sid),
// enabling the authorization server to perform global or single sign-out
// revocations.
func WithRefreshTokenSession(sessionID string) RefreshTokenClaimsOption {
	return func(c *RefreshTokenClaims) {
		c.SessionClaims = NewSessionClaims(sessionID)
	}
}
