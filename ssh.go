package civogo

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

const sshBasePath = "/v2/sshkeys"

// SSHKeyService is an interface for interfacing with the ssh
type SSHKeyService interface {
	List(ctx context.Context) ([]SSHKey, *Metadata, error)
	GetByID(ctx context.Context, id string) (*SSHKey, *Metadata, error)
	Find(ctx context.Context, value string) (*SSHKey, *Metadata, error)
	Create(ctx context.Context, createRequest *SSHKeyCreateRequest) (*SimpleResponse, *Metadata, error)
	Update(ctx context.Context, sshID string, updateRequest *SSHKeyUpdateRequest) (*SSHKey, *Metadata, error)
	Delete(ctx context.Context, sshID string) (*SimpleResponse, *Metadata, error)
}

// SSHKeyServiceOp Service used for communicating with the API
type SSHKeyServiceOp struct {
	client  *Client
}

type SSHKeyGetter interface {
	SSHKey() SSHKeyService
}

// newSSHKey returns a SSHKey
func newSSHKey(c *Client) *SSHKeyServiceOp {
	return &SSHKeyServiceOp{
		client:  c,
	}
}

// SSHKey represents an SSH public key, uploaded to access instances
type SSHKey struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Fingerprint string `json:"fingerprint"`
}

// SSHKeyCreateRequest represents a request to create a new key.
type SSHKeyCreateRequest struct {
	Name      string `json:"name"`
	PublicKey string `json:"public_key"`
}

// SSHKeyUpdateRequest represents a request to update an existing key.
type SSHKeyUpdateRequest struct {
	Name string `json:"name"`
}

// List list all SSH key for an account
func (c *SSHKeyServiceOp) List(ctx context.Context) ([]SSHKey, *Metadata, error) {
	path := sshBasePath
	req, err := c.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new([]SSHKey)
	resp, err := c.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return *root, resp, err
}

// GetByID get an SSH key by ID
func (c *SSHKeyServiceOp) GetByID(ctx context.Context, sshID string) (*SSHKey, *Metadata, error) {
	if sshID == "" {
		return nil, nil, errors.New("sshID cannot be empty")
	}

	path := fmt.Sprintf("%s/%s", sshBasePath, sshID)
	req, err := c.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(SSHKey)
	resp, err := c.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Find get an SSH  key by name or ID
func (c *SSHKeyServiceOp) Find(ctx context.Context, value string) (*SSHKey, *Metadata, error) {
	if value == "" {
		return nil, nil, errors.New("the search term cannot be empty")
	}

	allSSHKeys, meta, err := c.List(ctx)
	if err != nil {
		return nil, nil, err
	}

	for _, sshKey := range allSSHKeys {
		if sshKey.ID == value || sshKey.Name == value {
			return &sshKey, meta, nil
		}
	}

	return nil, nil, errors.New("no SSH key found")
}

// Create create a new SSH key
func (c *SSHKeyServiceOp) Create(ctx context.Context, createRequest *SSHKeyCreateRequest) (*SimpleResponse, *Metadata, error) {
	if createRequest == nil {
		return nil, nil, errors.New("createRequest is required")
	}

	path := sshBasePath
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

// Update update an existing SSH key
func (c *SSHKeyServiceOp) Update(ctx context.Context, sshID string, updateRequest *SSHKeyUpdateRequest) (*SSHKey, *Metadata, error) {
	if sshID == "" {
		return nil, nil, errors.New("sshID cannot be empty")
	}
	if updateRequest == nil {
		return nil, nil, errors.New("updateRequest is required")
	}

	path := fmt.Sprintf("%s/%s", sshBasePath, sshID)
	req, err := c.client.NewRequest(ctx, http.MethodPut, path, updateRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(SSHKey)
	resp, err := c.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Delete delete an existing SSH key
func (c *SSHKeyServiceOp) Delete(ctx context.Context, sshID string) (*SimpleResponse, *Metadata, error) {
	if sshID == "" {
		return nil, nil, errors.New("sshID cannot be empty")
	}

	path := fmt.Sprintf("%s/%s", sshBasePath, sshID)
	req, err := c.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(SimpleResponse)
	resp, err := c.client.Do(ctx, req, root)
	if err != nil {
		return nil, nil, err
	}

	return root, resp, err
}
