package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	AppPort  string
	DBHost   string
	DBPort   string
	DBUser   string
	DBPass   string
	DBName   string
	LogLevel string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables.")
	}

	return &Config{
		AppPort:  getEnv("APP_PORT", "8080"),
		DBHost:   getEnv("DB_HOST", "localhost"),
		DBPort:   getEnv("DB_PORT", "5432"),
		DBUser:   getEnv("DB_USER", "postgres"),
		DBPass:   getEnv("DB_PASSWORD", "pass"),
		DBName:   getEnv("DB_NAME", "song"),
		LogLevel: getEnv("LOG_LEVEL", "debug"),
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
