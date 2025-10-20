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

	Authorisation struct {
		AllowImplicitOwnerAccess bool
		IsEnabled                bool
	}

	AuthorisationRequest struct {
		Db       *gorm.DB
		User     *User
		Resource UserOwnedModel
		Action   ActionT
		Scope    *Scope
	}

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
)

func NewAuthorisation() *Authorisation {
	return &Authorisation{
		AllowImplicitOwnerAccess: true,
		IsEnabled:                false,
	}
}

func (a *Authorisation) getQueryString() string {
	return "accessor_type = ? AND accessor_id = ? AND resource_type = ? AND resource_id = ? AND action = ?"
}

func (c *Authorisation) GetResourcesForUser(authReq *AuthorisationRequest) ([]uint, error) {
	var (
		resources   []uint
		roleAccess  []RoleAccess
		err         error
		accessorIDs []uint
	)

	if c.IsEnabled == false {
		return resources, nil
	}

	accessorIDs = append(accessorIDs, authReq.User.GetID())

	userGroups, err := NewGroupFetcher(authReq.Db, authReq.User).GetGroups()
	if err != nil {
		return resources, err
	}
	for _, g := range userGroups {
		accessorIDs = append(accessorIDs, g.GetID())
	}

	if err = authReq.Db.Where(c.getQueryString(),
		"User",
		accessorIDs,
		authReq.Resource.GetConfig().Name,
		authReq.Action,
	).Find(&roleAccess).Error; err != nil {
		return resources, err
	}

	for _, ra := range roleAccess {
		resources = append(resources, ra.ResourceID)
	}

	return resources, nil
}

func (a *Authorisation) CanAccessResource(authReq *AuthorisationRequest) bool {

	var (
		err error
	)

	if a.IsEnabled == false {
		return true
	}

	// If the accessor is the owner of the resource, allow access
	if a.AllowImplicitOwnerAccess && authReq.Resource.GetUserID() == authReq.User.GetID() {
		return true
	}

	roleAccess := &RoleAccess{}
	authReq.Db.Where(a.getQueryString(),
		"User",
		authReq.User.GetID(),
		authReq.Resource.GetConfig().Name,
		authReq.Resource.GetID(),
		authReq.Action,
	).First(roleAccess)
	if roleAccess.ID != 0 {
		return true
	}

	possibleAccessorGroups := []*Group{}
	possibleResourceGroups := []*Group{}

	if possibleAccessorGroups, err = NewGroupFetcher(authReq.Db, authReq.User).GetGroups(); err != nil {
		return false
	}
	if possibleResourceGroups, err = NewGroupFetcher(authReq.Db, authReq.Resource).GetGroups(); err != nil {
		return false
	}

	return a.compareWithGroups(authReq.Db, possibleAccessorGroups, possibleResourceGroups, authReq.Action, authReq.Scope)
}

func (a *Authorisation) compareWithGroups(db *gorm.DB, accessorGroups, resourceGroups []*Group, action ActionT, scope *Scope) bool {
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
