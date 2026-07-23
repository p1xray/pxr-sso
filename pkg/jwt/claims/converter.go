package claims

import "github.com/go-jose/go-jose/v4/jwt"

// NumericDateToInt64 converts a *NumericDate pointer to a primitive int64.
// Returns 0 if the input pointer is nil.
func NumericDateToInt64(nd *jwt.NumericDate) int64 {
	if nd == nil {
		return 0
	}
	return int64(*nd)
}
