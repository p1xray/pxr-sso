package models

import (
	"encoding/json"
	"time"
)

type AuthorizedGrant struct {
	Session Session              `json:"session" required:"true"`
	Request AuthorizationRequest `json:"request" required:"true"`
}

//goland:noinspection GoMixedReceiverTypes
func (a AuthorizedGrant) MarshalBinary() ([]byte, error) {
	return json.Marshal(a)
}

//goland:noinspection GoMixedReceiverTypes
func (a *AuthorizedGrant) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &a)
}

type Session struct {
	ID               string    `json:"id" required:"true"`
	Subject          int64     `json:"subject" required:"true"`
	AuthTime         time.Time `json:"auth_time" required:"true"`
	IdentityProvider string    `json:"identity_provider" required:"true"`
}

//goland:noinspection GoMixedReceiverTypes
func (s Session) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

//goland:noinspection GoMixedReceiverTypes
func (s *Session) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &s)
}
