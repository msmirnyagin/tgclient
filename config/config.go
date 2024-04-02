package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Url      string
	Bot      string
	AppId    int
	AppHash  string
	Phone    string
	Password string
}

func ReadEnv() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Print("Error loading .env file")
	}
	return &Config{
		Url:      getEnv("URL_WEBHOOK", ""),
		Bot:      getEnv("BOT_TOKEN", ""),
		AppId:    getEnvAsInt("APP_ID", 0),
		AppHash:  getEnv("APP_HASH", ""),
		Phone:    getEnv("TG_PHONE", ""),
		Password: getEnv("TG_PASSWORD", ""),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}
