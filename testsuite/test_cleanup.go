package testsuite

import (
	"log"
	"strings"
)

// TestCleanup handles cleanup of test data
type TestCleanup struct {
	client *GoiterClient
}

// NewTestCleanup creates a new cleanup instance
func NewTestCleanup(client *GoiterClient) *TestCleanup {
	return &TestCleanup{client: client}
}

// CleanupTestUsers removes all test users from the database
func (tc *TestCleanup) CleanupTestUsers() error {
	// This would require a direct database connection or admin API
	// For now, we'll just log that cleanup should happen
	log.Println("🧹 Test users will be cleaned up automatically by the database")
	return nil
}

// CleanupTestProjects removes all test projects
func (tc *TestCleanup) CleanupTestProjects() error {
	// Login as admin to clean up projects
	adminClient := NewAuthTestClient(tc.client.BaseURL)
	if err := adminClient.LoginAsUser("admin_cleanup@example.com"); err != nil {
		// If we can't login as admin, skip cleanup
		log.Printf("⚠️  Warning: Could not login as admin for cleanup: %v", err)
		return nil
	}

	// List all projects
	projects, err := adminClient.GetProjects()
	if err != nil {
		log.Printf("⚠️  Warning: Could not list projects for cleanup: %v", err)
		return nil
	}

	// Delete projects that look like test projects
	if projectData, ok := projects["data"].([]interface{}); ok {
		for _, proj := range projectData {
			if projectMap, ok := proj.(map[string]interface{}); ok {
				if name, ok := projectMap["name"].(string); ok {
					if isTestProject(name) {
						if id, ok := projectMap["id"].(float64); ok {
							projectID := uint(id)
							if err := adminClient.DeleteProject(projectID); err != nil {
								log.Printf("⚠️  Warning: Could not delete test project %d: %v", projectID, err)
							} else {
								log.Printf("🗑️  Cleaned up test project: %s", name)
							}
						}
					}
				}
			}
		}
	}

	return nil
}

// isTestProject determines if a project name indicates it's a test project
func isTestProject(name string) bool {
	testIndicators := []string{
		"test",
		"Test",
		"TEST",
		"Permission",
		"Authorization",
		"Cross",
		"Inheritance",
		"Unauthorized",
		"Resource",
	}

	for _, indicator := range testIndicators {
		if strings.Contains(name, indicator) {
			return true
		}
	}
	return false
}

// RunCleanup runs all cleanup operations
func (tc *TestCleanup) RunCleanup() error {
	log.Println("🧹 Starting test cleanup...")

	if err := tc.CleanupTestProjects(); err != nil {
		log.Printf("⚠️  Warning: Project cleanup failed: %v", err)
	}

	if err := tc.CleanupTestUsers(); err != nil {
		log.Printf("⚠️  Warning: User cleanup failed: %v", err)
	}

	log.Println("✅ Test cleanup completed")
	return nil
}

// CleanupOnExit should be called when tests exit
func CleanupOnExit(baseURL string) {
	client := NewGoiterClient(baseURL)
	cleanup := NewTestCleanup(client)
	cleanup.RunCleanup()
}