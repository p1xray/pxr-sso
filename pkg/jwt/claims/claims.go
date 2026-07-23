package claims

import (
	"github.com/go-jose/go-jose/v4/jwt"
	"github.com/google/uuid"
	"time"
)

const defaultTokenTTL = 10 * time.Minute

// RegisteredClaims describes the token itself, its intended recipients, and its
// expiration date according to the JWT specification (RFC 7519).
type RegisteredClaims struct {
	// ID provides a unique identifier for the JWT.
	//
	// Defined in RFC 7519, Section 4.1.7. The "jti" (JWT ID) claim is a
	// case-sensitive string. Use this claim to prevent the JWT from being replayed.
	ID string `json:"jti,omitempty"`

	// Issuer represents the principal (e.g., authorization server) that issued the JWT.
	//
	// Defined in RFC 7519, Section 4.1.1. It is a case-sensitive string containing a
	// StringOrURI value. Use this claim to identify the issuer of the JWT uniquely.
	Issuer string `json:"iss,omitempty"`

	// Subject represents the principal that is the subject of the JWT.
	//
	// Defined in RFC 7519, Section 4.1.2. The "sub" value is a case-sensitive string
	// containing a StringOrURI value. This claim is used to identify the subject of
	// the JWT, which could be an end user or a device.
	Subject string `json:"sub,omitempty"`

	// Audience identifies the recipients that the JWT is intended for.
	//
	// Defined in RFC 7519, Section 4.1.3. It is generally a case-sensitive string or
	// an array of strings containing StringOrURI values. The audience claim ensures
	// that the JWT is sent to the intended recipients.
	Audience jwt.Audience `json:"aud,omitempty"`

	// Expiry specifies the expiration time on or after which the JWT must not be
	// accepted for processing.
	//
	// Defined in RFC 7519, Section 4.1.4. The "exp" claim is a NumericDate value.
	// Use this claim to define the validity period of the JWT.
	Expiry *jwt.NumericDate `json:"exp,omitempty"`

	// NotBefore defines a time before which the JWT MUST NOT be accepted for
	// processing.
	//
	// Defined in RFC 7519, Section 4.1.5. The "nbf" (Not Before) claim is a
	// NumericDate value. This claim helps in ensuring that a JWT is not accepted
	// before a certain time.
	NotBefore *jwt.NumericDate `json:"nbf,omitempty"`

	// IssuedAt indicates the time at which the JWT was issued.
	//
	// Defined in RFC 7519, Section 4.1.6. The "iat" (Issued At) claim is a
	// NumericDate value. This claim can be used to determine the age of the JWT.
	IssuedAt *jwt.NumericDate `json:"iat,omitempty"`
}

// NewRegisteredClaims returns a new RegisteredClaims structure with a unique
// identifier, the provided parameters, and standardized timestamps.
func NewRegisteredClaims(issuer, subject string, audience []string, ttl time.Duration) RegisteredClaims {
	utcNow := time.Now().UTC()

	expiry := utcNow.Add(defaultTokenTTL)
	if ttl > 0 {
		expiry = utcNow.Add(ttl)
	}

	registeredClaims := RegisteredClaims{
		ID:        uuid.New().String(),
		Issuer:    issuer,
		Subject:   subject,
		Audience:  audience,
		Expiry:    jwt.NewNumericDate(expiry),
		IssuedAt:  jwt.NewNumericDate(utcNow),
		NotBefore: jwt.NewNumericDate(utcNow),
	}

	return registeredClaims
}
