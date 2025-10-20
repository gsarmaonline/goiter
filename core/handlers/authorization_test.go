package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gsarmaonline/goiter/config"
	"github.com/gsarmaonline/goiter/core/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Authorization test setup
func setupAuthorizationTest() (*gorm.DB, *gin.Engine, *Handler) {
	os.Setenv("JWT_SECRET", "test_jwt_secret_key")

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to test database")
	}

	db.AutoMigrate(&models.User{}, &models.Account{}, &models.Plan{}, &models.Profile{}, &models.Group{})

	gin.SetMode(gin.TestMode)
	router := gin.New()

	cfg := &config.Config{}
	handler := NewHandler(router, db, cfg)

	return db, router, handler
}

func createAuthTestUser(db *gorm.DB, email string) (*models.User, string) {
	user := &models.User{
		Email:    email,
		GoogleID: "google_" + email,
	}
	db.Create(user)

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
	})
	tokenString, _ := token.SignedString([]byte("test_jwt_secret_key"))

	return user, tokenString
}

func makeAuthRequest(handler *Handler, method, path string, handlerFunc gin.HandlerFunc, token string, body interface{}) *httptest.ResponseRecorder {
	var reqBody []byte
	if body != nil {
		reqBody, _ = json.Marshal(body)
	}

	req := httptest.NewRequest(method, path, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Set user context if authenticated
	if token != "" {
		userEmail := extractUserEmailFromToken(token)
		if userEmail != "" {
			var user models.User
			handler.Db.Where("email = ?", userEmail).First(&user)
			c.Set("user", &user)
		}
	}

	handlerFunc(c)
	return w
}

func extractUserEmailFromToken(tokenString string) string {
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("test_jwt_secret_key"), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if email, exists := claims["email"]; exists {
			return email.(string)
		}
	}
	return ""
}

// Test authorization helper functions
func TestAuthorization_UserScopedDB(t *testing.T) {
	db, _, handler := setupAuthorizationTest()
	defer db.Migrator().DropTable(&models.User{}, &models.Account{}, &models.Plan{}, &models.Profile{})

	// Create test users
	user1, token1 := createAuthTestUser(db, "user1@example.com")
	user2, token2 := createAuthTestUser(db, "user2@example.com")

	// Create test plan
	testPlan := &models.Plan{
		Name:  "Test Plan",
		Price: 1000,
	}
	db.Create(testPlan)

	// Create accounts for each user
	account1 := &models.Account{
		Name:   "User1 Account",
		PlanID: testPlan.ID,
	}
	account1.UserID = user1.ID
	db.Create(account1)

	account2 := &models.Account{
		Name:   "User2 Account",
		PlanID: testPlan.ID,
	}
	account2.UserID = user2.ID
	db.Create(account2)

	t.Run("User can only access their own account", func(t *testing.T) {
		// Create account handlers
		accountHandler := NewAccountHandler(handler)

		// Test user1 access
		w1 := makeAuthRequest(handler, "GET", "/account", accountHandler.GetAccount, token1, nil)
		assert.Equal(t, http.StatusOK, w1.Code)

		var response1 map[string]interface{}
		err := json.Unmarshal(w1.Body.Bytes(), &response1)
		assert.NoError(t, err)
		assert.Equal(t, "User1 Account", response1["name"])

		// Test user2 access
		w2 := makeAuthRequest(handler, "GET", "/account", accountHandler.GetAccount, token2, nil)
		assert.Equal(t, http.StatusOK, w2.Code)

		var response2 map[string]interface{}
		err = json.Unmarshal(w2.Body.Bytes(), &response2)
		assert.NoError(t, err)
		assert.Equal(t, "User2 Account", response2["name"])

		// Verify they're different
		assert.NotEqual(t, response1["name"], response2["name"])
	})
}

func TestAuthorization_ProfileAccess(t *testing.T) {
	db, _, handler := setupAuthorizationTest()
	defer db.Migrator().DropTable(&models.User{}, &models.Account{}, &models.Plan{}, &models.Profile{})

	// Create test users
	user1, token1 := createAuthTestUser(db, "profile1@example.com")
	user2, token2 := createAuthTestUser(db, "profile2@example.com")

	// Create profiles
	profile1 := &models.Profile{
		CompanyName: "Company One",
		JobTitle:    "Developer",
	}
	profile1.UserID = user1.ID
	db.Create(profile1)

	profile2 := &models.Profile{
		CompanyName: "Company Two",
		JobTitle:    "Manager",
	}
	profile2.UserID = user2.ID
	db.Create(profile2)

	t.Run("Users can only access their own profiles", func(t *testing.T) {
		// User1 gets their profile
		w1 := makeAuthRequest(handler, "GET", "/profile", handler.handleGetProfile, token1, nil)
		assert.Equal(t, http.StatusOK, w1.Code)

		var response1 map[string]interface{}
		err := json.Unmarshal(w1.Body.Bytes(), &response1)
		assert.NoError(t, err)
		assert.Equal(t, "Company One", response1["company_name"])

		// User2 gets their profile
		w2 := makeAuthRequest(handler, "GET", "/profile", handler.handleGetProfile, token2, nil)
		assert.Equal(t, http.StatusOK, w2.Code)

		var response2 map[string]interface{}
		err = json.Unmarshal(w2.Body.Bytes(), &response2)
		assert.NoError(t, err)
		assert.Equal(t, "Company Two", response2["company_name"])

		// Verify isolation
		assert.NotEqual(t, response1["company_name"], response2["company_name"])
	})

	t.Run("Users can only update their own profiles", func(t *testing.T) {
		updateData := map[string]string{
			"company_name": "Updated Company",
		}

		// User1 updates their profile
		w1 := makeAuthRequest(handler, "PUT", "/profile", handler.handleUpdateProfile, token1, updateData)
		assert.Equal(t, http.StatusOK, w1.Code)

		// Verify update worked
		var updatedProfile models.Profile
		err := db.Where("user_id = ?", user1.ID).First(&updatedProfile).Error
		assert.NoError(t, err)
		assert.Equal(t, "Updated Company", updatedProfile.CompanyName)
		assert.Equal(t, "Developer", updatedProfile.JobTitle) // JobTitle unchanged

		// Verify user2's profile is unchanged
		var user2Profile models.Profile
		err = db.Where("user_id = ?", user2.ID).First(&user2Profile).Error
		assert.NoError(t, err)
		assert.Equal(t, "Company Two", user2Profile.CompanyName) // Should be unchanged
		assert.Equal(t, "Manager", user2Profile.JobTitle)
	})
}

func TestAuthorization_UnauthenticatedAccess(t *testing.T) {
	db, _, handler := setupAuthorizationTest()
	defer db.Migrator().DropTable(&models.User{}, &models.Account{}, &models.Plan{}, &models.Profile{})

	t.Run("UserScopedDB requires user context", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		// No user set in context

		// This should panic or fail when trying to get user from empty context
		assert.Panics(t, func() {
			handler.UserScopedDB(c)
		})
	})
}

func TestAuthorization_ResourceOwnership(t *testing.T) {
	db, _, handler := setupAuthorizationTest()
	defer db.Migrator().DropTable(&models.User{}, &models.Account{}, &models.Plan{}, &models.Profile{})

	// Create test users
	user1, _ := createAuthTestUser(db, "owner1@example.com")
	user2, _ := createAuthTestUser(db, "owner2@example.com")

	// Create test plan
	testPlan := &models.Plan{
		Name:  "Test Plan",
		Price: 1000,
	}
	db.Create(testPlan)

	// Create account owned by user1
	account := &models.Account{
		Name:   "Test Account",
		PlanID: testPlan.ID,
	}
	account.UserID = user1.ID
	db.Create(account)

	t.Run("Authorization helper correctly identifies resource ownership", func(t *testing.T) {
		// Test with owner
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("user", user1)

		canAccess := handler.authorisation.CanAccessResource(c, account)
		assert.True(t, canAccess, "Owner should be able to access their resource")

		// Test with non-owner
		c.Set("user", user2)
		canAccess = handler.authorisation.CanAccessResource(c, account)
		assert.False(t, canAccess, "Non-owner should not be able to access resource")
	})

	t.Run("UserScopedDB returns only user's resources", func(t *testing.T) {
		// Create another account for user2
		account2 := &models.Account{
			Name:   "User2 Account",
			PlanID: testPlan.ID,
		}
		account2.UserID = user2.ID
		db.Create(account2)

		// Test user1's scoped query
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("user", user1)

		scopedDB := handler.authorisation.UserScopedDB(c, db)
		var user1Accounts []models.Account
		err := scopedDB.Find(&user1Accounts).Error
		assert.NoError(t, err)
		assert.Len(t, user1Accounts, 1)
		assert.Equal(t, "Test Account", user1Accounts[0].Name)

		// Test user2's scoped query
		c.Set("user", user2)
		scopedDB = handler.authorisation.UserScopedDB(c, db)
		var user2Accounts []models.Account
		err = scopedDB.Find(&user2Accounts).Error
		assert.NoError(t, err)
		assert.Len(t, user2Accounts, 1)
		assert.Equal(t, "User2 Account", user2Accounts[0].Name)
	})
}

func TestAuthorization_UpdateWithUser(t *testing.T) {
	db, _, handler := setupAuthorizationTest()
	defer db.Migrator().DropTable(&models.User{}, &models.Account{}, &models.Plan{}, &models.Profile{})

	user, _ := createAuthTestUser(db, "update@example.com")

	t.Run("UpdateWithUser sets correct user ID", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/test", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("user", user)

		// Create a new profile
		profile := &models.Profile{
			CompanyName: "Test Company",
			JobTitle:    "Developer",
		}

		// Use authorization helper to set user ID
		err := handler.authorisation.UpdateWithUser(c, profile)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, profile.UserID)
	})
}

// Integration test for authorization across multiple resources
func TestAuthorization_DatabaseIntegration(t *testing.T) {
	db, _, handler := setupAuthorizationTest()
	defer db.Migrator().DropTable(&models.User{}, &models.Account{}, &models.Plan{}, &models.Profile{})

	// Create test users
	user1, _ := createAuthTestUser(db, "multi1@example.com")
	user2, _ := createAuthTestUser(db, "multi2@example.com")

	// Create test plan
	testPlan := &models.Plan{
		Name:  "Test Plan",
		Price: 1000,
	}
	db.Create(testPlan)

	// Create complete user setup (account + profile)
	account1 := &models.Account{
		Name:   "User1 Account",
		PlanID: testPlan.ID,
	}
	account1.UserID = user1.ID
	db.Create(account1)

	profile1 := &models.Profile{
		CompanyName: "Company One",
		JobTitle:    "Developer",
	}
	profile1.UserID = user1.ID
	db.Create(profile1)

	account2 := &models.Account{
		Name:   "User2 Account",
		PlanID: testPlan.ID,
	}
	account2.UserID = user2.ID
	db.Create(account2)

	profile2 := &models.Profile{
		CompanyName: "Company Two",
		JobTitle:    "Manager",
	}
	profile2.UserID = user2.ID
	db.Create(profile2)

	t.Run("Database queries are properly scoped by user", func(t *testing.T) {
		// Test user1's scoped queries
		req1 := httptest.NewRequest("GET", "/test", nil)
		w1 := httptest.NewRecorder()
		c1, _ := gin.CreateTestContext(w1)
		c1.Request = req1
		c1.Set("user", user1)

		user1DB := handler.UserScopedDB(c1)

		var user1Accounts []models.Account
		var user1Profiles []models.Profile

		err := user1DB.Find(&user1Accounts).Error
		assert.NoError(t, err)
		assert.Len(t, user1Accounts, 1)
		assert.Equal(t, "User1 Account", user1Accounts[0].Name)

		err = user1DB.Find(&user1Profiles).Error
		assert.NoError(t, err)
		assert.Len(t, user1Profiles, 1)
		assert.Equal(t, "Company One", user1Profiles[0].CompanyName)

		// Test user2's scoped queries
		req2 := httptest.NewRequest("GET", "/test", nil)
		w2 := httptest.NewRecorder()
		c2, _ := gin.CreateTestContext(w2)
		c2.Request = req2
		c2.Set("user", user2)

		user2DB := handler.UserScopedDB(c2)

		var user2Accounts []models.Account
		var user2Profiles []models.Profile

		err = user2DB.Find(&user2Accounts).Error
		assert.NoError(t, err)
		assert.Len(t, user2Accounts, 1)
		assert.Equal(t, "User2 Account", user2Accounts[0].Name)

		err = user2DB.Find(&user2Profiles).Error
		assert.NoError(t, err)
		assert.Len(t, user2Profiles, 1)
		assert.Equal(t, "Company Two", user2Profiles[0].CompanyName)

		// Ensure complete isolation
		assert.NotEqual(t, user1Accounts[0].Name, user2Accounts[0].Name)
		assert.NotEqual(t, user1Profiles[0].CompanyName, user2Profiles[0].CompanyName)
	})
}
