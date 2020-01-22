package user

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	hashedBytes, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(hashedBytes)
}
