package sha256

import (
	"crypto/sha256"
	"encoding/base64"
)

func Base64URLHash(v string) string {
	sum := sha256.Sum256([]byte(v))
	hash := base64.RawURLEncoding.EncodeToString(sum[:])
	return hash
}
