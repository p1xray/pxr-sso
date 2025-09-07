package entity

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const (
	invalidUserAgent   = "invalid user agent"
	invalidFingerprint = "invalid fingerprint"
)

func Test_Session_Validate(t *testing.T) {
	testCases := []struct {
		name                string
		userID              int64
		userAgent           string
		fingerprint         string
		expiresAt           time.Time
		expectedUserAgent   string
		expectedFingerprint string
		expectedError       error
	}{
		{
			name:                "successfully validating",
			userID:              userID,
			userAgent:           userAgent,
			fingerprint:         fingerprint,
			expiresAt:           time.Now().Add(time.Hour),
			expectedUserAgent:   userAgent,
			expectedFingerprint: fingerprint,
			expectedError:       nil,
		},
		{
			name:                "throws an error when userAgent is invalid",
			userID:              userID,
			userAgent:           userAgent,
			fingerprint:         fingerprint,
			expiresAt:           time.Now().Add(time.Hour),
			expectedUserAgent:   invalidUserAgent,
			expectedFingerprint: fingerprint,
			expectedError:       ErrInvalidSession,
		},
		{
			name:                "throws an error when fingerprint is invalid",
			userID:              userID,
			userAgent:           userAgent,
			fingerprint:         fingerprint,
			expiresAt:           time.Now().Add(time.Hour),
			expectedUserAgent:   userAgent,
			expectedFingerprint: invalidFingerprint,
			expectedError:       ErrInvalidSession,
		},
		{
			name:                "throws an error when userAgent and fingerprint is invalid",
			userID:              userID,
			userAgent:           userAgent,
			fingerprint:         fingerprint,
			expiresAt:           time.Now().Add(time.Hour),
			expectedUserAgent:   invalidUserAgent,
			expectedFingerprint: invalidFingerprint,
			expectedError:       ErrInvalidSession,
		},
		{
			name:                "throws an error when session is expired",
			userID:              userID,
			userAgent:           userAgent,
			fingerprint:         fingerprint,
			expiresAt:           time.Now().Add(-time.Hour),
			expectedUserAgent:   userAgent,
			expectedFingerprint: fingerprint,
			expectedError:       ErrRefreshTokenExpired,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			session, err := NewSession(tc.userID, tc.userAgent, tc.fingerprint, WithSessionExpiresAt(tc.expiresAt))
			require.NoError(t, err)

			err = session.Validate(tc.expectedUserAgent, tc.expectedFingerprint)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
