package claims

import (
	"github.com/go-jose/go-jose/v4/jwt"
	"time"
)

// OpenIDProfileClaims represents personal identifiable information claims that is
// included when providing "profile" scope.
//
// Defined in OpenID Connect Core 1.0, Section 5.1.
type OpenIDProfileClaims struct {
	// Name represents the full name of the user.
	//
	// Used in OpenID Connect for representing the user's full name in a single
	// string. It might include the first, middle, last, and other names.
	Name string `json:"name,omitempty"`

	// FamilyName represents the surname(s) or last name(s) of the user.
	//
	// Used in OpenID Connect. This claim focuses on the user's family name or
	// surname(s), excluding middle names.
	FamilyName string `json:"family_name,omitempty"`

	// GivenName represents the first or given name(s) of the user.
	//
	// Used in OpenID Connect. This claim is intended to refer to the user's first
	// name or given name(s). It allows for middle names if applicable.
	GivenName string `json:"given_name,omitempty"`

	// MiddleName represents the middle name(s) of the user.
	//
	// Used in OpenID Connect. This claim is intended for the user's middle name(s),
	// which might not be present for all users.
	MiddleName string `json:"middle_name,omitempty"`

	// Nickname represents the casual name of the user.
	//
	// Used in OpenID Connect. This claim is for the user's casual or informal name
	// that might differ from their legal name.
	Nickname string `json:"nickname,omitempty"`

	// PreferredUsername is the username preferred by the user, which may be
	// different from their actual or legal name.
	//
	// This claim is used to convey the user's preferred username within the system.
	// It allows the user to specify a nickname or alias that is used within the
	// application for display purposes, providing a more personalized user
	// experience.
	PreferredUsername string `json:"preferred_username,omitempty"`

	// Profile represents the URL of the user's profile page.
	//
	// Used in OpenID Connect. It is the URL of a web page containing information
	// about the user or a social profile page.
	Profile string `json:"profile,omitempty"`

	// Picture represents the URL of the user's profile picture.
	//
	// Used in OpenID Connect. This claim provides a URL pointing to a profile
	// picture or avatar of the user.
	Picture string `json:"picture,omitempty"`

	// Website represents the URL of the user's web page or blog.
	//
	// Used in OpenID Connect. This claim indicates the URL of the user's personal or
	// business website.
	Website string `json:"website,omitempty"`

	// Gender represents the user's gender.
	//
	// Used in OpenID Connect. This claim can be used to convey the user's gender.
	// The value is not strictly defined and can vary based on the user's preference
	// and the application's requirements.
	Gender string `json:"gender,omitempty"`

	// Birthdate represents the user's date of birth.
	//
	// Used in OpenID Connect. This claim is for the user's birthdate, typically
	// represented in the ISO 8601:2004 YYYY-MM-DD format.
	Birthdate string `json:"birthdate,omitempty"`

	// ZoneInfo represents the user's time zone.
	//
	// Used in OpenID Connect. This claim indicates the user's time zone,
	// facilitating localization and personalization.
	ZoneInfo string `json:"zoneinfo,omitempty"`

	// Locale represents the user's locale.
	//
	// Used in OpenID Connect. This claim specifies the user's preferred language and
	// optionally, region. Typically represented as a language tag, e.g., en-US.
	Locale string `json:"locale,omitempty"`

	// UpdatedAt indicates when the user's information was last updated.
	//
	// Used in OpenID Connect. This claim provides a Unix time stamp indicating when
	// the user's information was last updated.
	UpdatedAt *jwt.NumericDate `json:"updated_at,omitempty"`
}

// NewOpenIDProfileClaims returns a new OpenIDProfileClaims struct with optional
// configuration. It creates an empty profile instance and sequentially applies
// the provided configurations to populate standard OIDC user profile attributes
// (such as 'name', 'given_name', 'family_name', 'picture').
func NewOpenIDProfileClaims(opts ...OpenIDProfileClaimsOption) OpenIDProfileClaims {
	profileClaims := OpenIDProfileClaims{}

	for _, opt := range opts {
		opt(&profileClaims)
	}

	return profileClaims
}

// OpenIDProfileClaimsOption is how options for the OpenIDProfileClaims are set up.
type OpenIDProfileClaimsOption func(*OpenIDProfileClaims)

// WithProfileName applies standard full name, family, given, and middle name
// attributes to the profile claims according to OpenID Connect Core 1.0.
func WithProfileName(name, familyName, givenName, middleName string) OpenIDProfileClaimsOption {
	return func(c *OpenIDProfileClaims) {
		c.Name = name
		c.FamilyName = familyName
		c.GivenName = givenName
		c.MiddleName = middleName
	}
}

// WithProfileNameAlias applies casual and system-wide identification aliases
// using the 'nickname' and 'preferred_username' claims.
func WithProfileNameAlias(nickname, preferredUsername string) OpenIDProfileClaimsOption {
	return func(c *OpenIDProfileClaims) {
		c.Nickname = nickname
		c.PreferredUsername = preferredUsername
	}
}

// WithProfileURL applies external resource web links associated with the
// subject, populating their profile page, avatar picture, and main website URLs.
func WithProfileURL(profile, picture, website string) OpenIDProfileClaimsOption {
	return func(c *OpenIDProfileClaims) {
		c.Profile = profile
		c.Picture = picture
		c.Website = website
	}
}

// WithProfileGender applies the end-user's gender preference claim.
func WithProfileGender(gender string) OpenIDProfileClaimsOption {
	return func(c *OpenIDProfileClaims) {
		c.Gender = gender
	}
}

// WithProfileBirthdate applies the user's date of birth string, normally
// formatted as YYYY-MM-DD following ISO 8601.
func WithProfileBirthdate(birthdate string) OpenIDProfileClaimsOption {
	return func(c *OpenIDProfileClaims) {
		c.Birthdate = birthdate
	}
}

// WithProfileLocation applies regional preferences by setting the time zone name
// and BCP47 language/locale identifier.
func WithProfileLocation(zoneInfo, locale string) OpenIDProfileClaimsOption {
	return func(c *OpenIDProfileClaims) {
		c.ZoneInfo = zoneInfo
		c.Locale = locale
	}
}

// WithProfileUpdatedAt applies the exact time when the end-user's profile
// information was last modified.
func WithProfileUpdatedAt(updatedAt time.Time) OpenIDProfileClaimsOption {
	return func(c *OpenIDProfileClaims) {
		c.UpdatedAt = jwt.NewNumericDate(updatedAt)
	}
}
