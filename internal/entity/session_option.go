package entity

import (
	"fmt"
	"time"
)

// SessionOption is how options for the Session are set up.
type SessionOption func(*Session) error

// WithSessionID is an option which sets up the session ID for the user session entity.
func WithSessionID(id int64) SessionOption {
	return func(s *Session) error {
		s.ID = id

		return nil
	}
}

// WithSessionRefreshTokenID is an option which sets up the refresh token ID for the user session entity.
func WithSessionRefreshTokenID(refreshTokenID string) SessionOption {
	return func(s *Session) error {
		s.RefreshTokenID = refreshTokenID

		return nil
	}
}

// WithSessionExpiresAt is an option which sets up the time of expires session for the user session entity.
func WithSessionExpiresAt(expiresAt time.Time) SessionOption {
	return func(s *Session) error {
		s.ExpiresAt = expiresAt

		return nil
	}
}

// WithGeneratedTokens is an option which sets up the generated tokens for the user session entity.
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
