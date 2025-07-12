package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/gsarmaonline/goiter/core/models"
)

// handleGetProfile handles the profile read request
func (h *Handler) handleGetProfile(c *gin.Context) {
	var profile models.Profile
	if err := h.FirstWithUser(c, &profile, h.UserScopedDB(c)); err != nil {
		c.JSON(404, gin.H{"error": "Profile not found"})
		return
	}

	h.WriteSuccess(c, profile)
}

// handleUpdateProfile handles the profile update request
func (h *Handler) handleUpdateProfile(c *gin.Context) {
	updatedProfile := &models.Profile{}
	profile := &models.Profile{}

	if err := c.ShouldBindJSON(&updatedProfile); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}
	if err := h.FirstWithUser(c, profile, h.UserScopedDB(c)); err != nil {
		c.JSON(404, gin.H{"error": "Profile not found"})
		return
	}

	// Update profile
	if err := h.UpdateWithUser(c, profile, updatedProfile); err != nil {
		c.JSON(500, gin.H{"error": "Failed to update profile"})
		return
	}

	h.WriteSuccess(c, updatedProfile)
}
