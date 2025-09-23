package app

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gsarmaonline/goiter/core"
	"github.com/gsarmaonline/goiter/core/models"
)

type (
	App struct {
		*core.Server
	}

	ModelOne struct {
		models.BaseModelWithUser

		Name string `json:"name"`
	}

	ModelTwo struct {
		models.BaseModelWithUser
		Title string `json:"title"`
	}
)

func (c ModelOne) GetConfig() models.ModelConfig {
	return models.ModelConfig{
		Name:      "ModelOne",
		ScopeType: models.ProjectScopeType,
	}
}

func (c ModelTwo) GetConfig() models.ModelConfig {
	return models.ModelConfig{
		Name:      "ModelTwo",
		ScopeType: models.ProjectScopeType,
	}
}

func NewApp(srv *core.Server) (app *App, err error) {
	app = &App{
		Server: srv,
	}
	if err = app.DbMgr.RegisterModels([]models.UserOwnedModel{&ModelOne{}, &ModelTwo{}}); err != nil {
		log.Println(err)
		os.Exit(-1)
	}

	app.Handler.OpenRouteGroup.GET("/app_ping", app.Ping)
	app.Handler.ProtectedRouteGroup.GET("/app_protected_ping", app.Ping)
	app.Handler.ProtectedRouteGroup.GET("/model_ones", app.ListModelOnesHandler)
	app.Handler.ProtectedRouteGroup.POST("/model_ones", app.CreateModelOneHandler)
	return
}

func (app *App) Start() (err error) {
	if err = app.Server.Start(); err != nil {
		log.Fatalln(err)
		return
	}
	return
}

func (app *App) Ping(c *gin.Context) {
	app.Handler.WriteSuccess(c, "pong")
}
