package generator

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func PasswordHash(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("%s: %w", "generate password hash", err)
	}

	return string(passwordHash), nil
}
