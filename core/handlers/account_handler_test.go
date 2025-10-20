package handlers

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gsarmaonline/goiter/core/models"
)

func TestAccountHandler(t *testing.T) {
	handler, db := setupTestHandler(t)

	t.Run("GetAccount", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			// Create a user with account
			user, token := createTestUser(t, db, "accountuser@example.com")

			// Get the account that was created by the user AfterCreate hook
			var account models.Account
			err := db.Where("user_id = ?", user.ID).First(&account).Error
			require.NoError(t, err)

			// Update some account data to verify it's returned
			account.Name = "Updated Account Name"
			account.Description = "Test account description"
			err = db.Save(&account).Error
			require.NoError(t, err)

			// Make request
			w := makeAuthenticatedRequest(t, handler, "GET", "/account", nil, token)

			// Assert response
			assert.Equal(t, 200, w.Code)

			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			// Check if response has 'data' wrapper
			if data, hasData := response["data"]; hasData {
				accountData := data.(map[string]interface{})
				assert.Equal(t, float64(account.ID), accountData["id"])
				assert.Equal(t, float64(user.ID), accountData["user_id"])
				assert.Equal(t, "Updated Account Name", accountData["name"])
				assert.Equal(t, "Test account description", accountData["description"])
				assert.Contains(t, accountData, "plan") // Should include plan relationship
			}
		})

		t.Run("Unauthenticated", func(t *testing.T) {
			// Make request without authentication
			w := makeAuthenticatedRequest(t, handler, "GET", "/account", nil, "")
			assertErrorResponse(t, w, 401, "Not authenticated")
		})

		t.Run("Account not found", func(t *testing.T) {
			// Create user but delete the account to simulate missing account
			user, token := createTestUser(t, db, "noaccount@example.com")

			// Delete the account that was created by the user AfterCreate hook
			err := db.Where("user_id = ?", user.ID).Delete(&models.Account{}).Error
			require.NoError(t, err)

			// Make request
			w := makeAuthenticatedRequest(t, handler, "GET", "/account", nil, token)
			assertErrorResponse(t, w, 404, "Account not found")
		})

		t.Run("Invalid JWT", func(t *testing.T) {
			// Make request with invalid token
			w := makeAuthenticatedRequest(t, handler, "GET", "/account", nil, "invalid.jwt.token")
			assertErrorResponse(t, w, 401, "Invalid token")
		})
	})

	t.Run("UpdateAccount", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			// Create a user with account
			user, token := createTestUser(t, db, "updateaccount@example.com")

			// Get the account
			var account models.Account
			err := db.Where("user_id = ?", user.ID).First(&account).Error
			require.NoError(t, err)

			// Get a plan for testing plan_id update
			var plan models.Plan
			err = db.First(&plan).Error
			require.NoError(t, err)

			// Prepare update data
			updateData := map[string]interface{}{
				"name":        "Updated Account Name",
				"description": "Updated account description",
				"plan_id":     plan.ID,
			}

			// Make request
			w := makeAuthenticatedRequest(t, handler, "PUT", "/account", updateData, token)

			// Assert successful response
			assert.Equal(t, 200, w.Code)

			// Verify the account was updated in database
			var updatedAccount models.Account
			err = db.Where("user_id = ?", user.ID).First(&updatedAccount).Error
			require.NoError(t, err)

			assert.Equal(t, "Updated Account Name", updatedAccount.Name)
			assert.Equal(t, "Updated account description", updatedAccount.Description)
			assert.Equal(t, plan.ID, updatedAccount.PlanID)

			// Ensure the account ID hasn't changed
			assert.Equal(t, account.ID, updatedAccount.ID)
		})

		t.Run("Unauthenticated", func(t *testing.T) {
			updateData := map[string]interface{}{
				"name": "Should not work",
			}

			// Make request without authentication
			w := makeAuthenticatedRequest(t, handler, "PUT", "/account", updateData, "")
			assertErrorResponse(t, w, 401, "Not authenticated")
		})

		t.Run("Account not found", func(t *testing.T) {
			// Create user but delete the account
			user, token := createTestUser(t, db, "noaccountupdate@example.com")

			// Delete the account that was created by the user AfterCreate hook
			err := db.Where("user_id = ?", user.ID).Delete(&models.Account{}).Error
			require.NoError(t, err)

			updateData := map[string]interface{}{
				"name": "Should not work",
			}

			// Make request
			w := makeAuthenticatedRequest(t, handler, "PUT", "/account", updateData, token)
			assertErrorResponse(t, w, 404, "Account not found")
		})

		t.Run("Invalid JSON", func(t *testing.T) {
			// Create a user with account
			_, token := createTestUser(t, db, "invalidjsonaccount@example.com")

			// Make request with invalid JSON by using a string instead of map
			w := makeAuthenticatedRequest(t, handler, "PUT", "/account", "invalid json", token)
			assertErrorResponse(t, w, 400, "Invalid request body")
		})

		t.Run("Partial update", func(t *testing.T) {
			// Create a user with account
			user, token := createTestUser(t, db, "partialaccount@example.com")

			// Get and set some initial data
			var account models.Account
			err := db.Where("user_id = ?", user.ID).First(&account).Error
			require.NoError(t, err)

			account.Name = "Original name"
			account.Description = "Original description"
			err = db.Save(&account).Error
			require.NoError(t, err)

			// Prepare partial update data (only name)
			updateData := map[string]interface{}{
				"name": "Only name updated",
			}

			// Make request
			w := makeAuthenticatedRequest(t, handler, "PUT", "/account", updateData, token)

			// Assert successful response
			assert.Equal(t, 200, w.Code)

			// Verify only name was updated, other fields should remain
			var updatedAccount models.Account
			err = db.Where("user_id = ?", user.ID).First(&updatedAccount).Error
			require.NoError(t, err)

			assert.Equal(t, "Only name updated", updatedAccount.Name)
			assert.Equal(t, "Original description", updatedAccount.Description) // Should remain unchanged
		})
	})
}

func TestAccountHandler_Authorization(t *testing.T) {
	handler, db := setupTestHandler(t)

	t.Run("User can only access their own account", func(t *testing.T) {
		// Create two users with accounts
		user1, token1 := createTestUser(t, db, "accountuser1@example.com")
		user2, token2 := createTestUser(t, db, "accountuser2@example.com")

		// Update accounts with different data
		var account1, account2 models.Account
		err := db.Where("user_id = ?", user1.ID).First(&account1).Error
		require.NoError(t, err)
		err = db.Where("user_id = ?", user2.ID).First(&account2).Error
		require.NoError(t, err)

		account1.Name = "User 1 Account"
		account1.Description = "User 1 Description"
		account2.Name = "User 2 Account"
		account2.Description = "User 2 Description"
		err = db.Save(&account1).Error
		require.NoError(t, err)
		err = db.Save(&account2).Error
		require.NoError(t, err)

		// User 1 should get their own account
		w1 := makeAuthenticatedRequest(t, handler, "GET", "/account", nil, token1)
		assert.Equal(t, 200, w1.Code)

		// Verify user1 gets their own account data
		var response1 map[string]interface{}
		err = json.Unmarshal(w1.Body.Bytes(), &response1)
		require.NoError(t, err)

		if data1, hasData := response1["data"]; hasData {
			accountData1 := data1.(map[string]interface{})
			assert.Equal(t, float64(account1.ID), accountData1["id"])
			assert.Equal(t, "User 1 Account", accountData1["name"])
			assert.Equal(t, "User 1 Description", accountData1["description"])
		}

		// User 2 should get their own account
		w2 := makeAuthenticatedRequest(t, handler, "GET", "/account", nil, token2)
		assert.Equal(t, 200, w2.Code)

		// Verify user2 gets their own account data
		var response2 map[string]interface{}
		err = json.Unmarshal(w2.Body.Bytes(), &response2)
		require.NoError(t, err)

		if data2, hasData := response2["data"]; hasData {
			accountData2 := data2.(map[string]interface{})
			assert.Equal(t, float64(account2.ID), accountData2["id"])
			assert.Equal(t, "User 2 Account", accountData2["name"])
			assert.Equal(t, "User 2 Description", accountData2["description"])
		}
	})

	t.Run("User can only update their own account", func(t *testing.T) {
		// Create two users with accounts
		user1, token1 := createTestUser(t, db, "accountupdate1@example.com")
		user2, _ := createTestUser(t, db, "accountupdate2@example.com")

		// Set initial data
		var account1, account2 models.Account
		err := db.Where("user_id = ?", user1.ID).First(&account1).Error
		require.NoError(t, err)
		err = db.Where("user_id = ?", user2.ID).First(&account2).Error
		require.NoError(t, err)

		account1.Name = "User 1 original name"
		account2.Name = "User 2 original name"
		db.Save(&account1)
		db.Save(&account2)

		updateData := map[string]interface{}{
			"name": "Updated by user 1",
		}

		// User 1 updates their account
		w := makeAuthenticatedRequest(t, handler, "PUT", "/account", updateData, token1)
		assert.Equal(t, 200, w.Code)

		// Verify user1's account was updated
		var updatedAccount1 models.Account
		err = db.Where("user_id = ?", user1.ID).First(&updatedAccount1).Error
		require.NoError(t, err)
		assert.Equal(t, "Updated by user 1", updatedAccount1.Name)

		// Verify user2's account was NOT affected
		var unchangedAccount2 models.Account
		err = db.Where("user_id = ?", user2.ID).First(&unchangedAccount2).Error
		require.NoError(t, err)
		assert.Equal(t, "User 2 original name", unchangedAccount2.Name)
	})
}

func TestAccountHandler_PlanIntegration(t *testing.T) {
	handler, db := setupTestHandler(t)

	t.Run("Account should have default plan", func(t *testing.T) {
		// Create a user (which creates an account via AfterCreate hook)
		user, token := createTestUser(t, db, "defaultplan@example.com")

		// Make request to get account
		w := makeAuthenticatedRequest(t, handler, "GET", "/account", nil, token)
		assert.Equal(t, 200, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		// Check if response has 'data' wrapper
		if data, hasData := response["data"]; hasData {
			accountData := data.(map[string]interface{})

			// Verify plan relationship exists
			assert.Contains(t, accountData, "plan")
			planData := accountData["plan"].(map[string]interface{})
			assert.Contains(t, planData, "name")
			assert.Equal(t, "Free", planData["name"]) // Should be the default plan we created
		}

		// Verify in database that account has the default plan
		var account models.Account
		err = db.Preload("Plan").Where("user_id = ?", user.ID).First(&account).Error
		require.NoError(t, err)
		assert.NotNil(t, account.Plan)
		assert.Equal(t, "Free", account.Plan.Name)
	})

	t.Run("Update account with different plan", func(t *testing.T) {
		// Create another plan for testing - use Price: 0 to avoid Stripe integration
		premiumPlan := &models.Plan{
			Name:          "Premium",
			Price:         0, // Set to 0 to avoid Stripe API calls in tests
			BillingPeriod: "monthly",
			Description:   "Premium plan for testing",
		}
		err := db.Create(premiumPlan).Error
		require.NoError(t, err)

		// Create a user
		user, token := createTestUser(t, db, "updateplan@example.com")

		// Update account with premium plan
		updateData := map[string]interface{}{
			"name":    "Premium Account",
			"plan_id": premiumPlan.ID,
		}

		w := makeAuthenticatedRequest(t, handler, "PUT", "/account", updateData, token)
		assert.Equal(t, 200, w.Code)

		// Verify the account was updated with the new plan
		var updatedAccount models.Account
		err = db.Preload("Plan").Where("user_id = ?", user.ID).First(&updatedAccount).Error
		require.NoError(t, err)
		assert.Equal(t, premiumPlan.ID, updatedAccount.PlanID)
		assert.Equal(t, "Premium", updatedAccount.Plan.Name)
	})
}
