package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	KeyLength       int    `env:"KEY_LENGTH" envDefault:"5"`
	Letters         string `env:"LETTERS" envDefault:"0123456789abcdefghijklmnopqrstuvwxyz"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:":8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:""`
}

func NewConfig() Config {
	var c Config
	if err := env.Parse(&c); err != nil {
		log.Fatal(err)
	}
	return c
}

func (c *Config) ParseFlags() {
	flag.StringVar(&c.ServerAddress, "a", c.ServerAddress, "Server address")
	flag.StringVar(&c.BaseURL, "b", c.BaseURL, "Base URL")
	flag.StringVar(&c.FileStoragePath, "f", c.FileStoragePath, "File storage path")
	flag.Parse()
}
