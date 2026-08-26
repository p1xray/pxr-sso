package middleware

import (
	"context"
	"fmt"
	"github.com/p1xray/pxr-sso/pkg/jwt/claims"
	"net/http"
)

// Middleware is a middleware that validates JWTs and makes claims available in
// the request context. It wraps the core validation engine and provides
// HTTP-specific functionality like token extraction and error handling.
type Middleware struct {
	validateToken  ValidateToken
	errorHandler   ErrorHandler
	tokenExtractor TokenExtractor
}

// ValidateToken defines the function signature responsible for validating a raw
// JWT string within a context and returning its verified core claims.
type ValidateToken func(context.Context, string) (claims.RegisteredClaims, error)

// ContextKey is a unique type used as a key to store and retrieve validated JWT
// claims within the request context, avoiding key collisions.
type ContextKey struct{}

// New returns a new Middleware instance with the provided token validation function.
func New(validateToken ValidateToken) *Middleware {
	return &Middleware{
		validateToken:  validateToken,
		errorHandler:   DefaultErrorHandler,
		tokenExtractor: AuthHeaderTokenExtractor,
	}
}

// CheckJWT is the main JWTMiddleware function which performs the main logic. It
// is passed an http.Handler which will be called if the JWT passes validation.
func (m *Middleware) CheckJWT(next http.Handler) http.Handler {
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
			m.errorHandler(w, r, fmt.Errorf("%w: %w", ErrJWTInvalid, err))
			return
		}

		r = r.Clone(context.WithValue(r.Context(), ContextKey{}, validatedToken))
		next.ServeHTTP(w, r)
	})
}
