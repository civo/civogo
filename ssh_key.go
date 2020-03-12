package civogo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
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

// ListSSHKeys finds the default private SSH key for an account
func (c *Client) ListSSHKeys(page int, perPage int) (PaginatedSSHKeyList, error) {
	url := "/v2/sshkeys"
	if page != 0 && perPage != 0 {
		url = url + fmt.Sprintf("?page=%d&per_page=%d", page, perPage)
	}

	PaginatedSSHKeys := PaginatedSSHKeyList{}

	resp, err := c.SendGetRequest(url)
	if err != nil {
		return PaginatedSSHKeys, err
	}

	err = json.NewDecoder(bytes.NewReader(resp)).Decode(&PaginatedSSHKeys)
	if err != nil {
		return PaginatedSSHKeys, err
	}

	return PaginatedSSHKeys, nil
}

// GetDefaultSSHKey finds the default private SSH key for an account
func (c *Client) GetDefaultSSHKey() (*SSHKey, error) {
	keys, err := c.ListSSHKeys(1, 9999999)
	if err != nil {
		return nil, err
	}

	for _, key := range keys.Items {
		if key.Default {
			return &key, nil
		}
	}

	return nil, fmt.Errorf("default SSH key not found")
}

// FindSSHKey finds a SSH key by either part of the ID or part of the name
func (c *Client) FindSSHKey(search string) (*SSHKey, error) {
	keys, err := c.ListSSHKeys(1, 99999999)
	if err != nil {
		return nil, err
	}

	found := -1

	for i, key := range keys.Items {
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

	return &keys.Items[found], nil
}
