package entity

import "errors"

var (
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrUserExists           = errors.New("user already exists")
	ErrGeneratePasswordHash = errors.New("error generating password hash")
	ErrRefreshTokenExpired  = errors.New("refresh token expired")
	ErrValidateSession      = errors.New("error validating session")
	ErrInvalidSession       = errors.New("invalid session")
	ErrSessionNotFound      = errors.New("session not found")
	ErrCreateSession        = errors.New("error creating session")
	ErrCreateTokens         = errors.New("error creating tokens")
	ErrCreateAccessToken    = errors.New("error creating access token")
	ErrCreateRefreshToken   = errors.New("error creating refresh token")
)
