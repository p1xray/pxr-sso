package models

import "time"

// UserClientLink is data for user client link in storage.
type UserClientLink struct {
	ID        int64
	UserID    int64
	ClientID  int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserClientLinkOption func(*UserClientLink)

func UserClientLinkCreated() UserClientLinkOption {
	now := time.Now()
	return func(ucl *UserClientLink) {
		ucl.CreatedAt = now
		ucl.UpdatedAt = now
	}
}
