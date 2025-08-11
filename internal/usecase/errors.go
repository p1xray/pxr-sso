package usecase

import "errors"

var (
	ErrClientNotFound     = errors.New("client not found")
	ErrSessionNotFound    = errors.New("session not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
)
