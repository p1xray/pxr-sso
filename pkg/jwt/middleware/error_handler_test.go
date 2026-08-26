package middleware

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDefaultErrorHandler(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "successful handled a jwt missing error",
			err:            ErrJWTMissing,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"JWT is missing."}`,
		},
		{
			name:           "successful handled a jwt invalid error",
			err:            ErrJWTInvalid,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"message":"JWT is invalid."}`,
		},
		{
			name:           "successful handled a generic error",
			err:            assert.AnError,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":"Something went wrong while checking the JWT."}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/test", nil)

			DefaultErrorHandler(w, r, tc.err)

			// Check status code
			assert.Equal(t, tc.expectedStatus, w.Code)

			// Check Content-Type
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

			// Check response body
			body, err := io.ReadAll(w.Body)
			require.NoError(t, err)

			assert.Equal(t, tc.expectedBody, string(body))
		})
	}
}
