package testsuite

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// GetUser retrieves the current user information
func (c *GoiterClient) GetUser() (user map[string]interface{}, err error) {
	cliResp := &ClientResponse{}
	if cliResp, err = c.makeRequest(&ClientRequest{
		Method: "GET",
		URL:    "/me",
		Body:   nil,
	}); err != nil {
		return
	}
	user = cliResp.RespBody
	return
}

func (c *GoiterClient) shortCircuitLogin(baseURL, email string) (token string, err error) {
	var (
		resp *http.Response
	)

	reqBody := map[string]string{
		"email": email,
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
		err = fmt.Errorf("Failed to login: %s", err.Error())
		return
	}
	respBody, _ := io.ReadAll(resp.Body)
	var respData map[string]string
	if err = json.Unmarshal(respBody, &respData); err != nil {
		return
	}
	token = respData["token"]
	return
}

// Login initiates the Google OAuth flow
func (c *GoiterClient) Login(email string) (map[string]interface{}, error) {
	fmt.Println("🔐 Goiter Client Login")
	var (
		jwtToken string
		err      error
	)
	if jwtToken, err = c.shortCircuitLogin(c.BaseURL, email); err != nil {
		return nil, fmt.Errorf("failed to login: %v", err)
	}
	if jwtToken == "" {
		return nil, fmt.Errorf("no jwt token provided")
	}

	// Set the session cookie
	c.jwtToken = jwtToken

	// Test the authentication by making a request to /me
	fmt.Println("\n🔄 Testing authentication for user", email)
	user, err := c.GetUser()
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %v", err)
	}

	fmt.Println("✅ Login successful! Welcome", user["email"])
	return user, nil
}

// Logout clears the session
func (c *GoiterClient) Logout() (err error) {
	if c.jwtToken == "" {
		return nil
	}

	if _, err := c.makeRequest(&ClientRequest{
		Method: "POST",
		URL:    "/logout",
		Body:   nil,
	}); err != nil {
		return err
	}

	c.jwtToken = ""
	fmt.Println("Logged out successfully!")
	return nil
}

func (c *GoiterClient) RunUserSuite() (err error) {
	if _, err = c.Login(c.users["root"].Email); err != nil {
		log.Println("Login failed:", err)
		return
	}
	return
}
