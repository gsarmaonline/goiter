package handlers

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gsarmaonline/goiter/config"
	"github.com/gsarmaonline/goiter/core/helpers/authorisation"
	"github.com/gsarmaonline/goiter/core/middleware"
	"github.com/gsarmaonline/goiter/core/models"
	"gorm.io/gorm"
)

const (
	DefaultUrlKeyName = "id"

	// For FindWithUser query types
	NilQuery = ""
)

type (
	Handler struct {
		router        *gin.Engine
		Db            *gorm.DB
		middleware    *middleware.Middleware
		cfg           *config.Config
		authorisation *authorisation.Authorisation

		OpenRouteGroup      *gin.RouterGroup
		ProtectedRouteGroup *gin.RouterGroup
	}
)

func NewHandler(router *gin.Engine, db *gorm.DB, cfg *config.Config) (handler *Handler) {
	middleware := middleware.NewMiddleware(router, db)
	handler = &Handler{
		router:     router,
		Db:         db,
		middleware: middleware,
		cfg:        cfg,

		OpenRouteGroup:      router.Group("/"),
		ProtectedRouteGroup: router.Group(""),
	}
	// Setup routes
	handler.SetupRoutes()
	// Setup authorisation
	handler.authorisation = authorisation.NewAuthorisation(handler)
	return
}

func (h *Handler) SetAuthorisationState(status bool) {
	h.authorisation.IsEnabled = status
	return
}

func (h *Handler) SetupRoutes() {
	h.router.GET("/ping", h.handlePing)
	h.setupAuthRoutes()
}

func (h *Handler) setupAuthRoutes() {
	authOpenRoutes := h.router.Group("/auth")
	{
		// Public routes (no auth required)
		authOpenRoutes.POST("/shortcircuitlogin", h.handleShortCircuitLogin)
		authOpenRoutes.GET("/google", h.handleGoogleLogin)
		authOpenRoutes.GET("/google/callback", h.handleGoogleCallback)
	}

	// Protected routes (auth required)
	h.ProtectedRouteGroup.Use(h.middleware.AuthenticationMiddleware())
	h.ProtectedRouteGroup.Use(h.middleware.AuthorisationMiddleware())
	{
		h.ProtectedRouteGroup.GET("/me", h.handleGetUser)
		h.ProtectedRouteGroup.POST("/logout", h.handleLogout)
		h.ProtectedRouteGroup.GET("/profile", h.handleGetProfile)
		h.ProtectedRouteGroup.PUT("/profile", h.handleUpdateProfile)

		// Initialize handlers
		accountHandler := NewAccountHandler(h)
		billingHandler := NewBillingHandler(h)

		// Account routes
		accountRoutes := h.ProtectedRouteGroup.Group("/account")
		{
			accountRoutes.GET("", accountHandler.GetAccount)
			accountRoutes.PUT("", accountHandler.UpdateAccount)
		}

		// Billing routes
		billingRoutes := h.ProtectedRouteGroup.Group("/billing")
		{
			billingRoutes.POST("/subscriptions", billingHandler.CreateSubscription)
			billingRoutes.DELETE("/subscriptions", billingHandler.CancelSubscription)
			billingRoutes.GET("/subscriptions", billingHandler.GetSubscriptionStatus)
		}

		h.OpenRouteGroup.GET("/plans", h.GetPlans)
		h.OpenRouteGroup.POST("/webhook", billingHandler.HandleWebhook)

	}
}

func (h *Handler) handlePing(c *gin.Context) {
	// Test database connection
	sqlDB, err := h.Db.DB()
	if err != nil {
		c.JSON(500, gin.H{"error": "Database connection error"})
		return
	}

	if err := sqlDB.Ping(); err != nil {
		c.JSON(500, gin.H{"error": "Database ping failed"})
		return
	}

	c.JSON(200, gin.H{"message": "pong"})
}

// GetUserFromContext is a helper function to get the user from the context
func (handler *Handler) GetUserFromContext(c *gin.Context) (user *models.User) {
	var (
		exists bool
		cObj   interface{}
	)
	if cObj, exists = c.Get("user"); !exists {
		panic("user not found in context")
	}
	user = cObj.(*models.User)
	return
}

func (h *Handler) UserScopedDB(c *gin.Context) (db *gorm.DB) {
	db = h.authorisation.UserScopedDB(c, h.Db)
	return
}

func (h *Handler) FirstWithUser(c *gin.Context, userOwnedModel models.UserOwnedModel) (err error) {
	if err = h.UserScopedDB(c).First(userOwnedModel).Error; err != nil {
		return
	}
	return
}

func (h *Handler) FindWithUser(c *gin.Context) (userOwnedModels []models.UserOwnedModel, err error) {
	if err = h.UserScopedDB(c).Find(&userOwnedModels).Error; err != nil {
		return
	}
	return
}

func (h *Handler) CreateWithUser(c *gin.Context, model models.UserOwnedModel) (err error) {
	h.authorisation.UpdateWithUser(c, model)
	if err = h.Db.Create(model).Error; err != nil {
		return
	}
	return
}

func (h *Handler) UpdateWithUser(c *gin.Context, model models.UserOwnedModel, toUpdateWith interface{}) (err error) {
	h.authorisation.UpdateWithUser(c, model)
	err = h.Db.Model(model).Updates(toUpdateWith).Error
	return
}

func (h *Handler) DeleteWithUser(c *gin.Context, model models.UserOwnedModel) (err error) {
	if h.authorisation.CanAccessResource(c, model) == false {
		err = fmt.Errorf("unauthorised to delete resource")
		return
	}
	err = h.Db.Delete(model).Error
	return
}

func (h *Handler) WriteSuccess(c *gin.Context, data interface{}) {
	h.WriteJSON(c, 200, data)
	return
}

func (h *Handler) WriteError(c *gin.Context, err error, message string) {
	if err != nil {
		log.Println(err, message)
	}
	c.JSON(500, gin.H{
		"error": message,
	})
}

func (h *Handler) WriteJSON(c *gin.Context, status int, data interface{}) {
	c.JSON(status, gin.H{
		"data": data,
	})
}
