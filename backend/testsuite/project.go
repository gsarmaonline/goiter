package testsuite

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Project struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	AccountID   uint   `json:"account_id"`
	UserID      uint   `json:"user_id"`
}

// GetProjects retrieves all projects for the current user
func (c *GoiterClient) GetProjects() ([]Project, error) {
	resp, err := c.makeRequest("GET", "/projects", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get projects: %s", string(body))
	}

	var projects []Project
	if err := json.NewDecoder(resp.Body).Decode(&projects); err != nil {
		return nil, err
	}

	return projects, nil
}

func (c *GoiterClient) RunProjectSuite() (err error) {
	return
}
