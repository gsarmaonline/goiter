package handlers

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gsarmaonline/goiter/core/models"
	"github.com/gsarmaonline/goiter/core/services"
	"gorm.io/gorm"
)

type BillingHandler struct {
	db            *gorm.DB
	handler       *Handler
	stripeService *services.StripeService
}

func NewBillingHandler(handler *Handler) *BillingHandler {
	return &BillingHandler{
		handler:       handler,
		db:            handler.Db,
		stripeService: services.NewStripeService(handler.Db),
	}
}

// CreateSubscriptionRequest represents the request body for creating a subscription
type CreateSubscriptionRequest struct {
	PlanID          uint   `json:"plan_id" binding:"required"`
	PaymentMethodID string `json:"payment_method_id" binding:"required"`
}

// CreateSubscription creates a new subscription for the current user's account
func (h *BillingHandler) CreateSubscription(c *gin.Context) {
	var req CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the account for the current user
	var account models.Account
	if err := h.handler.UserScopedDB(c).First(&account).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	// Get the plan
	var plan models.Plan
	if err := h.db.Where("id = ?", req.PlanID).First(&plan).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan ID"})
		return
	}

	// Create the subscription
	subscription, err := account.CreateStripeSubscription(h.db, &plan, req.PaymentMethodID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"subscription_id": subscription.ID,
		"status":          subscription.Status,
		"message":         "Subscription created successfully",
	})
}

// CancelSubscription cancels the current user's subscription
func (h *BillingHandler) CancelSubscription(c *gin.Context) {
	// Get the account for the current user
	var account models.Account
	if err := h.handler.UserScopedDB(c).First(&account).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	// Cancel the subscription
	if err := account.CancelStripeSubscription(h.db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Subscription cancelled successfully",
	})
}

// GetSubscriptionStatus returns the current subscription status
func (h *BillingHandler) GetSubscriptionStatus(c *gin.Context) {
	// Get the account for the current user
	var account models.Account
	if err := h.handler.UserScopedDB(c).Preload("Plan").First(&account).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"subscription_id": account.StripeSubscriptionID,
		"status":          account.SubscriptionStatus,
		"plan":            account.Plan,
	})
}

// HandleWebhook processes Stripe webhooks
func (h *BillingHandler) HandleWebhook(c *gin.Context) {
	// Read the request body
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	// Get the signature from headers
	signature := c.GetHeader("Stripe-Signature")
	if signature == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing Stripe signature"})
		return
	}

	// Process the webhook
	if err := h.stripeService.ProcessWebhook(payload, signature); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Webhook processed successfully"})
}
