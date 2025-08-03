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

	"github.com/gsarmaonline/goiter/config"
	"github.com/gsarmaonline/goiter/core"
	"github.com/joho/godotenv"
)

type GoiterClient struct {
	BaseURL    string
	httpClient *http.Client

	jwtToken string
}

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
	}
	go gc.StartServer()
	time.Sleep(2 * time.Second) // Wait for the server to start
	return
}

func (c *GoiterClient) StartServer() {
	fmt.Println("Attempting to start server...")

	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found or error loading it: %v", err)
	}
	cfg := config.DefaultConfig()
	cfg.Mode = config.ModeDev
	cfg.Port = "8090"
	cfg.DBType = config.SqliteDbType

	server := core.NewServer(cfg)

	log.Printf("Starting server on :%s", cfg.Port)

	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// makeRequest makes an authenticated HTTP request
func (c *GoiterClient) makeRequest(method, endpoint string,
	body interface{}) (resp *http.Response, respBody map[string]interface{}, err error) {

	var (
		reqBody *bytes.Reader
		req     *http.Request
	)
	if c.jwtToken == "" {
		err = errors.New("jwt token is not set")
		return
	}
	bodyb, _ := json.Marshal(body)
	reqBody = bytes.NewReader(bodyb)

	if req, err = http.NewRequest(method, c.BaseURL+endpoint, reqBody); err != nil {
		return
	}

	// Add Authorization header
	req.Header.Set("Authorization", "Bearer "+c.jwtToken)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if resp, err = c.httpClient.Do(req); err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	respBody = make(map[string]interface{})
	respB, err := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK && err == nil {
		err = fmt.Errorf("request failed with status: %d with body", resp.StatusCode, string(respB))
		return
	}
	if err = json.Unmarshal(respB, &respBody); err != nil {
		return
	}
	if respBody["data"] != nil {
		if data, ok := respBody["data"].(map[string]interface{}); ok {
			respBody = data
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
	return
}

func Run() {
	c := NewGoiterClient()
	if err := c.Run(); err != nil {
		log.Fatalf("Test suite failed: %v", err)
	}
}
