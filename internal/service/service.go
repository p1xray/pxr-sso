package service

import (
	"errors"
)

var (
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrUserExists          = errors.New("user already exists")
	ErrSessionNotFound     = errors.New("session not found")
	ErrRefreshTokenExpired = errors.New("refresh token expired")
	ErrInvalidSession      = errors.New("invalid session")
	ErrUserNotFound        = errors.New("user not found")
)
