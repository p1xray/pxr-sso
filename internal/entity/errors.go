package entity

import "errors"

var (
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrUserExists           = errors.New("user already exists")
	ErrGeneratePasswordHash = errors.New("error generating password hash")
	ErrRefreshTokenExpired  = errors.New("refresh token expired")
	ErrInvalidSession       = errors.New("invalid session")
	ErrSessionNotFound      = errors.New("session not found")
)
