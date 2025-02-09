package config

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	Domain        string
	ServerAddress string
}

func New() *Config {

	return &Config{
		Domain:        os.Getenv("DOMAIN"),
		ServerAddress: os.Getenv("PORT"),
	}
}
