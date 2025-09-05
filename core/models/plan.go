package models

import (
	"fmt"
	"os"

	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/price"
	"github.com/stripe/stripe-go/v74/product"
	"gorm.io/gorm"
)

// Plan represents a subscription plan with features and limits
type (
	Plan struct {
		BaseModelWithoutUser

		Name        string `json:"name" gorm:"not null"`
		Description string `json:"description"`

		Price         float64 `json:"price" gorm:"not null;default:0"`
		BillingPeriod string  `json:"billing_period" gorm:"not null;default:'monthly'"`

		StripeProductID string `json:"stripe_product_id"`
		StripePriceID   string `json:"stripe_price_id"`

		Features []*Feature `json:"features" gorm:"many2many:plan_features;"`
		Accounts []*Account `json:"accounts" gorm:"foreignKey:PlanID"`
	}

	// Feature represents a feature that can be included in plans
	Feature struct {
		BaseModelWithoutUser

		Name        string  `json:"name" gorm:"not null"`
		Description string  `json:"description"`
		Limit       int     `json:"limit" gorm:"not null;default:-1"` // -1 means unlimited
		Plans       []*Plan `json:"plans" gorm:"many2many:plan_features;"`
	}

	// PlanFeature is the explicit join table for the many-to-many relationship
	// between Plan and Feature, ensuring correct table and constraint creation.
	PlanFeature struct {
		BaseModelWithoutUser

		PlanID    uint `gorm:"primaryKey"`
		FeatureID uint `gorm:"primaryKey"`
	}
)

func (p Plan) GetConfig() ModelConfig {
	return ModelConfig{
		Name:      "Plan",
		ScopeType: AccountScopeType,
	}
}

func (p Feature) GetConfig() ModelConfig {
	return ModelConfig{
		Name:      "Feature",
		ScopeType: AccountScopeType,
	}
}

func (p PlanFeature) GetConfig() ModelConfig {
	return ModelConfig{
		Name:      "PlanFeature",
		ScopeType: AccountScopeType,
	}
}

func GetDefaultPlan(db *gorm.DB) (plan *Plan, err error) {
	plan = &Plan{}
	if err = db.First(plan).Error; err == nil {
		// Return if no error
		return
	}
	return
}

func (plan *Plan) GetStripeInterval() (interval stripe.PriceRecurringInterval) {
	switch plan.BillingPeriod {
	case "monthly":
		interval = stripe.PriceRecurringIntervalMonth
	case "yearly":
		interval = stripe.PriceRecurringIntervalYear
	case "weekly":
		interval = stripe.PriceRecurringIntervalWeek
	case "daily":
		interval = stripe.PriceRecurringIntervalDay
	default:
		interval = stripe.PriceRecurringIntervalMonth // Default to monthly
	}
	return
}

// CreatePrice creates a Stripe Price object for a plan
func (plan *Plan) BeforeCreate(tx *gorm.DB) (err error) {
	var (
		stripePrice   *stripe.Price
		stripeProduct *stripe.Product
	)
	if plan.Price == 0 {
		return
	}

	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	if stripeProduct, err = plan.CreateStripeProduct(); err != nil {
		return
	}
	plan.StripeProductID = stripeProduct.ID

	if stripePrice, err = plan.CreateStripePrice(); err != nil {
		return
	}
	plan.StripePriceID = stripePrice.ID

	return
}

func (plan *Plan) CreateStripeProduct() (stripeProduct *stripe.Product, err error) {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	stripeProduct, err = product.New(&stripe.ProductParams{
		Name: stripe.String(plan.Name),
	})
	return
}

func (plan *Plan) CreateStripePrice() (stripePrice *stripe.Price, err error) {
	var (
		interval stripe.PriceRecurringInterval
	)
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	interval = plan.GetStripeInterval()

	// Create price parameters
	params := &stripe.PriceParams{
		Currency: stripe.String("usd"),
		Product:  stripe.String(plan.StripeProductID),
		Recurring: &stripe.PriceRecurringParams{
			Interval: stripe.String(string(interval)),
		},
		UnitAmount: stripe.Int64(int64(plan.Price * 100)), // Convert to cents
		Params: stripe.Params{
			Metadata: map[string]string{
				"plan_id":   fmt.Sprintf("%d", plan.ID),
				"plan_name": plan.Name,
			},
		},
	}

	// Create the price in Stripe
	if stripePrice, err = price.New(params); err != nil {
		return
	}
	return

}
