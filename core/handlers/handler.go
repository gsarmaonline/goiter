package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gsarmaonline/goiter/config"
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
		db            *gorm.DB
		middleware    *middleware.Middleware
		cfg           *config.Config
		authorisation *models.Authorisation

		OpenRouteGroup      *gin.RouterGroup
		ProtectedRouteGroup *gin.RouterGroup
	}
)

func NewHandler(router *gin.Engine, db *gorm.DB, cfg *config.Config) (handler *Handler) {
	middleware := middleware.NewMiddleware(router, db)
	handler = &Handler{
		router:     router,
		db:         db,
		middleware: middleware,
		cfg:        cfg,

		OpenRouteGroup:      router.Group("/"),
		ProtectedRouteGroup: router.Group(""),
	}
	// Setup routes
	handler.SetupRoutes()
	// Setup authorisation
	handler.authorisation = models.NewAuthorisation()
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
		groupHandler := NewGroupHandler(h)

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

		// Group routes
		groupRoutes := h.ProtectedRouteGroup.Group("/groups")
		{
			groupRoutes.POST("", groupHandler.CreateGroup)
			groupRoutes.GET("", groupHandler.ListGroups)
			groupRoutes.GET("/:id", groupHandler.GetGroup)
			groupRoutes.GET("/:id/ancestors", groupHandler.GetGroupAncestors)
			groupRoutes.DELETE("/:id", groupHandler.DeleteGroup)
			groupRoutes.POST("/:id/members", groupHandler.AddGroupMember)
			groupRoutes.DELETE("/:id/members", groupHandler.RemoveGroupMember)
		}

		h.OpenRouteGroup.GET("/plans", h.GetPlans)
		h.OpenRouteGroup.POST("/webhook", billingHandler.HandleWebhook)

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

func (h *Handler) GetModelFromUrl(c *gin.Context, userOwnedModel models.UserOwnedModel, urlKeyName string) (err error) {
	var (
		modelID uint64
	)
	if modelID, err = strconv.ParseUint(c.Param(urlKeyName), 10, 32); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}
	if err = h.db.First(&userOwnedModel, modelID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}
	if err = h.FirstWithUser(c, userOwnedModel, h.db.Where("id = ?", modelID)); err != nil {
		return
	}
	return
}

func (h *Handler) FirstWithUser(c *gin.Context, userOwnedModel models.UserOwnedModel, dbQuery *gorm.DB) (err error) {
	user := h.GetUserFromContext(c)
	err = dbQuery.First(&userOwnedModel).Error
	if false == h.authorisation.CanAccessResource(h.db, user, userOwnedModel, models.ReadAction, models.Scope{}) {
		err = fmt.Errorf("unauthorized access to resource: %s", h.GetTableName(userOwnedModel))
		return
	}
	return
}

func (h *Handler) FindWithUser(c *gin.Context, userOwnedModel interface{},
	query string) (userOwnedModels []models.UserOwnedModel, err error) {

	user := h.GetUserFromContext(c)
	err = h.db.Model(userOwnedModel).Where(query).Find(&userOwnedModels).Error
	for _, model := range userOwnedModels {
		if false == h.authorisation.CanAccessResource(h.db,
			user,
			model.(models.UserOwnedModel),
			models.ReadAction,
			models.Scope{},
		) {
			err = fmt.Errorf("unauthorized access to resource: %s", model.(models.UserOwnedModel).GetConfig().Name)
			return
		}
	}
	return
}

func (h *Handler) CreateWithUser(c *gin.Context, model models.UserOwnedModel) (err error) {
	user := h.GetUserFromContext(c)
	model.SetUserID(user.ID)
	if false == h.authorisation.CanAccessResource(h.db, user, model, models.CreateAction, models.Scope{}) {
		err = fmt.Errorf("unauthorized access to resource: %s", h.GetTableName(model))
		return
	}
	err = h.db.Create(model).Error
	return
}

func (h *Handler) UpdateWithUser(c *gin.Context, model models.UserOwnedModel, toUpdateWith interface{}) (err error) {
	user := h.GetUserFromContext(c)
	model.SetUserID(user.ID)
	if false == h.authorisation.CanAccessResource(h.db, user, model, models.UpdateAction, models.Scope{}) {
		err = fmt.Errorf("unauthorized access to resource: %s", h.GetTableName(model))
		return
	}
	err = h.db.Model(model).Updates(toUpdateWith).Error
	return
}

func (h *Handler) DeleteWithUser(c *gin.Context, model models.UserOwnedModel) (err error) {
	user := h.GetUserFromContext(c)
	model.SetUserID(user.ID)
	if false == h.authorisation.CanAccessResource(h.db, user, model, models.DeleteAction, models.Scope{}) {
		err = fmt.Errorf("unauthorized access to resource: %s", h.GetTableName(model))
		return
	}
	err = h.db.Delete(model).Error
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
