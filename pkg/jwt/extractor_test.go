package jwtmiddleware

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func Test_AuthHeaderTokenExtractor(t *testing.T) {
	const token = "token-example"

	testCases := []struct {
		name          string
		request       *http.Request
		expectedToken string
		expectedError string
	}{
		{
			name:    "no header",
			request: &http.Request{},
		},
		{
			name: "valid authorization header value with token in header",
			request: &http.Request{
				Header: http.Header{
					"Authorization": []string{"Bearer " + token},
				},
			},
			expectedToken: token,
		},
		{
			name: "valid authorization header value with empty token in header",
			request: &http.Request{
				Header: http.Header{
					"Authorization": []string{"Bearer "},
				},
			},
			expectedError: ErrInvalidHeaderFormat.Error(),
		},
		{
			name: "invalid authorization header value with no bearer",
			request: &http.Request{
				Header: http.Header{
					"Authorization": []string{token},
				},
			},
			expectedError: ErrInvalidHeaderFormat.Error(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gotToken, err := AuthHeaderTokenExtractor(tc.request)
			if tc.expectedError != "" {
				assert.EqualError(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tc.expectedToken, gotToken)
		})
	}
}
