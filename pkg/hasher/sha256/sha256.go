package sha256

import (
	"crypto/sha256"
	"encoding/base64"
)

func Base64URLHash(v string) string {
	sum := sha256.Sum256([]byte(v))
	hash := base64.URLEncoding.EncodeToString(sum[:])
	return hash
}
