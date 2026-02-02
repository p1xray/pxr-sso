package models

import "time"

// Audience is data for audience in storage.
type Audience struct {
	ID        int64
	ClientID  int64
	URL       string
	CreatedAt time.Time
	UpdatedAt time.Time
}
