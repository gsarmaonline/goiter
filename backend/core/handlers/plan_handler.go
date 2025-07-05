package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gsarmaonline/goiter/core/models"
)

func (h *Handler) GetPlans(c *gin.Context) {
	plans := []models.Plan{}
	err := h.db.Preload("Features").Find(&plans).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.WriteSuccess(c, plans)
}
