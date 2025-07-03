package testsuite

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type User struct {
	ID      uint   `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

// GetUser retrieves the current user information
func (c *GoiterClient) GetUser() (*User, error) {
	resp, err := c.makeRequest("GET", "/me", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get user: %s", string(body))
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (c *GoiterClient) shortCircuitLogin(baseURL string) (token string, err error) {
	var (
		resp *http.Response
	)

	reqBody := map[string]string{
		"email": "user1@gmail.com",
	}
	reqBodyJSON, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", c.BaseURL+"/auth/shortcircuitlogin", bytes.NewReader(reqBodyJSON))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")

	if resp, err = c.httpClient.Do(req); err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK || err != nil {
		err = fmt.Errorf("Failed to login: %s", err)
		return
	}
	respBody, _ := io.ReadAll(resp.Body)
	var respData map[string]string
	if err = json.Unmarshal(respBody, &respData); err != nil {
		return
	}
	token = respData["code"]
	return
}

// Login initiates the Google OAuth flow
func (c *GoiterClient) Login() error {
	fmt.Println("🔐 Goiter Client Login")
	var (
		sessionCookie string
		err           error
	)
	if sessionCookie, err = c.shortCircuitLogin(c.BaseURL); err != nil {
		return fmt.Errorf("failed to login: %v", err)
	}
	if sessionCookie == "" {
		return fmt.Errorf("no session cookie provided")
	}

	// Set the session cookie
	c.sessionID = sessionCookie

	// Test the authentication by making a request to /me
	fmt.Println("\n🔄 Testing authentication...")
	user, err := c.GetUser()
	if err != nil {
		return fmt.Errorf("authentication failed: %v", err)
	}

	fmt.Printf("✅ Login successful! Welcome, %s (%s)\n", user.Name, user.Email)
	return nil
}

func (c *GoiterClient) RunUserSuite() (err error) {
	return
}
