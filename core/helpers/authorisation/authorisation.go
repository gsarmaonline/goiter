package authorisation

import (
	"github.com/gin-gonic/gin"
	"github.com/gsarmaonline/goiter/core/models"
	"gorm.io/gorm"
)

type (
	Authorisation struct {
		handler                  HandlerInterface
		AllowImplicitOwnerAccess bool
		IsEnabled                bool
	}

	HandlerInterface interface {
		GetUserFromContext(c *gin.Context) *models.User
	}

	AuthorisationRequest struct {
		Db       *gorm.DB
		User     *models.User
		Resource models.UserOwnedModel
		Action   models.ActionT
	}
)

func NewAuthorisation(handler HandlerInterface) *Authorisation {
	return &Authorisation{
		handler:                  handler,
		AllowImplicitOwnerAccess: true,
		IsEnabled:                false,
	}
}

func (a *Authorisation) UserScopedDB(c *gin.Context, db *gorm.DB) *gorm.DB {
	db = db.Where("user_id = ?", a.handler.GetUserFromContext(c).ID)
	return db
}

func (a *Authorisation) UpdateWithUser(c *gin.Context, model models.UserOwnedModel) (err error) {
	model.SetUserID(a.handler.GetUserFromContext(c).ID)
	return
}

func (a *Authorisation) CanAccessResource(c *gin.Context, model models.UserOwnedModel) bool {
	if a.handler.GetUserFromContext(c).ID == model.GetUserID() {
		return true
	}
	return false
}
