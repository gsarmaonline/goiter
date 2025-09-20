package app

import (
	"github.com/gin-gonic/gin"
	"github.com/gsarmaonline/goiter/core/handlers"
	"github.com/gsarmaonline/goiter/core/models"
)

func (app *App) CreateModelOneHandler(c *gin.Context) {
	var modelOne ModelOne
	if err := c.ShouldBindJSON(&modelOne); err != nil {
		app.Handler.WriteError(c, err, "Invalid request payload")
		return
	}
	if err := app.Handler.CreateWithUser(c, &modelOne); err != nil {
		app.Handler.WriteError(c, err, "Failed to create ModelOne")
		return
	}
	app.Handler.WriteSuccess(c, modelOne)
}

func (app *App) ListModelOnesHandler(c *gin.Context) {
	var (
		models []models.UserOwnedModel
		err    error
	)
	if models, err = app.Handler.FindWithUser(c, &ModelOne{}, handlers.NilQuery); err != nil {
		app.Handler.WriteError(c, err, "Failed to list ModelOnes")
		return
	}
	app.Handler.WriteSuccess(c, models)
}
