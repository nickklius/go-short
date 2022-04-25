package utils

import (
	"github.com/nickklius/go-short/internal/config"
	"math/rand"
)

func GenerateKey() string {
	b := make([]byte, config.KeyLength)
	for i := range b {
		b[i] = config.Letters[rand.Intn(len(config.Letters))]
	}
	return string(b)
}
