package testsuite

import (
	"log"

	"github.com/gsarmaonline/goiter/core/models"
)

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

func (auth *AuthorisationScenario) Setup() (err error) {
	auth.gc.CreateAndLogin("direct_user_auth@example.com", false)
	auth.gc.CreateAndLogin("indirect_user_auth@example.com", false)
	return
}

func (auth *AuthorisationScenario) CheckForOwnerUserAndRoleAccess() (err error) {
	return
}

func (auth *AuthorisationScenario) CheckForDirectUserAndRoleAccess() (err error) {
	return
}

func (auth *AuthorisationScenario) CheckForIndirectUserAndDirectRoleAccess() (err error) {
	return
}

func (auth *AuthorisationScenario) CheckForDirectUserAndIndirectRoleAccess() (err error) {
	return
}

func (auth *AuthorisationScenario) CheckForIndirectUserAndIndirectRoleAccess() (err error) {
	return
}

func (auth *AuthorisationScenario) CheckForMultilevelIndirectUserAndIndirectRoleAccess() (err error) {
	return
}

func (c *GoiterClient) RunAuthorisationSuite() (err error) {
	log.Println("🔐 Running Authorisation tests...")

	authorisation := NewAuthorisationScenario(c)

	// Test cases
	if err = authorisation.CheckForDirectUserAndRoleAccess(); err != nil {
		return err
	}
	if err = authorisation.CheckForOwnerUserAndRoleAccess(); err != nil {
		return err
	}
	if err = authorisation.CheckForIndirectUserAndDirectRoleAccess(); err != nil {
		return err
	}
	if err = authorisation.CheckForDirectUserAndIndirectRoleAccess(); err != nil {
		return err
	}
	if err = authorisation.CheckForIndirectUserAndIndirectRoleAccess(); err != nil {
		return err
	}
	if err = authorisation.CheckForMultilevelIndirectUserAndIndirectRoleAccess(); err != nil {
		return err
	}
	log.Println("✅ Authorisation tests passed")
	return nil
}
