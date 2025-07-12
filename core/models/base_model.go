package models

import (
	"time"

	"gorm.io/gorm"
)

type (
	UserOwnedModel interface {
		GetID() uint
		GetUserID() uint
		SetUserID(uint)
	}
	BaseModel struct {
		gorm.Model

		ID        uint      `json:"id" gorm:"primary_key"`
		CreatedAt time.Time `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
		UpdatedAt time.Time `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	}
	BaseModelWithUser struct {
		BaseModel

		UserID uint  `json:"user_id" gorm:"not null"`
		User   *User `json:"user" gorm:"foreignKey:UserID"`
	}
	BaseModelWithoutUser struct {
		BaseModel
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
