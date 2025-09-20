package testsuite

import (
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

func (c *GoiterClient) NewServer() (server *core.Server) {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found or error loading it: %v", err)
	}

	server = core.NewServer(c.getServerConfig())
	return
}
