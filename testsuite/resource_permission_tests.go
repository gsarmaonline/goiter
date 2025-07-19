package testsuite

import (
	"github.com/gsarmaonline/goiter/core/models"
)

// ResourcePermissionTest defines a test for resource-specific permissions
type ResourcePermissionTest struct {
	ResourceType  string
	Action        models.ActionT
	MinLevel      models.PermissionLevel
	Description   string
	ShouldSucceed map[models.PermissionLevel]bool
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
