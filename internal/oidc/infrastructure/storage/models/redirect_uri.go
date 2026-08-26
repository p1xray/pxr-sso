package models

import "time"

// RedirectURI is data for redirect URI in storage.
type RedirectURI struct {
	ID        int64
	ClientID  int64
	URI       string
	CreatedAt time.Time
	UpdatedAt time.Time
}
