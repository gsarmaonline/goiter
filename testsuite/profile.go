package testsuite

import (
	"log"

	"github.com/gsarmaonline/goiter/core/models"
	"github.com/stretchr/testify/assert"
)

func (c *GoiterClient) UpdateProfile() (respBody map[string]interface{}, err error) {
	profile := &models.Profile{
		Address:    "123 Main St",
		City:       "Anytown",
		State:      "CA",
		PostalCode: "12345",
		Country:    "USA",

		CompanyName: "Goiter Inc.",
		JobTitle:    "Software Engineer",
		Department:  "Engineering",
	}
	cliResp := &ClientResponse{}
	if cliResp, err = c.makeRequest(&ClientRequest{
		Method: "PUT",
		URL:    "/profile",
		Body:   profile,
	}); err != nil {
		return
	}
	respBody = cliResp.RespBody
	return
}

func (c *GoiterClient) GetProfile() (respBody map[string]interface{}, err error) {
	cliResp := &ClientResponse{}
	if cliResp, err = c.makeRequest(&ClientRequest{
		Method: "GET",
		URL:    "/profile",
		Body:   nil,
	}); err != nil {
		return nil, err
	}
	respBody = cliResp.RespBody
	return
}

func (c *GoiterClient) RunProfileSuite() (err error) {
	log.Println("Running profile suite...")

	_, err = c.UpdateProfile()
	assert.Nil(c, err, "UpdateProfile failed")

	respBody, err := c.GetProfile()
	assert.Nil(c, err, "GetProfile failed")
	assert.Equal(c, "123 Main St", respBody["address"], "Address should be updated")
	assert.Equal(c, "Anytown", respBody["city"], "City should be updated")
	assert.Equal(c, "CA", respBody["state"], "State should be updated")
	assert.Equal(c, "12345", respBody["postal_code"], "Postal code should be updated")
	assert.Equal(c, "USA", respBody["country"], "Country should be updated")
	assert.Equal(c, "Goiter Inc.", respBody["company_name"], "Company name should be updated")
	assert.Equal(c, "Software Engineer", respBody["job_title"], "Job title should be updated")
	assert.Equal(c, "Engineering", respBody["department"], "Department should be updated")
	return
}
