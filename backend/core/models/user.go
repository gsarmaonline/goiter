package models

import (
	"gorm.io/gorm"
)

const (
	ActiveUser   UserStatus = "active"
	InactiveUser UserStatus = "inactive"
)

type UserStatus string

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`

	Name       string `json:"name"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`

	Picture string `json:"picture"`
}

type User struct {
	BaseModelWithoutUser

	GoogleID string `json:"-" gorm:"uniqueIndex"`
	Email    string `json:"email" gorm:"uniqueIndex"`
	Name     string `json:"name"`

	Picture     string `json:"picture"`
	AccessToken string `json:"-"`

	UserStatus UserStatus `json:"user_status" gorm:"type:varchar(20);not null;default:'active'"`

	CreatedFrom string `json:"-" gorm:"type:varchar(20);not null;default:'login'"`

	Profile Profile `json:"-" gorm:"foreignKey:UserID"`
}

// AfterCreate hook to create a profile and account when a new user is created
func (u *User) AfterCreate(tx *gorm.DB) error {
	// Create profile
	profile := Profile{
		BaseModelWithUser: BaseModelWithUser{
			UserID: u.ID,
		},
	}
	if err := tx.Create(&profile).Error; err != nil {
		return err
	}

	// Create account
	account := Account{
		BaseModelWithUser: BaseModelWithUser{
			UserID: u.ID,
		},
		Name: u.Name + "'s Account",
	}
	if err := tx.Create(&account).Error; err != nil {
		return err
	}

	return nil
}
