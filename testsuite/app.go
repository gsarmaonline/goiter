package testsuite

import (
	"fmt"
	"log"

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

		Name string `json:"name" gorm:"not null"`
	}

	ModelTwo struct {
		models.BaseModelWithUser
		Title string `json:"title" gorm:"not null"`
	}
)

func NewApp(srv *core.Server) (app *App, err error) {
	app = &App{
		Server: srv,
	}
	app.DbMgr.RegisterModels("model_one", &ModelOne{})
	app.DbMgr.RegisterModels("model_two", &ModelTwo{})

	app.Handler.OpenRouteGroup.GET("/app_ping", app.Ping)
	app.Handler.ProtectedRouteGroup.GET("/app_protected_ping", app.Ping)
	return
}

func (app *App) Ping(c *gin.Context) {
	app.Handler.WriteSuccess(c, "pong")
}

func (c *GoiterClient) RunAppTestSuite() (err error) {
	log.Println("Running app test suite...")

	// Test the /app_ping endpoint
	cliResp := &ClientResponse{}
	cliResp, err = c.makeRequest(&ClientRequest{
		Method: "GET",
		URL:    "/app_ping",
		Body:   nil,
	})
	if err != nil {
		return err
	}

	if cliResp.RespBody["data"] != "pong" {
		return fmt.Errorf("expected 'pong', got '%s'", cliResp.RespBody["data"])
	}

	log.Println("App test suite completed successfully.")
	return nil
}
