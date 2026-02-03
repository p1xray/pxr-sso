package models

import (
	"encoding/json"
	"github.com/p1xray/pxr-sso/internal/oauth/infrastructure/redis/builder"
)

type Flow struct {
	ID                  string `json:"id"`
	ClientID            string `json:"client_id" required:"true"`
	RedirectURI         string `json:"redirect_uri" required:"true"`
	CodeChallenge       string `json:"code_challenge" required:"true"`
	CodeChallengeMethod string `json:"code_challenge_method" required:"true"`
	State               string `json:"state" required:"true"`
	AuthorizationCode   string `json:"authorization_code"`
}

func (f *Flow) RedisKey() string {
	key := builder.BuildRedisFlowKey(f.ID)
	return key
}

func (f *Flow) MarshalBinary() ([]byte, error) {
	return json.Marshal(f)
}

func (f *Flow) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &f)
}
