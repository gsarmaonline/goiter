package testsuite

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type GoiterClient struct {
	BaseURL    string
	httpClient *http.Client

	jwtToken string
}

// NewGoiterClient creates a new client instance
func NewGoiterClient(baseURL string) (gc *GoiterClient) {
	if baseURL == "" {
		baseURL = "http://localhost:8090"
	}

	gc = &GoiterClient{
		BaseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
	return
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
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && err == nil {
		err = fmt.Errorf("request failed with status: %d", resp.StatusCode)
		return
	}

	respBody = make(map[string]interface{})
	respB, err := io.ReadAll(resp.Body)
	if err = json.Unmarshal(respB, &respBody); err != nil {
		return
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

func Run() {

	// Initialize client
	baseURL := os.Getenv("GOITER_BASE_URL")
	client := NewGoiterClient(baseURL)

	if err := client.RunUserSuite(); err != nil {
		log.Fatalf("User suite failed: %v", err)
	}
	if err := client.RunProfileSuite(); err != nil {
		log.Fatalf("Profile suite failed: %v", err)
	}
	if err := client.RunAccountSuite(); err != nil {
		log.Fatalf("Account suite failed: %v", err)
	}
	if err := client.RunProjectSuite(); err != nil {
		log.Fatalf("Project suite failed: %v", err)
	}

}
