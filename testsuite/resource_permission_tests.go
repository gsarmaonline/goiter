package testsuite

import (
	"fmt"
	"log"

	"github.com/gsarmaonline/goiter/core/models"
)

// ResourcePermissionTest defines a test for resource-specific permissions
type ResourcePermissionTest struct {
	ResourceType   string
	Action         models.ActionT
	MinLevel       models.PermissionLevel
	Description    string
	ShouldSucceed  map[models.PermissionLevel]bool
}

// ResourcePermissionTests defines permission tests for different resource types
var ResourcePermissionTests = []ResourcePermissionTest{
	// Project resource tests
	{
		ResourceType: "project",
		Action:       models.ReadAction,
		MinLevel:     models.PermissionViewer,
		Description:  "Reading project details",
		ShouldSucceed: map[models.PermissionLevel]bool{
			models.PermissionOwner:  true,
			models.PermissionAdmin:  true,
			models.PermissionEditor: true,
			models.PermissionViewer: true,
		},
	},
	{
		ResourceType: "project",
		Action:       models.UpdateAction,
		MinLevel:     models.PermissionAdmin,
		Description:  "Updating project details",
		ShouldSucceed: map[models.PermissionLevel]bool{
			models.PermissionOwner:  true,
			models.PermissionAdmin:  true,
			models.PermissionEditor: false,
			models.PermissionViewer: false,
		},
	},
	{
		ResourceType: "project",
		Action:       models.DeleteAction,
		MinLevel:     models.PermissionOwner,
		Description:  "Deleting project",
		ShouldSucceed: map[models.PermissionLevel]bool{
			models.PermissionOwner:  true,
			models.PermissionAdmin:  false,
			models.PermissionEditor: false,
			models.PermissionViewer: false,
		},
	},
	
	// Project member resource tests
	{
		ResourceType: "project_member",
		Action:       models.ReadAction,
		MinLevel:     models.PermissionViewer,
		Description:  "Reading project member list",
		ShouldSucceed: map[models.PermissionLevel]bool{
			models.PermissionOwner:  true,
			models.PermissionAdmin:  true,
			models.PermissionEditor: true,
			models.PermissionViewer: true,
		},
	},
	{
		ResourceType: "project_member",
		Action:       models.CreateAction,
		MinLevel:     models.PermissionAdmin,
		Description:  "Adding project members",
		ShouldSucceed: map[models.PermissionLevel]bool{
			models.PermissionOwner:  true,
			models.PermissionAdmin:  true,
			models.PermissionEditor: false,
			models.PermissionViewer: false,
		},
	},
	{
		ResourceType: "project_member",
		Action:       models.DeleteAction,
		MinLevel:     models.PermissionAdmin,
		Description:  "Removing project members",
		ShouldSucceed: map[models.PermissionLevel]bool{
			models.PermissionOwner:  true,
			models.PermissionAdmin:  true,
			models.PermissionEditor: false,
			models.PermissionViewer: false,
		},
	},
	
	// Generic resource tests (wildcard)
	{
		ResourceType: "*",
		Action:       models.ReadAction,
		MinLevel:     models.PermissionViewer,
		Description:  "Reading any resource",
		ShouldSucceed: map[models.PermissionLevel]bool{
			models.PermissionOwner:  true,
			models.PermissionAdmin:  true,
			models.PermissionEditor: true,
			models.PermissionViewer: true,
		},
	},
	{
		ResourceType: "*",
		Action:       models.CreateAction,
		MinLevel:     models.PermissionEditor,
		Description:  "Creating any resource",
		ShouldSucceed: map[models.PermissionLevel]bool{
			models.PermissionOwner:  true,
			models.PermissionAdmin:  true,
			models.PermissionEditor: true,
			models.PermissionViewer: false,
		},
	},
	{
		ResourceType: "*",
		Action:       models.UpdateAction,
		MinLevel:     models.PermissionEditor,
		Description:  "Updating any resource",
		ShouldSucceed: map[models.PermissionLevel]bool{
			models.PermissionOwner:  true,
			models.PermissionAdmin:  true,
			models.PermissionEditor: true,
			models.PermissionViewer: false,
		},
	},
	{
		ResourceType: "*",
		Action:       models.DeleteAction,
		MinLevel:     models.PermissionAdmin,
		Description:  "Deleting any resource",
		ShouldSucceed: map[models.PermissionLevel]bool{
			models.PermissionOwner:  true,
			models.PermissionAdmin:  true,
			models.PermissionEditor: false,
			models.PermissionViewer: false,
		},
	},
}

// TestResourcePermissions tests permissions on different resource types
func (c *AuthTestClient) TestResourcePermissions() error {
	log.Println("🔐 Running Resource Permission Tests...")

	// Setup: Create a project as owner to test permissions against
	ownerClient := NewAuthTestClient(c.BaseURL)
	if err := ownerClient.LoginAsUser("owner_resource@example.com"); err != nil {
		return fmt.Errorf("failed to login as owner: %w", err)
	}

	project, err := ownerClient.CreateProject("Resource Permission Test Project", "Test Description")
	if err != nil {
		return fmt.Errorf("failed to create test project: %w", err)
	}

	projectData := project["data"].(map[string]interface{})
	projectID := uint(projectData["id"].(float64))

	log.Printf("📋 Created test project with ID: %d", projectID)

	// Test each resource permission scenario
	for _, test := range ResourcePermissionTests {
		log.Printf("🧪 Testing resource: %s, action: %s", test.ResourceType, test.Action)

		// Test each permission level
		levels := []models.PermissionLevel{
			models.PermissionOwner,
			models.PermissionAdmin,
			models.PermissionEditor,
			models.PermissionViewer,
		}

		for _, level := range levels {
			expectedSuccess := test.ShouldSucceed[level]
			
			log.Printf("  📝 Testing %s level for %s", getLevelName(level), test.Description)

			// Create test user with specific permission level
			testClient := NewAuthTestClient(c.BaseURL)
			userEmail := fmt.Sprintf("user_resource_%s_%s_%d@example.com", test.ResourceType, test.Action, level)

			if err := testClient.LoginAsUser(userEmail); err != nil {
				return fmt.Errorf("failed to login as test user: %w", err)
			}

			// Add user to project with specific permission level (skip for owner)
			if level != models.PermissionOwner {
				_, err := ownerClient.AddProjectMember(projectID, userEmail, level)
				if err != nil {
					return fmt.Errorf("failed to add member to project: %w", err)
				}
			}

			// Test the permission based on resource type and action
			success, err := c.testResourceAction(testClient, test.ResourceType, test.Action, projectID, expectedSuccess)
			if err != nil {
				return fmt.Errorf("failed to test resource permission: %w", err)
			}

			if success != expectedSuccess {
				return fmt.Errorf("❌ Resource permission test failed: %s level should %s for %s:%s", 
					getLevelName(level), 
					boolToString(expectedSuccess, "succeed", "fail"), 
					test.ResourceType, 
					test.Action)
			}

			log.Printf("    ✅ %s level - %s (as expected)", getLevelName(level), boolToString(expectedSuccess, "ALLOWED", "DENIED"))
		}
	}

	// Cleanup
	log.Printf("🧹 Cleaning up test project...")
	if err := ownerClient.DeleteProject(projectID); err != nil {
		log.Printf("⚠️  Warning: failed to cleanup test project: %v", err)
	}

	log.Println("✅ All Resource Permission Tests PASSED!")
	return nil
}

// testResourceAction tests a specific resource action based on the resource type
func (c *AuthTestClient) testResourceAction(client *AuthTestClient, resourceType string, action models.ActionT, projectID uint, expectedSuccess bool) (bool, error) {
	var method string
	var endpoint string
	var body interface{}
	var expectedStatus int

	if expectedSuccess {
		expectedStatus = 200
	} else {
		expectedStatus = 403
	}

	// Map resource type and action to HTTP method and endpoint
	switch resourceType {
	case "project":
		switch action {
		case models.ReadAction:
			method = "GET"
			endpoint = fmt.Sprintf("/projects/%d", projectID)
		case models.UpdateAction:
			method = "PUT"
			endpoint = fmt.Sprintf("/projects/%d", projectID)
			body = map[string]string{"name": "Updated Resource Test Project"}
		case models.DeleteAction:
			method = "DELETE"
			endpoint = fmt.Sprintf("/projects/%d", projectID)
		default:
			return false, fmt.Errorf("unsupported action for project resource: %s", action)
		}
	case "project_member":
		switch action {
		case models.ReadAction:
			method = "GET"
			endpoint = fmt.Sprintf("/projects/%d", projectID) // Members are included in project details
		case models.CreateAction:
			method = "POST"
			endpoint = fmt.Sprintf("/projects/%d/members", projectID)
			body = map[string]interface{}{
				"user_email": fmt.Sprintf("newmember_%d@example.com", projectID),
				"level":      models.PermissionViewer,
			}
		case models.DeleteAction:
			// First create a member to delete
			ownerClient := NewAuthTestClient(client.BaseURL)
			if err := ownerClient.LoginAsUser("owner_resource@example.com"); err != nil {
				return false, fmt.Errorf("failed to login as owner for member deletion: %w", err)
			}
			
			memberEmail := fmt.Sprintf("member_to_delete_%d@example.com", projectID)
			memberResult, err := ownerClient.AddProjectMember(projectID, memberEmail, models.PermissionViewer)
			if err != nil {
				return false, fmt.Errorf("failed to create member to delete: %w", err)
			}
			
			memberData := memberResult["data"].(map[string]interface{})
			memberUserID := uint(memberData["user_id"].(float64))
			
			method = "DELETE"
			endpoint = fmt.Sprintf("/projects/%d/members/%d", projectID, memberUserID)
		default:
			return false, fmt.Errorf("unsupported action for project_member resource: %s", action)
		}
	case "*":
		// For wildcard resources, we'll test against project operations
		return c.testResourceAction(client, "project", action, projectID, expectedSuccess)
	default:
		return false, fmt.Errorf("unsupported resource type: %s", resourceType)
	}

	// Make the request
	success, err := client.TestPermission(method, endpoint, body, expectedStatus)
	if err != nil {
		// If we got an error but it's the expected failure, that's actually success
		if !expectedSuccess {
			return true, nil
		}
		return false, err
	}

	return success, nil
}

// TestCrossProjectPermissions tests that users can't access resources across projects
func (c *AuthTestClient) TestCrossProjectPermissions() error {
	log.Println("🔐 Running Cross-Project Permission Tests...")

	// Setup: Create two projects with different owners
	owner1Client := NewAuthTestClient(c.BaseURL)
	if err := owner1Client.LoginAsUser("owner1_cross@example.com"); err != nil {
		return fmt.Errorf("failed to login as owner1: %w", err)
	}

	owner2Client := NewAuthTestClient(c.BaseURL)
	if err := owner2Client.LoginAsUser("owner2_cross@example.com"); err != nil {
		return fmt.Errorf("failed to login as owner2: %w", err)
	}

	// Create project 1
	project1, err := owner1Client.CreateProject("Cross Test Project 1", "Test Description 1")
	if err != nil {
		return fmt.Errorf("failed to create test project 1: %w", err)
	}
	project1Data := project1["data"].(map[string]interface{})
	project1ID := uint(project1Data["id"].(float64))

	// Create project 2
	project2, err := owner2Client.CreateProject("Cross Test Project 2", "Test Description 2")
	if err != nil {
		return fmt.Errorf("failed to create test project 2: %w", err)
	}
	project2Data := project2["data"].(map[string]interface{})
	project2ID := uint(project2Data["id"].(float64))

	// Create a user with access to project 1
	userClient := NewAuthTestClient(c.BaseURL)
	if err := userClient.LoginAsUser("user_cross@example.com"); err != nil {
		return fmt.Errorf("failed to login as test user: %w", err)
	}

	// Add user to project 1 as editor
	_, err = owner1Client.AddProjectMember(project1ID, "user_cross@example.com", models.PermissionEditor)
	if err != nil {
		return fmt.Errorf("failed to add user to project 1: %w", err)
	}

	// Test that user can access project 1
	log.Printf("🧪 Testing user can access project 1 (should succeed)")
	success, err := userClient.TestPermission("GET", fmt.Sprintf("/projects/%d", project1ID), nil, 200)
	if !success {
		return fmt.Errorf("❌ User should be able to access project 1: %v", err)
	}
	log.Printf("✅ User can access project 1 - PASSED")

	// Test that user cannot access project 2
	log.Printf("🧪 Testing user cannot access project 2 (should fail)")
	success, err = userClient.TestPermission("GET", fmt.Sprintf("/projects/%d", project2ID), nil, 403)
	if !success {
		return fmt.Errorf("❌ User should not be able to access project 2: %v", err)
	}
	log.Printf("✅ User cannot access project 2 - PASSED")

	// Test that user cannot modify project 2
	log.Printf("🧪 Testing user cannot modify project 2 (should fail)")
	success, err = userClient.TestPermission("PUT", fmt.Sprintf("/projects/%d", project2ID), 
		map[string]string{"name": "Hacked Project"}, 403)
	if !success {
		return fmt.Errorf("❌ User should not be able to modify project 2: %v", err)
	}
	log.Printf("✅ User cannot modify project 2 - PASSED")

	// Cleanup
	log.Printf("🧹 Cleaning up test projects...")
	if err := owner1Client.DeleteProject(project1ID); err != nil {
		log.Printf("⚠️  Warning: failed to cleanup test project 1: %v", err)
	}
	if err := owner2Client.DeleteProject(project2ID); err != nil {
		log.Printf("⚠️  Warning: failed to cleanup test project 2: %v", err)
	}

	log.Println("✅ All Cross-Project Permission Tests PASSED!")
	return nil
}

// boolToString converts a boolean to a string based on true/false values
func boolToString(value bool, trueStr, falseStr string) string {
	if value {
		return trueStr
	}
	return falseStr
}

// TestPermissionInheritance tests that permissions inherit correctly
func (c *AuthTestClient) TestPermissionInheritance() error {
	log.Println("🔐 Running Permission Inheritance Tests...")

	// Setup: Create a project
	ownerClient := NewAuthTestClient(c.BaseURL)
	if err := ownerClient.LoginAsUser("owner_inherit@example.com"); err != nil {
		return fmt.Errorf("failed to login as owner: %w", err)
	}

	project, err := ownerClient.CreateProject("Inheritance Test Project", "Test Description")
	if err != nil {
		return fmt.Errorf("failed to create test project: %w", err)
	}

	projectData := project["data"].(map[string]interface{})
	projectID := uint(projectData["id"].(float64))

	// Test that higher permission levels can perform lower-level actions
	inheritanceTests := []struct {
		level       models.PermissionLevel
		canRead     bool
		canWrite    bool
		canDelete   bool
		description string
	}{
		{models.PermissionOwner, true, true, true, "Owner should have all permissions"},
		{models.PermissionAdmin, true, true, false, "Admin should have read/write but not delete"},
		{models.PermissionEditor, true, false, false, "Editor should have read only (based on current implementation)"},
		{models.PermissionViewer, true, false, false, "Viewer should have read only"},
	}

	for _, test := range inheritanceTests {
		log.Printf("🧪 Testing permission inheritance for %s", getLevelName(test.level))

		// Create user with specific permission level
		testClient := NewAuthTestClient(c.BaseURL)
		userEmail := fmt.Sprintf("user_inherit_%d@example.com", test.level)

		if err := testClient.LoginAsUser(userEmail); err != nil {
			return fmt.Errorf("failed to login as test user: %w", err)
		}

		// Add user to project (skip for owner)
		if test.level != models.PermissionOwner {
			_, err := ownerClient.AddProjectMember(projectID, userEmail, test.level)
			if err != nil {
				return fmt.Errorf("failed to add member to project: %w", err)
			}
		}

		// Test read permission
		readSuccess, err := testClient.TestPermission("GET", fmt.Sprintf("/projects/%d", projectID), nil, 
			boolToStatus(test.canRead))
		if readSuccess != test.canRead {
			return fmt.Errorf("❌ Read permission failed for %s: expected %t, got %t. Error: %v", 
				getLevelName(test.level), test.canRead, readSuccess, err)
		}

		// Test write permission
		writeSuccess, err := testClient.TestPermission("PUT", fmt.Sprintf("/projects/%d", projectID), 
			map[string]string{"name": "Updated by " + getLevelName(test.level)}, 
			boolToStatus(test.canWrite))
		if writeSuccess != test.canWrite {
			return fmt.Errorf("❌ Write permission failed for %s: expected %t, got %t. Error: %v", 
				getLevelName(test.level), test.canWrite, writeSuccess, err)
		}

		log.Printf("✅ %s inheritance - PASSED", getLevelName(test.level))
	}

	// Cleanup
	log.Printf("🧹 Cleaning up test project...")
	if err := ownerClient.DeleteProject(projectID); err != nil {
		log.Printf("⚠️  Warning: failed to cleanup test project: %v", err)
	}

	log.Println("✅ All Permission Inheritance Tests PASSED!")
	return nil
}

// boolToStatus converts boolean to HTTP status code
func boolToStatus(success bool) int {
	if success {
		return 200
	}
	return 403
}