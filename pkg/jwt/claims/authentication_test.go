package claims

import (
	"github.com/go-jose/go-jose/v4/jwt"
	"github.com/google/go-cmp/cmp"
	"testing"
	"time"
)

const (
	testNonce  = "_Xm8-pWq29_LzR-k"
	testACR    = "urn:pxr:loa:1fa:any"
	testAZP    = "ELuXIbyntYQwIgUb8kFRmAotDzoGdPQZ"
	testAtHash = "U4osGiwnaxsnr4zCZBs2pA"
	testCHash  = "vF9dft4qmTGHjlOhiHcTqw"
)

var (
	testAuthTime = time.Date(2026, time.July, 21, 12, 0, 0, 0, time.UTC)
	testAMR      = []string{"pwd", "mfa"}
)

func TestNewAuthenticationClaims(t *testing.T) {
	testCases := []struct {
		name           string
		authTime       time.Time
		opts           []AuthenticationClaimsOption
		expectedClaims AuthenticationClaims
	}{
		{
			name:     "successful creation of new authentication claims with required parameters",
			authTime: testAuthTime,
			expectedClaims: AuthenticationClaims{
				AuthTime: jwt.NewNumericDate(testAuthTime),
			},
		},
		{
			name:     "successful creation of new authentication claims with all required and optional parameters",
			authTime: testAuthTime,
			opts: []AuthenticationClaimsOption{
				WithNonce(testNonce),
				WithReferences(testACR, testAMR),
				WithAuthorizedParty(testAZP),
				WithHash(testAtHash, testCHash),
			},
			expectedClaims: AuthenticationClaims{
				AuthTime: jwt.NewNumericDate(testAuthTime),
				Nonce:    testNonce,
				ACR:      testACR,
				AMR:      testAMR,
				AZP:      testAZP,
				AtHash:   testAtHash,
				CHash:    testCHash,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			authenticationClaims := NewAuthenticationClaims(tc.authTime, tc.opts...)

			if !cmp.Equal(tc.expectedClaims, authenticationClaims) {
				t.Fatal(cmp.Diff(tc.expectedClaims, authenticationClaims))
			}
		})
	}
}
