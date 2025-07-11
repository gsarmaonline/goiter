package main

import (
	"log"
	"os"

	"github.com/gsarmaonline/goiter/config"
	"github.com/gsarmaonline/goiter/core"
	"github.com/joho/godotenv"
)

func main() {
	// Try to load .env file, but don't fail if it doesn't exist
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found or error loading it: %v", err)
	}

	// Verify environment variables are loaded
	requiredEnvVars := []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME"}
	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			log.Fatalf("Required environment variable %s is not set", envVar)
		}
	}

	server := core.NewServer(config.DefaultConfig())

	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
