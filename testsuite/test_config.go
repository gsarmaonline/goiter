package testsuite

import (
	"log"
	"os"
)

// TestConfig holds configuration for running tests
type TestConfig struct {
	BaseURL        string
	EnableAuth     bool
	EnableVerbose  bool
	TestTimeout    int
}

// GetTestConfig returns the test configuration
func GetTestConfig() *TestConfig {
	config := &TestConfig{
		BaseURL:        getEnvOrDefault("GOITER_BASE_URL", "http://localhost:8090"),
		EnableAuth:     getEnvOrDefault("GOITER_TEST_AUTH", "true") == "true",
		EnableVerbose:  getEnvOrDefault("GOITER_TEST_VERBOSE", "false") == "true",
		TestTimeout:    30, // seconds
	}

	if config.EnableVerbose {
		log.Println("📋 Test Configuration:")
		log.Printf("  Base URL: %s", config.BaseURL)
		log.Printf("  Auth Tests: %t", config.EnableAuth)
		log.Printf("  Verbose: %t", config.EnableVerbose)
		log.Printf("  Timeout: %d seconds", config.TestTimeout)
	}

	return config
}

// getEnvOrDefault returns environment variable value or default
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// RunTestSuite runs the complete test suite with configuration
func RunTestSuite() {
	config := GetTestConfig()
	
	// Initialize clients
	client := NewGoiterClient(config.BaseURL)
	authClient := NewAuthTestClient(config.BaseURL)

	// Setup cleanup to run on exit
	defer func() {
		if config.EnableVerbose {
			log.Println("🧹 Running cleanup...")
		}
		CleanupOnExit(config.BaseURL)
	}()

	log.Println("🚀 Starting Goiter Test Suite...")

	// Run basic functional tests
	log.Println("📋 Running Basic Functional Tests...")
	if err := client.RunUserSuite(); err != nil {
		log.Fatalf("❌ User suite failed: %v", err)
	}
	if err := client.RunProfileSuite(); err != nil {
		log.Fatalf("❌ Profile suite failed: %v", err)
	}
	if err := client.RunAccountSuite(); err != nil {
		log.Fatalf("❌ Account suite failed: %v", err)
	}
	if err := client.RunProjectSuite(); err != nil {
		log.Fatalf("❌ Project suite failed: %v", err)
	}

	// Run authorization tests if enabled
	if config.EnableAuth {
		log.Println("🔐 Running Authorization Tests...")
		
		if err := authClient.RunProjectPermissionTests(); err != nil {
			log.Fatalf("❌ Project permission tests failed: %v", err)
		}

		if err := authClient.RunUnauthorizedAccessTests(); err != nil {
			log.Fatalf("❌ Unauthorized access tests failed: %v", err)
		}

		if err := authClient.TestResourcePermissions(); err != nil {
			log.Fatalf("❌ Resource permission tests failed: %v", err)
		}

		if err := authClient.TestCrossProjectPermissions(); err != nil {
			log.Fatalf("❌ Cross-project permission tests failed: %v", err)
		}

		if err := authClient.TestPermissionInheritance(); err != nil {
			log.Fatalf("❌ Permission inheritance tests failed: %v", err)
		}
	} else {
		log.Println("⚠️  Authorization tests disabled (set GOITER_TEST_AUTH=true to enable)")
	}

	log.Println("🎉 All tests passed! Your authorization system is working correctly.")
}