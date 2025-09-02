package entity

import (
	"fmt"
	"time"
)

type SessionOption func(*Session) error

func WithSessionID(id int64) SessionOption {
	return func(s *Session) error {
		s.ID = id

		return nil
	}
}

func WithSessionRefreshTokenID(refreshTokenID string) SessionOption {
	return func(s *Session) error {
		s.RefreshTokenID = refreshTokenID

		return nil
	}
}

func WithSessionExpiresAt(expiresAt time.Time) SessionOption {
	return func(s *Session) error {
		s.ExpiresAt = expiresAt

		return nil
	}
}

func WithGeneratedTokens(data SessionWithGeneratedTokensParams) SessionOption {
	return func(s *Session) error {
		createTokensParams := CreateTokensParams{
			UserID:          s.UserID,
			Permissions:     data.UserPermissions,
			Audiences:       data.Audiences,
			SecretKey:       data.ClientSecretKey,
			Issuer:          data.Issuer,
			AccessTokenTTL:  data.AccessTokenTTL,
			RefreshTokenTTL: data.RefreshTokenTTL,
		}
		tokens, err := NewTokens(createTokensParams)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrCreateTokens, err)
		}

		s.Tokens = tokens
		s.RefreshTokenID = tokens.RefreshTokenID
		s.ExpiresAt = time.Now().Add(data.RefreshTokenTTL)

		return nil
	}
}
