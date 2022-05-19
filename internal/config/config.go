package config

import (
	"flag"
	"math"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	KeyLength         int    `env:"KEY_LENGTH" envDefault:"5"`
	Letters           string `env:"LETTERS" envDefault:"0123456789abcdefghijklmnopqrstuvwxyz"`
	ShortenerCapacity int    `env:"CAPACITY" envDefault:"10"`
	BaseURL           string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	ServerAddress     string `env:"SERVER_ADDRESS" envDefault:":8080"`
	FileStoragePath   string `env:"FILE_STORAGE_PATH" envDefault:"storage.json"`
}

func NewConfig() (Config, error) {
	var c Config
	if err := env.Parse(&c); err != nil {
		return c, err
	}

	flag.StringVar(&c.ServerAddress, "a", c.ServerAddress, "Server address")
	flag.StringVar(&c.BaseURL, "b", c.BaseURL, "Base URL")
	flag.StringVar(&c.FileStoragePath, "f", c.FileStoragePath, "File storage path")
	flag.Parse()

	c.ShortenerCapacity = int(math.Pow(float64(len(c.Letters)), float64(c.KeyLength)))

	return c, nil
}
