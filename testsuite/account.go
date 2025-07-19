package testsuite

import (
	"log"

	"github.com/stretchr/testify/assert"
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
	account, err := c.GetAccount()
	assert.Nil(c, err, "Failed to get account")
	assert.NotEqual(c, "", account["id"], "Account ID should not be empty")

	// Update the account
	updatedAccount, err := c.UpdateAccount("Updated Test Account", "This is an updated test account.")
	assert.Nil(c, err, "Failed to update account")
	assert.Equal(c, "Updated Test Account", updatedAccount["name"], "Account name should be updated")
	assert.Equal(c, "This is an updated test account.", updatedAccount["description"], "Account description should be updated")

	// Verify the account is updated
	account, err = c.GetAccount()
	assert.Nil(c, err, "Failed to get account")
	assert.NotEqual(c, "", account["id"], "Account ID should not be empty")

	return
}
