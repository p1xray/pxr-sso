package jwtmiddleware

import (
	"context"
	"fmt"
	jwtclaims "github.com/p1xray/pxr-sso/pkg/jwt/claims"
	"net/http"
)

type JWTMiddleware struct {
	validateToken  ValidateToken
	errorHandler   ErrorHandler
	tokenExtractor TokenExtractor
}

type ValidateToken func(context.Context, string) (jwtclaims.ValidatedClaims, error)

type ContextKey struct{}

func New(validateToken ValidateToken) *JWTMiddleware {
	return &JWTMiddleware{
		validateToken:  validateToken,
		errorHandler:   DefaultErrorHandler,
		tokenExtractor: AuthHeaderTokenExtractor,
	}
}

// ParseJWT is the main JWTMiddleware function which performs the main logic.
func (m *JWTMiddleware) ParseJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := m.tokenExtractor(r)
		if err != nil {
			m.errorHandler(w, r, fmt.Errorf("error extracting token: %w", err))
			return
		}

		if token == "" {
			m.errorHandler(w, r, ErrJWTMissing)
			return
		}

		validatedToken, err := m.validateToken(r.Context(), token)
		if err != nil {
			m.errorHandler(w, r, err)
			return
		}

		r = r.Clone(context.WithValue(r.Context(), ContextKey{}, validatedToken))
		next.ServeHTTP(w, r)
	})
}
