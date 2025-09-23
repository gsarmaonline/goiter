package app

import (
	"github.com/gin-gonic/gin"
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
		modelOnes []ModelOne
		err       error
	)
	if err = app.Handler.Db.Find(&modelOnes).Error; err != nil {
		app.Handler.WriteError(c, err, "Failed to list ModelOnes")
		return
	}
	app.Handler.WriteSuccess(c, modelOnes)
}
