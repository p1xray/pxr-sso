package jwt

import (
	"github.com/go-jose/go-jose/v4/jwt"
	"github.com/google/go-cmp/cmp"
	"github.com/p1xray/pxr-sso/pkg/jwt/claims"
	"github.com/p1xray/pxr-sso/pkg/jwt/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestParseToken(t *testing.T) {
	testCases := []struct {
		name                       string
		rawToken                   string
		allowedSignatureAlgorithms []crypto.SignatureAlgorithm
		expectedClaims             map[string]interface{}
		expectedError              error
	}{
		{
			name:     "successfully parse a token claims",
			rawToken: testValidMinimalAccessToken,
			expectedClaims: map[string]interface{}{
				"aud": testAudience,
				"exp": float64(testLargeNumericDate),
				"iat": float64(testSmallNumericDate),
				"iss": testIssuer,
				"jti": testAccessTokenID,
				"nbf": float64(testSmallNumericDate),
				"sub": testSubject,
			},
			allowedSignatureAlgorithms: []crypto.SignatureAlgorithm{crypto.HS256},
		},
		{
			name:                       "throws an error when parsing a token with not allowed signature algorithm",
			rawToken:                   testValidMinimalAccessToken,
			allowedSignatureAlgorithms: []crypto.SignatureAlgorithm{crypto.RS256},
			expectedError:              ErrParseToken,
		},
		{
			name:                       "throws an error when parsing an invalid token",
			rawToken:                   testInvalidAccessToken,
			allowedSignatureAlgorithms: []crypto.SignatureAlgorithm{crypto.HS256},
			expectedError:              ErrParseToken,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			token, err := ParseToken(tc.rawToken, tc.allowedSignatureAlgorithms)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)

				tokenClaims := make(map[string]interface{})
				err = token.UnsafeClaimsWithoutVerification(&tokenClaims)
				require.NoError(t, err)

				if !cmp.Equal(tc.expectedClaims, tokenClaims) {
					t.Fatal(cmp.Diff(tc.expectedClaims, tokenClaims))
				}
			}
		})
	}
}

func TestUnsafeParseTokenClaims(t *testing.T) {
	testCases := []struct {
		name           string
		token          *jwt.JSONWebToken
		expectedClaims claims.RegisteredClaims
		expectedError  error
	}{
		{
			name:  "successfully unsafe parse a token claims",
			token: NewTestTokenWithRegisteredClaims(),
			expectedClaims: claims.RegisteredClaims{
				ID:        testAccessTokenID,
				Issuer:    testIssuer,
				Subject:   testSubject,
				Audience:  jwt.Audience{testAudience},
				Expiry:    jwt.NewNumericDate(time.Unix(testLargeNumericDate, 0)),
				IssuedAt:  jwt.NewNumericDate(time.Unix(testSmallNumericDate, 0)),
				NotBefore: jwt.NewNumericDate(time.Unix(testSmallNumericDate, 0)),
			},
		},
		{
			name:          "throws an error when unsafe parsing an invalid token",
			token:         &jwt.JSONWebToken{},
			expectedError: ErrParseClaims,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tokenClaims, err := UnsafeParseTokenClaims(tc.token)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)

				if !cmp.Equal(tc.expectedClaims, tokenClaims) {
					t.Fatal(cmp.Diff(tc.expectedClaims, tokenClaims))
				}
			}
		})
	}
}

func TestParseTokenClaims(t *testing.T) {
	testCases := []struct {
		name           string
		token          *jwt.JSONWebToken
		key            []byte
		expectedClaims claims.RegisteredClaims
		expectedError  error
	}{
		{
			name:  "successfully parse a token claims",
			token: NewTestTokenWithRegisteredClaims(),
			key:   []byte(testValidKey),
			expectedClaims: claims.RegisteredClaims{
				ID:        testAccessTokenID,
				Issuer:    testIssuer,
				Subject:   testSubject,
				Audience:  jwt.Audience{testAudience},
				Expiry:    jwt.NewNumericDate(time.Unix(testLargeNumericDate, 0)),
				IssuedAt:  jwt.NewNumericDate(time.Unix(testSmallNumericDate, 0)),
				NotBefore: jwt.NewNumericDate(time.Unix(testSmallNumericDate, 0)),
			},
		},
		{
			name:          "throws an error when parsing a token claims with invalid key",
			token:         NewTestTokenWithRegisteredClaims(),
			key:           []byte(testInvalidKey),
			expectedError: ErrParseClaims,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tokenClaims, err := ParseTokenClaims(tc.token, tc.key)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)

				if !cmp.Equal(tc.expectedClaims, tokenClaims) {
					t.Fatal(cmp.Diff(tc.expectedClaims, tokenClaims))
				}
			}
		})
	}
}

func TestParseAccessTokenClaims(t *testing.T) {
	testCases := []struct {
		name           string
		token          *jwt.JSONWebToken
		key            []byte
		expectedClaims claims.AccessTokenClaims
		expectedError  error
	}{
		{
			name:           "successfully parse a access token claims",
			token:          NewTestAccessToken(),
			key:            []byte(testValidKey),
			expectedClaims: NewTestFullAccessTokenClaims(),
		},
		{
			name:          "throws an error when parsing a access token claims with invalid key",
			token:         NewTestAccessToken(),
			key:           []byte(testInvalidKey),
			expectedError: ErrParseClaims,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tokenClaims, err := ParseAccessTokenClaims(tc.token, tc.key)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)

				if !cmp.Equal(tc.expectedClaims, tokenClaims) {
					t.Fatal(cmp.Diff(tc.expectedClaims, tokenClaims))
				}
			}
		})
	}
}

func TestParseRefreshTokenClaims(t *testing.T) {
	testCases := []struct {
		name           string
		token          *jwt.JSONWebToken
		key            []byte
		expectedClaims claims.RefreshTokenClaims
		expectedError  error
	}{
		{
			name:           "successfully parse a refresh token claims",
			token:          NewTestRefreshToken(),
			key:            []byte(testValidKey),
			expectedClaims: NewTestFullRefreshTokenClaims(),
		},
		{
			name:          "throws an error when parsing a refresh token claims with invalid key",
			token:         NewTestRefreshToken(),
			key:           []byte(testInvalidKey),
			expectedError: ErrParseClaims,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tokenClaims, err := ParseRefreshTokenClaims(tc.token, tc.key)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)

				if !cmp.Equal(tc.expectedClaims, tokenClaims) {
					t.Fatal(cmp.Diff(tc.expectedClaims, tokenClaims))
				}
			}
		})
	}
}

func TestParseIDTokenClaims(t *testing.T) {
	testCases := []struct {
		name           string
		token          *jwt.JSONWebToken
		key            []byte
		expectedClaims claims.IDTokenClaims
		expectedError  error
	}{
		{
			name:           "successfully parse an id token claims",
			token:          NewTestIDToken(),
			key:            []byte(testValidKey),
			expectedClaims: NewTestFullIDTokenClaims(),
		},
		{
			name:          "throws an error when parsing an id token claims with invalid key",
			token:         NewTestIDToken(),
			key:           []byte(testInvalidKey),
			expectedError: ErrParseClaims,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tokenClaims, err := ParseIDTokenClaims(tc.token, tc.key)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)

				if !cmp.Equal(tc.expectedClaims, tokenClaims) {
					t.Fatal(cmp.Diff(tc.expectedClaims, tokenClaims))
				}
			}
		})
	}
}
