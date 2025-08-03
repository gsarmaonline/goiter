package testsuite

import (
	"fmt"
	"log"

	"github.com/gsarmaonline/goiter/config"
	"github.com/gsarmaonline/goiter/core"
	"github.com/joho/godotenv"
)

func (c *GoiterClient) getServerConfig() (cfg *config.Config) {
	cfg = config.DefaultConfig()
	cfg.Mode = config.ModeDev
	cfg.Port = "8090"
	cfg.DBType = config.SqliteDbType
	return
}

func (c *GoiterClient) StartServer() {
	fmt.Println("Attempting to start server...")

	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found or error loading it: %v", err)
	}

	server := core.NewServer(c.getServerConfig())
	NewApp(server)

	log.Printf("Starting server on :%s", server.Cfg.Port)

	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
