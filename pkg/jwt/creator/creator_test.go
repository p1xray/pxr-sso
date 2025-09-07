package jwtcreator

import (
	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	"time"
)

const (
	validKey   = "05c5328f-17cb-4b42-a085-4089c03b86f8"
	invalidKey = "invalid_key"
)

func Test_NewAccessToken(t *testing.T) {
	testCases := []struct {
		name          string
		data          AccessTokenCreateData
		expectedError error
	}{
		{
			name: "successfully creates a new token",
			data: AccessTokenCreateData{
				TTL: time.Duration(30) * time.Minute,
				Key: []byte(validKey),
			},
		},
		{
			name: "successfully creates a new token with data",
			data: AccessTokenCreateData{
				Subject:   "1",
				Audiences: []string{"testAudience"},
				Issuer:    "testIssuer",
				Scopes:    []string{"test.read", "test.write"},
				TTL:       time.Duration(30) * time.Minute,
				Key:       []byte(validKey),
			},
		},
		{
			name: "successfully creates a new token with custom claims",
			data: AccessTokenCreateData{
				Subject:   "1",
				Audiences: []string{"testAudience"},
				Issuer:    "testIssuer",
				Scopes:    []string{"test.read", "test.write"},
				CustomClaims: map[string]interface{}{
					"custom1": "value1",
					"custom2": "value2",
				},
				TTL: time.Duration(30) * time.Minute,
				Key: []byte(validKey),
			},
		},
		{
			name: "throws an error when creating a token signed by invalid key",
			data: AccessTokenCreateData{
				Key: []byte(invalidKey),
			},
			expectedError: ErrTokenSerialize,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tokenStr, err := NewAccessToken(tc.data)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)

				token, err := jwt.ParseSigned(tokenStr, []jose.SignatureAlgorithm{jose.HS256})
				require.NoError(t, err)

				claims := make(map[string]interface{})
				err = token.Claims(tc.data.Key, &claims)
				require.NoError(t, err)

				checkAccessTokenClaims(t, claims, tc.data)
			}
		})
	}
}

func Test_NewRefreshToken(t *testing.T) {
	testCases := []struct {
		name          string
		key           []byte
		ttl           time.Duration
		expectedError error
	}{
		{
			name: "successfully creates a new token",
			key:  []byte(validKey),
			ttl:  time.Duration(12) * time.Hour,
		},
		{
			name:          "throws an error when creating a token signed by invalid key",
			key:           []byte(invalidKey),
			ttl:           time.Duration(12) * time.Hour,
			expectedError: ErrTokenSerialize,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tokenStr, tokenID, err := NewRefreshToken(tc.key, tc.ttl)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)

				token, err := jwt.ParseSigned(tokenStr, []jose.SignatureAlgorithm{jose.HS256})
				require.NoError(t, err)

				claims := make(map[string]interface{})
				err = token.Claims(tc.key, &claims)
				require.NoError(t, err)

				checkRefreshTokenClaims(t, claims, tokenID)
			}
		})
	}
}

func checkAccessTokenClaims(t *testing.T, claims map[string]interface{}, expectedData AccessTokenCreateData) {
	checkJtiClaim(t, claims)
	checkSubClaim(t, claims, expectedData.Subject)
	checkIssClaim(t, claims, expectedData.Issuer)
	checkAudClaim(t, claims, expectedData.Audiences)
	checkExpClaim(t, claims)
	checkNbfClaim(t, claims)
	checkIatClaim(t, claims)
	checkScopeClaim(t, claims, expectedData.Scopes)
	checkCustomClaims(t, claims, expectedData.CustomClaims)
}

func checkRefreshTokenClaims(t *testing.T, claims map[string]interface{}, tokenID string) {
	checkJtiClaim(t, claims)
	checkExpClaim(t, claims)

	err := uuid.Validate(tokenID)
	assert.NoError(t, err)

	jti := parseJtiClaim(t, claims)
	assert.Equal(t, jti, tokenID)
}

func checkJtiClaim(t *testing.T, claims map[string]interface{}) {
	jti := parseJtiClaim(t, claims)

	err := uuid.Validate(jti)
	assert.NoError(t, err)
}

func checkSubClaim(t *testing.T, claims map[string]interface{}, expectedSubject string) {
	if expectedSubject == "" {
		return
	}

	sub, ok := claims["sub"]
	require.True(t, ok)

	assert.Equal(t, expectedSubject, sub.(string))
}

func checkIssClaim(t *testing.T, claims map[string]interface{}, expectedIssuer string) {
	if expectedIssuer == "" {
		return
	}

	iss, ok := claims["iss"]
	require.True(t, ok)

	assert.Equal(t, expectedIssuer, iss.(string))
}

func checkAudClaim(t *testing.T, claims map[string]interface{}, expectedAudiences []string) {
	if len(expectedAudiences) == 0 {
		return
	}

	aud, ok := claims["aud"]
	require.True(t, ok)

	audiences := make([]string, 0)
	switch audValue := aud.(type) {
	case string:
		audiences = []string{audValue}
	case []interface{}:
		audiences = make([]string, len(audValue))
		for i, v := range audValue {
			audValueStr, ok := v.(string)
			require.True(t, ok)
			audiences[i] = audValueStr
		}
	default:
		t.Error("invalid type of aud claim")
	}
	
	assert.Equal(t, expectedAudiences, audiences)
}

func checkExpClaim(t *testing.T, claims map[string]interface{}) {
	now := time.Now().UTC()
	leeway := time.Duration(60) * time.Second

	exp, ok := claims["exp"]
	require.True(t, ok)

	expTime := time.Unix(int64(exp.(float64)), 0).UTC()
	expValid := expTime.After(now.Add(-leeway))
	assert.True(t, expValid)
}

func checkNbfClaim(t *testing.T, claims map[string]interface{}) {
	now := time.Now().UTC()
	leeway := time.Duration(60) * time.Second

	nbf, ok := claims["nbf"]
	require.True(t, ok)

	nbfTime := time.Unix(int64(nbf.(float64)), 0).UTC()
	nbfValid := nbfTime.Before(now.Add(leeway))
	assert.True(t, nbfValid)
}

func checkIatClaim(t *testing.T, claims map[string]interface{}) {
	now := time.Now().UTC()
	leeway := time.Duration(60) * time.Second

	iat, ok := claims["iat"]
	require.True(t, ok)

	iatTime := time.Unix(int64(iat.(float64)), 0).UTC()
	iatValid := iatTime.Before(now.Add(leeway))
	assert.True(t, iatValid)
}

func checkScopeClaim(t *testing.T, claims map[string]interface{}, expectedScopes []string) {
	if len(expectedScopes) == 0 {
		return
	}

	scope, ok := claims["scope"]
	require.True(t, ok)

	assert.Equal(t, strings.Join(expectedScopes, " "), scope.(string))
}

func checkCustomClaims(t *testing.T, claims map[string]interface{}, expectedCustom map[string]interface{}) {
	if len(expectedCustom) == 0 {
		return
	}

	for k, v := range expectedCustom {
		value, ok := claims[k]
		assert.True(t, ok)

		assert.Equal(t, v, value)
	}
}

func parseJtiClaim(t *testing.T, claims map[string]interface{}) string {
	jti, ok := claims["jti"]
	require.True(t, ok)

	return jti.(string)
}
