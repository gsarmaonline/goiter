package models

import (
	"fmt"
	"os"

	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/customer"
	"github.com/stripe/stripe-go/v74/paymentmethod"
	"github.com/stripe/stripe-go/v74/subscription"
	"gorm.io/gorm"
)

// Account represents an organization or workspace that can contain multiple projects
type Account struct {
	BaseModelWithUser

	Name        string `json:"name" gorm:"not null"`
	Description string `json:"description"`

	Projects []Project `json:"projects" gorm:"foreignKey:AccountID"`

	PlanID uint `json:"plan_id"`
	Plan   Plan `json:"plan" gorm:"foreignKey:PlanID"`

	StripeCustomerID     string `json:"-"`
	StripeSubscriptionID string `json:"-"`
	SubscriptionStatus   string `json:"subscription_status" gorm:"default:'active'"`
}

// TableName specifies the table name for the Account model
func (Account) TableName() string {
	return "accounts"
}

// BeforeCreate is a GORM hook that ensures new accounts have the free plan
func (a *Account) BeforeCreate(tx *gorm.DB) error {
	plan, err := GetDefaultPlan(tx)
	if err != nil {
		return err
	}
	a.PlanID = plan.ID
	return nil
}

func (account *Account) BeforeUpdate(tx *gorm.DB) (err error) {
	if account.PlanID == 1 && account.StripeSubscriptionID != "" {
		err = account.CancelStripeSubscription(tx)
	}
	return
}

// BeforeDelete is a GORM hook that handles cleanup before account deletion
func (a *Account) BeforeDelete(tx *gorm.DB) error {
	// Delete all projects associated with this account
	if err := tx.Where("account_id = ?", a.ID).Delete(&Project{}).Error; err != nil {
		return err
	}
	return nil
}

func (account *Account) CreateStripeCustomer(tx *gorm.DB) (*stripe.Customer, error) {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	user := &User{}
	if err := tx.First(user, "id = ?", account.UserID).Error; err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	params := &stripe.CustomerParams{
		Email: stripe.String(user.Email),
		Name:  stripe.String(user.Name),
		Params: stripe.Params{
			Metadata: map[string]string{
				"user_id": fmt.Sprintf("%d", user.ID),
			},
		},
	}

	customer, err := customer.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create Stripe customer: %v", err)
	}

	// Update user with Stripe customer ID
	if err := tx.Model(account).Update("stripe_customer_id", customer.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to update user with Stripe customer ID: %v", err)
	}

	return customer, nil
}

// CreateStripeSubscription creates a subscription for an account
func (account *Account) CreateStripeSubscription(tx *gorm.DB, plan *Plan, paymentMethodID string) (*stripe.Subscription, error) {
	// Get or create Stripe customer
	var stripeCustomerID string
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	if account.HasActiveSubscription(tx) {
		return nil, fmt.Errorf("account already has an active subscription")
	}

	if account.StripeCustomerID == "" {
		customer, err := account.CreateStripeCustomer(tx)
		if err != nil {
			return nil, err
		}
		stripeCustomerID = customer.ID
	} else {
		stripeCustomerID = account.StripeCustomerID
	}

	// Attach payment method to customer
	attachParams := &stripe.PaymentMethodAttachParams{
		Customer: stripe.String(stripeCustomerID),
	}

	_, err := paymentmethod.Attach(paymentMethodID, attachParams)
	if err != nil {
		return nil, fmt.Errorf("failed to attach payment method to customer: %v", err)
	}

	// Create subscription
	params := &stripe.SubscriptionParams{
		Customer: stripe.String(stripeCustomerID),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Price: stripe.String(plan.StripePriceID),
			},
		},
		DefaultPaymentMethod: stripe.String(paymentMethodID),
		Params: stripe.Params{
			Metadata: map[string]string{
				"account_id": fmt.Sprintf("%d", account.ID),
				"plan_id":    fmt.Sprintf("%d", plan.ID),
			},
		},
	}

	sub, err := subscription.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create subscription: %v", err)
	}

	// Update account with subscription details
	if err := tx.Model(account).Updates(map[string]interface{}{
		"stripe_subscription_id": sub.ID,
		"plan_id":                plan.ID,
	}).Error; err != nil {
		return nil, fmt.Errorf("failed to update account: %v", err)
	}

	return sub, nil
}

func (account *Account) CancelStripeSubscription(tx *gorm.DB) error {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	if account.StripeSubscriptionID == "" {
		return fmt.Errorf("no subscription found for account")
	}

	params := &stripe.SubscriptionParams{
		CancelAtPeriodEnd: stripe.Bool(true),
	}

	_, err := subscription.Update(account.StripeSubscriptionID, params)
	if err != nil {
		return fmt.Errorf("failed to cancel subscription: %v", err)
	}

	// Update account to reflect cancellation
	if err := tx.Model(account).Update("subscription_status", "canceling").Error; err != nil {
		return fmt.Errorf("failed to update account status: %v", err)
	}

	return nil
}

func (account *Account) HasActiveSubscription(tx *gorm.DB) (hasActiveSubscription bool) {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	if account.StripeSubscriptionID == "" {
		return
	}
	sub, err := subscription.Get(account.StripeSubscriptionID, nil)
	if err != nil {
		return
	}
	hasActiveSubscription = sub.Status == "active"
	return
}
