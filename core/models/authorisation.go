package models

import "gorm.io/gorm"

const (
	// Action types
	ReadAction   ActionT = "read"
	CreateAction ActionT = "create"
	UpdateAction ActionT = "update"
	DeleteAction ActionT = "delete"

	// Scope types
	ProjectScopeType ScopeTypeT = "Project"
	AccountScopeType ScopeTypeT = "Account"
)

type (
	ScopeTypeT string
	ElementT   string
	ActionT    string

	RoleAccess struct {
		gorm.Model

		AccessorType ElementT
		AccessorID   uint

		ResourceType ElementT
		ResourceID   uint

		Scope

		Action ActionT
	}

	Scope struct {
		ScopeType string
		ScopeID   uint
	}

	Authorisation struct {
		AllowImplicitOwnerAccess bool
	}
)

func NewAuthorisation() *Authorisation {
	return &Authorisation{
		AllowImplicitOwnerAccess: true,
	}
}

func (a *Authorisation) getQueryString() string {
	return "accessor_type = ? AND accessor_id = ? AND resource_type = ? AND resource_id = ? AND action = ? AND scope_type = ? AND scope_id = ?"
}

func (a *Authorisation) CanAccessResource(db *gorm.DB,
	accessor *User, resource UserOwnedModel, action ActionT, scope Scope) bool {

	var (
		err error
	)

	// If the accessor is the owner of the resource, allow access
	if a.AllowImplicitOwnerAccess && resource.GetUserID() == accessor.GetID() {
		return true
	}

	roleAccess := &RoleAccess{}
	db.Where(a.getQueryString(),
		"User", accessor.GetID(), resource.GetConfig().Name,
		resource.GetID(), action, scope.ScopeType, scope.ScopeID,
	).First(roleAccess)
	if roleAccess.ID != 0 {
		return true
	}

	possibleAccessorGroups := []*Group{}
	possibleResourceGroups := []*Group{}

	if possibleAccessorGroups, err = NewGroupFetcher(db, accessor).GetGroups(); err != nil {
		return false
	}
	if possibleResourceGroups, err = NewGroupFetcher(db, resource).GetGroups(); err != nil {
		return false
	}

	return a.compareWithGroups(db, possibleAccessorGroups, possibleResourceGroups, action, scope)
}

func (a *Authorisation) compareWithGroups(db *gorm.DB, accessorGroups, resourceGroups []*Group, action ActionT, scope Scope) bool {
	for _, ag := range accessorGroups {
		for _, rg := range resourceGroups {
			roleAccess := RoleAccess{}
			// Check role access for this accessor group and resource group
			// If they match, allow access
			db.Where(a.getQueryString(),
				"Group", ag.GetID(), "Group", rg.GetID(),
				action, scope.ScopeType, scope.ScopeID,
			).First(&roleAccess)

			// If role access is found, allow access
			if roleAccess.ID != 0 {
				return true
			}
		}
	}
	return false
}
