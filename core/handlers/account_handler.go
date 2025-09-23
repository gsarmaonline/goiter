package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gsarmaonline/goiter/core/models"
	"gorm.io/gorm"
)

type (
	AccountHandler struct {
		db      *gorm.DB
		handler *Handler
	}

	AccountUpdateRequest struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		PlanID      uint   `json:"plan_id"`
	}
)

func NewAccountHandler(handler *Handler) *AccountHandler {
	return &AccountHandler{handler: handler, db: handler.Db}
}

// GetAccount retrieves the account for the current user
func (h *AccountHandler) GetAccount(c *gin.Context) {

	var account models.Account
	if err := h.handler.UserScopedDB(c).Preload("Plan").First(&account).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	h.handler.WriteSuccess(c, account)
}

// UpdateAccount updates the account for the current user
func (h *AccountHandler) UpdateAccount(c *gin.Context) {

	var (
		account    models.Account
		updateData AccountUpdateRequest
	)

	if err := h.handler.UserScopedDB(c).First(&account).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Update the account fields with the new data
	account.Name = updateData.Name
	account.Description = updateData.Description
	account.PlanID = updateData.PlanID

	if err := h.handler.UpdateWithUser(c, &account, &account); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update account"})
		return
	}

	h.handler.WriteSuccess(c, account)
}
