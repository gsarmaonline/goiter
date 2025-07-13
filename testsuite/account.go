package testsuite

import (
	"fmt"
	"log"
)

// GetAccount retrieves the current user's account
func (c *GoiterClient) GetAccount() (respBody map[string]interface{}, err error) {
	if _, respBody, err = c.makeRequest("GET", "/account", nil); err != nil {
		return
	}
	return
}

// UpdateAccount updates the current user's account
func (c *GoiterClient) UpdateAccount(name, description string) (respBody map[string]interface{}, err error) {
	body := map[string]string{
		"name":        name,
		"description": description,
	}
	if _, respBody, err = c.makeRequest("PUT", "/account", body); err != nil {
		return
	}
	return
}

func (c *GoiterClient) RunAccountSuite() (err error) {
	log.Println("Running Account test suite...")

	// Get the account
	log.Println("Getting the account...")
	account, err := c.GetAccount()
	if err != nil {
		return fmt.Errorf("failed to get account: %v", err)
	}
	log.Println("Fetched account details:", account)

	// Update the account
	log.Println("Updating the account...")
	updatedAccount, err := c.UpdateAccount("Updated Test Account", "This is an updated test account.")
	if err != nil {
		return fmt.Errorf("failed to update account: %v", err)
	}
	log.Println("Account updated:", updatedAccount)

	// Get the account again to verify the changes
	log.Println("Getting the account again...")
	account, err = c.GetAccount()
	if err != nil {
		return fmt.Errorf("failed to get account: %v", err)
	}
	log.Println("Fetched account details:", account)

	return
}
