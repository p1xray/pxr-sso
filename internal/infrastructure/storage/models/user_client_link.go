package models

import "time"

type UserClientLink struct {
	ID        int64
	UserID    int64
	ClientID  int64
	CreatedAt time.Time
	UpdatedAt time.Time
}
