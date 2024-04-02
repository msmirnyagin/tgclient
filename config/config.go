package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Url string
	Bot string
}

func ReadEnv() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Print("Error loading .env file")
	}
	return &Config{
		Url: os.Getenv("URL_WEBHOOK"),
		Bot: os.Getenv("BOT_TOKEN"),
	}
}
