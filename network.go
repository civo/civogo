package civogo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
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

type networkConfig struct {
	Label string `json:"label"`
}

// NetworkResult represents the result from a network create/update call
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

// NewNetwork creates a new private network
func (c *Client) NewNetwork(label string) (*NetworkResult, error) {
	nc := networkConfig{Label: label}
	body, err := c.SendPostRequest("/v2/networks", nc)
	if err != nil {
		return nil, err
	}

	var result = &NetworkResult{}
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}

// ListNetworks list all private networks
func (c *Client) ListNetworks() ([]Network, error) {
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

// FindNetwork finds a network by either part of the ID or part of the name
func (c *Client) FindNetwork(search string) (*Network, error) {
	networks, err := c.ListNetworks()
	if err != nil {
		return nil, err
	}

	found := -1

	for i, network := range networks {
		if strings.Contains(network.ID, search) || strings.Contains(network.Name, search) {
			if found != -1 {
				return nil, fmt.Errorf("unable to find %s because there were multiple matches", search)
			}
			found = i
		}
	}

	if found == -1 {
		return nil, fmt.Errorf("unable to find %s, zero matches", search)
	}

	return &networks[found], nil
}

// RenameNetwork renames an existing private network
func (c *Client) RenameNetwork(label, id string) (*NetworkResult, error) {
	nc := networkConfig{Label: label}
	body, err := c.SendPutRequest("/v2/networks/"+id, nc)
	if err != nil {
		return nil, err
	}

	var result = &NetworkResult{}
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}

// DeleteNetwork deletes a private network
func (c *Client) DeleteNetwork(id string) (*SimpleResponse, error) {
	resp, err := c.SendDeleteRequest(fmt.Sprintf("/v2/networks/%s", id))
	if err != nil {
		return nil, err
	}

	return c.DecodeSimpleResponse(resp)
}
