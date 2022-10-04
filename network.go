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
	Name    string `json:"name,omitempty"`
	Default bool   `json:"default,omitempty"`
	CIDR    string `json:"cidr,omitempty"`
	Label   string `json:"label,omitempty"`
	Status  string `json:"status,omitempty"`
}

type networkConfig struct {
	Label   string     `json:"label"`
	Region  string     `json:"region"`
	Default bool       `json:"default"`
	IPv4    ipv4Config `json:"ipv4" schema:"ipv4"`
	IPv6    ipv6Config `json:"ipv6" schema:"ipv6"`
}

type ipv4Config struct {
	Enabled     *bool    `json:"enabled" schema:"enabled"` // enabled by default
	Nameservers []string `json:"nameservers" schema:"nameservers"`
	Subnet      string   `json:"subnets" schema:"subnets"` // default is: 192.168.1.0/24
}

type ipv6Config struct {
	Enabled     *bool    `json:"enabled" schema:"enabled"` // enabled by default
	Nameservers []string `json:"nameservers" schema:"nameservers"`
}

// NetworkResult represents the result from a network create/update call
type NetworkResult struct {
	ID     string `json:"id"`
	Label  string `json:"label"`
	Result string `json:"result"`
}

// Subnet represents an ipv6 allocation for an existing ipv6-enabled network
type Subnet struct {
	ID        string `json:"id"`
	Name      string `json:"name,omitempty"`
	NetworkID string `json:"networkId"`
	Label     string `json:"label,omitempty"`
	Status    string `json:"status,omitempty"`
}

// GetDefaultNetwork finds the default private network for an account
func (c *Client) GetDefaultNetwork() (*Network, error) {
	resp, err := c.SendGetRequest("/v2/networks")
	if err != nil {
		return nil, decodeError(err)
	}

	networks := make([]Network, 0)
	json.NewDecoder(bytes.NewReader(resp)).Decode(&networks)
	for _, network := range networks {
		if network.Default {
			return &network, nil
		}
	}

	return nil, errors.New("no default network found")
}

// GetNetwork gets a network with ID
func (c *Client) GetNetwork(id string) (*Network, error) {
	resp, err := c.SendGetRequest("/v2/networks/" + id)
	if err != nil {
		return nil, decodeError(err)
	}

	network := Network{}
	err = json.NewDecoder(bytes.NewReader(resp)).Decode(&network)
	return &network, err
}

// NewNetwork creates a new private network
func (c *Client) NewNetwork(label string) (*NetworkResult, error) {
	nc := networkConfig{Label: label, Region: c.Region}
	body, err := c.SendPostRequest("/v2/networks", nc)
	if err != nil {
		return nil, decodeError(err)
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
		return nil, decodeError(err)
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
		return nil, decodeError(err)
	}

	exactMatch := false
	partialMatchesCount := 0
	result := Network{}

	for _, value := range networks {
		if value.Name == search || value.ID == search || value.Label == search {
			exactMatch = true
			result = value
		} else if strings.Contains(value.Name, search) || strings.Contains(value.ID, search) || strings.Contains(value.Label, search) {
			if !exactMatch {
				result = value
				partialMatchesCount++
			}
		}
	}

	if exactMatch || partialMatchesCount == 1 {
		return &result, nil
	} else if partialMatchesCount > 1 {
		err := fmt.Errorf("unable to find %s because there were multiple matches", search)
		return nil, MultipleMatchesError.wrap(err)
	} else {
		err := fmt.Errorf("unable to find %s, zero matches", search)
		return nil, ZeroMatchesError.wrap(err)
	}
}

// RenameNetwork renames an existing private network
func (c *Client) RenameNetwork(label, id string) (*NetworkResult, error) {
	nc := networkConfig{Label: label, Region: c.Region}
	body, err := c.SendPutRequest("/v2/networks/"+id, nc)
	if err != nil {
		return nil, decodeError(err)
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
		return nil, decodeError(err)
	}

	return c.DecodeSimpleResponse(resp)
}

// ListSubnets list all ipv6 subnets linked to a given ipv6-enabled network
func (c *Client) ListSubnets(networkID string) ([]Subnet, error) {
	resp, err := c.SendGetRequest(fmt.Sprintf("/v2/networks/%s/subnets", networkID))
	if err != nil {
		return nil, decodeError(err)
	}

	subnets := make([]Subnet, 0)
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&subnets); err != nil {
		return nil, err
	}

	return subnets, nil
}

// GetSubnetByID fetch the ipv6 subnet details linked to a given ipv6-enabled network given its ID
func (c *Client) GetSubnetByID(networkID string, subnetID string) (Subnet, error) {
	resp, err := c.SendGetRequest(fmt.Sprintf("/v2/networks/%s/subnets/%s", networkID, subnetID))
	if err != nil {
		return Subnet{}, decodeError(err)
	}

	subnet := Subnet{}
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&subnet); err != nil {
		return Subnet{}, err
	}

	return subnet, nil
}

// GetSubnetByID fetch the ipv6 subnet details linked to a given ipv6-enabled network given its name
func (c *Client) GetSubnetByName(networkID string, subnetName string) (Subnet, error) {
	resp, err := c.SendGetRequest(fmt.Sprintf("/v2/networks/%s/subnets?name=%s", networkID, subnetName))
	if err != nil {
		return Subnet{}, decodeError(err)
	}

	subnet := Subnet{}
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&subnet); err != nil {
		return Subnet{}, err
	}

	return subnet, nil
}

// DeleteSubnet deletes a ipv6 subnet from an ipv6-enabled network
func (c *Client) DeleteSubnet(networkID string, subnetID string) (*SimpleResponse, error) {
	resp, err := c.SendDeleteRequest(fmt.Sprintf("/v2/networks/%s/subnets/%s", networkID, subnetID))
	if err != nil {
		return nil, decodeError(err)
	}

	return c.DecodeSimpleResponse(resp)
}
