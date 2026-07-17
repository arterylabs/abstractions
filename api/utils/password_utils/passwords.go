package password_utils

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

// BcryptCost is the work factor used when hashing passwords.
// OWASP recommends >= 10; 12 is a reasonable balance of security and latency.
const BcryptCost = 12

func GeneratePasswordHash(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, BcryptCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func ComparePasswords(hashedPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
