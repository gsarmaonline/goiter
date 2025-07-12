package services

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gsarmaonline/goiter/core/models"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/webhook"
	"gorm.io/gorm"
)

type StripeService struct {
	db *gorm.DB
}

func NewStripeService(db *gorm.DB) *StripeService {
	// Set Stripe API key
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	return &StripeService{db: db}
}

// ProcessWebhook processes Stripe webhooks
func (s *StripeService) ProcessWebhook(payload []byte, signature string) error {
	event, err := webhook.ConstructEvent(payload, signature, os.Getenv("STRIPE_WEBHOOK_SECRET"))
	if err != nil {
		return fmt.Errorf("failed to verify webhook signature: %v", err)
	}

	switch event.Type {
	case "customer.subscription.created":
		return s.handleSubscriptionCreated(event)
	case "customer.subscription.updated":
		return s.handleSubscriptionUpdated(event)
	case "customer.subscription.deleted":
		return s.handleSubscriptionDeleted(event)
	case "invoice.payment_succeeded":
		return s.handlePaymentSucceeded(event)
	case "invoice.payment_failed":
		return s.handlePaymentFailed(event)
	}

	return nil
}

func (s *StripeService) handleSubscriptionCreated(event stripe.Event) error {
	var sub stripe.Subscription
	err := json.Unmarshal(event.Data.Raw, &sub)
	if err != nil {
		return fmt.Errorf("failed to unmarshal subscription: %v", err)
	}

	// Update account subscription status
	if err := s.db.Model(&models.Account{}).
		Where("stripe_subscription_id = ?", sub.ID).
		Update("subscription_status", sub.Status).Error; err != nil {
		return fmt.Errorf("failed to update account subscription status: %v", err)
	}

	return nil
}

func (s *StripeService) handleSubscriptionUpdated(event stripe.Event) error {
	var sub stripe.Subscription
	err := json.Unmarshal(event.Data.Raw, &sub)
	if err != nil {
		return fmt.Errorf("failed to unmarshal subscription: %v", err)
	}

	// Update account subscription status
	if err := s.db.Model(&models.Account{}).
		Where("stripe_subscription_id = ?", sub.ID).
		Update("subscription_status", sub.Status).Error; err != nil {
		return fmt.Errorf("failed to update account subscription status: %v", err)
	}

	return nil
}

func (s *StripeService) handleSubscriptionDeleted(event stripe.Event) error {
	var sub stripe.Subscription
	err := json.Unmarshal(event.Data.Raw, &sub)
	if err != nil {
		return fmt.Errorf("failed to unmarshal subscription: %v", err)
	}

	// Reset account to free plan
	if err := s.db.Model(&models.Account{}).
		Where("stripe_subscription_id = ?", sub.ID).
		Updates(map[string]interface{}{
			"plan_id":                "free",
			"stripe_subscription_id": "",
			"subscription_status":    "canceled",
		}).Error; err != nil {
		return fmt.Errorf("failed to reset account to free plan: %v", err)
	}

	return nil
}

func (s *StripeService) handlePaymentSucceeded(event stripe.Event) error {
	var invoice stripe.Invoice
	err := json.Unmarshal(event.Data.Raw, &invoice)
	if err != nil {
		return fmt.Errorf("failed to unmarshal invoice: %v", err)
	}

	// Update account subscription status
	if err := s.db.Model(&models.Account{}).
		Where("stripe_subscription_id = ?", invoice.Subscription).
		Update("subscription_status", "active").Error; err != nil {
		return fmt.Errorf("failed to update account subscription status: %v", err)
	}

	return nil
}

func (s *StripeService) handlePaymentFailed(event stripe.Event) error {
	var invoice stripe.Invoice
	err := json.Unmarshal(event.Data.Raw, &invoice)
	if err != nil {
		return fmt.Errorf("failed to unmarshal invoice: %v", err)
	}

	// Update account subscription status
	if err := s.db.Model(&models.Account{}).
		Where("stripe_subscription_id = ?", invoice.Subscription).
		Update("subscription_status", "past_due").Error; err != nil {
		return fmt.Errorf("failed to update account subscription status: %v", err)
	}

	return nil
}
