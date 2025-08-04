package testsuite

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gsarmaonline/goiter/core/models"
)

type (
	GoiterClient struct {
		BaseURL    string
		httpClient *http.Client

		jwtToken string

		users map[string]*models.User
	}

	ClientRequest struct {
		Method   string
		URL      string
		Body     interface{}
		SkipAuth bool
	}
	ClientResponse struct {
		Resp     *http.Response
		RespBody map[string]interface{}
	}
)

func (c *GoiterClient) Errorf(format string, args ...interface{}) {
	log.Printf("ERROR: "+format, args...)
}

// NewGoiterClient creates a new client instance
func NewGoiterClient() (gc *GoiterClient) {
	gc = &GoiterClient{
		BaseURL: "http://localhost:8090",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		users: make(map[string]*models.User),
	}
	gc.users["root"] = &models.User{
		Email: "root@example.com",
	}
	go gc.StartServer()
	time.Sleep(2 * time.Second) // Wait for the server to start
	return
}

// makeRequest makes an authenticated HTTP request
func (c *GoiterClient) makeRequest(cliReq *ClientRequest) (cliResp *ClientResponse, err error) {

	var (
		reqBody *bytes.Reader
		req     *http.Request
	)
	if c.jwtToken == "" {
		err = errors.New("jwt token is not set")
		return
	}
	bodyb, _ := json.Marshal(cliReq.Body)
	reqBody = bytes.NewReader(bodyb)

	if req, err = http.NewRequest(cliReq.Method, c.BaseURL+cliReq.URL, reqBody); err != nil {
		return
	}

	// Add Authorization header
	req.Header.Set("Authorization", "Bearer "+c.jwtToken)

	if cliReq.Body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	cliResp = &ClientResponse{}

	if cliResp.Resp, err = c.httpClient.Do(req); err != nil {
		log.Println(err)
		return
	}
	defer cliResp.Resp.Body.Close()

	cliResp.RespBody = make(map[string]interface{})
	respB, err := io.ReadAll(cliResp.Resp.Body)

	if cliResp.Resp.StatusCode != http.StatusOK && err == nil {
		err = fmt.Errorf("request failed with status: %d with body", cliResp.Resp.StatusCode, string(respB))
		return
	}
	if err = json.Unmarshal(respB, &cliResp.RespBody); err != nil {
		return
	}
	if cliResp.RespBody["data"] != nil {
		if data, ok := cliResp.RespBody["data"].(map[string]interface{}); ok {
			cliResp.RespBody = data
			return
		}
	}

	return
}

// Ping tests the connection to the server
func (c *GoiterClient) Ping() error {
	resp, err := http.Get(c.BaseURL + "/ping")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status: %d", resp.StatusCode)
	}

	return nil
}

func (c *GoiterClient) Run() (err error) {
	log.Println("🚀 Starting Goiter Test Suite...")

	// Run basic functional tests
	log.Println("📋 Running Basic Functional Tests...")
	if err := c.RunUserSuite(); err != nil {
		log.Fatalf("❌ User suite failed: %v", err)
	}
	if err := c.RunProfileSuite(); err != nil {
		log.Fatalf("❌ Profile suite failed: %v", err)
	}
	if err := c.RunAccountSuite(); err != nil {
		log.Fatalf("❌ Account suite failed: %v", err)
	}
	if err := c.RunProjectSuite(); err != nil {
		log.Fatalf("❌ Project suite failed: %v", err)
	}
	if err := c.RunAuthorisationSuite(); err != nil {
		log.Fatalf("❌ Authorisation suite failed: %v", err)
	}
	if err := c.RunAppTestSuite(); err != nil {
		log.Fatalf("❌ App test suite failed: %v", err)
	}
	return
}

func Run() {
	c := NewGoiterClient()
	if err := c.Run(); err != nil {
		log.Fatalf("Test suite failed: %v", err)
	}
}
