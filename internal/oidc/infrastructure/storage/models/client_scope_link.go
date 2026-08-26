package models

import "time"

// ClientScopeLink is data for client scope link in storage.
type ClientScopeLink struct {
	ID        int64
	ClientID  int64
	ScopeID   int64
	CreatedAt time.Time
	UpdatedAt time.Time

	Scope Scope
}
