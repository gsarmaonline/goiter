package testsuite

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type GoiterClient struct {
	BaseURL    string
	httpClient *http.Client
	sessionID  string
}

// NewGoiterClient creates a new client instance
func NewGoiterClient(baseURL string) *GoiterClient {
	if baseURL == "" {
		baseURL = "http://localhost:8090"
	}

	return &GoiterClient{
		BaseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// makeRequest makes an authenticated HTTP request
func (c *GoiterClient) makeRequest(method, endpoint string, body io.Reader) (*http.Response, error) {
	if c.sessionID == "" {
		return nil, fmt.Errorf("not authenticated - please login first")
	}

	req, err := http.NewRequest(method, c.BaseURL+endpoint, body)
	if err != nil {
		return nil, err
	}

	// Add session cookie
	req.AddCookie(&http.Cookie{
		Name:  "session",
		Value: c.sessionID,
	})

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return c.httpClient.Do(req)
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

// Logout clears the session
func (c *GoiterClient) Logout() error {
	if c.sessionID == "" {
		return nil
	}

	resp, err := c.makeRequest("POST", "/logout", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	c.sessionID = ""
	fmt.Println("Logged out successfully!")
	return nil
}

// Example usage and CLI interface
func Run() {

	// Initialize client
	baseURL := os.Getenv("GOITER_BASE_URL")
	client := NewGoiterClient(baseURL)

	client.RunUserSuite()
	client.RunProfileSuite()
	client.RunAccountSuite()
	client.RunProjectSuite()

}
