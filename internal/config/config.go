package config

import (
	"github.com/caarlos0/env/v6"
	"log"
)

type Config struct {
	KeyLength       int    `env:"KEY_LENGTH" envDefault:"5"`
	Letters         string `env:"LETTERS" envDefault:"0123456789abcdefghijklmnopqrstuvwxyz"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:":8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:"storage.json"`
}

func New() Config {
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}
	return cfg
}
