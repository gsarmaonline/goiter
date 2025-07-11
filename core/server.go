package core

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gsarmaonline/goiter/config"
	"github.com/gsarmaonline/goiter/core/handlers"
	"github.com/gsarmaonline/goiter/core/models"
)

type (
	ModeT  string
	Server struct {
		router  *gin.Engine
		dbMgr   *models.DbManager
		handler *handlers.Handler

		cfg *config.Config
	}
)

func NewServer(cfg *config.Config) *Server {
	// Initialize router
	router := gin.Default()

	// Initialize config
	if cfg == nil {
		cfg = config.DefaultConfig()
	}

	gin.SetMode(cfg.GinMode)

	// Get CORS origin from environment variable
	corsOrigin := os.Getenv("FRONTEND_URL")

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{corsOrigin},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Initialize database
	dbMgr, err := models.NewDbManager(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Create server instance
	server := &Server{
		router:  router,
		dbMgr:   dbMgr,
		handler: handlers.NewHandler(router, dbMgr.Db, cfg),
		cfg:     cfg,
	}

	return server
}

func (s *Server) Start() error {
	return s.router.Run(fmt.Sprintf(":%s", s.cfg.Port))
}
