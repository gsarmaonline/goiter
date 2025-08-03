package testsuite

import "github.com/gsarmaonline/goiter/core/models"

// Testing Authorisation
// - Define types of users
// - Define the resources to access

type (
	AuthorisationScenario struct {
		Users []*models.User
	}
)

func (c *GoiterClient) RunAuthorisationSuite() error {
	// Implement the logic to run the authorisation scenario
	return nil
}
