package claims

import "time"

// AccessTokenClaims represents the complete set of access token claims in strict
// compliance with the RFC 9068 specification.
//
// The structure combines basic JWT metadata, authentication session parameters,
// client identifiers, and authorization contexts.
type AccessTokenClaims struct {
	RegisteredClaims     `json:",inline"`
	AuthenticationClaims `json:",inline"`
	AuthorizationClaims  `json:",inline"`
	ClientClaims         `json:",inline"`
}

// NewAccessTokenClaims returns a new AccessTokenClaims struct in strict
// compliance with the RFC 9068 profile by embedding required registered JWT
// metadata and applying optional attributes. It constructs the underlying
// RegisteredClaims with the core identity and expiration fields, then evaluates
// functional options to populate authorization, authentication session, or
// client context states.
func NewAccessTokenClaims(issuer, subject string, audience []string, ttl time.Duration, opts ...AccessTokenClaimsOption) AccessTokenClaims {
	registeredClaims := NewRegisteredClaims(issuer, subject, audience, ttl)

	accessTokenClaims := AccessTokenClaims{
		RegisteredClaims: registeredClaims,
	}

	for _, opt := range opts {
		opt(&accessTokenClaims)
	}

	return accessTokenClaims
}

// AccessTokenClaimsOption is how options for the AccessTokenClaims are set up.
type AccessTokenClaimsOption func(*AccessTokenClaims)

// WithAccessTokenAuthentication applies authentication session details into the
// access token by initializing and applying OpenID Connect identity validation
// options.
func WithAccessTokenAuthentication(authTime time.Time, opts ...AuthenticationClaimsOption) AccessTokenClaimsOption {
	return func(c *AccessTokenClaims) {
		c.AuthenticationClaims = NewAuthenticationClaims(authTime, opts...)
	}
}

// WithAccessTokenAuthorization applies access control boundaries into the token
// by mapping scopes and evaluating additional options like roles, groups, or
// granular entitlements.
func WithAccessTokenAuthorization(scopes []string, opts ...AuthorizationClaimsOption) AccessTokenClaimsOption {
	return func(c *AccessTokenClaims) {
		c.AuthorizationClaims = NewAuthorizationClaims(scopes, opts...)
	}
}

// WithAccessTokenClient applies a specific registered OAuth 2.0 client
// application using its identifier.
func WithAccessTokenClient(clientID string) AccessTokenClaimsOption {
	return func(c *AccessTokenClaims) {
		c.ClientClaims = NewClientClaims(clientID)
	}
}
