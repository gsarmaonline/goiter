package testsuite

import (
	"fmt"

	"github.com/gsarmaonline/goiter/core/models"
)

// Testing Authorisation
// - Define types of users
// - Define the resources to access

type (
	AuthorisationScenario struct {
		gc    *GoiterClient
		users []*models.User
	}
)

func NewAuthorisationScenario(gc *GoiterClient) *AuthorisationScenario {
	return &AuthorisationScenario{
		gc:    gc,
		users: make([]*models.User, 0),
	}
}

func (as *AuthorisationScenario) createAuthUsers() (err error) {
	as.users = append(as.users, &models.User{
		Email: "auth_root@sample.com",
	}, &models.User{
		Email: "auth_user@sample.com",
	})
	for _, user := range as.users {
		if _, err = as.gc.Login(user.Email); err != nil {
			return fmt.Errorf("failed to create user %s: %w", user.Email, err)
		}
	}
	return
}

func (c *GoiterClient) RunAuthorisationSuite() (err error) {

	as := NewAuthorisationScenario(c)

	if err = as.createAuthUsers(); err != nil {
		c.Errorf("Failed to create auth users: %v", err)
		return err
	}
	return nil
}
