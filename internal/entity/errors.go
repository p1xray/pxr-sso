package entity

import "errors"

var (
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrUserExists           = errors.New("user already exists")
	ErrGeneratePasswordHash = errors.New("error generating password hash")
)
