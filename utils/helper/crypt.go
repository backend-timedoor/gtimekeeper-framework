package helper

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func Hash(param string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(param), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password %w", err)
	}

	return string(hash), nil
}

func HashCheck(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
