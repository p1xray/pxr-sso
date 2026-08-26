package claims

import "strings"

// AuthorizationClaims represents which access control included in
// the access token.
//
// The structure integrates various authorization models: application privilege
// differentiation, role-based access control, group security policies, and
// granular permissions.
type AuthorizationClaims struct {
	// Scope represents the scope associated with an access token.
	//
	// Specifies the permissions or access rights granted to an access token,
	// defining what actions the token bearer is authorized to perform. This claim is
	// fundamental in controlling access to resources.
	Scope string `json:"scope,omitempty"`

	// Roles represents the roles associated with the subject, often used in System
	// for Cross-domain Identity Management (SCIM).
	//
	// This claim is used to convey the roles attributed to a subject, facilitating
	// role-based access control (RBAC) and other authorization decisions in systems
	// implementing SCIM or similar identity management protocols.
	Roles []string `json:"roles,omitempty"`

	// Groups represents the groups that the subject belongs to, typically used in
	// identity and access management.
	//
	// This claim is used to convey group membership information, supporting
	// group-based access control and enabling systems to make authorization
	// decisions based on the groups a subject is associated with.
	Groups []string `json:"groups,omitempty"`

	// Entitlements represents specific entitlements or permissions granted to the subject.
	//
	// This claim is used to specify granular permissions or entitlements granted to
	// a subject, allowing for precise control over access rights and enabling
	// fine-grained authorization policies.
	Entitlements []string `json:"entitlements,omitempty"`
}

// NewAuthorizationClaims returns a new AuthorizationClaims struct using the
// required scopes and applies optional configurations.
func NewAuthorizationClaims(scopes []string, opts ...AuthorizationClaimsOption) AuthorizationClaims {
	authorizationClaims := AuthorizationClaims{
		Scope: strings.Join(scopes, " "),
	}

	for _, opt := range opts {
		opt(&authorizationClaims)
	}

	return authorizationClaims
}

// AuthorizationClaimsOption is how options for the AuthorizationClaims are set up.
type AuthorizationClaimsOption func(*AuthorizationClaims)

// WithRoles applies a list of assigned security roles into the claims, enabling
// Role-Based Access Control (RBAC) and SCIM-compatible authorization checks.
func WithRoles(roles []string) AuthorizationClaimsOption {
	return func(c *AuthorizationClaims) {
		c.Roles = roles
	}
}

// WithGroups applies a list of identity management group memberships into the
// claims to support group-based security policies and access control.
func WithGroups(groups []string) AuthorizationClaimsOption {
	return func(c *AuthorizationClaims) {
		c.Groups = groups
	}
}

// WithEntitlements applies a set of fine-grained permissions or explicit
// entitlements to the claims, allowing for precise and granular authorization
// decisions.
func WithEntitlements(entitlements []string) AuthorizationClaimsOption {
	return func(c *AuthorizationClaims) {
		c.Entitlements = entitlements
	}
}
