package models

import (
	"encoding/json"
)

type Flow struct {
	ID                  string `json:"id" redis:"id" required:"true"`
	ClientID            string `json:"client_id" redis:"client_id" required:"true"`
	RedirectURI         string `json:"redirect_uri" redis:"redirect_uri" required:"true"`
	CodeChallenge       string `json:"code_challenge" redis:"code_challenge" required:"true"`
	CodeChallengeMethod string `json:"code_challenge_method" redis:"code_challenge_method" required:"true"`
	State               string `json:"state" redis:"state" required:"true"`
	AuthorizationCode   string `json:"authorization_code" redis:"authorization_code"`
}

//goland:noinspection GoMixedReceiverTypes
func (f Flow) MarshalBinary() ([]byte, error) {
	return json.Marshal(f)
}

//goland:noinspection GoMixedReceiverTypes
func (f *Flow) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &f)
}
