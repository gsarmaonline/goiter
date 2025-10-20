package testutils

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
	"github.com/gsarmaonline/goiter/core/models"
)

// AuthTestHelper provides authentication-related test utilities
type AuthTestHelper struct {
	env *TestEnvironment
}

// NewAuthTestHelper creates a new authentication test helper
func NewAuthTestHelper(env *TestEnvironment) *AuthTestHelper {
	return &AuthTestHelper{env: env}
}

// CreateJWTForUser creates a valid JWT token for the given user
func (h *AuthTestHelper) CreateJWTForUser(t *testing.T, user *models.User) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"email": user.Email,
			"exp":   time.Now().Add(time.Hour * 24).Unix(),
		})

	secret := "test-secret-key-for-testing-only"
	tokenString, err := token.SignedString([]byte(secret))
	require.NoError(t, err)

	return tokenString
}

// CreateExpiredJWT creates an expired JWT token for testing
func (h *AuthTestHelper) CreateExpiredJWT(t *testing.T, email string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"email": email,
			"exp":   time.Now().Add(-time.Hour).Unix(), // Expired 1 hour ago
		})

	secret := "test-secret-key-for-testing-only"
	tokenString, err := token.SignedString([]byte(secret))
	require.NoError(t, err)

	return tokenString
}

// CreateInvalidJWT creates an invalid JWT token for testing
func (h *AuthTestHelper) CreateInvalidJWT(t *testing.T) string {
	return "invalid.jwt.token"
}

// CreateUserWithProfile creates a user and ensures they have a profile
func (h *AuthTestHelper) CreateUserWithProfile(t *testing.T, email string) (*TestUser, *models.Profile) {
	testUser := h.env.CreateTestUser(t, email)
	
	// Verify profile was created by the AfterCreate hook
	var profile models.Profile
	err := h.env.DB.Where("user_id = ?", testUser.User.ID).First(&profile).Error
	require.NoError(t, err)

	return testUser, &profile
}

// CreateUserWithAccount creates a user and ensures they have an account
func (h *AuthTestHelper) CreateUserWithAccount(t *testing.T, email string) (*TestUser, *models.Account) {
	testUser := h.env.CreateTestUser(t, email)
	
	// Verify account was created by the AfterCreate hook
	var account models.Account
	err := h.env.DB.Where("user_id = ?", testUser.User.ID).First(&account).Error
	require.NoError(t, err)

	return testUser, &account
}

// TestAuthenticationScenarios contains common authentication test scenarios
type TestAuthenticationScenarios struct {
	ValidUser       *TestUser
	ValidToken      string
	ExpiredToken    string
	InvalidToken    string
	MalformedToken  string
	NonexistentUser *TestUser
}

// SetupAuthenticationScenarios prepares all common authentication test scenarios
func (h *AuthTestHelper) SetupAuthenticationScenarios(t *testing.T) *TestAuthenticationScenarios {
	// Create a valid user
	validUser := h.env.CreateTestUser(t, "valid@example.com")
	
	// Create a user that doesn't exist in the database (for JWT with nonexistent user)
	nonexistentUser := &TestUser{
		User: &models.User{
			Email: "nonexistent@example.com",
		},
		Token: h.CreateJWTForUser(t, &models.User{Email: "nonexistent@example.com"}),
	}

	return &TestAuthenticationScenarios{
		ValidUser:       validUser,
		ValidToken:      validUser.Token,
		ExpiredToken:    h.CreateExpiredJWT(t, validUser.User.Email),
		InvalidToken:    h.CreateInvalidJWT(t),
		MalformedToken:  "Bearer invalid-format",
		NonexistentUser: nonexistentUser,
	}
}

// DBTestHelper provides database-related test utilities
type DBTestHelper struct {
	env *TestEnvironment
}

// NewDBTestHelper creates a new database test helper
func NewDBTestHelper(env *TestEnvironment) *DBTestHelper {
	return &DBTestHelper{env: env}
}

// CleanupDB clears all data from test database tables
func (h *DBTestHelper) CleanupDB(t *testing.T) {
	tables := []string{
		"accounts", "profiles", "users", "plans", "groups",
	}
	
	for _, table := range tables {
		err := h.env.DB.Exec(fmt.Sprintf("DELETE FROM %s", table)).Error
		require.NoError(t, err)
	}
}

// SeedTestData seeds the database with common test data
func (h *DBTestHelper) SeedTestData(t *testing.T) {
	// Create test plans
	basicPlan := &models.Plan{
		Name:          "Basic",
		Price:         9.99,
		BillingPeriod: "monthly",
		Description:   "Basic plan for testing",
	}
	
	premiumPlan := &models.Plan{
		Name:          "Premium", 
		Price:         29.99,
		BillingPeriod: "monthly",
		Description:   "Premium plan for testing",
	}

	err := h.env.DB.Create(basicPlan).Error
	require.NoError(t, err)
	
	err = h.env.DB.Create(premiumPlan).Error
	require.NoError(t, err)
}

// GetUserCount returns the number of users in the database
func (h *DBTestHelper) GetUserCount(t *testing.T) int64 {
	var count int64
	err := h.env.DB.Model(&models.User{}).Count(&count).Error
	require.NoError(t, err)
	return count
}

// AssertUserExists asserts that a user with the given email exists
func (h *DBTestHelper) AssertUserExists(t *testing.T, email string) *models.User {
	var user models.User
	err := h.env.DB.Where("email = ?", email).First(&user).Error
	require.NoError(t, err, "User with email %s should exist", email)
	return &user
}

// AssertUserDoesNotExist asserts that a user with the given email does not exist
func (h *DBTestHelper) AssertUserDoesNotExist(t *testing.T, email string) {
	var user models.User
	err := h.env.DB.Where("email = ?", email).First(&user).Error
	require.Error(t, err, "User with email %s should not exist", email)
}