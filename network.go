package civogo

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

const networkBasePath = "/v2/networks"

// NetworkService is an interface for interfacing with the network
type NetworkService interface {
	List(ctx context.Context) ([]Network, *Metadata, error)
	GetDefault(ctx context.Context) (*Network, *Metadata, error)
	GetByID(ctx context.Context, networkID string) (*Network, *Metadata, error)
	Find(ctx context.Context, value string) (*Network, *Metadata, error)
	Create(ctx context.Context, createRequest *NetworkCreateRequest) (*SimpleResponse, *Metadata, error)
	Update(ctx context.Context, networkID string, updateRequest *NetworkUpdateRequest) (*SimpleResponse, *Metadata, error)
	Delete(ctx context.Context, networkID string) (*SimpleResponse, *Metadata, error)
}

// SSH Service used for communicating with the API
type NetworkServiceOp struct {
	client *Client
}

var _ NetworkService = &NetworkServiceOp{}

// Network represents a private network for instances to connect to
type Network struct {
	ID      string `json:"id"`
	Name    string `json:"name,omitempty"`
	Default bool   `json:"default,omitempty"`
	CIDR    string `json:"cidr,omitempty"`
	Label   string `json:"label,omitempty"`
	Status  string `json:"status,omitempty"`
}

type NetworkCreateRequest struct {
	Label  string `json:"label"`
	Region string `json:"region"`
}

type NetworkUpdateRequest struct {
	Label string `json:"label"`
	Region string `json:"region"`
}

// List list all network for an account
func (c *NetworkServiceOp) List(ctx context.Context) ([]Network, *Metadata, error) {
	path := networkBasePath
	req, err := c.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new([]Network)
	resp, err := c.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return *root, resp, err
}

// GetDefault returns the default network for an account
func (c *NetworkServiceOp) GetDefault(ctx context.Context) (*Network, *Metadata, error) {
	networks, metadata, err := c.List(ctx)
	if err != nil {
		return nil, metadata, err
	}

	for _, network := range networks {
		if network.Default {
			return &network, metadata, nil
		}
	}

	return nil, metadata, nil
}

// GetByID returns a network by ID
func (c *NetworkServiceOp) GetByID(ctx context.Context, networkID string) (*Network, *Metadata, error) {
	if networkID == "" {
		return nil, nil, errors.New("networkID cannot be empty")
	}

	path := fmt.Sprintf("%s/%s", networkBasePath, networkID)
	req, err := c.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Network)
	resp, err := c.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Find returns a network by name or ID
func (c *NetworkServiceOp) Find(ctx context.Context, value string) (*Network, *Metadata, error) {
	if value == "" {
		return nil, nil, errors.New("value cannot be empty")
	}

	networks, metadata, err := c.List(ctx)
	if err != nil {
		return nil, metadata, err
	}

	for _, network := range networks {
		if network.ID == value || network.Name == value || network.Label == value {
			return &network, metadata, nil
		}
	}

	return nil, metadata, nil
}

// Create creates a new network
func (c *NetworkServiceOp) Create(ctx context.Context, createRequest *NetworkCreateRequest) (*SimpleResponse, *Metadata, error) {
	if createRequest == nil {
		return nil, nil, errors.New("createRequest is required")
	}

	path := networkBasePath
	req, err := c.client.NewRequest(ctx, http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(SimpleResponse)
	resp, err := c.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Update updates a network
func (c *NetworkServiceOp) Update(ctx context.Context, networkID string, updateRequest *NetworkUpdateRequest) (*SimpleResponse, *Metadata, error) {
	if networkID == "" {
		return nil, nil, errors.New("networkID is required")
	}

	if updateRequest == nil {
		return nil, nil, errors.New("updateRequest is required")
	}

	path := fmt.Sprintf("%s/%s", networkBasePath, networkID)
	req, err := c.client.NewRequest(ctx, http.MethodPut, path, updateRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(SimpleResponse)
	resp, err := c.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Delete deletes a network
func (c *NetworkServiceOp) Delete(ctx context.Context, networkID string) (*SimpleResponse, *Metadata, error) {
	if networkID == "" {
		return nil, nil, errors.New("networkID is required")
	}

	path := fmt.Sprintf("%s/%s", networkBasePath, networkID)
	req, err := c.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(SimpleResponse)
	resp, err := c.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}
