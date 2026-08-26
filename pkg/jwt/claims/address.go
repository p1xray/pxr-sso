package claims

// OpenIDAddressClaims represents user's postal address information claims that is
// included when providing "address" scope.
//
// Defined in OpenID Connect Core 1.0, Section 5.1.
type OpenIDAddressClaims struct {
	// Address represents the user's postal address.
	//
	// Used in OpenID Connect. This JSON structured claim contains components of the
	// user's address such as street address, locality, region, postal code, and
	// country.
	Address *AddressClaim `json:"address,omitempty"`
}

// AddressClaim represents a physical mailing address.
//
// Defined in OpenID Connect Core 1.0, Section 5.1.1.
type AddressClaim struct {
	// Formatted represents the full mailing address, formatted for display or use on
	// a mailing label.
	//
	// This claim may contain multiple lines, separated by newlines. Newlines can be
	// represented either as a carriage return/line feed pair ("\r\n") or as a single
	// line feed character ("\n").
	Formatted string `json:"formatted,omitempty"`

	// StreetAddress represents the full street address component, which may include
	// house number, street name, Post Office Box, and multi-line extended street
	// address information.
	//
	// This claim may contain multiple lines, separated by newlines. Newlines can be
	// represented either as a carriage return/line feed pair ("\r\n") or as a single
	// line feed character ("\n").
	StreetAddress string `json:"street_address,omitempty"`

	// Locality represents the city or locality component.
	Locality string `json:"locality,omitempty"`

	// Region represents the state, province, prefecture, or region component.
	Region string `json:"region,omitempty"`

	// PostalCode represents the zip code or postal code component.
	PostalCode string `json:"postal_code,omitempty"`

	// Country represents the country name component.
	Country string `json:"country,omitempty"`
}

// NewOpenIDAddressClaims returns a new OpenIDAddressClaims struct by
// constructing a nested AddressClaim with the provided physical location details
// to support the OIDC 'address' scope.
func NewOpenIDAddressClaims(
	formatted,
	streetAddress,
	locality,
	region,
	postalCode,
	country string,
) OpenIDAddressClaims {
	return OpenIDAddressClaims{
		Address: &AddressClaim{
			Formatted:     formatted,
			StreetAddress: streetAddress,
			Locality:      locality,
			Region:        region,
			PostalCode:    postalCode,
			Country:       country,
		},
	}
}
