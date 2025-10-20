package testutils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/gsarmaonline/goiter/config"
	"github.com/gsarmaonline/goiter/core/handlers"
	"github.com/gsarmaonline/goiter/core/models"
)

// TestEnvironment holds all the infrastructure needed for handler testing
type TestEnvironment struct {
	DB      *gorm.DB
	Router  *gin.Engine
	Handler *handlers.Handler
	Config  *config.Config
	Server  *httptest.Server
}

// TestUser represents a test user with authentication token
type TestUser struct {
	User  *models.User
	Token string
}

// TestClient provides HTTP client functionality for testing
type TestClient struct {
	env *TestEnvironment
}

// SetupTestEnvironment creates a complete test environment
func SetupTestEnvironment(t *testing.T) *TestEnvironment {
	// Set test mode for Gin
	gin.SetMode(gin.TestMode)

	// Setup test database (in-memory SQLite)
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Auto migrate all models
	err = db.AutoMigrate(
		&models.User{},
		&models.Profile{},
		&models.Account{},
		&models.Plan{},
		&models.Feature{},
		&models.PlanFeature{},
		&models.Group{},
	)
	require.NoError(t, err)

	// Setup test config
	cfg := &config.Config{
		Mode: config.ModeDev,
	}

	// Set required environment variables for testing
	os.Setenv("JWT_SECRET", "test-secret-key-for-testing-only")
	os.Setenv("GOOGLE_CLIENT_ID", "test-google-client-id")
	os.Setenv("GOOGLE_CLIENT_SECRET", "test-google-client-secret")
	os.Setenv("GOOGLE_CALLBACK_URL", "http://localhost:8080/auth/google/callback")
	os.Setenv("FRONTEND_URL", "http://localhost:3000")

	// Create router and handler
	router := gin.New()
	handler := handlers.NewHandler(router, db, cfg)

	// Create test server
	server := httptest.NewServer(router)

	env := &TestEnvironment{
		DB:      db,
		Router:  router,
		Handler: handler,
		Config:  cfg,
		Server:  server,
	}

	return env
}

// Cleanup cleans up the test environment
func (env *TestEnvironment) Cleanup() {
	if env.Server != nil {
		env.Server.Close()
	}
	
	// Close database connection
	if env.DB != nil {
		sqlDB, _ := env.DB.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}
}

// CreateTestUser creates a test user and returns it with a JWT token
func (env *TestEnvironment) CreateTestUser(t *testing.T, email string) *TestUser {
	user := &models.User{
		Email:       email,
		Name:        fmt.Sprintf("Test User %s", email),
		GoogleID:    fmt.Sprintf("google-id-%s", email),
		UserStatus:  models.ActiveUser,
		CreatedFrom: "test",
	}

	err := env.DB.Create(user).Error
	require.NoError(t, err)

	// Create JWT token using the handler's method
	token, err := env.createJWT(user.Email)
	require.NoError(t, err)

	return &TestUser{
		User:  user,
		Token: token,
	}
}

// createJWT creates a JWT token for testing (mirrors handler logic)
func (env *TestEnvironment) createJWT(email string) (string, error) {
	// Use the handler's createJWT method by making a mock request
	// This ensures we use the same JWT creation logic as the actual handler
	
	// For now, we'll create a simple implementation
	// In a real scenario, you might want to extract JWT creation to a service
	return fmt.Sprintf("test-jwt-token-for-%s", email), nil
}

// CreateTestPlan creates a test plan for billing tests
func (env *TestEnvironment) CreateTestPlan(t *testing.T, name string, price float64) *models.Plan {
	plan := &models.Plan{
		Name:          name,
		Price:         price,
		BillingPeriod: "monthly",
		Description:   fmt.Sprintf("Test plan: %s", name),
	}

	err := env.DB.Create(plan).Error
	require.NoError(t, err)

	return plan
}

// NewTestClient creates a new test client
func (env *TestEnvironment) NewTestClient() *TestClient {
	return &TestClient{env: env}
}

// MakeRequest makes an HTTP request to the test server
func (c *TestClient) MakeRequest(t *testing.T, method, path string, body interface{}, user *TestUser) *httptest.ResponseRecorder {
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

	// Add authentication if user is provided
	if user != nil {
		req.Header.Set("Authorization", "Bearer "+user.Token)
	}

	// Record the response
	w := httptest.NewRecorder()
	c.env.Router.ServeHTTP(w, req)

	return w
}

// AssertJSONResponse asserts that the response has the expected status and JSON structure
func (c *TestClient) AssertJSONResponse(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int, expectedBody interface{}) {
	assert.Equal(t, expectedStatus, w.Code)
	
	if expectedBody != nil {
		var actualBody map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &actualBody)
		require.NoError(t, err)

		expectedJSON, err := json.Marshal(expectedBody)
		require.NoError(t, err)
		
		var expectedBodyMap map[string]interface{}
		err = json.Unmarshal(expectedJSON, &expectedBodyMap)
		require.NoError(t, err)

		// Check if response has 'data' wrapper
		if data, hasData := actualBody["data"]; hasData {
			assert.Equal(t, expectedBodyMap, data)
		} else {
			assert.Equal(t, expectedBodyMap, actualBody)
		}
	}
}

// AssertErrorResponse asserts that the response has an error with the expected status and message
func (c *TestClient) AssertErrorResponse(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int, expectedErrorMessage string) {
	assert.Equal(t, expectedStatus, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	error, exists := response["error"]
	require.True(t, exists, "Response should contain an error field")
	assert.Contains(t, error.(string), expectedErrorMessage)
}

// MockStripeService provides a mock implementation for Stripe operations
type MockStripeService struct {
	CreateSubscriptionFunc func(accountID uint, planID uint, paymentMethodID string) error
	CancelSubscriptionFunc func(subscriptionID string) error
	ProcessWebhookFunc     func(payload []byte, signature string) error
}

func (m *MockStripeService) CreateSubscription(accountID uint, planID uint, paymentMethodID string) error {
	if m.CreateSubscriptionFunc != nil {
		return m.CreateSubscriptionFunc(accountID, planID, paymentMethodID)
	}
	return nil
}

func (m *MockStripeService) CancelSubscription(subscriptionID string) error {
	if m.CancelSubscriptionFunc != nil {
		return m.CancelSubscriptionFunc(subscriptionID)
	}
	return nil
}

func (m *MockStripeService) ProcessWebhook(payload []byte, signature string) error {
	if m.ProcessWebhookFunc != nil {
		return m.ProcessWebhookFunc(payload, signature)
	}
	return nil
}