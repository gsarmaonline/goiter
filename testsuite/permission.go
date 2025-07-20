package testsuite

import (
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
