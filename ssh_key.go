package civogo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

// SSHKey represents an SSH public key, uploaded to access instances
type SSHKey struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Fingerprint string `json:"fingerprint"`
}

// ListSSHKeys list all SSH key for an account
func (c *Client) ListSSHKeys() ([]SSHKey, error) {
	resp, err := c.SendGetRequest("/v2/sshkeys")
	if err != nil {
		return nil, err
	}

	sshKeys := make([]SSHKey, 0)
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&sshKeys); err != nil {
		return nil, err
	}

	return sshKeys, nil
}

// NewSSHKey creates a new SSH key record
func (c *Client) NewSSHKey(name string, publicKey string) (*SimpleResponse, error) {
	resp, err := c.SendPostRequest("/v2/sshkeys", map[string]string{
		"name":       name,
		"public_key": publicKey,
	})
	if err != nil {
		return nil, err
	}

	return c.DecodeSimpleResponse(resp)
}

// UpdateSSHKey update a SSH key record
func (c *Client) UpdateSSHKey(name string, sshKeyID string) (*SSHKey, error) {
	resp, err := c.SendPutRequest(fmt.Sprintf("/v2/sshkeys/%s", sshKeyID), map[string]string{
		"name": name,
	})
	if err != nil {
		return nil, err
	}

	result := &SSHKey{}
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}

// FindSSHKey finds an SSH key by either part of the ID or part of the name
func (c *Client) FindSSHKey(search string) (*SSHKey, error) {
	keys, err := c.ListSSHKeys()
	if err != nil {
		return nil, err
	}

	found := -1

	for i, key := range keys {
		if strings.Contains(key.ID, search) || strings.Contains(key.Name, search) {
			if found != -1 {
				return nil, fmt.Errorf("unable to find %s because there were multiple matches", search)
			}
			found = i
		}
	}

	if found == -1 {
		return nil, fmt.Errorf("unable to find %s, zero matches", search)
	}

	return &keys[found], nil
}

// DeleteSSHKey deletes an SSH key
func (c *Client) DeleteSSHKey(id string) (*SimpleResponse, error) {
	resp, err := c.SendDeleteRequest(fmt.Sprintf("/v2/sshkeys/%s", id))
	if err != nil {
		return nil, err
	}

	return c.DecodeSimpleResponse(resp)
}
