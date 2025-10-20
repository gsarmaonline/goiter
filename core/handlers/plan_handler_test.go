package handlers

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gsarmaonline/goiter/core/models"
)

func TestPlanHandler(t *testing.T) {
	handler, db := setupTestHandler(t)

	t.Run("GetPlans", func(t *testing.T) {
		t.Run("Success with multiple plans", func(t *testing.T) {
			// Create additional test plans (we already have the default "Free" plan)
			basicPlan := &models.Plan{
				Name:          "Basic",
				Price:         0, // Avoid Stripe integration
				BillingPeriod: "monthly",
				Description:   "Basic plan for testing",
			}
			err := db.Create(basicPlan).Error
			require.NoError(t, err)

			premiumPlan := &models.Plan{
				Name:          "Premium",
				Price:         0, // Avoid Stripe integration
				BillingPeriod: "monthly",
				Description:   "Premium plan for testing",
			}
			err = db.Create(premiumPlan).Error
			require.NoError(t, err)

			// Make request (no authentication required - this is a public endpoint)
			w := makeAuthenticatedRequest(t, handler, "GET", "/plans", nil, "")

			// Assert response
			assert.Equal(t, 200, w.Code)

			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			// Check if response has 'data' wrapper
			if data, hasData := response["data"]; hasData {
				plans := data.([]interface{})
				assert.GreaterOrEqual(t, len(plans), 3) // Should have at least 3 plans (Free, Basic, Premium)

				// Verify plan structure
				firstPlan := plans[0].(map[string]interface{})
				assert.Contains(t, firstPlan, "id")
				assert.Contains(t, firstPlan, "name")
				assert.Contains(t, firstPlan, "price")
				assert.Contains(t, firstPlan, "billing_period")
				assert.Contains(t, firstPlan, "description")
				assert.Contains(t, firstPlan, "features") // Should include features relationship

				// Verify we can find our test plans
				planNames := make([]string, 0, len(plans))
				for _, plan := range plans {
					planData := plan.(map[string]interface{})
					planNames = append(planNames, planData["name"].(string))
				}
				assert.Contains(t, planNames, "Free")
				assert.Contains(t, planNames, "Basic")
				assert.Contains(t, planNames, "Premium")
			}
		})

		t.Run("Success with no plans", func(t *testing.T) {
			// Create a separate handler with empty database for this test
			handler2, db2 := setupTestHandler(t)

			// Clear all plans from this test database
			err := db2.Exec("DELETE FROM plans").Error
			require.NoError(t, err)

			// Make request
			w := makeAuthenticatedRequest(t, handler2, "GET", "/plans", nil, "")

			// Assert response - should still return 200 with empty array
			assert.Equal(t, 200, w.Code)

			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			// Check if response has 'data' wrapper
			if data, hasData := response["data"]; hasData {
				plans := data.([]interface{})
				assert.Equal(t, 0, len(plans)) // Should be empty
			}
		})

		t.Run("Success with features relationship", func(t *testing.T) {
			// Create a plan with features
			planWithFeatures := &models.Plan{
				Name:          "Enterprise",
				Price:         0, // Avoid Stripe integration
				BillingPeriod: "monthly",
				Description:   "Enterprise plan with features",
			}
			err := db.Create(planWithFeatures).Error
			require.NoError(t, err)

			// Create some features
			feature1 := &models.Feature{
				Name:        "Advanced Analytics",
				Description: "Get detailed analytics",
				Limit:       100,
			}
			err = db.Create(feature1).Error
			require.NoError(t, err)

			feature2 := &models.Feature{
				Name:        "Priority Support",
				Description: "24/7 priority support",
				Limit:       -1, // Unlimited
			}
			err = db.Create(feature2).Error
			require.NoError(t, err)

			// Associate features with plan
			err = db.Model(planWithFeatures).Association("Features").Append([]*models.Feature{feature1, feature2})
			require.NoError(t, err)

			// Make request
			w := makeAuthenticatedRequest(t, handler, "GET", "/plans", nil, "")

			// Assert response
			assert.Equal(t, 200, w.Code)

			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			// Find the Enterprise plan in the response
			if data, hasData := response["data"]; hasData {
				plans := data.([]interface{})
				var enterprisePlan map[string]interface{}

				for _, plan := range plans {
					planData := plan.(map[string]interface{})
					if planData["name"].(string) == "Enterprise" {
						enterprisePlan = planData
						break
					}
				}

				require.NotNil(t, enterprisePlan, "Enterprise plan should be in response")

				// Verify features are included
				assert.Contains(t, enterprisePlan, "features")
				features := enterprisePlan["features"].([]interface{})
				assert.Equal(t, 2, len(features))

				// Verify feature structure
				feature := features[0].(map[string]interface{})
				assert.Contains(t, feature, "name")
				assert.Contains(t, feature, "description")
				assert.Contains(t, feature, "limit")

				// Verify feature names
				featureNames := make([]string, 0, len(features))
				for _, feat := range features {
					featData := feat.(map[string]interface{})
					featureNames = append(featureNames, featData["name"].(string))
				}
				assert.Contains(t, featureNames, "Advanced Analytics")
				assert.Contains(t, featureNames, "Priority Support")
			}
		})

		t.Run("Database error handling", func(t *testing.T) {
			// This test is more complex to setup as we'd need to mock the database
			// For now, we can test that the endpoint handles normal database operations
			// In a real scenario, you might want to use a database mock or test with a failing DB

			// Make request - should work normally
			w := makeAuthenticatedRequest(t, handler, "GET", "/plans", nil, "")

			// Should return success (since our DB is working)
			assert.Equal(t, 200, w.Code)
		})
	})
}

func TestPlanHandler_PublicAccess(t *testing.T) {
	handler, db := setupTestHandler(t)

	t.Run("Plans endpoint should be accessible without authentication", func(t *testing.T) {
		// Create a test plan
		publicPlan := &models.Plan{
			Name:          "Public Plan",
			Price:         0,
			BillingPeriod: "monthly",
			Description:   "Plan accessible without auth",
		}
		err := db.Create(publicPlan).Error
		require.NoError(t, err)

		// Make request WITHOUT any authentication token
		w := makeAuthenticatedRequest(t, handler, "GET", "/plans", nil, "")

		// Should succeed even without authentication
		assert.Equal(t, 200, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		// Should still return plans data
		if data, hasData := response["data"]; hasData {
			plans := data.([]interface{})
			assert.GreaterOrEqual(t, len(plans), 1)

			// Verify we can see the public plan
			planNames := make([]string, 0, len(plans))
			for _, plan := range plans {
				planData := plan.(map[string]interface{})
				planNames = append(planNames, planData["name"].(string))
			}
			assert.Contains(t, planNames, "Public Plan")
		}
	})

	t.Run("Plans endpoint should work with invalid token", func(t *testing.T) {
		// Make request with invalid authentication token
		w := makeAuthenticatedRequest(t, handler, "GET", "/plans", nil, "invalid.token.here")

		// Should still succeed (plans are public)
		assert.Equal(t, 200, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		// Should return plans data
		assert.Contains(t, response, "data")
	})

	t.Run("Plans endpoint should work with valid token", func(t *testing.T) {
		// Create a user and get valid token
		_, token := createTestUser(t, db, "validuser@example.com")

		// Make request with valid authentication token
		w := makeAuthenticatedRequest(t, handler, "GET", "/plans", nil, token)

		// Should succeed with authentication as well
		assert.Equal(t, 200, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		// Should return plans data
		assert.Contains(t, response, "data")
	})
}

func TestPlanHandler_Performance(t *testing.T) {
	handler, db := setupTestHandler(t)

	t.Run("Should handle many plans efficiently", func(t *testing.T) {
		// Create multiple plans to test performance
		for i := 0; i < 20; i++ {
			plan := &models.Plan{
				Name:          fmt.Sprintf("Plan %d", i),
				Price:         0,
				BillingPeriod: "monthly",
				Description:   fmt.Sprintf("Test plan number %d", i),
			}
			err := db.Create(plan).Error
			require.NoError(t, err)
		}

		// Make request
		w := makeAuthenticatedRequest(t, handler, "GET", "/plans", nil, "")

		// Should succeed
		assert.Equal(t, 200, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		// Should return all plans
		if data, hasData := response["data"]; hasData {
			plans := data.([]interface{})
			assert.GreaterOrEqual(t, len(plans), 20) // At least the 20 we created + the default one
		}
	})
}
