package testsuite

import (
	"log"

	"github.com/gsarmaonline/goiter/core/models"
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
	if _, respBody, err = c.makeRequest("PUT", "/profile", profile); err != nil {
		return
	}
	log.Println("Profile updated successfully:", respBody)
	return
}

func (c *GoiterClient) GetProfile() (respBody map[string]interface{}, err error) {
	if _, respBody, err = c.makeRequest("GET", "/profile", nil); err != nil {
		return nil, err
	}
	return
}

func (c *GoiterClient) RunProfileSuite() (err error) {
	respBody := make(map[string]interface{})
	if _, err = c.UpdateProfile(); err != nil {
		return
	}
	if respBody, err = c.GetProfile(); err != nil {
		return
	}
	log.Println(respBody)
	return
}
