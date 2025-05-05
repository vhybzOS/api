package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	JWTSecret string
	DBPath    string
}

var AppConfig Config

func LoadConfig() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Set default values
	AppConfig = Config{
		Port:      getEnv("PORT", "8080"),
		JWTSecret: getEnv("JWT_SECRET", "your-default-secret-key"),
		DBPath:    getEnv("DB_PATH", "app.db"),
	}

	// Validate required configurations
	if AppConfig.JWTSecret == "your-default-secret-key" {
		log.Println("Warning: Using default JWT secret key. Please set JWT_SECRET in your environment variables.")
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
