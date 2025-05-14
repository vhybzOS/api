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
	// Azure OpenAI Configuration
	AzureOpenAIEndpoint          string
	AzureOpenAIKey               string
	AzureOpenAIDeployment        string
	AzureOpenAIDeploymentVersion string
	// Gemini API Configuration
	GeminiAPIKey string
}

var AppConfig Config

func LoadConfig() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
		if err := godotenv.Load("../.env"); err != nil {
			log.Println("No .env file found, using environment variables")
			os.Exit(1)
		}
	}

	// Set default values
	AppConfig = Config{
		Port:                         getEnv("PORT", "8080"),
		JWTSecret:                    getEnv("JWT_SECRET", "your-default-secret-key"),
		DBPath:                       getEnv("DB_PATH", "app.db"),
		AzureOpenAIEndpoint:          getEnv("AZURE_OPENAI_ENDPOINT", ""),
		AzureOpenAIKey:               getEnv("AZURE_OPENAI_KEY", ""),
		AzureOpenAIDeployment:        getEnv("AZURE_OPENAI_DEPLOYMENT", "gpt-4o"),
		AzureOpenAIDeploymentVersion: getEnv("AZURE_OPENAI_DEPLOYMENT_VERSION", "gpt-4o"),
		GeminiAPIKey:                 getEnv("GEMINI_API_KEY", ""),
	}

	// Validate required configurations
	if AppConfig.JWTSecret == "your-default-secret-key" {
		log.Println("Warning: Using default JWT secret key. Please set JWT_SECRET in your environment variables.")
	}

	// Validate Azure OpenAI configuration
	if AppConfig.AzureOpenAIEndpoint == "" || AppConfig.AzureOpenAIKey == "" || AppConfig.AzureOpenAIDeployment == "" {
		log.Println("Warning: Azure OpenAI configuration is incomplete. Please set AZURE_OPENAI_ENDPOINT, AZURE_OPENAI_KEY, and AZURE_OPENAI_DEPLOYMENT in your environment variables.")
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
