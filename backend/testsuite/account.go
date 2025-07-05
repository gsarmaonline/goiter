package testsuite

import "log"

// GetAccount retrieves the current user's account
func (c *GoiterClient) GetAccount() (respBody map[string]interface{}, err error) {
	if _, respBody, err = c.makeRequest("GET", "/account", nil); err != nil {
		return
	}
	return
}

func (c *GoiterClient) RunAccountSuite() (err error) {
	respBody := make(map[string]interface{})
	if respBody, err = c.GetAccount(); err != nil {
		return
	}
	log.Println("Fetched account details:", respBody)
	return
}
