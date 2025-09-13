package testsuite

import (
	"log"
)

func (c *GoiterClient) ListGroups() (respBody map[string]interface{}, err error) {
	cliResp := &ClientResponse{}
	if cliResp, err = c.makeRequest(&ClientRequest{
		Method: "GET",
		URL:    "/groups",
		Body:   nil,
	}); err != nil {
		return nil, err
	}
	respBody = cliResp.RespBody
	return
}

func (c *GoiterClient) RunGroupSuite() (err error) {
	log.Println("Running group suite...")

	return
}
