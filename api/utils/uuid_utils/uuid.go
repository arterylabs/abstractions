package uuid_utils

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
)

func IsValidUUID(id string) bool {
	_, err := uuid.Parse(id)
	if err != nil {
		return false
	}
	return true
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandomShortID(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
