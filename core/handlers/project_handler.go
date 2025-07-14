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

// CheckProjectPermission checks if a user has the required permission level for a project
func (h *ProjectHandler) CheckProjectPermission(c *gin.Context, projectID uint, requiredLevel models.PermissionLevel) bool {
	user := h.handler.GetUserFromContext(c)

	var permission models.Permission

	// Check if user has the required permission
	err := h.handler.UserScopedDB(c).Where("project_id = ?", projectID).First(&permission).Error
	if err != nil {
		// If no permission record found, check if user is the owner
		var project models.Project
		if err := h.db.First(&project, projectID).Error; err != nil {
			return false
		}
		return project.UserID == user.ID
	}

	// Check permission level
	switch requiredLevel {
	case models.PermissionOwner:
		return permission.Level == models.PermissionOwner
	case models.PermissionAdmin:
		return permission.Level == models.PermissionOwner || permission.Level == models.PermissionAdmin
	case models.PermissionEditor:
		return permission.Level == models.PermissionOwner || permission.Level == models.PermissionAdmin || permission.Level == models.PermissionEditor
	case models.PermissionViewer:
		return true // Everyone can view
	default:
		return false
	}
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

	// Check if user has at least viewer permission
	if !h.CheckProjectPermission(c, uint(projectID), models.PermissionViewer) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to view this project"})
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

	// Check if user has admin or owner permission
	if !h.CheckProjectPermission(c, uint(projectID), models.PermissionAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to update this project"})
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

	// Only owner can delete the project
	if !h.CheckProjectPermission(c, uint(projectID), models.PermissionOwner) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only the project owner can delete the project"})
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

	// Only owner or admin can add members
	if !h.CheckProjectPermission(c, uint(projectID), models.PermissionAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to add members to this project"})
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

	// Only owner or admin can remove members
	if !h.CheckProjectPermission(c, uint(projectID), models.PermissionAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to remove members from this project"})
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
