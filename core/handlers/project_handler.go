package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gsarmaonline/goiter/core/models"
	"gorm.io/gorm"
)

type ProjectHandler struct {
	handler *Handler
	db      *gorm.DB
}

func NewProjectHandler(handler *Handler) *ProjectHandler {
	return &ProjectHandler{handler: handler, db: handler.db}
}

// ListProjects retrieves all projects for the current user
func (h *ProjectHandler) ListProjects(c *gin.Context) {
	userID := h.handler.GetUserFromContext(c).ID
	var projects []models.Project

	// Get projects where user is owner or member
	if err := h.handler.UserScopedDB(c).Where("id IN (SELECT project_id FROM permissions WHERE user_id = ?)", userID).
		Preload("Members").
		Find(&projects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list projects"})
		return
	}

	h.handler.WriteSuccess(c, projects)
}

// CreateProject creates a new project
func (h *ProjectHandler) CreateProject(c *gin.Context) {
	project := &models.Project{}

	if err := c.ShouldBindJSON(project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user from context
	user := h.handler.GetUserFromContext(c)
	project.UserID = user.ID

	// Start a transaction
	tx := h.db.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	// Get account owned by user
	var account models.Account
	if err := h.handler.UserScopedDB(c).First(&account).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user account"})
		return
	}

	// Set account ID on project
	project.AccountID = account.ID

	// Create the project
	if err := h.handler.CreateWithUser(c, project); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project"})
		return
	}

	// Create owner permission
	permission := models.Permission{
		ProjectID: project.ID,
		Level:     models.PermissionOwner,
	}
	permission.UserID = user.ID
	if err := h.handler.CreateWithUser(c, &permission); err != nil {
		log.Println("Failed to create project permission:", err)
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set project permissions"})
		return
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	h.handler.WriteSuccess(c, project)
}

// GetProject retrieves a project by ID
func (h *ProjectHandler) GetProject(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var project models.Project
	if err := h.db.Preload("User").Preload("Members").First(&project, projectID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	h.handler.WriteSuccess(c, project)
}

// UpdateProject updates a project
func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var project models.Project
	if err := h.db.First(&project, projectID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	// Update only allowed fields
	var updateData struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields
	if updateData.Name != "" {
		project.Name = updateData.Name
	}
	if updateData.Description != "" {
		project.Description = updateData.Description
	}

	if err := h.db.Save(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update project"})
		return
	}

	h.handler.WriteSuccess(c, project)
}

// DeleteProject deletes a project
func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	if err := h.db.Delete(&models.Project{}, projectID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete project"})
		return
	}

	h.handler.WriteSuccess(c, nil)
}

// AddProjectMember adds a user to a project with a specific permission level
func (h *ProjectHandler) AddProjectMember(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var input struct {
		UserID    uint                   `json:"user_id"`
		UserEmail string                 `json:"user_email" binding:"required"`
		Level     models.PermissionLevel `json:"level" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create the permission
	permission := models.Permission{
		ProjectID: uint(projectID),
		UserEmail: input.UserEmail,
		Level:     input.Level,
	}
	permission.UserID = input.UserID

	if err := h.db.Create(&permission).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.handler.WriteSuccess(c, permission)
}

// RemoveProjectMember removes a user from a project
func (h *ProjectHandler) RemoveProjectMember(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	memberID, err := strconv.ParseUint(c.Param("user_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Don't allow removing the owner
	var project models.Project
	if err := h.db.First(&project, projectID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}
	if uint(memberID) == project.UserID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot remove the project owner"})
		return
	}

	if err := h.handler.UserScopedDB(c).Where("project_id = ?", projectID).Delete(&models.Permission{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove project member"})
		return
	}

	h.handler.WriteSuccess(c, nil)
}
