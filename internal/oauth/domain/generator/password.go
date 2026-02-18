package generator

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func PasswordHash(password string) (string, error) {
	const op = "password hash"

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	return string(passwordHash), nil
}
