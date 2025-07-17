package testsuite

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// AuthTestClient wraps GoiterClient with authorization testing capabilities
type AuthTestClient struct {
	*GoiterClient
	userID uint
}

// NewAuthTestClient creates a new authorization test client
func NewAuthTestClient(baseURL string) *AuthTestClient {
	return &AuthTestClient{
		GoiterClient: NewGoiterClient(baseURL),
	}
}

// LoginAsUser logs in as a specific user for testing
func (c *AuthTestClient) LoginAsUser(email string) error {
	token, err := c.shortCircuitLoginWithEmail(email)
	if err != nil {
		return err
	}
	c.jwtToken = token

	// Get user info to store user ID
	user, err := c.GetUser()
	if err != nil {
		return err
	}
	
	// Handle the case where data might be nil
	if user["data"] == nil {
		return fmt.Errorf("user data is nil")
	}
	
	userData, ok := user["data"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("user data is not a map")
	}
	
	userID, ok := userData["id"].(float64)
	if !ok {
		return fmt.Errorf("user id is not a number")
	}
	
	c.userID = uint(userID)
	return nil
}

// shortCircuitLoginWithEmail performs login with a specific email
func (c *AuthTestClient) shortCircuitLoginWithEmail(email string) (token string, err error) {
	var resp *http.Response

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

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("failed to login with status: %d", resp.StatusCode)
		return
	}

	var respData map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return
	}

	if tokenData, ok := respData["token"]; ok {
		token = tokenData.(string)
	} else {
		err = fmt.Errorf("no token in response")
	}
	return
}

// TestPermission tests if a user can perform an action on a resource
func (c *AuthTestClient) TestPermission(method, endpoint string, body interface{}, expectedStatus int) (bool, error) {
	resp, _, err := c.makeRequestWithStatus(method, endpoint, body)
	if err != nil && resp == nil {
		return false, err
	}

	actualStatus := resp.StatusCode
	if actualStatus == expectedStatus {
		return true, nil
	}

	return false, fmt.Errorf("expected status %d but got %d", expectedStatus, actualStatus)
}

// makeRequestWithStatus is like makeRequest but returns response with status code
func (c *GoiterClient) makeRequestWithStatus(method, endpoint string, body interface{}) (*http.Response, map[string]interface{}, error) {
	var (
		reqBody *bytes.Reader
		req     *http.Request
	)

	if c.jwtToken == "" {
		return nil, nil, fmt.Errorf("jwt token is not set")
	}

	bodyb, _ := json.Marshal(body)
	reqBody = bytes.NewReader(bodyb)

	req, err := http.NewRequest(method, c.BaseURL+endpoint, reqBody)
	if err != nil {
		return nil, nil, err
	}

	// Add Authorization header
	req.Header.Set("Authorization", "Bearer "+c.jwtToken)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	respBody := make(map[string]interface{})
	if err = json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		// If we can't decode JSON, it's still a valid response
		respBody = map[string]interface{}{"raw_response": "could not decode"}
	}

	return resp, respBody, nil
}

// GetUserID returns the current user's ID
func (c *AuthTestClient) GetUserID() uint {
	return c.userID
}

// SetupTestUser creates a test user and logs them in
func (c *AuthTestClient) SetupTestUser(email string) error {
	return c.LoginAsUser(email)
}

// CleanupTestUser performs any cleanup for the test user
func (c *AuthTestClient) CleanupTestUser() error {
	return c.Logout()
}
