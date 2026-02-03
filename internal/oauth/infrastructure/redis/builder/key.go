package builder

import (
	"fmt"
)

// RedisObjectTypeNameFlow is a redis object type name for flow.
const redisObjectTypeNameFlow = "flow"

func BuildRedisFlowKey(id string) string {
	key := fmt.Sprintf("%s:%s", redisObjectTypeNameFlow, id)
	return key
}
