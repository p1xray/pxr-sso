package storage

import (
	"errors"
)

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrClientNotFound       = errors.New("client not found")
	ErrSessionExists        = errors.New("session already exists")
	ErrUserClientLinkExists = errors.New("user's client link already exists")
	ErrSessionNotFound      = errors.New("session not found")
)
