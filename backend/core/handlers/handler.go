package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gsarmaonline/goiter/core/middleware"
	"github.com/gsarmaonline/goiter/core/models"
	"gorm.io/gorm"
)

type (
	Handler struct {
		router     *gin.Engine
		db         *gorm.DB
		middleware *middleware.Middleware
	}
)

func NewHandler(router *gin.Engine, db *gorm.DB) (handler *Handler) {
	middleware := middleware.NewMiddleware(router, db)
	handler = &Handler{
		router:     router,
		db:         db,
		middleware: middleware,
	}
	// Setup routes
	handler.SetupRoutes()
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
		authOpenRoutes.GET("/google", h.handleGoogleLogin)
		authOpenRoutes.GET("/google/callback", h.handleGoogleCallback)
	}

	// Protected routes (auth required)
	protected := h.router.Group("")
	protected.Use(h.middleware.AuthenticationMiddleware())
	protected.Use(h.middleware.AuthorisationMiddleware())
	{
		protected.GET("/me", h.handleGetUser)
		protected.POST("/logout", h.handleLogout)
		protected.GET("/profile", h.handleGetProfile)
		protected.PUT("/profile", h.handleUpdateProfile)

		// Initialize handlers
		projectHandler := NewProjectHandler(h)
		accountHandler := NewAccountHandler(h)
		billingHandler := NewBillingHandler(h)

		// Project routes
		projectRoutes := protected.Group("/projects")
		{
			projectRoutes.POST("", projectHandler.CreateProject)
			projectRoutes.GET("", projectHandler.ListProjects)
			projectRoutes.GET("/:id", projectHandler.GetProject)
			projectRoutes.PUT("/:id", projectHandler.UpdateProject)
			projectRoutes.DELETE("/:id", projectHandler.DeleteProject)
			projectRoutes.POST("/:id/members", projectHandler.AddProjectMember)
			projectRoutes.DELETE("/:id/members/:user_id", projectHandler.RemoveProjectMember)
		}

		// Account routes
		accountRoutes := protected.Group("/account")
		{
			accountRoutes.GET("", accountHandler.GetAccount)
			accountRoutes.PUT("", accountHandler.UpdateAccount)
		}

		// Billing routes
		billingRoutes := protected.Group("/billing")
		{
			billingRoutes.POST("/subscriptions", billingHandler.CreateSubscription)
			billingRoutes.DELETE("/subscriptions", billingHandler.CancelSubscription)
			billingRoutes.GET("/subscriptions", billingHandler.GetSubscriptionStatus)
		}

		openRoutes := h.router.Group("/")
		{
			openRoutes.GET("/plans", h.GetPlans)
			openRoutes.POST("/webhook", billingHandler.HandleWebhook)
		}

	}
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

func (h *Handler) handlePing(c *gin.Context) {
	// Test database connection
	sqlDB, err := h.db.DB()
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

func (h *Handler) GetTableName(model models.UserOwnedModel) (tableName string) {
	stmt := &gorm.Statement{DB: h.db}
	stmt.Parse(model)
	return stmt.Schema.Table
}

func (h *Handler) UserScopedDB(c *gin.Context) (db *gorm.DB) {
	user := h.GetUserFromContext(c)
	db = h.db.Where("user_id", user.ID)
	return
}

func (h *Handler) FirstWithUser(c *gin.Context, userOwnedModel models.UserOwnedModel, dbQuery *gorm.DB) (err error) {
	user := h.GetUserFromContext(c)
	err = dbQuery.First(&userOwnedModel).Error
	if false == models.CanAccessResource(h.db, h.GetTableName(userOwnedModel), userOwnedModel.GetID(), user, models.ReadAction) {
		err = fmt.Errorf("unauthorized access to resource: %s", h.GetTableName(userOwnedModel))
		return
	}
	return
}

func (h *Handler) FindWithUser(c *gin.Context, userOwnedModels []models.UserOwnedModel, query string) (err error) {
	user := h.GetUserFromContext(c)
	err = h.db.Where(query).Find(&userOwnedModels).Error
	for _, model := range userOwnedModels {
		if false == models.CanAccessResource(h.db, h.GetTableName(model), model.GetID(), user, models.ReadAction) {
			err = fmt.Errorf("unauthorized access to resource: %s", h.GetTableName(model))
			return
		}
	}
	return
}

func (h *Handler) CreateWithUser(c *gin.Context, model models.UserOwnedModel) (err error) {
	user := h.GetUserFromContext(c)
	model.SetUserID(user.ID)
	if false == models.CanAccessResource(h.db, h.GetTableName(model), model.GetID(), user, models.CreateAction) {
		err = fmt.Errorf("unauthorized access to resource: %s", h.GetTableName(model))
		return
	}
	err = h.db.Create(model).Error
	return
}

func (h *Handler) UpdateWithUser(c *gin.Context, model models.UserOwnedModel) (err error) {
	user := h.GetUserFromContext(c)
	model.SetUserID(user.ID)
	if false == models.CanAccessResource(h.db, h.GetTableName(model), model.GetID(), user, models.UpdateAction) {
		err = fmt.Errorf("unauthorized access to resource: %s", h.GetTableName(model))
		return
	}
	err = h.db.Save(model).Error
	return
}

func (h *Handler) DeleteWithUser(c *gin.Context, model models.UserOwnedModel) (err error) {
	user := h.GetUserFromContext(c)
	model.SetUserID(user.ID)
	if false == models.CanAccessResource(h.db, h.GetTableName(model), model.GetID(), user, models.DeleteAction) {
		err = fmt.Errorf("unauthorized access to resource: %s", h.GetTableName(model))
		return
	}
	err = h.db.Delete(model).Error
	return
}
