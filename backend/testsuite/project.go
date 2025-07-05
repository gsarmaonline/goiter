package testsuite

import "log"

// GetProjects retrieves all projects for the current user
func (c *GoiterClient) GetProjects() (projects map[string]interface{}, err error) {
	if _, projects, err = c.makeRequest("GET", "/projects", nil); err != nil {
		return nil, err
	}
	return
}

func (c *GoiterClient) RunProjectSuite() (err error) {
	respBody := make(map[string]interface{})
	if respBody, err = c.GetProjects(); err != nil {
		return
	}
	log.Println("Fetched projects:", respBody)
	return
}
