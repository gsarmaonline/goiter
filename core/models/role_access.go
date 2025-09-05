package models

import (
	"log"

	"gorm.io/gorm"
)

// Pending Items
// - Implement scope of the RoleAccess. It can be for a project or account

const (
	// Actions
	ReadAction   ActionT = "read"
	WriteAction  ActionT = "write"
	DeleteAction ActionT = "delete"
	UpdateAction ActionT = "update"
	CreateAction ActionT = "create"

	// Wildcards
	WildcardResourceType         = "*"
	WildcardResourceID           = 0
	WildcardAction       ActionT = "*"

	// Scope types
	AccountScopeType ScopeTypeT = "account_scope"
	ProjectScopeType ScopeTypeT = "account_scope"
)

type (
	ActionT    string
	ScopeTypeT string
	// RoleAccess represents a user's permission level in a project
	// A User can access the resource if they have the permission level for the resource.
	// The permission level of the user is defined by the ProjectPermission table.
	// By default, if there is no corresponding role access entry, then the user has no access to the resource.
	RoleAccess struct {
		BaseModelWithoutUser

		// If the ResourceType is empty, then it becomes applicable to all resources unless
		// there is another entry for a specific resource type.
		ResourceType string `json:"resource_type" gorm:"not null"`
		// If the ResourceID is 0 and the ResourceType is not empty, then it becomes applicable to all resources of the given type unless
		// there is another entry for a specific resource ID.
		ResourceID uint `json:"resource_id"`

		Level PermissionLevel `json:"level" gorm:"not null;default:1"`

		// ProjectID can be used to scope the access to a specific project.
		// If the ProjectID is 0, then it becomes applicable to all projects unless
		ProjectID uint `json:"project_id"`

		// Scope is used to identify the actual scope where the rules for access
		// will be applied to. For example, a scope can be for the entire account,
		// or a project, or any other group
		ScopeType ScopeTypeT `json:"scope_type"`
		ScopeID   uint       `json:"scope_id"`

		Action ActionT `json:"action" gorm:"not null"`
	}
)

func (r RoleAccess) GetConfig() ModelConfig {
	return ModelConfig{
		Name:      "RoleAccess",
		ScopeType: AccountScopeType,
	}
}

func CanAccessResource(db *gorm.DB,
	resourceType string,
	resourceID uint,
	user *User,
	action ActionT,
) bool {
	// Get User's level in the project
	var (
		projectPermission Permission
	)
	db.Where("user_id = ?", user.ID).First(&projectPermission)

	for _, permissionScope := range formPermissionScopes(resourceType, resourceID, action) {
		var (
			allowAccess   bool
			proceedToNext bool
		)
		if allowAccess, proceedToNext = formRoleAccessQuery(
			db,
			permissionScope.ResourceType,
			permissionScope.ResourceID,
			permissionScope.Action,
			projectPermission.Level,
		); allowAccess {
			// If we found a matching role access entry, we can return true
			log.Println("Access granted for resource:",
				permissionScope.ResourceType,
				permissionScope.ResourceID,
				"with action:",
				permissionScope.Action,
			)
			return true
		}
		if !proceedToNext {
			// If we found a matching role access entry but it doesn't allow access, we can
			// stop checking further scopes
			return false
		}
	}

	return false
}

func formPermissionScopes(resourceType string, resourceID uint, action ActionT) []RoleAccess {
	return []RoleAccess{
		{
			ResourceType: resourceType,
			ResourceID:   resourceID,
			Action:       action,
		},
		{
			ResourceType: resourceType,
			ResourceID:   WildcardResourceID,
			Action:       action,
		},
		{
			ResourceType: WildcardResourceType,
			ResourceID:   WildcardResourceID,
			Action:       action,
		},
		{
			ResourceType: WildcardResourceType,
			ResourceID:   WildcardResourceID,
			Action:       WildcardAction,
		},
	}
}

func formRoleAccessQuery(db *gorm.DB,
	resourceType string,
	resourceID uint,
	action ActionT,
	level PermissionLevel,
) (allowAccess bool, proceedToNext bool) {

	proceedToNext = true
	roleAccess := &RoleAccess{}

	db.Where("resource_type = ? AND resource_id = ? AND action = ?", resourceType, resourceID, action).First(&roleAccess)

	if roleAccess.ID != 0 {
		allowAccess = roleAccess.Level >= level
		proceedToNext = false
		return
	}
	return
}
