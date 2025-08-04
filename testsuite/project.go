package testsuite

import (
	"fmt"
	"log"

	"github.com/gsarmaonline/goiter/core/models"
	"github.com/stretchr/testify/assert"
)

// GetProjects retrieves all projects for the current user
func (c *GoiterClient) GetProjects() (projects []interface{}, err error) {
	cliResp := &ClientResponse{}
	if cliResp, err = c.makeRequest(&ClientRequest{
		Method: "GET",
		URL:    "/projects",
		Body:   nil,
	}); err != nil {
		return nil, err
	}
	projects = cliResp.RespBody["data"].([]interface{})
	return
}

// CreateProject creates a new project
func (c *GoiterClient) CreateProject(name, description string) (project map[string]interface{}, err error) {
	body := map[string]string{
		"name":        name,
		"description": description,
	}
	cliResp := &ClientResponse{}
	if cliResp, err = c.makeRequest(&ClientRequest{
		Method: "POST",
		URL:    "/projects",
		Body:   body,
	}); err != nil {
		return
	}
	project = cliResp.RespBody
	return
}

// GetProject retrieves a single project by its ID
func (c *GoiterClient) GetProject(id uint) (project map[string]interface{}, err error) {
	cliResp := &ClientResponse{}
	if cliResp, err = c.makeRequest(&ClientRequest{
		Method: "GET",
		URL:    fmt.Sprintf("/projects/%d", id),
		Body:   nil,
	}); err != nil {
		return nil, err
	}
	project = cliResp.RespBody
	return
}

// UpdateProject updates a project's details
func (c *GoiterClient) UpdateProject(id uint, name, description string) (project map[string]interface{}, err error) {
	body := map[string]string{
		"name":        name,
		"description": description,
	}
	cliResp := &ClientResponse{}
	if cliResp, err = c.makeRequest(&ClientRequest{
		Method: "PUT",
		URL:    fmt.Sprintf("/projects/%d", id),
		Body:   body,
	}); err != nil {
		return nil, err
	}
	project = cliResp.RespBody
	return
}

// DeleteProject deletes a project by its ID
func (c *GoiterClient) DeleteProject(id uint) (err error) {
	if _, err = c.makeRequest(&ClientRequest{
		Method: "DELETE",
		URL:    fmt.Sprintf("/projects/%d", id),
		Body:   nil,
	}); err != nil {
		return err
	}
	return
}

// AddProjectMember adds a member to a project
func (c *GoiterClient) AddProjectMember(projectID uint, userEmail string, level models.PermissionLevel) (permission map[string]interface{}, err error) {
	body := map[string]interface{}{
		"user_email": userEmail,
		"level":      level,
	}
	cliResp := &ClientResponse{}
	if cliResp, err = c.makeRequest(&ClientRequest{
		Method: "POST",
		URL:    fmt.Sprintf("/projects/%d/members", projectID),
		Body:   body,
	}); err != nil {
		return nil, err
	}
	permission = cliResp.RespBody
	return
}

// RemoveProjectMember removes a member from a project
func (c *GoiterClient) RemoveProjectMember(projectID, userID uint) (err error) {
	if _, err = c.makeRequest(&ClientRequest{
		Method: "DELETE",
		URL:    fmt.Sprintf("/projects/%d/members/%d", projectID, userID),
		Body:   nil,
	}); err != nil {
		return err
	}
	return
}

func (c *GoiterClient) RunProjectCRUDSuite() (err error) {
	return
}

func (c *GoiterClient) RunProjectMembersSuite() (err error) {
	return
}

func (c *GoiterClient) RunProjectSuite() (err error) {
	log.Println("Running Project test suite...")

	createdProject, err := c.CreateProject("Test Project", "This is a test project.")
	assert.Nil(c, err, "Failed to create project")
	assert.Equal(c, "Test Project", createdProject["name"], "Project name should be 'Test Project'")

	projectID := uint(createdProject["id"].(float64))

	// Get the project
	project, err := c.GetProject(projectID)
	assert.Nil(c, err, "Failed to get project")
	assert.Equal(c, "Test Project", project["name"], "Project name should be 'Test Project'")

	// Update the project
	updatedProject, err := c.UpdateProject(projectID, "Updated Test Project", "This is an updated test project.")
	assert.Nil(c, err, "Failed to update project")
	assert.Equal(c, "Updated Test Project", updatedProject["name"], "Project name should be 'Updated Test Project'")

	// Add a member to the project
	addedMember, err := c.AddProjectMember(projectID, "testuser@example.com", models.PermissionEditor)
	assert.Nil(c, err, "Failed to add member to project")
	userID := uint(addedMember["user_id"].(float64))

	// List projects
	projects, err := c.GetProjects()
	assert.Nil(c, err, "Failed to list projects")
	assert.NotEmpty(c, projects, "Projects list should not be empty")

	// Remove a member from the project
	err = c.RemoveProjectMember(projectID, userID)
	assert.Nil(c, err, "Failed to remove member from project")

	// Delete the project
	err = c.DeleteProject(projectID)
	assert.Nil(c, err, "Failed to delete project")

	return
}
