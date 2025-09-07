package entity

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_NewTokens(t *testing.T) {
	testCases := []struct {
		name          string
		data          CreateTokensParams
		expectedError error
	}{
		{
			name: "successfully create tokens",
			data: CreateTokensParams{
				UserID:          userID,
				SecretKey:       secretKey,
				AccessTokenTTL:  accessTokenTTL,
				RefreshTokenTTL: refreshTokenTTL,
			},
			expectedError: nil,
		},
		{
			name: "throws an error when secret key is empty",
			data: CreateTokensParams{
				UserID:          userID,
				SecretKey:       "",
				AccessTokenTTL:  accessTokenTTL,
				RefreshTokenTTL: refreshTokenTTL,
			},
			expectedError: ErrCreateAccessToken,
		},
		{
			name: "throws an error when secret key is invalid",
			data: CreateTokensParams{
				UserID:          userID,
				SecretKey:       invalidSecretKey,
				AccessTokenTTL:  accessTokenTTL,
				RefreshTokenTTL: refreshTokenTTL,
			},
			expectedError: ErrCreateAccessToken,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tokens, err := NewTokens(tc.data)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)

				assert.NotEmpty(t, tokens.AccessToken)
				assert.NotEmpty(t, tokens.RefreshToken)
				assert.NotEmpty(t, tokens.RefreshTokenID)
			}
		})
	}
}
