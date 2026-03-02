package models

import (
	"encoding/json"
)

type Authorization struct {
	Code          string `json:"code" required:"true"`
	Username      string `json:"username" required:"true"`
	ClientID      string `json:"client_id" required:"true"`
	RedirectURI   string `json:"redirect_uri" required:"true"`
	CodeChallenge string `json:"code_challenge" required:"true"`
	Scope         string `json:"scope" required:"true"`
}

//goland:noinspection GoMixedReceiverTypes
func (a Authorization) MarshalBinary() ([]byte, error) {
	return json.Marshal(a)
}

//goland:noinspection GoMixedReceiverTypes
func (a *Authorization) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &a)
}
