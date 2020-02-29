package civogo

import (
	"bytes"
	"encoding/json"
)

// Region represents a geographical/DC region for Civo resources
type Region struct {
	Code    string `json:"code"`
	Name    string `json:"name"`
	Default bool   `json:"default"`
}

// ListRegions returns all load balancers owned by the calling API account
func (c *Client) ListRegions() ([]Region, error) {
	resp, err := c.SendGetRequest("/v2/regions")
	if err != nil {
		return nil, err
	}

	regions := make([]Region, 0)
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&regions); err != nil {
		return nil, err
	}

	return regions, nil
}
