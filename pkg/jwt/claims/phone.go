package claims

// OpenIDPhoneNumberClaims represents user's phone number information claims that
// is included when providing "phone" scope.
//
// Defined in OpenID Connect Core 1.0, Section 5.1.
type OpenIDPhoneNumberClaims struct {
	// PhoneNumber represents the user's phone number.
	//
	// Used in OpenID Connect. This claim provides the user's preferred phone number.
	// The format of the number can vary, and it's not guaranteed to be in a standard
	// format.
	PhoneNumber string `json:"phone_number,omitempty"`

	// PhoneNumberVerified indicates whether the user's phone number has been
	// verified.
	//
	// Used in OpenID Connect as a boolean. True if the user's phone number has been
	// verified; otherwise false.
	PhoneNumberVerified *bool `json:"phone_number_verified,omitempty"`
}

// NewOpenIDPhoneNumberClaims returns a new OpenIDPhoneNumberClaims struct with
// the user's phone number and its verification status to support the OIDC
// 'phone' scope.
func NewOpenIDPhoneNumberClaims(number string, verified bool) OpenIDPhoneNumberClaims {
	return OpenIDPhoneNumberClaims{
		PhoneNumber:         number,
		PhoneNumberVerified: &verified,
	}
}
