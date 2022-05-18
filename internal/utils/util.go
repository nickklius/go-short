package utils

import (
	"math/rand"
	"time"
)

func GenerateKey(letters string, keyLength int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, keyLength)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
