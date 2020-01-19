package civogo

import (
	"bytes"
	"encoding/json"
	"errors"
)

// Network represents a private network for instances to connect to
type Network struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Region  string `json:"region"`
	Default bool   `json:"default"`
	CIDR    string `json:"cidr"`
	Label   string `json:"label"`
}

type NetworkConfig struct {
	Label string `form:"label"`
}

type NetworkResult struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Result string `json:"result"`
}

// GetDefaultNetwork finds the default private network for an account
func (c *Client) GetDefaultNetwork() (*Network, error) {
	resp, err := c.SendGetRequest("/v2/networks")
	if err != nil {
		return nil, err
	}

	networks := make([]Network, 0)
	err = json.NewDecoder(bytes.NewReader(resp)).Decode(&networks)
	for _, network := range networks {
		if network.Default {
			return &network, nil
		}
	}

	return nil, errors.New("No default network found")
}

// NewVolumes creates a new volume
func (c *Client) NewNetwork(r *NetworkConfig) (*NetworkResult, error) {
	body, err := c.SendPostRequest("/v2/networks/", r)
	if err != nil {
		return nil, err
	}

	var result = &NetworkResult{}
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}
