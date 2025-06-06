package jwtmiddleware

import (
	"errors"
	"net/http"
)

var (
	// ErrJWTMissing is returned when the JWT is missing.
	ErrJWTMissing = errors.New("jwt missing")

	// ErrJWTInvalid is returned when the JWT is invalid.
	ErrJWTInvalid = errors.New("jwt invalid")
)

// ErrorHandler is a handler which is called when an error occurs in the JWTMiddleware.
type ErrorHandler func(w http.ResponseWriter, r *http.Request, err error)

// DefaultErrorHandler is the default error handler implementation for the JWTMiddleware.
func DefaultErrorHandler(w http.ResponseWriter, _ *http.Request, err error) {
	w.Header().Set("Content-Type", "application/json")

	switch {
	case errors.Is(err, ErrJWTMissing):
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"message":"JWT is missing."}`))
	case errors.Is(err, ErrJWTInvalid):
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"message":"JWT is invalid."}`))
	default:
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"message":"Something went wrong while checking the JWT."}`))
	}
}
