package models

import (
	"encoding/json"
)

type AuthorizationRequest struct {
	ResponseType        string        `json:"response_type" required:"true"`
	Prompt              string        `json:"prompt" required:"true"`
	ClientID            string        `json:"client_id" required:"true"`
	RedirectURI         string        `json:"redirect_uri" required:"true"`
	CodeChallenge       string        `json:"code_challenge" required:"true"`
	CodeChallengeMethod string        `json:"code_challenge_method" required:"true"`
	State               string        `json:"state" required:"true"`
	Audience            string        `json:"audience" required:"true"`
	Scopes              ScopesRequest `json:"scopes" required:"true"`
}

//goland:noinspection GoMixedReceiverTypes
func (a AuthorizationRequest) MarshalBinary() ([]byte, error) {
	return json.Marshal(a)
}

//goland:noinspection GoMixedReceiverTypes
func (a *AuthorizationRequest) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &a)
}

type ScopesRequest struct {
	Pending []string `json:"pending" required:"true"`
	Granted []string `json:"granted" required:"true"`
}

//goland:noinspection GoMixedReceiverTypes
func (s ScopesRequest) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

//goland:noinspection GoMixedReceiverTypes
func (s *ScopesRequest) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &s)
}
