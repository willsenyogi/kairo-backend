package configs

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	AppEnv      string
	ServerPort  string
	AutoMigrate bool
}

func LoadAppConfig() *AppConfig {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	return &AppConfig{
		AppEnv:      getEnv("APP_ENV", "development"),
		ServerPort:  getEnv("APP_PORT", "8080"),
		AutoMigrate: getEnvAsBool("AUTOMIGRATE", false),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return fallback
}

func getEnvAsBool(key string, fallback bool) bool {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		log.Printf("Invalid boolean value for %s, using fallback: %t", key, fallback)
		return fallback
	}

	return parsed
}
