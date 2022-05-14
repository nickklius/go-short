package utils

import (
	"github.com/nickklius/go-short/internal/config"
	"math/rand"
)

func GenerateKey() string {
	c := config.New()

	b := make([]byte, c.KeyLength)
	for i := range b {
		b[i] = c.Letters[rand.Intn(len(c.Letters))]
	}
	return string(b)
}
