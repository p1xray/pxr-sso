package jwtparser

import (
	"context"
	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
	"github.com/google/go-cmp/cmp"
	jwtclaims "github.com/p1xray/pxr-sso/pkg/jwt/claims"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const (
	validKey            = "98649a5c-2137-4a78-a63f-fbab416a7f9e"
	invalidKey          = "invalid_key"
	validRefreshToken   = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTA3OTMyMjgsImp0aSI6ImNiYThhOWNiLTc5MDktNGQxYy05ZWM2LTljYTBhNDIyNjVmYiJ9.rEdD8JT8yo9xG2N-iFkK1dzGLFNukTiA7kMNd2ctyZc"
	invalidRefreshToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.ey"
)

type testCustomClaims struct {
	TestCustom  string `json:"test_custom,omitempty"`
	Test2Custom string `json:"test2_custom,omitempty"`
}

func (tcc testCustomClaims) Validate(context.Context) error {
	return nil
}

type testCustomClaimsWithoutTags struct {
	TestCustom  string
	Test2Custom string
}

func (tcc testCustomClaimsWithoutTags) Validate(context.Context) error {
	return nil
}

func Test_ParseAccessToken(t *testing.T) {
	now := time.Now()
	tokenClaims := jwtclaims.AccessTokenClaims{
		Claims: jwt.Claims{
			ID:        "5f7a093e-9301-4fb5-9eeb-c7b529f16ce8",
			Subject:   "1",
			Issuer:    "testIssuer",
			Audience:  []string{"testAudience"},
			Expiry:    jwt.NewNumericDate(now.Add(time.Hour)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		RegisteredCustomClaims: jwtclaims.RegisteredCustomClaims{
			Scope: "test.read test.write",
		},
	}

	customClaims := testCustomClaims{
		TestCustom:  "testCustom",
		Test2Custom: "testCustom2",
	}

	testCases := []struct {
		name                 string
		keyToCreate          []byte
		keyToParse           []byte
		customClaimsFunc     func() jwtclaims.CustomClaims
		customClaims         interface{}
		expectedTokenClaims  jwtclaims.AccessTokenClaims
		expectedCustomClaims interface{}
		expectedError        error
	}{
		{
			name:                "successfully parse a token",
			keyToCreate:         []byte(validKey),
			keyToParse:          []byte(validKey),
			expectedTokenClaims: tokenClaims,
			customClaims:        nil,
		},
		{
			name:                 "successfully parse a token with custom claims",
			keyToCreate:          []byte(validKey),
			keyToParse:           []byte(validKey),
			expectedTokenClaims:  tokenClaims,
			customClaims:         customClaims,
			expectedCustomClaims: &customClaims,
			customClaimsFunc: func() jwtclaims.CustomClaims {
				return &testCustomClaims{}
			},
		},
		{
			name:                "successfully parse a token even if customClaims function returns nil",
			keyToCreate:         []byte(validKey),
			keyToParse:          []byte(validKey),
			expectedTokenClaims: tokenClaims,
			customClaimsFunc: func() jwtclaims.CustomClaims {
				return nil
			},
		},
		{
			name:                "throws an error if key is invalid",
			keyToCreate:         []byte(validKey),
			keyToParse:          []byte(invalidKey),
			expectedTokenClaims: tokenClaims,
			expectedError:       ErrParseTokenClaims,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			sig, err := jose.NewSigner(
				jose.SigningKey{Algorithm: jose.HS256, Key: tc.keyToCreate},
				(&jose.SignerOptions{}).WithType("JWT"))
			require.NoError(t, err)

			tokenBuilder := jwt.Signed(sig)
			tokenBuilder = tokenBuilder.Claims(tc.expectedTokenClaims)

			if tc.customClaims != nil {
				tokenBuilder = tokenBuilder.Claims(tc.customClaims)
			}

			token, err := tokenBuilder.Token()
			require.NoError(t, err)

			registeredClaims, customClaims, err := ParseAccessToken(token, tc.keyToParse, tc.customClaimsFunc)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)

				if !cmp.Equal(tc.expectedTokenClaims, registeredClaims) {
					t.Fatal(cmp.Diff(tc.expectedTokenClaims, registeredClaims))
				}

				if tc.expectedCustomClaims != nil {
					if !cmp.Equal(tc.expectedCustomClaims, customClaims) {
						t.Fatal(cmp.Diff(tc.expectedCustomClaims, customClaims))
					}
				}
			}
		})
	}
}

func Test_ParseRefreshToken(t *testing.T) {
	exp := jwt.NumericDate(1750793228)

	testCases := []struct {
		name                string
		tokenStr            string
		key                 []byte
		expectedTokenClaims jwtclaims.RefreshTokenClaims
		expectedError       error
	}{
		{
			name:     "successfully parse a token",
			tokenStr: validRefreshToken,
			key:      []byte(validKey),
			expectedTokenClaims: jwtclaims.RefreshTokenClaims{
				ID:     "cba8a9cb-7909-4d1c-9ec6-9ca0a42265fb",
				Expiry: &exp,
			},
		},
		{
			name:          "throws an error if key is invalid",
			tokenStr:      validRefreshToken,
			key:           []byte(invalidKey),
			expectedError: ErrParseTokenClaims,
		},
		{
			name:          "throws an error if token is empty",
			tokenStr:      "",
			key:           []byte(validKey),
			expectedError: ErrParseToken,
		},
		{
			name:          "throws an error if token is invalid",
			tokenStr:      invalidRefreshToken,
			key:           []byte(validKey),
			expectedError: ErrParseToken,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tokenClaims, err := ParseRefreshToken(tc.tokenStr, tc.key)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)

				if !cmp.Equal(tc.expectedTokenClaims, tokenClaims) {
					t.Fatal(cmp.Diff(tc.expectedTokenClaims, tokenClaims))
				}
			}
		})
	}
}
