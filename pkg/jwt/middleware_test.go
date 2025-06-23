package jwtmiddleware

import (
	"context"
	"github.com/go-jose/go-jose/v4/jwt"
	"github.com/google/go-cmp/cmp"
	jwtclaims "github.com/p1xray/pxr-sso/pkg/jwt/claims"
	"github.com/p1xray/pxr-sso/pkg/jwt/validator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_ParseJWT(t *testing.T) {
	const (
		validToken   = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ0ZXN0IiwiZXhwIjoxNzUwNzExODQ0LCJpYXQiOjE3NTA3MDgyNDQsImlzcyI6Imh0dHA6Ly9sb2NhbGhvc3Q6NjAwNCIsImp0aSI6ImI2MWE0NjI2LWQ4MjItNDE0Yy04YWE1LTdiYjVmNjcwMGJhZSIsIm5iZiI6MTc1MDcwODI0NCwic2NvcGUiOiJwcm9maWxlLnJlYWQiLCJzdWIiOiIxIn0.F4noN66vHF5-jCFHMpta6ENobWeKnwFwy0kkoy5Ow1U"
		invalidToken = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJpbnZhbGlkVGVzdElzc3VlciIsImlhdCI6MTc1MDQ5MzQwNiwiZXhwIjoxNzgyMDI5NDA2LCJhdWQiOiIiLCJzdWIiOiIifQ.osd2XpJtwCHsGTJkCIu_yZKDG1TuGk9IyZi9mjdpe3A"
		issuer       = "http://localhost:6004"
		audience     = "test"
	)

	exp := jwt.NumericDate(1750711844)
	nbf := jwt.NumericDate(1750708244)
	iat := jwt.NumericDate(1750708244)

	expectedTokenClaims := jwtclaims.ValidatedClaims{
		RegisteredClaims: jwtclaims.AccessTokenClaims{
			Claims: jwt.Claims{
				ID:        "b61a4626-d822-414c-8aa5-7bb5f6700bae",
				Subject:   "1",
				Issuer:    issuer,
				Audience:  []string{audience},
				Expiry:    &exp,
				NotBefore: &nbf,
				IssuedAt:  &iat,
			},
			RegisteredCustomClaims: jwtclaims.RegisteredCustomClaims{
				Scope: "profile.read",
			},
		},
	}

	keyFunc := func(context.Context) ([]byte, error) {
		return []byte("98649a5c-2137-4a78-a63f-fbab416a7f9e"), nil
	}

	jwtValidator, err := validator.New(keyFunc, issuer, []string{audience})
	require.NoError(t, err)

	testCases := []struct {
		name                string
		validateToken       ValidateToken
		path                string
		expectedStatusCode  int
		token               string
		expectedTokenClaims interface{}
		expectedBody        string
	}{
		{
			name:                "successfully validate a token",
			validateToken:       jwtValidator.ValidateToken,
			expectedStatusCode:  http.StatusOK,
			token:               validToken,
			expectedTokenClaims: expectedTokenClaims,
			expectedBody:        `{"message":"Authenticated."}`,
		},
		{
			name:               "fails to validate a token with a invalid format",
			validateToken:      jwtValidator.ValidateToken,
			expectedStatusCode: http.StatusInternalServerError,
			token:              "invalid token",
			expectedBody:       `{"message":"Something went wrong while checking the JWT."}`,
		},
		{
			name:               "fails to validate an empty token",
			validateToken:      jwtValidator.ValidateToken,
			expectedStatusCode: http.StatusBadRequest,
			token:              "",
			expectedBody:       `{"message":"JWT is missing."}`,
		},
		{
			name:               "fails to validate an invalid token",
			validateToken:      jwtValidator.ValidateToken,
			expectedStatusCode: http.StatusUnauthorized,
			token:              invalidToken,
			expectedBody:       `{"message":"JWT is invalid."}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			middleware := New(tc.validateToken)

			var tokenClaims interface{}
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				tokenClaims = r.Context().Value(ContextKey{})

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"message":"Authenticated."}`))
			})

			testServer := httptest.NewServer(middleware.ParseJWT(handler))
			defer testServer.Close()

			url := testServer.URL + tc.path
			request, err := http.NewRequest(http.MethodGet, url, nil)
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
