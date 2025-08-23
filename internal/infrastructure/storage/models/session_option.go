package models

import "time"

type SessionOption func(*Session)

func SessionCreated() SessionOption {
	now := time.Now()
	return func(s *Session) {
		s.CreatedAt = now
		s.UpdatedAt = now
	}
}

func SessionUpdated() SessionOption {
	return func(s *Session) {
		s.UpdatedAt = time.Now()
	}
}
