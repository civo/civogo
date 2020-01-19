package civogo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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
	Label  string `json:"label"`
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
	body, err := c.SendPostRequest("/v2/networks", r)
	if err != nil {
		return nil, err
	}

	var result = &NetworkResult{}
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}

// ListNetwork list all networks
func (c *Client) ListNetwork() ([]Network, error) {
	resp, err := c.SendGetRequest("/v2/networks")
	if err != nil {
		return nil, err
	}

	networks := make([]Network, 0)
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&networks); err != nil {
		return nil, err
	}

	return networks, nil
}

// DeleteNetwork deletes an network
func (c *Client) DeleteNetwork(id string) (*SimpleResponse, error) {
	resp, err := c.SendDeleteRequest(fmt.Sprintf("/v2/networks/%s", id))
	if err != nil {
		return nil, err
	}

	return c.DecodeSimpleResponse(resp)
}
