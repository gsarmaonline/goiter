package models

import (
	"gorm.io/gorm"
)

const (
	PermissionOwner    PermissionLevel = 20 // Full access, can manage members and resources
	PermissionAdmin    PermissionLevel = 19 // Can manage resources and some project settings
	PermissionEditor   PermissionLevel = 18 // Can edit resources
	PermissionViewer   PermissionLevel = 17 // Can only view resources
	PermissionEveryone PermissionLevel = 10 // Can only view resources

)

type (
	// PermissionLevel represents the different levels of access a user can have in a project
	// TODO: PermissionLevel should be a table
	PermissionLevel uint8

	// Project represents a group of resources that can be shared among users
	Project struct {
		BaseModelWithUser

		Name        string `json:"name" gorm:"not null"`
		Description string `json:"description"`

		AccountID uint    `json:"account_id" gorm:"not null"`
		Account   Account `json:"account" gorm:"foreignKey:AccountID"`

		Members []User `json:"members" gorm:"many2many:project_permissions;"`
	}

	// ProjectPermission represents a user's permission level in a project
	ProjectPermission struct {
		BaseModelWithUser

		ProjectID uint    `json:"project_id" gorm:"not null"`
		Project   Project `json:"project" gorm:"foreignKey:ProjectID"`

		UserEmail string `json:"user_email" gorm:"not null"`

		Level PermissionLevel `json:"level" gorm:"not null;default:10"`
	}
)

func (p *ProjectPermission) TableName() string {
	return "project_permissions"
}

func (p *Project) BeforeCreate(tx *gorm.DB) (err error) {
	// Set default role access for the project
	for _, roleAccess := range p.GetDefaultRoleAccessRows() {
		roleAccess.ProjectID = p.ID
		if err = tx.Create(roleAccess).Error; err != nil {
			return err
		}
	}

	return nil
}

// BeforeDelete is a GORM hook that handles cleanup before project deletion
func (p *Project) BeforeDelete(tx *gorm.DB) error {
	// Remove all project permissions
	if err := tx.Exec("DELETE FROM project_permissions WHERE project_id = ?", p.ID).Error; err != nil {
		return err
	}
	return nil
}

func (p *ProjectPermission) BeforeCreate(tx *gorm.DB) (err error) {
	if p.UserID != 0 {
		return
	}

	user := &User{}
	if err = tx.Where("email = ?", p.UserEmail).First(user).Error; err == nil {
		p.UserID = user.ID
		return
	}

	if user.ID == 0 {
		user = &User{
			Email:      p.UserEmail,
			UserStatus: InactiveUser,

			CreatedFrom: "project_permission",
		}
		if err = tx.Create(user).Error; err != nil {
			return
		}
	}

	p.UserID = user.ID
	return
}

func (p *Project) GetDefaultRoleAccessRows() (roleAcceses []*RoleAccess) {
	roleAcceses = []*RoleAccess{
		{
			ResourceType: WildcardResourceType,
			ResourceID:   WildcardResourceID,
			Action:       WildcardAction,
			Level:        PermissionViewer,
		},
	}
	return
}
