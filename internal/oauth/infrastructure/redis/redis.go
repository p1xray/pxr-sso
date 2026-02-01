package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/p1xray/pxr-sso/internal/oauth"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/dto"
)

type Redis struct {
	data map[string][]byte
}

func New() *Redis {
	return &Redis{
		data: make(map[string][]byte),
	}
}

func (r *Redis) SaveFlow(ctx context.Context, flow dto.Flow) error {
	redisFlow := Flow{
		ID:                  flow.ID().String(),
		ClientID:            flow.ClientID(),
		RedirectURI:         flow.RedirectURI(),
		CodeChallenge:       flow.CodeChallenge(),
		CodeChallengeMethod: flow.CodeChallengeMethod(),
		State:               flow.State(),
		AuthorizationCode:   flow.AuthorizationCode(),
	}

	redisFlowKey := redisFlow.RedisKey()

	redisFlowBinary, err := redisFlow.MarshalBinary()
	if err != nil {
		return err
	}

	r.data[redisFlowKey] = redisFlowBinary
	return nil
}

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
	key := fmt.Sprintf("%s:%s", oauth.RedisObjectTypeNameFlow, f.ID)
	return key
}

func (f *Flow) MarshalBinary() ([]byte, error) {
	return json.Marshal(f)
}

func (f *Flow) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &f)
}
