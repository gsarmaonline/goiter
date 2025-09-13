package models

import (
	"time"
)

type (
	UserOwnedModel interface {
		GetID() uint
		GetUserID() uint
		SetUserID(uint)
		GetConfig() ModelConfig
	}
	BaseModel struct {
		ID        uint      `json:"id" gorm:"primary_key"`
		CreatedAt time.Time `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
		UpdatedAt time.Time `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP"`

		DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index;default:NULL"`
	}
	BaseModelWithUser struct {
		BaseModel

		UserID uint  `json:"user_id" gorm:"not null"`
		User   *User `json:"user" gorm:"foreignKey:UserID"`
	}
	BaseModelWithoutUser struct {
		BaseModel
	}
	ModelConfig struct {
		Name      string     `json:"name"`
		ScopeType ScopeTypeT `json:"scope_type"`
	}
)

func (b *BaseModelWithUser) GetUserID() uint {
	return b.UserID
}

func (b *BaseModelWithUser) SetUserID(userID uint) {
	b.UserID = userID
}

func (b *BaseModelWithoutUser) GetUserID() uint {
	return 0
}

func (b *BaseModelWithoutUser) SetUserID(userID uint) {
	return
}

func (b *BaseModel) GetID() uint {
	return b.ID
}

func (b *BaseModel) GetConfig() ModelConfig {
	return ModelConfig{
		ScopeType: ProjectScopeType,
	}
}
