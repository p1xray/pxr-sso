package validator

import (
	"context"
	"errors"
	"github.com/go-jose/go-jose/v4/jwt"
	jwtclaims "github.com/p1xray/pxr-sso/pkg/jwt/claims"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	// JWT claim Issuer (iss).
	issuer = "testIssuer"

	// JWT claim Audience (aud).
	audience = "testAudience"
)

type testCustomClaims struct {
	TestCustom  string `json:"test_custom,omitempty"`
	Test2Custom string `json:"test2_custom,omitempty"`
	ReturnError error
}

func (tcc testCustomClaims) Validate(context.Context) error {
	return tcc.ReturnError
}

func Test_New(t *testing.T) {
	t.Run("throws an error when the keyFunc is nil", func(t *testing.T) {
		_, err := New(nil, issuer, []string{audience})
		assert.EqualError(t, err, ErrEmptyKeyFunc.Error())
	})

	t.Run("throws an error when the issuer is empty", func(t *testing.T) {
		_, err := New(validKeyFunc, "", []string{audience})
		assert.EqualError(t, err, ErrEmptyIssuer.Error())
	})

	t.Run("throws an error when the audience is nil", func(t *testing.T) {
		_, err := New(validKeyFunc, issuer, nil)
		assert.EqualError(t, err, ErrEmptyAudience.Error())
	})

	t.Run("throws an error when the audience is empty", func(t *testing.T) {
		_, err := New(validKeyFunc, issuer, []string{})
		assert.EqualError(t, err, ErrEmptyAudience.Error())
	})
}

func Test_ValidateToken(t *testing.T) {

	const (
		// JWT claim ID (jti).
		id = "0472dd7c-4821-4910-a6bc-29d4c5157f22"

		// JWT claim Subject (sub).
		subject = "1"

		// Валидный токен с заполненными стандартными клэймами.
		validDefaultToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ0ZXN0QXVkaWVuY2UiLCJleHAiOjMzMzA1NjQ0ODAwLCJpYXQiOjE3NTA0OTcyOTUsImlzcyI6InRlc3RJc3N1ZXIiLCJqdGkiOiIwNDcyZGQ3Yy00ODIxLTQ5MTAtYTZiYy0yOWQ0YzUxNTdmMjIiLCJuYmYiOjE3NTA0OTcyOTUsInN1YiI6IjEifQ.NoYm5P_1hlqC3e-C0fsdkvCX9TqKg8Wr1bYunQpNJsE"

		// Токен с другим алгоритмом подписи.
		tokenWithAnotherSignatureAlgorithm = "eyJhbGciOiJIUzM4NCIsInR5cCI6IkpXVCJ9.e30.CIcOiI5WfaY7BTYlGcMtE24fKDAqSvpnv8jCb9-3VYQuXlh4qh0ssPvu3QNycFuD"

		// Валидный минимальный токен, заполнены только Issuer и Audience.
		validMinimalToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ0ZXN0QXVkaWVuY2UiLCJpc3MiOiJ0ZXN0SXNzdWVyIn0.d2qDNajH-3dsEf1ZVbYepR6KgrNkaPOSfrKwH-c2tAE"

		// Токен с заполненными стандартными клэймами (Issuer невалидный).
		tokenWithInvalidIssuer = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ0ZXN0QXVkaWVuY2UiLCJpc3MiOiJpbnZhbGlkVGVzdElzc3VlciJ9.mKXEhugZ9DBwD8ELzCYBug3nQQbJZ2p4kj65SAx63Fw"

		// Токен с заполненными стандартными клэймами (Audience невалидный).
		tokenWithInvalidAudience = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJpbnZhbGlkVGVzdEF1ZGllbmNlIiwiaXNzIjoidGVzdElzc3VlciJ9.GcWshNfR-Fy9zAQ1Kh_tchxADIbofK58BQBvyrwZDpU"

		// Токен с заполненными стандартными клэймами (NotBefore невалидный).
		tokenWithInvalidNotBefore = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ0ZXN0QXVkaWVuY2UiLCJleHAiOjMzMzA1NjQ0ODAwLCJpYXQiOjE3NTA0OTcyOTUsImlzcyI6InRlc3RJc3N1ZXIiLCJuYmYiOjMzMzA1NjQ0ODAwfQ.BmMPHQhYhF5dBu4L1UB3ffkIQofwubeogr5JVPWDllI"

		// Токен с заполненными стандартными клэймами (Expiry невалидный).
		tokenWithInvalidExpiry = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ0ZXN0QXVkaWVuY2UiLCJleHAiOjE3NTA0OTcyOTUsImlhdCI6MTc1MDQ5NzI5NSwiaXNzIjoidGVzdElzc3VlciIsIm5iZiI6MTc1MDQ5NzI5NX0.PG-3-ZtvmXlNXid5ctW9Kg0mvaxI8WUrjYIK0ef3acA"

		// Токен с заполненными стандартными клэймами (IssuedAt невалидный).
		tokenWithInvalidIssuedAt = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ0ZXN0QXVkaWVuY2UiLCJleHAiOjMzMzA1NjQ0ODAwLCJpYXQiOjMzMzA1NjQ0ODAwLCJpc3MiOiJ0ZXN0SXNzdWVyIiwibmJmIjoxNzUwNDk3Mjk1fQ.sAOuHEUvLkP-GzLeLUOz961IVZxnUvmUvTZqL_yh0fE"

		// Токен с заполненными стандартными клэймами + заполнен Scope.
		tokenWithRegisteredCustomClaims = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ0ZXN0QXVkaWVuY2UiLCJleHAiOjMzMzA1NjQ0ODAwLCJpYXQiOjE3NTA0OTcyOTUsImlzcyI6InRlc3RJc3N1ZXIiLCJqdGkiOiIwNDcyZGQ3Yy00ODIxLTQ5MTAtYTZiYy0yOWQ0YzUxNTdmMjIiLCJuYmYiOjE3NTA0OTcyOTUsInNjb3BlIjoidGVzdC5yZWFkIHRlc3Qud3JpdGUiLCJzdWIiOiIxIn0.iih5bOFv6qcG8VrgcAxufA6AeGW62Pb_wOiEnnswkEc"

		// Токен с заполненными стандартными клэймами + добавлены кастомные клэймы.
		tokenWithCustomClaims = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ0ZXN0QXVkaWVuY2UiLCJleHAiOjMzMzA1NjQ0ODAwLCJpYXQiOjE3NTA0OTcyOTUsImlzcyI6InRlc3RJc3N1ZXIiLCJqdGkiOiIwNDcyZGQ3Yy00ODIxLTQ5MTAtYTZiYy0yOWQ0YzUxNTdmMjIiLCJuYmYiOjE3NTA0OTcyOTUsInNjb3BlIjoidGVzdC5yZWFkIHRlc3Qud3JpdGUiLCJzdWIiOiIxIiwidGVzdDJfY3VzdG9tIjoidGVzdCBjdXN0b20gY2xhaW0gMiIsInRlc3RfY3VzdG9tIjoidGVzdCBjdXN0b20gY2xhaW0ifQ.9O_WaXeDWacQuBmzWpbptBP6GV4Wz6aN4UbW5L6qVuY"
	)

	var (
		// Large numeric date value for Expiry, NotBefore or IssuedAt JWT claims (exp, nbf, iat).
		largeNumericDate = jwt.NumericDate(33305644800)

		// Before now time numeric date value for Expiry, NotBefore or IssuedAt JWT claims (exp, nbf, iat).
		beforeNowNumericDate = jwt.NumericDate(1750497295)

		// Token claims for default valid token.
		defaultTokenClaims = jwtclaims.ValidatedClaims{
			RegisteredClaims: jwtclaims.AccessTokenClaims{
				Claims: jwt.Claims{
					ID:        id,
					Subject:   subject,
					Issuer:    issuer,
					Audience:  []string{audience},
					Expiry:    &largeNumericDate,
					NotBefore: &beforeNowNumericDate,
					IssuedAt:  &beforeNowNumericDate,
				},
			},
		}
	)

	testCases := []struct {
		name                string
		token               string
		keyFunc             func(context.Context) ([]byte, error)
		customClaims        func() jwtclaims.CustomClaims
		expectedTokenClaims jwtclaims.ValidatedClaims
		expectedError       error
	}{
		{
			name:          "throws an error when token is empty",
			token:         "",
			keyFunc:       validKeyFunc,
			expectedError: ErrParsingToken,
		},
		{
			name:          "throws an error when token has a different signing algorithm than the validator",
			token:         tokenWithAnotherSignatureAlgorithm,
			keyFunc:       validKeyFunc,
			expectedError: ErrParsingToken,
		},
		{
			name:          "throws an error when it fails to fetch the keys from the key func",
			token:         validDefaultToken,
			keyFunc:       keyFuncReturnsNil,
			expectedError: ErrGettingKey,
		},
		{
			name:          "throws an error when parsing the token by invalid key",
			token:         validDefaultToken,
			keyFunc:       invalidKeyFunc,
			expectedError: ErrParsingToken,
		},
		{
			name:    "successfully validates a minimal token",
			token:   validMinimalToken,
			keyFunc: validKeyFunc,
			expectedTokenClaims: jwtclaims.ValidatedClaims{
				RegisteredClaims: jwtclaims.AccessTokenClaims{
					Claims: jwt.Claims{
						Issuer:   issuer,
						Audience: []string{audience},
					},
				},
			},
		},
		{
			name:                "successfully validates a default token",
			token:               validDefaultToken,
			keyFunc:             validKeyFunc,
			expectedTokenClaims: defaultTokenClaims,
		},
		{
			name:          "throws an error when token issuer is invalid",
			token:         tokenWithInvalidIssuer,
			keyFunc:       validKeyFunc,
			expectedError: jwt.ErrInvalidIssuer,
		},
		{
			name:          "throws an error when token audience is invalid",
			token:         tokenWithInvalidAudience,
			keyFunc:       validKeyFunc,
			expectedError: jwt.ErrInvalidAudience,
		},
		{
			name:          "throws an error when token is not valid yet",
			token:         tokenWithInvalidNotBefore,
			keyFunc:       validKeyFunc,
			expectedError: jwt.ErrNotValidYet,
		},
		{
			name:          "throws an error when token is expired",
			token:         tokenWithInvalidExpiry,
			keyFunc:       validKeyFunc,
			expectedError: jwt.ErrExpired,
		},
		{
			name:          "throws an error when token is issued in the future",
			token:         tokenWithInvalidIssuedAt,
			keyFunc:       validKeyFunc,
			expectedError: jwt.ErrIssuedInTheFuture,
		},
		{
			name:    "successfully validates a token with registered custom claims",
			token:   tokenWithRegisteredCustomClaims,
			keyFunc: validKeyFunc,
			expectedTokenClaims: jwtclaims.ValidatedClaims{
				RegisteredClaims: jwtclaims.AccessTokenClaims{
					Claims: jwt.Claims{
						ID:        id,
						Subject:   subject,
						Issuer:    issuer,
						Audience:  []string{audience},
						Expiry:    &largeNumericDate,
						NotBefore: &beforeNowNumericDate,
						IssuedAt:  &beforeNowNumericDate,
					},
					RegisteredCustomClaims: jwtclaims.RegisteredCustomClaims{
						Scope: "test.read test.write",
					},
				},
			},
		},
		{
			name:    "successfully validates a token even if customClaims function returns nil",
			token:   validDefaultToken,
			keyFunc: validKeyFunc,
			customClaims: func() jwtclaims.CustomClaims {
				return nil
			},
			expectedTokenClaims: defaultTokenClaims,
		},
		{
			name:    "successfully validates a token with custom claims",
			token:   tokenWithCustomClaims,
			keyFunc: validKeyFunc,
			customClaims: func() jwtclaims.CustomClaims {
				return &testCustomClaims{}
			},
			expectedTokenClaims: jwtclaims.ValidatedClaims{
				RegisteredClaims: jwtclaims.AccessTokenClaims{
					Claims: jwt.Claims{
						ID:        id,
						Subject:   subject,
						Issuer:    issuer,
						Audience:  []string{audience},
						Expiry:    &largeNumericDate,
						NotBefore: &beforeNowNumericDate,
						IssuedAt:  &beforeNowNumericDate,
					},
					RegisteredCustomClaims: jwtclaims.RegisteredCustomClaims{
						Scope: "test.read test.write",
					},
				},
				CustomClaims: &testCustomClaims{
					TestCustom:  "test custom claim",
					Test2Custom: "test custom claim 2",
				},
			},
		},
		{
			name:    "throws an error when it fails to validate the custom claims",
			token:   tokenWithCustomClaims,
			keyFunc: validKeyFunc,
			customClaims: func() jwtclaims.CustomClaims {
				return &testCustomClaims{
					ReturnError: errors.New("error validating custom claims"),
				}
			},
			expectedError: ErrValidatingCustomClaims,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			validator, err := New(
				tc.keyFunc,
				issuer,
				[]string{audience},
				WithCustomClaims(tc.customClaims))
			require.NoError(t, err)

			tokenClaims, err := validator.ValidateToken(context.Background(), tc.token)
			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
				assert.Exactly(t, tc.expectedTokenClaims, tokenClaims)
			}
		})
	}
}

func validKeyFunc(context.Context) ([]byte, error) {
	return []byte("05c5328f-17cb-4b42-a085-4089c03b86f8"), nil
}

func invalidKeyFunc(context.Context) ([]byte, error) {
	return []byte("fae35d9e-3696-499d-ae4a-9786b4273e68"), nil
}

func keyFuncReturnsNil(context.Context) ([]byte, error) {
	return nil, ErrGettingKey
}
