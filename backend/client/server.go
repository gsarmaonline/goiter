package client

import (
	"fmt" // Added for printing
	"log"

	"github.com/gsarmaonline/goiter/config"
	"github.com/gsarmaonline/goiter/core"
	"github.com/joho/godotenv"
)

func StartServer() {
	fmt.Println("Attempting to start server...")

	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found or error loading it: %v", err)
	}
	cfg := config.DefaultConfig()
	cfg.Mode = config.ModeDev
	cfg.Port = "8090"
	cfg.DBType = config.SqliteDbType

	server := core.NewServer(cfg)

	log.Printf("Starting server on :%s", cfg.Port)

	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
