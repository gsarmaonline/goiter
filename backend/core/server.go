package core

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gsarmaonline/goiter/core/handlers"
	"github.com/gsarmaonline/goiter/core/models"
)

type (
	Server struct {
		router  *gin.Engine
		dbMgr   *models.DbManager
		handler *handlers.Handler
	}
)

func NewServer() *Server {
	// Initialize router
	router := gin.Default()

	// Get CORS origin from environment variable
	corsOrigin := os.Getenv("FRONTEND_URL")
	if corsOrigin == "" {
		corsOrigin = "http://localhost:3000" // Default to localhost for development
	}
	log.Printf("CORS origin: %s", corsOrigin)

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{corsOrigin},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Initialize database
	dbMgr, err := models.NewDbManager()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Create server instance
	server := &Server{
		router:  router,
		dbMgr:   dbMgr,
		handler: handlers.NewHandler(router, dbMgr.Db),
	}

	return server
}

func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}
