package testsuite

import (
	"fmt"
	"log"
	"net/http"
)

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

func (c *GoiterClient) CreateModelOne(name string) (modelOne map[string]interface{}, err error) {
	body := map[string]string{
		"name": name,
	}
	cliResp := &ClientResponse{}
	if cliResp, err = c.makeRequest(&ClientRequest{
		Method: "POST",
		URL:    "/model_ones",
		Body:   body,
	}); err != nil {
		return
	}
	modelOne = cliResp.RespBody
	return
}

func (c *GoiterClient) ListModelOnes() (models []interface{}, err error) {
	cliResp := &ClientResponse{}
	if cliResp, err = c.makeRequest(&ClientRequest{
		Method: "GET",
		URL:    "/model_ones",
		Body:   nil,
	}); err != nil {
		return
	}
	models = cliResp.RespBody["data"].([]interface{})
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
	if _, err = c.CreateModelOne("Test Model One"); err != nil {
		return
	}
	//if _, err = c.ListModelOnes(); err != nil {
	//	return
	//}

	log.Println("App test suite completed successfully.")
	return nil
}
