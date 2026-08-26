package claims

import "time"

// IDTokenClaims represents the complete set of claims for an ID token as defined
// in the OpenID Connect Core 1.0 specification.
//
// The structure combines required protocol security metadata, the context of
// the current authentication session, and standard user profile data requested
// through the corresponding OIDC scopes.
type IDTokenClaims struct {
	RegisteredClaims        `json:",inline"`
	AuthenticationClaims    `json:",inline"`
	OpenIDProfileClaims     `json:",inline"`
	OpenIDPhoneNumberClaims `json:",inline"`
	OpenIDEmailClaims       `json:",inline"`
	OpenIDAddressClaims     `json:",inline"`
}

// NewIDTokenClaims returns a new IDTokenClaims struct in strict compliance with
// the OpenID Connect Core 1.0 profile by embedding required registered JWT
// metadata and lifecycle rules. It sets up core identification and expiration
// fields first, then evaluates functional options to append session
// authentication data and user-specific claims derived from authorized OIDC
// scopes.
func NewIDTokenClaims(
	issuer string,
	subject string,
	audience []string,
	ttl time.Duration,
	opts ...IDTokenClaimsOption,
) IDTokenClaims {
	registeredClaims := NewRegisteredClaims(issuer, subject, audience, ttl)

	idTokenClaims := IDTokenClaims{
		RegisteredClaims: registeredClaims,
	}

	for _, opt := range opts {
		opt(&idTokenClaims)
	}

	return idTokenClaims
}

// IDTokenClaimsOption is how options for the IDTokenClaims are set up.
type IDTokenClaimsOption func(*IDTokenClaims)

// WithIDTokenAuthentication applies authentication session details and context
// references into the ID Token, supplying required OIDC attributes like
// 'auth_time' and 'nonce'.
func WithIDTokenAuthentication(authTime time.Time, opts ...AuthenticationClaimsOption) IDTokenClaimsOption {
	return func(c *IDTokenClaims) {
		c.AuthenticationClaims = NewAuthenticationClaims(authTime, opts...)
	}
}

// WithIDTokenProfile applies standard user identity attributes (such as names,
// locale, URLs, and picture) into the ID Token to satisfy requests containing
// the 'profile' scope.
func WithIDTokenProfile(opts ...OpenIDProfileClaimsOption) IDTokenClaimsOption {
	return func(c *IDTokenClaims) {
		c.OpenIDProfileClaims = NewOpenIDProfileClaims(opts...)
	}
}

// WithIDTokenPhoneNumber applies the user's phone number and its verified state
// to the ID Token, satisfying requests containing the 'phone' scope.
func WithIDTokenPhoneNumber(number string, verified bool) IDTokenClaimsOption {
	return func(c *IDTokenClaims) {
		c.OpenIDPhoneNumberClaims = NewOpenIDPhoneNumberClaims(number, verified)
	}
}

// WithIDTokenEmail applies the user's communication email address and its
// verified state to the ID Token, satisfying requests containing the 'email'
// scope.
func WithIDTokenEmail(email string, verified bool) IDTokenClaimsOption {
	return func(c *IDTokenClaims) {
		c.OpenIDEmailClaims = NewOpenIDEmailClaims(email, verified)
	}
}

// WithIDTokenAddress applies structured and formatted physical mailing details
// into the ID Token, satisfying requests containing the 'address' scope.
func WithIDTokenAddress(
	formatted,
	streetAddress,
	locality,
	region,
	postalCode,
	country string,
) IDTokenClaimsOption {
	return func(c *IDTokenClaims) {
		c.OpenIDAddressClaims = NewOpenIDAddressClaims(
			formatted,
			streetAddress,
			locality,
			region,
			postalCode,
			country,
		)
	}
}
