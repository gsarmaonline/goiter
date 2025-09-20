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
		Router  *gin.Engine
		DbMgr   *models.DbManager
		Handler *handlers.Handler

		Cfg *config.Config
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

	if dbMgr.Db == nil {
		log.Fatalf("Database connection is nil")
	}

	// Create server instance
	server := &Server{
		Router:  router,
		DbMgr:   dbMgr,
		Handler: handlers.NewHandler(router, dbMgr.Db, cfg),
		Cfg:     cfg,
	}

	return server
}

func (s *Server) Start() (err error) {
	if err = s.DbMgr.Migrate(); err != nil {
		return
	}
	return s.Router.Run(fmt.Sprintf(":%s", s.Cfg.Port))
}
