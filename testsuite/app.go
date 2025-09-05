package testsuite

import (
	"fmt"
	"log"
	"net/http"
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

		Name string `json:"name" gorm:"not null"`
	}

	ModelTwo struct {
		models.BaseModelWithUser
		Title string `json:"title" gorm:"not null"`
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
	return
}

func (app *App) Ping(c *gin.Context) {
	app.Handler.WriteSuccess(c, "pong")
}

func (c *GoiterClient) PingOpenRoute() (err error) {
	// Test the /app_ping endpoint
	cliResp := &ClientResponse{}
	cliResp, err = c.makeRequest(&ClientRequest{
		Method:   "GET",
		URL:      "/app_ping",
		Body:     nil,
		SkipAuth: true,
	})
	if err != nil {
		return err
	}

	if cliResp.RespBody["data"] != "pong" {
		return fmt.Errorf("expected 'pong', got '%s'", cliResp.RespBody["data"])
	}

	return
}

func (c *GoiterClient) PingProtectedRoute() (err error) {
	// Test the /app_protected_ping endpoint
	cliResp := &ClientResponse{}
	cliResp, err = c.makeRequest(&ClientRequest{
		Method:   "GET",
		URL:      "/app_protected_ping",
		Body:     nil,
		SkipAuth: true,
	})
	if cliResp.Resp.StatusCode != http.StatusUnauthorized {
		err = fmt.Errorf("expected status 401 Unauthorized, got %d", cliResp.Resp.StatusCode)
	}

	// Enable auth and try
	cliResp, err = c.makeRequest(&ClientRequest{
		Method:   "GET",
		URL:      "/app_protected_ping",
		Body:     nil,
		SkipAuth: false,
	})
	if cliResp.Resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("expected status 200, got %d", cliResp.Resp.StatusCode)
	}

	return
}

func (c *GoiterClient) RunAppTestSuite() (err error) {
	log.Println("Running app test suite...")
	if err = c.PingOpenRoute(); err != nil {
		return
	}
	if err = c.PingProtectedRoute(); err != nil {
		return
	}

	log.Println("App test suite completed successfully.")
	return nil
}
