package claims

import (
	"github.com/google/go-cmp/cmp"
	"strings"
	"testing"
)

var (
	testScopes       = []string{"openid profile email users:read"}
	testRoles        = []string{"member", "manager"}
	testGroups       = []string{"admin", "moderator"}
	testEntitlements = []string{"users:read", "users:write"}
)

func TestNewAuthorizationClaims(t *testing.T) {
	testCases := []struct {
		name           string
		scopes         []string
		opts           []AuthorizationClaimsOption
		expectedClaims AuthorizationClaims
	}{
		{
			name:   "successful creation of new authorization claims with required parameters",
			scopes: testScopes,
			expectedClaims: AuthorizationClaims{
				Scope: strings.Join(testScopes, " "),
			},
		},
		{
			name:   "successful creation of new authentication claims with all required and optional parameters",
			scopes: testScopes,
			opts: []AuthorizationClaimsOption{
				WithRoles(testRoles),
				WithGroups(testGroups),
				WithEntitlements(testEntitlements),
			},
			expectedClaims: AuthorizationClaims{
				Scope:        strings.Join(testScopes, " "),
				Roles:        testRoles,
				Groups:       testGroups,
				Entitlements: testEntitlements,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			authorizationClaims := NewAuthorizationClaims(tc.scopes, tc.opts...)

			if !cmp.Equal(tc.expectedClaims, authorizationClaims) {
				t.Fatal(cmp.Diff(tc.expectedClaims, authorizationClaims))
			}
		})
	}
}
