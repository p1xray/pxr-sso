package entity

import (
	"fmt"
	"time"
)

type Session struct {
	UserID         int64
	RefreshTokenID string
	UserAgent      string
	Fingerprint    string
	ExpiresAt      time.Time

	Tokens Tokens
}

func NewSession(data CreateNewSessionParams) (Session, error) {
	const op = "entity.NewSession"

	createTokensParams := CreateTokensParams{
		UserID:          data.UserID,
		Permissions:     data.UserPermissions,
		ClientCode:      data.ClientCode,
		SecretKey:       data.ClientSecretKey,
		Issuer:          data.Issuer,
		AccessTokenTTL:  data.AccessTokenTTL,
		RefreshTokenTTL: data.RefreshTokenTTL,
	}
	tokens, err := NewTokens(createTokensParams)
	if err != nil {
		return Session{}, fmt.Errorf("%s: %w", op, err)
	}

	return Session{
		UserID:         data.UserID,
		RefreshTokenID: tokens.RefreshTokenID,
		UserAgent:      data.UserAgent,
		Fingerprint:    data.Fingerprint,
		ExpiresAt:      time.Now().Add(data.RefreshTokenTTL),

		Tokens: tokens,
	}, nil
}
