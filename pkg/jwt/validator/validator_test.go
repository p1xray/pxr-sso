package validator

import (
	"context"
	"errors"
	josejwt "github.com/go-jose/go-jose/v4/jwt"
	"github.com/p1xray/pxr-sso/pkg/jwt"
	"github.com/p1xray/pxr-sso/pkg/jwt/claims"
	"github.com/p1xray/pxr-sso/pkg/jwt/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	const (
		issuer    = "https://example.com/"
		audience  = "test_audience"
		algorithm = crypto.HS256
	)

	var keyFunc = func(context.Context) (any, error) {
		return []byte("secret_key"), nil
	}

	testCases := []struct {
		name          string
		opts          []Option
		expectedError error
	}{
		// successful tests:
		{
			name: "successful creation with all required options",
			opts: []Option{
				WithKeyFunc(keyFunc),
				WithAlgorithm(algorithm),
				WithIssuer(issuer),
				WithAudience(audience),
			},
		},
		{
			name: "successful creation with WithAudiences",
			opts: []Option{
				WithKeyFunc(keyFunc),
				WithAlgorithm(algorithm),
				WithIssuer(issuer),
				WithAudiences([]string{audience, "test_audience_2"}),
			},
		},
		{
			name: "successful creation with optional parameters",
			opts: []Option{
				WithKeyFunc(keyFunc),
				WithAlgorithm(algorithm),
				WithIssuer(issuer),
				WithAudience(audience),
				WithAllowedClockSkew(30 * time.Second),
			},
		},
		{
			name: "successful creation with multiple algorithms",
			opts: []Option{
				WithKeyFunc(keyFunc),
				WithAlgorithms([]crypto.SignatureAlgorithm{algorithm, crypto.RS256}),
				WithIssuer(issuer),
				WithAudience(audience),
			},
		},
		{
			name: "successful creation with multiple issuers",
			opts: []Option{
				WithKeyFunc(keyFunc),
				WithAlgorithm(algorithm),
				WithIssuers([]string{issuer, "https://other_example.com/"}),
				WithAudience(audience),
			},
		},
		{
			name: "successful creation with multiple audiences",
			opts: []Option{
				WithKeyFunc(keyFunc),
				WithAlgorithm(algorithm),
				WithIssuer(issuer),
				WithAudiences([]string{audience, "other_audience"}),
			},
		},

		// failed keyFunc tests:
		{
			name: "throws an error when the keyFunc is nil",
			opts: []Option{
				WithKeyFunc(nil),
				WithAlgorithm(algorithm),
				WithIssuer(issuer),
				WithAudience(audience),
			},
			expectedError: ErrEmptyKeyFunc,
		},
		{
			name: "throws an error when keyFunc is missing",
			opts: []Option{
				WithAlgorithm(algorithm),
				WithIssuer(issuer),
				WithAudience(audience),
			},
			expectedError: ErrKeyFuncMissing,
		},

		// failed signature algorithm tests:
		{
			name: "throws an error when the signature algorithm is empty",
			opts: []Option{
				WithKeyFunc(keyFunc),
				WithAlgorithms([]crypto.SignatureAlgorithm{}),
				WithIssuer(issuer),
				WithAudience(audience),
			},
			expectedError: ErrEmptyAlgorithms,
		},
		{
			name: "throws an error when the signature algorithm is unsupported",
			opts: []Option{
				WithKeyFunc(keyFunc),
				WithAlgorithm("unsupported"),
				WithIssuer(issuer),
				WithAudience(audience),
			},
			expectedError: ErrUnsupportedAlgorithm,
		},
		{
			name: "throws an error when algorithm is missing",
			opts: []Option{
				WithKeyFunc(keyFunc),
				WithIssuer(issuer),
				WithAudience(audience),
			},
			expectedError: ErrSignatureAlgorithmMissing,
		},
		{
			name: "throws an error when both algorithm options are set up",
			opts: []Option{
				WithKeyFunc(keyFunc),
				WithAlgorithm(algorithm),
				WithAlgorithms([]crypto.SignatureAlgorithm{algorithm}),
				WithIssuer(issuer),
				WithAudience(audience),
			},
			expectedError: ErrUseBothAlgorithmOption,
		},

		// failed issuer tests:
		{
			name: "throws an error when the issuer is empty",
			opts: []Option{
				WithKeyFunc(keyFunc),
				WithAlgorithm(algorithm),
				WithIssuer(""),
				WithAudience(audience),
			},
			expectedError: ErrEmptyIssuer,
		},
		{
			name: "throws an error when the issuer URL is invalid",
			opts: []Option{
				WithKeyFunc(keyFunc),
				WithAlgorithm(algorithm),
				WithIssuer("htt$p://invalid issuer url"),
				WithAudience(audience),
			},
			expectedError: ErrInvalidIssuerURL,
		},
		{
			name: "throws an error when issuer is missing",
			opts: []Option{
				WithKeyFunc(keyFunc),
				WithAlgorithm(algorithm),
				WithAudience(audience),
			},
			expectedError: ErrIssuerMissing,
		},
		{
			name: "throws an error when the issuers list is empty",
			opts: []Option{
				WithKeyFunc(keyFunc),
				WithAlgorithm(algorithm),
				WithIssuers([]string{}),
				WithAudience(audience),
			},
			expectedError: ErrEmptyIssuers,
		},
		{
			name: "throws an error when the issuers list contains empty string",
			opts: []Option{
				WithKeyFunc(keyFunc),
				WithAlgorithm(algorithm),
				WithIssuers([]string{""}),
				WithAudience(audience),
			},
			expectedError: ErrEmptyIssuer,
		},
		{
			name: "throws an error when the issuers list contains invalid issuer URL",
			opts: []Option{
				WithKeyFunc(keyFunc),
				WithAlgorithm(algorithm),
				WithIssuers([]string{"htt$p://invalid issuer url"}),
				WithAudience(audience),
			},
			expectedError: ErrInvalidIssuerURL,
		},

		// failed audience tests:
		{
			name: "throws an error when the audience is empty",
			opts: []Option{
				WithKeyFunc(keyFunc),
				WithAlgorithm(algorithm),
				WithIssuer(issuer),
				WithAudience(""),
			},
			expectedError: ErrEmptyAudience,
		},
		{
			name: "throws an error when audience is missing",
			opts: []Option{
				WithKeyFunc(keyFunc),
				WithAlgorithm(algorithm),
				WithIssuer(issuer),
			},
			expectedError: ErrAudienceMissing,
		},
		{
			name: "throws an error when the audiences list is empty",
			opts: []Option{
				WithKeyFunc(keyFunc),
				WithAlgorithm(algorithm),
				WithIssuer(issuer),
				WithAudiences([]string{}),
			},
			expectedError: ErrEmptyAudiences,
		},
		{
			name: "throws an error when the audiences list contains empty string",
			opts: []Option{
				WithKeyFunc(keyFunc),
				WithAlgorithm(algorithm),
				WithIssuer(issuer),
				WithAudiences([]string{""}),
			},
			expectedError: ErrEmptyAudience,
		},

		// failed clock skew tests:
		{
			name: "throws an error when clock skew is negative",
			opts: []Option{
				WithKeyFunc(keyFunc),
				WithAlgorithm(algorithm),
				WithIssuer(issuer),
				WithAudience(audience),
				WithAllowedClockSkew(-30 * time.Second),
			},
			expectedError: ErrNegativeClockSkew,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			v, err := New(tc.opts...)
			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, v)
			}
		})
	}
}

func TestValidateToken(t *testing.T) {
	const (
		issuer    = "https://example.com/"
		audience  = "test_audience"
		algorithm = crypto.HS256

		validToken                    = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ0ZXN0X2F1ZGllbmNlIiwiZXhwIjozMzMwNTY0NDgwMCwiaWF0IjoxNzUwNDk3Mjk1LCJpc3MiOiJodHRwczovL2V4YW1wbGUuY29tLyIsImp0aSI6IjQ2M2FjMDJmLTAzOWMtNGRlZS1iNDhhLTExZThlMjQ5ZGVmOSIsIm5iZiI6MTc1MDQ5NzI5NSwic3ViIjoiMSJ9.9nu4HlbShynL16YJerbsRdr8-w2KecdncvlfvDbELdg"
		invalidIssuerToken            = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ0ZXN0X2F1ZGllbmNlIiwiZXhwIjozMzMwNTY0NDgwMCwiaWF0IjoxNzUwNDk3Mjk1LCJpc3MiOiJodHRwczovL290aGVyX2V4YW1wbGUuY29tLyIsImp0aSI6IjQ2M2FjMDJmLTAzOWMtNGRlZS1iNDhhLTExZThlMjQ5ZGVmOSIsIm5iZiI6MTc1MDQ5NzI5NSwic3ViIjoiMSJ9.Jb6gG6kuQ69MzEK5GMJDbuRS5oefi5iu5Ge0FULVSvg"
		invalidAudienceToken          = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJvdGhlcl9hdWRpZW5jZSIsImV4cCI6MzMzMDU2NDQ4MDAsImlhdCI6MTc1MDQ5NzI5NSwiaXNzIjoiaHR0cHM6Ly9leGFtcGxlLmNvbS8iLCJqdGkiOiI0NjNhYzAyZi0wMzljLTRkZWUtYjQ4YS0xMWU4ZTI0OWRlZjkiLCJuYmYiOjE3NTA0OTcyOTUsInN1YiI6IjEifQ.h63ONV-p5Ud1oV6bMakm_69Vq_4EnfkwXiY8X4PYyBc"
		invalidNotValidYetToken       = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ0ZXN0X2F1ZGllbmNlIiwiZXhwIjozMzMwNTY0NDgwMCwiaWF0IjoxNzUwNDk3Mjk1LCJpc3MiOiJodHRwczovL2V4YW1wbGUuY29tLyIsImp0aSI6IjQ2M2FjMDJmLTAzOWMtNGRlZS1iNDhhLTExZThlMjQ5ZGVmOSIsIm5iZiI6MzMzMDU2NDQ4MDAsInN1YiI6IjEifQ.iR5JIW184mOm0hTQdslnNBw2aBagJA59JAetggg5uHA"
		invalidExpiredToken           = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ0ZXN0X2F1ZGllbmNlIiwiZXhwIjoxNzUwNDk3Mjk1LCJpYXQiOjE3NTA0OTcyOTUsImlzcyI6Imh0dHBzOi8vZXhhbXBsZS5jb20vIiwianRpIjoiNDYzYWMwMmYtMDM5Yy00ZGVlLWI0OGEtMTFlOGUyNDlkZWY5IiwibmJmIjoxNzUwNDk3Mjk1LCJzdWIiOiIxIn0.WxVkBGfth_follr0Cy6_5LuGnXoR5fgsSe9CxxpaZjk"
		invalidIssuedInTheFutureToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ0ZXN0X2F1ZGllbmNlIiwiZXhwIjozMzMwNTY0NDgwMCwiaWF0IjozMzMwNTY0NDgwMCwiaXNzIjoiaHR0cHM6Ly9leGFtcGxlLmNvbS8iLCJqdGkiOiI0NjNhYzAyZi0wMzljLTRkZWUtYjQ4YS0xMWU4ZTI0OWRlZjkiLCJuYmYiOjE3NTA0OTcyOTUsInN1YiI6IjEifQ.HD8oblJj91TiihhlrmKCdlNJeVnTiJXyboL42YJ5noA"
	)

	var validKeyFunc = func(context.Context) (any, error) {
		return []byte("05c5328f-17cb-4b42-a085-4089c03b86f8"), nil
	}

	var invalidKeyFunc = func(context.Context) (any, error) {
		return []byte("fae35d9e-3696-499d-ae4a-9786b4273e68"), nil
	}

	var keyFuncReturnsNil = func(context.Context) (any, error) {
		return nil, errors.New("failed to fetch key")
	}

	testCases := []struct {
		name           string
		token          string
		opts           []Option
		expectedClaims claims.RegisteredClaims
		expectedError  error
	}{
		{
			name:  "successfully validates the token",
			token: validToken,
			opts: []Option{
				WithKeyFunc(validKeyFunc),
				WithAlgorithm(algorithm),
				WithIssuer(issuer),
				WithAudience(audience),
			},
			expectedClaims: claims.RegisteredClaims{
				ID:        "463ac02f-039c-4dee-b48a-11e8e249def9",
				Issuer:    issuer,
				Subject:   "1",
				Audience:  []string{audience},
				Expiry:    josejwt.NewNumericDate(time.Unix(33305644800, 0)),
				IssuedAt:  josejwt.NewNumericDate(time.Unix(1750497295, 0)),
				NotBefore: josejwt.NewNumericDate(time.Unix(1750497295, 0)),
			},
		},
		{
			name:  "throws an error when token is empty",
			token: "",
			opts: []Option{
				WithKeyFunc(validKeyFunc),
				WithAlgorithm(algorithm),
				WithIssuer(issuer),
				WithAudience(audience),
			},
			expectedError: ErrEmptyToken,
		},
		{
			name:  "throws an error when token has a different signing algorithm than the validator",
			token: validToken,
			opts: []Option{
				WithKeyFunc(validKeyFunc),
				WithAlgorithm(crypto.RS256),
				WithIssuer(issuer),
				WithAudience(audience),
			},
			expectedError: jwt.ErrParseToken,
		},
		{
			name:  "throws an error when it fails to fetch the keys from the key func",
			token: validToken,
			opts: []Option{
				WithKeyFunc(keyFuncReturnsNil),
				WithAlgorithm(algorithm),
				WithIssuer(issuer),
				WithAudience(audience),
			},
			expectedError: ErrKeyFunc,
		},
		{
			name:  "throws an error when parsing the token by invalid key",
			token: validToken,
			opts: []Option{
				WithKeyFunc(invalidKeyFunc),
				WithAlgorithm(algorithm),
				WithIssuer(issuer),
				WithAudience(audience),
			},
			expectedError: jwt.ErrParseClaims,
		},
		{
			name:  "throws an error when token issuer is invalid",
			token: invalidIssuerToken,
			opts: []Option{
				WithKeyFunc(validKeyFunc),
				WithAlgorithm(algorithm),
				WithIssuer(issuer),
				WithAudience(audience),
			},
			expectedError: ErrInvalidIssuerClaim,
		},
		{
			name:  "throws an error when token audience is invalid",
			token: invalidAudienceToken,
			opts: []Option{
				WithKeyFunc(validKeyFunc),
				WithAlgorithm(algorithm),
				WithIssuer(issuer),
				WithAudience(audience),
			},
			expectedError: ErrInvalidAudienceClaim,
		},
		{
			name:  "throws an error when token is not valid yet",
			token: invalidNotValidYetToken,
			opts: []Option{
				WithKeyFunc(validKeyFunc),
				WithAlgorithm(algorithm),
				WithIssuer(issuer),
				WithAudience(audience),
			},
			expectedError: ErrNotValidYet,
		},
		{
			name:  "throws an error when token is expired",
			token: invalidExpiredToken,
			opts: []Option{
				WithKeyFunc(validKeyFunc),
				WithAlgorithm(algorithm),
				WithIssuer(issuer),
				WithAudience(audience),
			},
			expectedError: ErrExpired,
		},
		{
			name:  "throws an error when token is issued in the future",
			token: invalidIssuedInTheFutureToken,
			opts: []Option{
				WithKeyFunc(validKeyFunc),
				WithAlgorithm(algorithm),
				WithIssuer(issuer),
				WithAudience(audience),
			},
			expectedError: ErrIssuedInTheFuture,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			v, err := New(tc.opts...)
			require.NoError(t, err)

			tokenClaims, err := v.ValidateToken(context.Background(), tc.token)
			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
				assert.Exactly(t, tc.expectedClaims, tokenClaims)
			}
		})
	}
}
