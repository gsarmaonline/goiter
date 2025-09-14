package testsuite

import (
	"fmt"
	"log"

	"github.com/stretchr/testify/assert"
)

// GetGroups retrieves all groups for the current user
func (c *GoiterClient) GetGroups() (groups []interface{}, err error) {
	cliResp := &ClientResponse{}
	if cliResp, err = c.makeRequest(&ClientRequest{
		Method: "GET",
		URL:    "/groups",
		Body:   nil,
	}); err != nil {
		return nil, err
	}
	groups = cliResp.RespBody["data"].([]interface{})
	return
}

// CreateGroup creates a new group
func (c *GoiterClient) CreateGroup(name, description string) (group map[string]interface{}, err error) {
	body := map[string]string{
		"name":        name,
		"description": description,
	}
	cliResp := &ClientResponse{}
	if cliResp, err = c.makeRequest(&ClientRequest{
		Method: "POST",
		URL:    "/groups",
		Body:   body,
	}); err != nil {
		return
	}
	group = cliResp.RespBody
	return
}

// GetGroup retrieves a single group by its ID
func (c *GoiterClient) GetGroup(id uint) (group map[string]interface{}, err error) {
	cliResp := &ClientResponse{}
	if cliResp, err = c.makeRequest(&ClientRequest{
		Method: "GET",
		URL:    fmt.Sprintf("/groups/%d", id),
		Body:   nil,
	}); err != nil {
		return nil, err
	}
	group = cliResp.RespBody
	return
}

// UpdateGroup updates a group's details
func (c *GoiterClient) UpdateGroup(id uint, name, description string) (group map[string]interface{}, err error) {
	body := map[string]string{
		"name":        name,
		"description": description,
	}
	cliResp := &ClientResponse{}
	if cliResp, err = c.makeRequest(&ClientRequest{
		Method: "PUT",
		URL:    fmt.Sprintf("/groups/%d", id),
		Body:   body,
	}); err != nil {
		return nil, err
	}
	group = cliResp.RespBody
	return
}

// DeleteGroup deletes a group by its ID
func (c *GoiterClient) DeleteGroup(id uint) (err error) {
	if _, err = c.makeRequest(&ClientRequest{
		Method: "DELETE",
		URL:    fmt.Sprintf("/groups/%d", id),
		Body:   nil,
	}); err != nil {
		return err
	}
	return
}

// AddGroupMember adds a member to a group
func (c *GoiterClient) AddGroupMember(groupID uint, memberType string, memberID uint) (permission map[string]interface{}, err error) {
	cliResp := &ClientResponse{}
	if cliResp, err = c.makeRequest(&ClientRequest{
		Method: "POST",
		URL:    fmt.Sprintf("/groups/%d/members", groupID),
		Body: map[string]interface{}{
			"group_id":    groupID,
			"member_id":   memberID,
			"member_type": memberType,
		},
	}); err != nil {
		return nil, err
	}
	permission = cliResp.RespBody
	return
}

// RemoveGroupMember removes a member from a group
func (c *GoiterClient) RemoveGroupMember(groupID, memberID uint) (err error) {
	if _, err = c.makeRequest(&ClientRequest{
		Method: "DELETE",
		URL:    fmt.Sprintf("/groups/%d/members", groupID),
		Body: map[string]interface{}{
			"group_id":    groupID,
			"member_id":   memberID,
			"member_type": "User",
		},
	}); err != nil {
		return err
	}
	return
}

// GetGroupAncestors fetches all the ancestors of the group
func (c *GoiterClient) GetGroupAncestors(groupID uint) (groups []map[string]interface{}, err error) {
	cliResp := &ClientResponse{}
	if cliResp, err = c.makeRequest(&ClientRequest{
		Method: "GET",
		URL:    fmt.Sprintf("/groups/%d/ancestors", groupID),
		Body:   nil,
	}); err != nil {
		return nil, err
	}
	log.Println(cliResp.RespBody)
	return
}

func (c *GoiterClient) RunGroupSuite() (err error) {
	log.Println("Running Group test suite...")

	createdGroup, err := c.CreateGroup("Test Group", "This is a test group.")
	assert.Nil(c, err, "Failed to create group")
	assert.Equal(c, "Test Group", createdGroup["name"], "group name should be 'Test Group'")

	groupID := uint(createdGroup["id"].(float64))

	// Get the group
	group, err := c.GetGroup(groupID)
	assert.Nil(c, err, "Failed to get group")
	assert.Equal(c, "Test Group", group["name"], "group name should be 'Test group'")

	// Add a member to the group
	userHash := make(map[string]interface{})
	// Login also creates the user. So just logging in and getting the user ID
	userHash, _ = c.Login("random_group_member_1@example.com")

	c.Login(c.users["root"].Email)

	_, err = c.AddGroupMember(groupID, "User", uint(userHash["id"].(float64)))
	assert.Nil(c, err, "Failed to add member to group")

	// List groups
	groups, err := c.GetGroups()
	assert.Nil(c, err, "Failed to list groups")
	assert.NotEmpty(c, groups, "Groups list should not be empty")

	// Get group ancestors (should be empty)
	ancestors, err := c.GetGroupAncestors(groupID)
	assert.Nil(c, err, "Failed to get group ancestors")
	assert.Empty(c, ancestors, "Group ancestors should not be empty")

	// Get group ancestors when ancestors are present
	parentGroup, err := c.CreateGroup("Parent Group", "This is a parent group.")
	grandParentGroup, err := c.CreateGroup("Grandparent Group", "This is a grandparent group.")
	_, err = c.AddGroupMember(uint(parentGroup["id"].(float64)), "Group", groupID)
	_, err = c.AddGroupMember(uint(grandParentGroup["id"].(float64)), "Group", uint(parentGroup["id"].(float64)))

	ancestors, err = c.GetGroupAncestors(groupID)
	assert.Nil(c, err, "Failed to get group ancestors")
	assert.Equal(c, 2, len(ancestors), "Group should have 2 ancestors")

	// Remove a member from the group
	//err = c.RemoveGroupMember(groupID, uint(userHash["id"].(float64)))
	//assert.Nil(c, err, "Failed to remove member from group")

	// Delete the group
	//err = c.DeleteGroup(groupID)
	//assert.Nil(c, err, "Failed to delete group")

	return
}
