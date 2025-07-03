package testsuite

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Account struct {
	ID                 uint   `json:"id"`
	Name               string `json:"name"`
	Description        string `json:"description"`
	PlanID             uint   `json:"plan_id"`
	SubscriptionStatus string `json:"subscription_status"`
}

// GetAccount retrieves the current user's account
func (c *GoiterClient) GetAccount() (*Account, error) {
	resp, err := c.makeRequest("GET", "/account", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get account: %s", string(body))
	}

	var account Account
	if err := json.NewDecoder(resp.Body).Decode(&account); err != nil {
		return nil, err
	}

	return &account, nil
}

func (c *GoiterClient) RunAccountSuite() (err error) {
	return
}
