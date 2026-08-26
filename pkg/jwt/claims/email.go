package claims

// OpenIDEmailClaims represents user's email address information claims that is
// included when providing "email" scope.
//
// Defined in OpenID Connect Core 1.0, Section 5.1.
type OpenIDEmailClaims struct {
	// Email represents the user's preferred email address.
	//
	// Used in OpenID Connect. This claim is for the user's preferred email. Note
	// that this email address might not be unique.
	Email string `json:"email,omitempty"`

	// EmailVerified represents whether the user's email address has been verified.
	//
	// Used in OpenID Connect as a boolean. True if the user's email address has been
	// verified; otherwise false.
	EmailVerified *bool `json:"email_verified,omitempty"`
}

// NewOpenIDEmailClaims returns a new OpenIDEmailClaims struct with the user's
// preferred email address and its verification status to support the OIDC
// 'email' scope.
func NewOpenIDEmailClaims(email string, verified bool) OpenIDEmailClaims {
	return OpenIDEmailClaims{
		Email:         email,
		EmailVerified: &verified,
	}
}
