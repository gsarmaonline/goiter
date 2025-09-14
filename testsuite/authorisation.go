package testsuite

import (
	"log"

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

func (gc *GoiterClient) CheckForDirectUserAndRoleAccess() (err error) {
	return
}

func (gc *GoiterClient) CheckForIndirectUserAndDirectRoleAccess() (err error) {
	return
}

func (gc *GoiterClient) CheckForDirectUserAndIndirectRoleAccess() (err error) {
	return
}

func (gc *GoiterClient) CheckForIndirectUserAndIndirectRoleAccess() (err error) {
	return
}

func (gc *GoiterClient) CheckForMultilevelIndirectUserAndIndirectRoleAccess() (err error) {
	return
}

func (c *GoiterClient) RunAuthorisationSuite() (err error) {
	log.Println("🔐 Running Authorisation tests...")

	// Test cases
	if err = c.CheckForDirectUserAndRoleAccess(); err != nil {
		return err
	}
	if err = c.CheckForIndirectUserAndDirectRoleAccess(); err != nil {
		return err
	}
	if err = c.CheckForDirectUserAndIndirectRoleAccess(); err != nil {
		return err
	}
	if err = c.CheckForIndirectUserAndIndirectRoleAccess(); err != nil {
		return err
	}
	if err = c.CheckForMultilevelIndirectUserAndIndirectRoleAccess(); err != nil {
		return err
	}
	log.Println("✅ Authorisation tests passed")
	return nil
}
