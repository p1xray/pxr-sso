package models

import (
	"encoding/json"
)

type Flow struct {
	ID                  string `json:"id" required:"true"`
	ClientID            string `json:"client_id" required:"true"`
	ResponseType        string `json:"response_type" required:"true"`
	RedirectURI         string `json:"redirect_uri" required:"true"`
	CodeChallenge       string `json:"code_challenge" required:"true"`
	CodeChallengeMethod string `json:"code_challenge_method" required:"true"`
	State               string `json:"state" required:"true"`
	Scope               string `json:"scope" required:"true"`
	Audience            string `json:"audience" required:"true"`
	Username            string `json:"username"`
}

//goland:noinspection GoMixedReceiverTypes
func (f Flow) MarshalBinary() ([]byte, error) {
	return json.Marshal(f)
}

//goland:noinspection GoMixedReceiverTypes
func (f *Flow) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &f)
}
