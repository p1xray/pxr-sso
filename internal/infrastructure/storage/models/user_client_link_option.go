package models

import "time"

type UserClientLinkOption func(*UserClientLink)

func UserClientLinkCreated() UserClientLinkOption {
	now := time.Now()
	return func(ucl *UserClientLink) {
		ucl.CreatedAt = now
		ucl.UpdatedAt = now
	}
}

func UserClientLinkUpdated() UserClientLinkOption {
	return func(ucl *UserClientLink) {
		ucl.UpdatedAt = time.Now()
	}
}
