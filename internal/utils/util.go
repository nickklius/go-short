package utils

import (
	"math/rand"
)

func GenerateKey(letters string, keyLength int) string {
	b := make([]byte, keyLength)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
