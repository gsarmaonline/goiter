package handlers

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/gsarmaonline/goiter/config"
	"github.com/gsarmaonline/goiter/core/models"
)

// Test helper functions
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(
		&models.User{},
		&models.Profile{},
		&models.Account{},
		&models.Plan{},
		&models.Feature{},
		&models.PlanFeature{},
	)
	require.NoError(t, err)

	// Create a default plan to satisfy the Account BeforeCreate hook
	defaultPlan := &models.Plan{
		Name:          "Free",
		Price:         0,
		BillingPeriod: "monthly",
		Description:   "Default free plan for testing",
	}
	err = db.Create(defaultPlan).Error
	require.NoError(t, err)

	return db
}

func setupTestHandler(t *testing.T) (*Handler, *gorm.DB) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB(t)

	// Set required environment variables
	os.Setenv("JWT_SECRET", "test-secret-key")

	cfg := &config.Config{Mode: config.ModeDev}
	router := gin.New()
	handler := NewHandler(router, db, cfg)

	return handler, db
}

func createTestUser(t *testing.T, db *gorm.DB, email string) (*models.User, string) {
	user := &models.User{
		Email:       email,
		Name:        "Test User",
		GoogleID:    "test-google-id-" + email, // Make GoogleID unique
		UserStatus:  models.ActiveUser,
		CreatedFrom: "test",
	}

	err := db.Create(user).Error
	require.NoError(t, err)

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"email": user.Email,
			"exp":   time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString([]byte("test-secret-key"))
	require.NoError(t, err)

	return user, tokenString
}

func makeAuthenticatedRequest(t *testing.T, handler *Handler, method, path string, body interface{}, token string) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	if body != nil {
		jsonBody, err := json.Marshal(body)
		require.NoError(t, err)
		reqBody = bytes.NewBuffer(jsonBody)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req := httptest.NewRequest(method, path, reqBody)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	handler.router.ServeHTTP(w, req)

	return w
}

func TestProfileHandler(t *testing.T) {
	handler, db := setupTestHandler(t)

	t.Run("GetProfile", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			// Create a user with profile
			user, token := createTestUser(t, db, "testuser@example.com")

			// Get the profile that was created by the user AfterCreate hook
			var profile models.Profile
			err := db.Where("user_id = ?", user.ID).First(&profile).Error
			require.NoError(t, err)

			// Update some profile data to verify it's returned
			profile.Address = "123 Test St"
			profile.City = "Test City"
			profile.CompanyName = "Test Company"
			err = db.Save(&profile).Error
			require.NoError(t, err)

			// Make request
			w := makeAuthenticatedRequest(t, handler, "GET", "/profile", nil, token)

			// Assert response
			assert.Equal(t, 200, w.Code)

			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			// Check if response has 'data' wrapper
			if data, hasData := response["data"]; hasData {
				profileData := data.(map[string]interface{})
				assert.Equal(t, float64(profile.ID), profileData["id"])
				assert.Equal(t, float64(user.ID), profileData["user_id"])
				assert.Equal(t, "123 Test St", profileData["address"])
				assert.Equal(t, "Test City", profileData["city"])
				assert.Equal(t, "Test Company", profileData["company_name"])
			}
		})

		t.Run("Unauthenticated", func(t *testing.T) {
			// Make request without authentication
			w := makeAuthenticatedRequest(t, handler, "GET", "/profile", nil, "")

			// Assert error response
			assert.Equal(t, 401, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)
			assert.Contains(t, response["error"].(string), "Not authenticated")
		})

		t.Run("Profile not found", func(t *testing.T) {
			// Create user but delete the profile to simulate missing profile
			user, token := createTestUser(t, db, "noprofile@example.com")

			// Delete the profile that was created by the user AfterCreate hook
			err := db.Where("user_id = ?", user.ID).Delete(&models.Profile{}).Error
			require.NoError(t, err)

			// Make request
			w := makeAuthenticatedRequest(t, handler, "GET", "/profile", nil, token)

			// Assert error response
			assert.Equal(t, 404, w.Code)

			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)
			assert.Contains(t, response["error"].(string), "Profile not found")
		})

		t.Run("Invalid JWT", func(t *testing.T) {
			// Make request with invalid token
			w := makeAuthenticatedRequest(t, handler, "GET", "/profile", nil, "invalid.jwt.token")

			// Assert error response
			assert.Equal(t, 401, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)
			assert.Contains(t, response["error"].(string), "Invalid token")
		})
	})

	t.Run("UpdateProfile", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			// Create a user with profile
			user, token := createTestUser(t, db, "update@example.com")

			// Prepare update data using actual Profile model fields
			updateData := map[string]interface{}{
				"address":      "456 Update St",
				"city":         "Update City",
				"company_name": "Updated Company",
				"job_title":    "Senior Developer",
			}

			// Make request
			w := makeAuthenticatedRequest(t, handler, "PUT", "/profile", updateData, token)

			// Assert successful response
			assert.Equal(t, 200, w.Code)

			// Verify the profile was updated in database
			var updatedProfile models.Profile
			err := db.Where("user_id = ?", user.ID).First(&updatedProfile).Error
			require.NoError(t, err)

			assert.Equal(t, "456 Update St", updatedProfile.Address)
			assert.Equal(t, "Update City", updatedProfile.City)
			assert.Equal(t, "Updated Company", updatedProfile.CompanyName)
			assert.Equal(t, "Senior Developer", updatedProfile.JobTitle)
		})

		t.Run("Unauthenticated", func(t *testing.T) {
			updateData := map[string]interface{}{
				"address": "Should not work",
			}

			// Make request without authentication
			w := makeAuthenticatedRequest(t, handler, "PUT", "/profile", updateData, "")

			// Assert error response
			assert.Equal(t, 401, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)
			assert.Contains(t, response["error"].(string), "Not authenticated")
		})

		t.Run("Profile not found", func(t *testing.T) {
			// Create user but delete the profile
			user, token := createTestUser(t, db, "noprofileupdate@example.com")

			// Delete the profile that was created by the user AfterCreate hook
			err := db.Where("user_id = ?", user.ID).Delete(&models.Profile{}).Error
			require.NoError(t, err)

			updateData := map[string]interface{}{
				"address": "Should not work",
			}

			// Make request
			w := makeAuthenticatedRequest(t, handler, "PUT", "/profile", updateData, token)

			// Assert error response
			assert.Equal(t, 404, w.Code)

			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)
			assert.Contains(t, response["error"].(string), "Profile not found")
		})

		t.Run("Partial update", func(t *testing.T) {
			// Create a user with profile
			user, token := createTestUser(t, db, "partial@example.com")

			// Get and set some initial data
			var profile models.Profile
			err := db.Where("user_id = ?", user.ID).First(&profile).Error
			require.NoError(t, err)

			profile.Address = "Original address"
			profile.City = "Original city"
			profile.CompanyName = "Original company"
			err = db.Save(&profile).Error
			require.NoError(t, err)

			// Prepare partial update data (only address)
			updateData := map[string]interface{}{
				"address": "Only address updated",
			}

			// Make request
			w := makeAuthenticatedRequest(t, handler, "PUT", "/profile", updateData, token)

			// Assert successful response
			assert.Equal(t, 200, w.Code)

			// Verify only address was updated, other fields should remain
			var updatedProfile models.Profile
			err = db.Where("user_id = ?", user.ID).First(&updatedProfile).Error
			require.NoError(t, err)

			assert.Equal(t, "Only address updated", updatedProfile.Address)
			assert.Equal(t, "Original city", updatedProfile.City)           // Should remain unchanged
			assert.Equal(t, "Original company", updatedProfile.CompanyName) // Should remain unchanged
		})
	})
}

func TestProfileHandler_Authorization(t *testing.T) {
	handler, db := setupTestHandler(t)

	t.Run("User can only access their own profile", func(t *testing.T) {
		// Create two users with profiles
		user1, token1 := createTestUser(t, db, "user1@example.com")
		user2, token2 := createTestUser(t, db, "user2@example.com")

		// Update profiles with different data
		var profile1, profile2 models.Profile
		err := db.Where("user_id = ?", user1.ID).First(&profile1).Error
		require.NoError(t, err)
		err = db.Where("user_id = ?", user2.ID).First(&profile2).Error
		require.NoError(t, err)

		profile1.Address = "User 1 address"
		profile2.Address = "User 2 address"
		err = db.Save(&profile1).Error
		require.NoError(t, err)
		err = db.Save(&profile2).Error
		require.NoError(t, err)

		// User 1 should get their own profile
		w1 := makeAuthenticatedRequest(t, handler, "GET", "/profile", nil, token1)
		assert.Equal(t, 200, w1.Code)

		// Verify user1 gets their own profile data
		var response1 map[string]interface{}
		err = json.Unmarshal(w1.Body.Bytes(), &response1)
		require.NoError(t, err)

		if data1, hasData := response1["data"]; hasData {
			profileData1 := data1.(map[string]interface{})
			assert.Equal(t, float64(profile1.ID), profileData1["id"])
			assert.Equal(t, "User 1 address", profileData1["address"])
		}

		// User 2 should get their own profile
		w2 := makeAuthenticatedRequest(t, handler, "GET", "/profile", nil, token2)
		assert.Equal(t, 200, w2.Code)

		// Verify user2 gets their own profile data
		var response2 map[string]interface{}
		err = json.Unmarshal(w2.Body.Bytes(), &response2)
		require.NoError(t, err)

		if data2, hasData := response2["data"]; hasData {
			profileData2 := data2.(map[string]interface{})
			assert.Equal(t, float64(profile2.ID), profileData2["id"])
			assert.Equal(t, "User 2 address", profileData2["address"])
		}
	})

	t.Run("User can only update their own profile", func(t *testing.T) {
		// Create two users with profiles
		user1, token1 := createTestUser(t, db, "update1@example.com")
		user2, _ := createTestUser(t, db, "update2@example.com")

		// Set initial data
		var profile1, profile2 models.Profile
		err := db.Where("user_id = ?", user1.ID).First(&profile1).Error
		require.NoError(t, err)
		err = db.Where("user_id = ?", user2.ID).First(&profile2).Error
		require.NoError(t, err)

		profile1.Address = "User 1 original address"
		profile2.Address = "User 2 original address"
		db.Save(&profile1)
		db.Save(&profile2)

		updateData := map[string]interface{}{
			"address": "Updated by user 1",
		}

		// User 1 updates their profile
		w := makeAuthenticatedRequest(t, handler, "PUT", "/profile", updateData, token1)
		assert.Equal(t, 200, w.Code)

		// Verify user1's profile was updated
		var updatedProfile1 models.Profile
		err = db.Where("user_id = ?", user1.ID).First(&updatedProfile1).Error
		require.NoError(t, err)
		assert.Equal(t, "Updated by user 1", updatedProfile1.Address)

		// Verify user2's profile was NOT affected
		var unchangedProfile2 models.Profile
		err = db.Where("user_id = ?", user2.ID).First(&unchangedProfile2).Error
		require.NoError(t, err)
		assert.Equal(t, "User 2 original address", unchangedProfile2.Address)
	})
}
