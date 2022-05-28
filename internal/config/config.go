package config

import (
	"flag"
	"math"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	KeyLength         int    `env:"KEY_LENGTH" envDefault:"5"`
	Letters           string `env:"LETTERS" envDefault:"0123456789abcdefghijklmnopqrstuvwxyz"`
	ShortenerCapacity int    `env:"SHORTENER_CAPACITY" envDefault:"1"`
	BaseURL           string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	ServerAddress     string `env:"SERVER_ADDRESS" envDefault:":8080"`
	FileStoragePath   string `env:"FILE_STORAGE_PATH" envDefault:"storage.json"`
	DatabaseDSN       string `env:"DATABASE_DSN" envDefault:"user=goshort password=1dc3sfdf host=localhost port=5432 dbname=goshort"`
}

func NewConfig() (Config, error) {
	var c Config
	if err := env.Parse(&c); err != nil {
		return c, err
	}

	flag.StringVar(&c.ServerAddress, "a", c.ServerAddress, "Server address")
	flag.StringVar(&c.BaseURL, "b", c.BaseURL, "Base URL")
	flag.StringVar(&c.FileStoragePath, "f", c.FileStoragePath, "File storage path")
	flag.StringVar(&c.DatabaseDSN, "d", c.DatabaseDSN, "PG conn address")
	flag.Parse()

	c.ShortenerCapacity = int(math.Pow(float64(len(c.Letters)), float64(c.KeyLength)))
	return c, nil
}
