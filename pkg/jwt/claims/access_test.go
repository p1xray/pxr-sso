package claims

import (
	"github.com/go-jose/go-jose/v4/jwt"
	"github.com/google/go-cmp/cmp"
	"strings"
	"testing"
	"time"
)

func TestNewAccessTokenClaims(t *testing.T) {
	testCases := []struct {
		name           string
		issuer         string
		subject        string
		audience       jwt.Audience
		ttl            time.Duration
		opts           []AccessTokenClaimsOption
		expectedClaims AccessTokenClaims
	}{
		{
			name:     "successful creation of new access token claims with required parameters",
			issuer:   testIssuer,
			subject:  testSubject,
			audience: testAudience,
			ttl:      testTTL,
			expectedClaims: AccessTokenClaims{
				RegisteredClaims: RegisteredClaims{
					Issuer:   testIssuer,
					Subject:  testSubject,
					Audience: testAudience,
				},
			},
		},
		{
			name:     "successful creation of new access token claims with all required and optional parameters",
			issuer:   testIssuer,
			subject:  testSubject,
			audience: testAudience,
			ttl:      testTTL,
			opts: []AccessTokenClaimsOption{
				WithAccessTokenAuthentication(
					testAuthTime,
					WithNonce(testNonce),
					WithReferences(testACR, testAMR),
					WithAuthorizedParty(testAZP),
					WithHash(testAtHash, testCHash),
				),
				WithAccessTokenAuthorization(
					testScopes,
					WithRoles(testRoles),
					WithGroups(testGroups),
					WithEntitlements(testEntitlements),
				),
				WithAccessTokenClient(testClientID),
			},
			expectedClaims: AccessTokenClaims{
				RegisteredClaims: RegisteredClaims{
					Issuer:   testIssuer,
					Subject:  testSubject,
					Audience: testAudience,
				},
				AuthenticationClaims: AuthenticationClaims{
					AuthTime: jwt.NewNumericDate(testAuthTime),
					Nonce:    testNonce,
					ACR:      testACR,
					AMR:      testAMR,
					AZP:      testAZP,
					AtHash:   testAtHash,
					CHash:    testCHash,
				},
				AuthorizationClaims: AuthorizationClaims{
					Scope:        strings.Join(testScopes, " "),
					Roles:        testRoles,
					Groups:       testGroups,
					Entitlements: testEntitlements,
				},
				ClientClaims: ClientClaims{
					ClientID: testClientID,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tokenClaims := NewAccessTokenClaims(tc.issuer, tc.subject, tc.audience, tc.ttl, tc.opts...)

			checkRegisteredClaims(
				t,
				tokenClaims.RegisteredClaims,
				tc.expectedClaims.Issuer,
				tc.expectedClaims.Subject,
				tc.expectedClaims.Audience,
				leeway,
			)

			if !cmp.Equal(tc.expectedClaims.AuthenticationClaims, tokenClaims.AuthenticationClaims) {
				t.Fatal(cmp.Diff(tc.expectedClaims.AuthenticationClaims, tokenClaims.AuthenticationClaims))
			}

			if !cmp.Equal(tc.expectedClaims.AuthorizationClaims, tokenClaims.AuthorizationClaims) {
				t.Fatal(cmp.Diff(tc.expectedClaims.AuthorizationClaims, tokenClaims.AuthorizationClaims))
			}

			if !cmp.Equal(tc.expectedClaims.ClientClaims, tokenClaims.ClientClaims) {
				t.Fatal(cmp.Diff(tc.expectedClaims.ClientClaims, tokenClaims.ClientClaims))
			}
		})
	}
}
