package middleware

import (
	"context"
	"github.com/go-jose/go-jose/v4/jwt"
	"github.com/google/go-cmp/cmp"
	"github.com/p1xray/pxr-sso/pkg/jwt/claims"
	"github.com/p1xray/pxr-sso/pkg/jwt/crypto"
	"github.com/p1xray/pxr-sso/pkg/jwt/validator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func Test_ParseJWT(t *testing.T) {
	const (
		validKey = "05c5328f-17cb-4b42-a085-4089c03b86f8"
		issuer   = "https://example.com/"
		audience = "testAudience"
		subject  = "1"

		validToken   = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ0ZXN0QXVkaWVuY2UiLCJleHAiOjMzMzA1NjQ0ODAwLCJpYXQiOjE3NTA0OTcyOTUsImlzcyI6Imh0dHBzOi8vZXhhbXBsZS5jb20vIiwianRpIjoiYTBjNDNhYjMtMzA2Ny00MjExLTgwODYtZjZjN2YzMDA5YTgyIiwibmJmIjoxNzUwNDk3Mjk1LCJzdWIiOiIxIn0.C6WSp4EaynQBbdgLEJ1hpoHxwEoGAW8RYnCa1YpwNfE"
		invalidToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhd"
	)

	var (
		expiry    = jwt.NewNumericDate(time.Unix(33305644800, 0))
		issuedAt  = jwt.NewNumericDate(time.Unix(1750497295, 0))
		notBefore = jwt.NewNumericDate(time.Unix(1750497295, 0))
	)

	validTokenClaims := claims.RegisteredClaims{
		ID:        "a0c43ab3-3067-4211-8086-f6c7f3009a82",
		Issuer:    issuer,
		Subject:   subject,
		Audience:  []string{audience},
		Expiry:    expiry,
		IssuedAt:  issuedAt,
		NotBefore: notBefore,
	}

	keyFunc := func(context.Context) (any, error) {
		return []byte(validKey), nil
	}

	jwtValidator, err := validator.New(
		validator.WithKeyFunc(keyFunc),
		validator.WithIssuer(issuer),
		validator.WithAudience(audience),
		validator.WithAlgorithm(crypto.HS256))
	require.NoError(t, err)

	testCases := []struct {
		name                string
		method              string
		path                string
		token               string
		expectedStatusCode  int
		expectedTokenClaims any
		expectedBody        string
	}{
		{
			name:                "successfully validate a token",
			method:              http.MethodGet,
			token:               validToken,
			expectedStatusCode:  http.StatusOK,
			expectedTokenClaims: validTokenClaims,
			expectedBody:        `{"message":"Authenticated."}`,
		},
		{
			name:                "successfully validate a token on options method",
			method:              http.MethodOptions,
			token:               validToken,
			expectedStatusCode:  http.StatusOK,
			expectedTokenClaims: validTokenClaims,
			expectedBody:        `{"message":"Authenticated."}`,
		},
		{
			name:               "fails to validate a token with an invalid format",
			token:              "invalid token",
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       `{"message":"Something went wrong while checking the JWT."}`,
		},
		{
			name:               "fails to validate an empty token",
			token:              "",
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"message":"JWT is missing."}`,
		},
		{
			name:               "fails to validate an invalid token",
			token:              invalidToken,
			expectedStatusCode: http.StatusUnauthorized,
			expectedBody:       `{"message":"JWT is invalid."}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			middleware := New(jwtValidator.ValidateToken)

			var tokenClaims interface{}
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				tokenClaims = r.Context().Value(ContextKey{})

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"message":"Authenticated."}`))
			})

			testServer := httptest.NewServer(middleware.CheckJWT(handler))
			defer testServer.Close()

			url := testServer.URL + tc.path
			request, err := http.NewRequest(tc.method, url, nil)
			require.NoError(t, err)

			if tc.token != "" {
				request.Header.Add("Authorization", "Bearer "+tc.token)
			}

			response, err := testServer.Client().Do(request)
			require.NoError(t, err)

			body, err := io.ReadAll(response.Body)
			require.NoError(t, err)
			defer response.Body.Close()

			assert.Equal(t, tc.expectedStatusCode, response.StatusCode)
			assert.Equal(t, "application/json", response.Header.Get("Content-Type"))
			assert.Equal(t, tc.expectedBody, string(body))

			if !cmp.Equal(tc.expectedTokenClaims, tokenClaims) {
				t.Fatal(cmp.Diff(tc.expectedTokenClaims, tokenClaims))
			}
		})
	}
}
