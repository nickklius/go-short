package utils

import "math/rand"

const (
	keyLength = 5
	letters   = "0123456789abcdefghijklmnopqrstuvwxyz"
)

func GenerateKey() string {
	b := make([]byte, keyLength)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
