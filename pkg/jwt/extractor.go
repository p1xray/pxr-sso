package jwtmiddleware

import (
	"errors"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	bearerPrefix        = "Bearer"
)

var (
	ErrInvalidHeaderFormat = errors.New("authorization header format must be Bearer {token}")
)

// TokenExtractor is a function that takes a request as input and returns either a token or an error.
type TokenExtractor func(r *http.Request) (string, error)

// AuthHeaderTokenExtractor is a TokenExtractor that takes a request and extracts the token
// from the Authorization header.
func AuthHeaderTokenExtractor(r *http.Request) (string, error) {
	authHeader := r.Header.Get(authorizationHeader)
	if authHeader == "" {
		return "", nil
	}

	authHeaderParts := strings.Fields(authHeader)
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != bearerPrefix {
		return "", ErrInvalidHeaderFormat
	}

	return authHeaderParts[1], nil
}
