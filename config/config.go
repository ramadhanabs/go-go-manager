package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string
	DbHost  string
	DbPort  string
	DbUser  string
	DbPass  string
	DbName  string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	return &Config{
		AppPort: getEnv("APP_PORT", "8080"),
		DbHost:  getEnv("DB_HOST", "localhost"),
		DbPort:  getEnv("DB_PORT", "5432"),
		DbUser:  getEnv("DB_USER", "postgres"),
		DbPass:  getEnv("DB_PASS", "password"),
		DbName:  getEnv("DB_NAME", "mydb"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
