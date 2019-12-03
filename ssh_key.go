package civogo

import (
	"bytes"
	"encoding/json"
	"errors"
)

// SSHKey represents an SSH public key, uploaded to access instances
type SSHKey struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Region  string `json:"region"`
	Default bool   `json:"default"`
	CIDR    string `json:"cidr"`
	Label   string `json:"label"`
}

// PaginatedSSHKeyList returns a paginated list of SSHKey objects
type PaginatedSSHKeyList struct {
	Page    int      `json:"page"`
	PerPage int      `json:"per_page"`
	Pages   int      `json:"pages"`
	Items   []SSHKey `json:"items"`
}

// GetDefaultSSHKey finds the default private SSH key for an account
func (c *Client) GetDefaultSSHKey() (*SSHKey, error) {
	resp, err := c.SendGetRequest("/v2/sshkeys")
	if err != nil {
		return nil, err
	}

	PaginatedSSHKeys := PaginatedSSHKeyList{}
	err = json.NewDecoder(bytes.NewReader(resp)).Decode(&PaginatedSSHKeys)
	if len(PaginatedSSHKeys.Items) > 0 {
		return &PaginatedSSHKeys.Items[0], nil
	}

	return nil, errors.New("No SSH key found")
}
