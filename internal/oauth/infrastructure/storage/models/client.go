package models

import "time"

// Client is data for client in storage.
type Client struct {
	ID        int64
	Name      string
	Code      string
	SecretKey string
	Deleted   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
