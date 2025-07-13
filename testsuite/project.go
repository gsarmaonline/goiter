package testsuite

import (
	"fmt"
	"log"

	"github.com/gsarmaonline/goiter/core/models"
)

// GetProjects retrieves all projects for the current user
func (c *GoiterClient) GetProjects() (projects map[string]interface{}, err error) {
	if _, projects, err = c.makeRequest("GET", "/projects", nil); err != nil {
		return nil, err
	}
	return
}

// CreateProject creates a new project
func (c *GoiterClient) CreateProject(name, description string) (project map[string]interface{}, err error) {
	body := map[string]string{
		"name":        name,
		"description": description,
	}
	_, project, err = c.makeRequest("POST", "/projects", body)
	return
}

// GetProject retrieves a single project by its ID
func (c *GoiterClient) GetProject(id uint) (project map[string]interface{}, err error) {
	_, project, err = c.makeRequest("GET", fmt.Sprintf("/projects/%d", id), nil)
	return
}

// UpdateProject updates a project's details
func (c *GoiterClient) UpdateProject(id uint, name, description string) (project map[string]interface{}, err error) {
	body := map[string]string{
		"name":        name,
		"description": description,
	}
	_, project, err = c.makeRequest("PUT", fmt.Sprintf("/projects/%d", id), body)
	return
}

// DeleteProject deletes a project by its ID
func (c *GoiterClient) DeleteProject(id uint) (err error) {
	_, _, err = c.makeRequest("DELETE", fmt.Sprintf("/projects/%d", id), nil)
	return
}

// AddProjectMember adds a member to a project
func (c *GoiterClient) AddProjectMember(projectID uint, userEmail string, level models.PermissionLevel) (permission map[string]interface{}, err error) {
	body := map[string]interface{}{
		"user_email": userEmail,
		"level":      level,
	}
	_, permission, err = c.makeRequest("POST", fmt.Sprintf("/projects/%d/members", projectID), body)
	return
}

// RemoveProjectMember removes a member from a project
func (c *GoiterClient) RemoveProjectMember(projectID, userID uint) (err error) {
	_, _, err = c.makeRequest("DELETE", fmt.Sprintf("/projects/%d/members/%d", projectID, userID), nil)
	return
}

func (c *GoiterClient) RunProjectSuite() (err error) {
	log.Println("Running Project test suite...")

	// Create a new project
	log.Println("Creating a new project...")
	createdProject, err := c.CreateProject("Test Project", "This is a test project.")
	if err != nil {
		return fmt.Errorf("failed to create project: %v", err)
	}
	log.Println("Project created:", createdProject)
	projectData := createdProject["data"].(map[string]interface{})
	projectID := uint(projectData["id"].(float64))

	// Get the project
	log.Println("Getting the project...")
	project, err := c.GetProject(projectID)
	if err != nil {
		return fmt.Errorf("failed to get project: %v", err)
	}
	log.Println("Project details:", project)

	// Update the project
	log.Println("Updating the project...")
	updatedProject, err := c.UpdateProject(projectID, "Updated Test Project", "This is an updated test project.")
	if err != nil {
		return fmt.Errorf("failed to update project: %v", err)
	}
	log.Println("Project updated:", updatedProject)

	// Add a member to the project
	log.Println("Adding a member to the project...")
	addedMember, err := c.AddProjectMember(projectID, "testuser@example.com", models.PermissionEditor)
	if err != nil {
		return fmt.Errorf("failed to add member to project: %v", err)
	}
	log.Println("Member added:", addedMember)
	memberData := addedMember["data"].(map[string]interface{})
	userID := uint(memberData["user_id"].(float64))

	// List projects
	log.Println("Listing projects...")
	projects, err := c.GetProjects()
	if err != nil {
		return fmt.Errorf("failed to list projects: %v", err)
	}
	log.Println("Fetched projects:", projects)

	// Remove a member from the project
	log.Println("Removing a member from the project...")
	if err := c.RemoveProjectMember(projectID, userID); err != nil {
		return fmt.Errorf("failed to remove member from project: %v", err)
	}
	log.Println("Member removed successfully.")

	// Delete the project
	log.Println("Deleting the project...")
	if err := c.DeleteProject(projectID); err != nil {
		return fmt.Errorf("failed to delete project: %v", err)
	}
	log.Println("Project deleted successfully.")

	return nil
}
