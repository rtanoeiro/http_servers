package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashed_pass, errHash := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if errHash != nil {
		return "", errors.New("Unable to Hash password")
	}

	return string(hashed_pass), nil
}

func CheckPasswordHash(password, hashDB string) error {
	checkError := bcrypt.CompareHashAndPassword([]byte(hashDB), []byte(password))
	if checkError != nil {
		return errors.New("Passwords do not match")
	}
	return nil
}
