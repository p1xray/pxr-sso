package claims

import (
	"github.com/go-jose/go-jose/v4/jwt"
	"github.com/google/go-cmp/cmp"
	"strings"
	"testing"
	"time"
)

func TestNewRefreshTokenClaims(t *testing.T) {
	testCases := []struct {
		name           string
		issuer         string
		subject        string
		audience       jwt.Audience
		ttl            time.Duration
		opts           []RefreshTokenClaimsOption
		expectedClaims RefreshTokenClaims
	}{
		{
			name:     "successful creation of new refresh token claims with required parameters",
			issuer:   testIssuer,
			subject:  testSubject,
			audience: testAudience,
			ttl:      testTTL,
			expectedClaims: RefreshTokenClaims{
				RegisteredClaims: RegisteredClaims{
					Issuer:   testIssuer,
					Subject:  testSubject,
					Audience: testAudience,
				},
			},
		},
		{
			name:     "successful creation of new refresh token claims with all required and optional parameters",
			issuer:   testIssuer,
			subject:  testSubject,
			audience: testAudience,
			ttl:      testTTL,
			opts: []RefreshTokenClaimsOption{
				WithRefreshTokenAuthentication(
					testAuthTime,
					WithNonce(testNonce),
					WithReferences(testACR, testAMR),
					WithAuthorizedParty(testAZP),
					WithHash(testAtHash, testCHash),
				),
				WithRefreshTokenAuthorization(
					testScopes,
					WithRoles(testRoles),
					WithGroups(testGroups),
					WithEntitlements(testEntitlements),
				),
				WithRefreshTokenClient(testClientID),
				WithRefreshTokenSession(testSessionID),
			},
			expectedClaims: RefreshTokenClaims{
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
				SessionClaims: SessionClaims{
					SessionID: testSessionID,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tokenClaims := NewRefreshTokenClaims(tc.issuer, tc.subject, tc.audience, tc.ttl, tc.opts...)

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

			if !cmp.Equal(tc.expectedClaims.SessionClaims, tokenClaims.SessionClaims) {
				t.Fatal(cmp.Diff(tc.expectedClaims.SessionClaims, tokenClaims.SessionClaims))
			}
		})
	}
}
