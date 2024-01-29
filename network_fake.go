package civogo

import (
	"fmt"
	"strings"
)

// GetDefaultNetwork implemented in a fake way for automated tests
func (c *FakeClient) GetDefaultNetwork() (*Network, error) {
	for _, network := range c.Networks {
		if network.Default {
			return &network, nil
		}
	}

	err := fmt.Errorf("unable to find default network, zero matches")
	return nil, ZeroMatchesError.wrap(err)
}

// NewNetwork implemented in a fake way for automated tests
func (c *FakeClient) NewNetwork(label string) (*NetworkResult, error) {
	network := Network{
		ID:   c.generateID(),
		Name: label,
	}
	c.Networks = append(c.Networks, network)

	return &NetworkResult{
		ID:     network.ID,
		Label:  network.Name,
		Result: "success",
	}, nil

}

// ListNetworks implemented in a fake way for automated tests
func (c *FakeClient) ListNetworks() ([]Network, error) {
	return c.Networks, nil
}

// FindNetwork implemented in a fake way for automated tests
func (c *FakeClient) FindNetwork(search string) (*Network, error) {
	for _, network := range c.Networks {
		if strings.Contains(network.Name, search) {
			return &network, nil
		}
	}

	err := fmt.Errorf("unable to find default network, zero matches")
	return nil, ZeroMatchesError.wrap(err)
}

// RenameNetwork implemented in a fake way for automated tests
func (c *FakeClient) RenameNetwork(label, id string) (*NetworkResult, error) {
	for i, network := range c.Networks {
		if network.ID == id {
			c.Networks[i].Label = label
			return &NetworkResult{
				ID:     network.ID,
				Label:  network.Label,
				Result: "success",
			}, nil
		}
	}

	err := fmt.Errorf("unable to find default network, zero matches")
	return nil, ZeroMatchesError.wrap(err)
}

// DeleteNetwork implemented in a fake way for automated tests
func (c *FakeClient) DeleteNetwork(id string) (*SimpleResponse, error) {
	for i, network := range c.Networks {
		if network.ID == id {
			c.Networks[len(c.Networks)-1], c.Networks[i] = c.Networks[i], c.Networks[len(c.Networks)-1]
			c.Networks = c.Networks[:len(c.Networks)-1]
			return &SimpleResponse{Result: "success"}, nil
		}
	}

	return &SimpleResponse{Result: "failed"}, nil
}
