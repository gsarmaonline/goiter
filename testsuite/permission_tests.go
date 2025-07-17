package testsuite

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gsarmaonline/goiter/core/models"
)

// PermissionTestScenario defines a test scenario for permissions
type PermissionTestScenario struct {
	Name           string
	UserLevel      models.PermissionLevel
	Action         string
	Method         string
	Endpoint       string
	Body           interface{}
	ExpectedStatus int
	Description    string
}

// ProjectPermissionTests defines test scenarios for project permissions
var ProjectPermissionTests = []PermissionTestScenario{
	// Owner tests
	{
		Name:           "Owner can read project",
		UserLevel:      models.PermissionOwner,
		Action:         "read",
		Method:         "GET",
		Endpoint:       "/projects/%d",
		ExpectedStatus: 200,
		Description:    "Project owner should be able to read project details",
	},
	{
		Name:           "Owner can update project",
		UserLevel:      models.PermissionOwner,
		Action:         "update",
		Method:         "PUT",
		Endpoint:       "/projects/%d",
		Body:           map[string]string{"name": "Updated Project by Owner"},
		ExpectedStatus: 200,
		Description:    "Project owner should be able to update project details",
	},
	{
		Name:           "Owner can delete project",
		UserLevel:      models.PermissionOwner,
		Action:         "delete",
		Method:         "DELETE",
		Endpoint:       "/projects/%d",
		ExpectedStatus: 200,
		Description:    "Project owner should be able to delete the project",
	},

	// Admin tests
	{
		Name:           "Admin can read project",
		UserLevel:      models.PermissionAdmin,
		Action:         "read",
		Method:         "GET",
		Endpoint:       "/projects/%d",
		ExpectedStatus: 200,
		Description:    "Admin should be able to read project details",
	},
	{
		Name:           "Admin can update project",
		UserLevel:      models.PermissionAdmin,
		Action:         "update",
		Method:         "PUT",
		Endpoint:       "/projects/%d",
		Body:           map[string]string{"name": "Updated Project by Admin"},
		ExpectedStatus: 200,
		Description:    "Admin should be able to update project details",
	},
	{
		Name:           "Admin cannot delete project",
		UserLevel:      models.PermissionAdmin,
		Action:         "delete",
		Method:         "DELETE",
		Endpoint:       "/projects/%d",
		ExpectedStatus: 403,
		Description:    "Admin should not be able to delete the project",
	},

	// Editor tests
	{
		Name:           "Editor can read project",
		UserLevel:      models.PermissionEditor,
		Action:         "read",
		Method:         "GET",
		Endpoint:       "/projects/%d",
		ExpectedStatus: 200,
		Description:    "Editor should be able to read project details",
	},
	{
		Name:           "Editor cannot update project",
		UserLevel:      models.PermissionEditor,
		Action:         "update",
		Method:         "PUT",
		Endpoint:       "/projects/%d",
		Body:           map[string]string{"name": "Updated Project by Editor"},
		ExpectedStatus: 403,
		Description:    "Editor should not be able to update project details",
	},
	{
		Name:           "Editor cannot delete project",
		UserLevel:      models.PermissionEditor,
		Action:         "delete",
		Method:         "DELETE",
		Endpoint:       "/projects/%d",
		ExpectedStatus: 403,
		Description:    "Editor should not be able to delete the project",
	},

	// Viewer tests
	{
		Name:           "Viewer can read project",
		UserLevel:      models.PermissionViewer,
		Action:         "read",
		Method:         "GET",
		Endpoint:       "/projects/%d",
		ExpectedStatus: 200,
		Description:    "Viewer should be able to read project details",
	},
	{
		Name:           "Viewer cannot update project",
		UserLevel:      models.PermissionViewer,
		Action:         "update",
		Method:         "PUT",
		Endpoint:       "/projects/%d",
		Body:           map[string]string{"name": "Updated Project by Viewer"},
		ExpectedStatus: 403,
		Description:    "Viewer should not be able to update project details",
	},
	{
		Name:           "Viewer cannot delete project",
		UserLevel:      models.PermissionViewer,
		Action:         "delete",
		Method:         "DELETE",
		Endpoint:       "/projects/%d",
		ExpectedStatus: 403,
		Description:    "Viewer should not be able to delete the project",
	},
}

// ProjectMemberPermissionTests defines test scenarios for project member management
var ProjectMemberPermissionTests = []PermissionTestScenario{
	// Owner tests
	{
		Name:           "Owner can add project member",
		UserLevel:      models.PermissionOwner,
		Action:         "create",
		Method:         "POST",
		Endpoint:       "/projects/%d/members",
		Body:           map[string]interface{}{"user_email": "newmember@example.com", "level": models.PermissionViewer},
		ExpectedStatus: 200,
		Description:    "Project owner should be able to add members",
	},
	{
		Name:           "Owner can remove project member",
		UserLevel:      models.PermissionOwner,
		Action:         "delete",
		Method:         "DELETE",
		Endpoint:       "/projects/%d/members/%d",
		ExpectedStatus: 200,
		Description:    "Project owner should be able to remove members",
	},

	// Admin tests
	{
		Name:           "Admin can add project member",
		UserLevel:      models.PermissionAdmin,
		Action:         "create",
		Method:         "POST",
		Endpoint:       "/projects/%d/members",
		Body:           map[string]interface{}{"user_email": "newmember2@example.com", "level": models.PermissionViewer},
		ExpectedStatus: 200,
		Description:    "Admin should be able to add members",
	},
	{
		Name:           "Admin can remove project member",
		UserLevel:      models.PermissionAdmin,
		Action:         "delete",
		Method:         "DELETE",
		Endpoint:       "/projects/%d/members/%d",
		ExpectedStatus: 200,
		Description:    "Admin should be able to remove members",
	},

	// Editor tests
	{
		Name:           "Editor cannot add project member",
		UserLevel:      models.PermissionEditor,
		Action:         "create",
		Method:         "POST",
		Endpoint:       "/projects/%d/members",
		Body:           map[string]interface{}{"user_email": "newmember3@example.com", "level": models.PermissionViewer},
		ExpectedStatus: 403,
		Description:    "Editor should not be able to add members",
	},
	{
		Name:           "Editor cannot remove project member",
		UserLevel:      models.PermissionEditor,
		Action:         "delete",
		Method:         "DELETE",
		Endpoint:       "/projects/%d/members/%d",
		ExpectedStatus: 403,
		Description:    "Editor should not be able to remove members",
	},

	// Viewer tests
	{
		Name:           "Viewer cannot add project member",
		UserLevel:      models.PermissionViewer,
		Action:         "create",
		Method:         "POST",
		Endpoint:       "/projects/%d/members",
		Body:           map[string]interface{}{"user_email": "newmember4@example.com", "level": models.PermissionViewer},
		ExpectedStatus: 403,
		Description:    "Viewer should not be able to add members",
	},
	{
		Name:           "Viewer cannot remove project member",
		UserLevel:      models.PermissionViewer,
		Action:         "delete",
		Method:         "DELETE",
		Endpoint:       "/projects/%d/members/%d",
		ExpectedStatus: 403,
		Description:    "Viewer should not be able to remove members",
	},
}

// RunProjectPermissionTests runs all project permission tests
func (c *AuthTestClient) RunProjectPermissionTests() error {
	log.Println("🔐 Running Project Permission Tests...")

	// Setup: Create an owner client
	ownerClient := NewAuthTestClient(c.BaseURL)
	if err := ownerClient.LoginAsUser("owner@example.com"); err != nil {
		return fmt.Errorf("failed to login as owner: %w", err)
	}

	// Test each permission level for project operations
	for _, test := range ProjectPermissionTests {
		log.Printf("🧪 Testing: %s", test.Name)

		// Create a fresh project for each test
		project, err := ownerClient.CreateProject("Test Project for Permissions", "Test Description")
		if err != nil {
			return fmt.Errorf("failed to create test project: %w", err)
		}

		projectData := project["data"].(map[string]interface{})
		projectID := uint(projectData["id"].(float64))

		// Create user with specific permission level
		testClient := NewAuthTestClient(c.BaseURL)
		userEmail := fmt.Sprintf("user_project_%d_%s@example.com", test.UserLevel, test.Action)

		if err := testClient.LoginAsUser(userEmail); err != nil {
			return fmt.Errorf("failed to login as test user: %w", err)
		}

		// Add user to project with specific permission level (skip for owner)
		if test.UserLevel != models.PermissionOwner {
			_, err := ownerClient.AddProjectMember(projectID, userEmail, test.UserLevel)
			if err != nil {
				return fmt.Errorf("failed to add member to project: %w", err)
			}
		}

		// Test the permission
		endpoint := fmt.Sprintf(test.Endpoint, projectID)
		success, err := testClient.TestPermission(test.Method, endpoint, test.Body, test.ExpectedStatus)

		if !success {
			return fmt.Errorf("❌ Permission test failed: %s - %s. Error: %v", test.Name, test.Description, err)
		}

		log.Printf("✅ %s - PASSED", test.Name)

		// Cleanup: Delete the project (unless it's already deleted by the test)
		if test.Action != "delete" {
			if err := ownerClient.DeleteProject(projectID); err != nil {
				log.Printf("⚠️  Warning: failed to cleanup test project %d: %v", projectID, err)
			}
		}
	}

	log.Println("🔐 Running Project Member Permission Tests...")

	// Test project member management permissions
	for _, test := range ProjectMemberPermissionTests {
		log.Printf("🧪 Testing: %s", test.Name)

		// Create a fresh project for each test
		project, err := ownerClient.CreateProject("Test Project for Member Permissions", "Test Description")
		if err != nil {
			return fmt.Errorf("failed to create test project: %w", err)
		}

		projectData := project["data"].(map[string]interface{})
		projectID := uint(projectData["id"].(float64))

		// Create user with specific permission level
		testClient := NewAuthTestClient(c.BaseURL)
		userEmail := fmt.Sprintf("user_member_%d_%s@example.com", test.UserLevel, test.Action)

		if err := testClient.LoginAsUser(userEmail); err != nil {
			return fmt.Errorf("failed to login as test user: %w", err)
		}

		// Add user to project with specific permission level (skip for owner)
		if test.UserLevel != models.PermissionOwner {
			_, err := ownerClient.AddProjectMember(projectID, userEmail, test.UserLevel)
			if err != nil {
				return fmt.Errorf("failed to add member to project: %w", err)
			}
		}

		// For member removal tests, we need to create a member first
		var memberUserID uint
		if test.Action == "delete" {
			memberEmail := fmt.Sprintf("member_to_remove_%d@example.com", test.UserLevel)
			memberResult, err := ownerClient.AddProjectMember(projectID, memberEmail, models.PermissionViewer)
			if err != nil {
				return fmt.Errorf("failed to create member to remove: %w", err)
			}
			memberData := memberResult["data"].(map[string]interface{})
			memberUserID = uint(memberData["user_id"].(float64))
		}

		// Test the permission
		var endpoint string
		if test.Action == "delete" {
			endpoint = fmt.Sprintf(test.Endpoint, projectID, memberUserID)
		} else {
			endpoint = fmt.Sprintf(test.Endpoint, projectID)
		}

		success, err := testClient.TestPermission(test.Method, endpoint, test.Body, test.ExpectedStatus)

		if !success {
			return fmt.Errorf("❌ Member permission test failed: %s - %s. Error: %v", test.Name, test.Description, err)
		}

		log.Printf("✅ %s - PASSED", test.Name)

		// Cleanup: Delete the project
		if err := ownerClient.DeleteProject(projectID); err != nil {
			log.Printf("⚠️  Warning: failed to cleanup test project %d: %v", projectID, err)
		}
	}

	log.Println("✅ All Project Permission Tests PASSED!")
	return nil
}

// RunUnauthorizedAccessTests tests access without proper authentication
func (c *AuthTestClient) RunUnauthorizedAccessTests() error {
	log.Println("🔐 Running Unauthorized Access Tests...")

	// Create a project first
	ownerClient := NewAuthTestClient(c.BaseURL)
	if err := ownerClient.LoginAsUser("owner_unauth@example.com"); err != nil {
		return fmt.Errorf("failed to login as owner: %w", err)
	}

	project, err := ownerClient.CreateProject("Unauthorized Test Project", "Test Description")
	if err != nil {
		return fmt.Errorf("failed to create test project: %w", err)
	}

	projectData := project["data"].(map[string]interface{})
	projectID := uint(projectData["id"].(float64))

	// Test unauthorized access
	unauthorizedClient := NewGoiterClient(c.BaseURL)
	// Don't log in - no JWT token

	endpoints := []struct {
		method   string
		endpoint string
		name     string
	}{
		{"GET", fmt.Sprintf("/projects/%d", projectID), "Read project without auth"},
		{"PUT", fmt.Sprintf("/projects/%d", projectID), "Update project without auth"},
		{"DELETE", fmt.Sprintf("/projects/%d", projectID), "Delete project without auth"},
		{"GET", "/projects", "List projects without auth"},
		{"POST", "/projects", "Create project without auth"},
	}

	for _, test := range endpoints {
		log.Printf("🧪 Testing: %s", test.name)

		resp, err := unauthorizedClient.httpClient.Do(func() *http.Request {
			req, _ := http.NewRequest(test.method, unauthorizedClient.BaseURL+test.endpoint, nil)
			return req
		}())

		if err != nil {
			return fmt.Errorf("failed to make request: %w", err)
		}
		resp.Body.Close()

		if resp.StatusCode != 401 {
			return fmt.Errorf("expected 401 Unauthorized but got %d for %s", resp.StatusCode, test.name)
		}

		log.Printf("✅ %s - PASSED (correctly returned 401)", test.name)
	}

	// Cleanup
	if err := ownerClient.DeleteProject(projectID); err != nil {
		log.Printf("⚠️  Warning: failed to cleanup test project: %v", err)
	}

	log.Println("✅ All Unauthorized Access Tests PASSED!")
	return nil
}

// getLevelName returns a human-readable name for a permission level
func getLevelName(level models.PermissionLevel) string {
	switch level {
	case models.PermissionOwner:
		return "Owner"
	case models.PermissionAdmin:
		return "Admin"
	case models.PermissionEditor:
		return "Editor"
	case models.PermissionViewer:
		return "Viewer"
	default:
		return "Unknown"
	}
}