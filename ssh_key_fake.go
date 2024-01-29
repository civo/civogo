package civogo

import (
	"fmt"
	"strings"
)

// ListSSHKeys implemented in a fake way for automated tests
func (c *FakeClient) ListSSHKeys() ([]SSHKey, error) {
	return c.SSHKeys, nil
}

// NewSSHKey implemented in a fake way for automated tests
func (c *FakeClient) NewSSHKey(name string, publicKey string) (*SimpleResponse, error) {
	sshKey := SSHKey{
		Name:        name,
		Fingerprint: publicKey, // This is weird, but we're just storing a value
	}
	c.SSHKeys = append(c.SSHKeys, sshKey)
	return &SimpleResponse{Result: "success"}, nil
}

// UpdateSSHKey implemented in a fake way for automated tests
func (c *FakeClient) UpdateSSHKey(name string, sshKeyID string) (*SSHKey, error) {
	for i, sshKey := range c.SSHKeys {
		if sshKey.ID == sshKeyID {
			c.SSHKeys[i].Name = name
			return &sshKey, nil
		}
	}

	err := fmt.Errorf("unable to find SSH key %s, zero matches", sshKeyID)
	return nil, ZeroMatchesError.wrap(err)
}

// FindSSHKey implemented in a fake way for automated tests
func (c *FakeClient) FindSSHKey(search string) (*SSHKey, error) {
	for _, sshKey := range c.SSHKeys {
		if strings.Contains(sshKey.Name, search) {
			return &sshKey, nil
		}
	}

	err := fmt.Errorf("unable to find SSH key %s, zero matches", search)
	return nil, ZeroMatchesError.wrap(err)
}

// DeleteSSHKey implemented in a fake way for automated tests
func (c *FakeClient) DeleteSSHKey(id string) (*SimpleResponse, error) {
	for i, sshKey := range c.SSHKeys {
		if sshKey.ID == id {
			c.SSHKeys[len(c.SSHKeys)-1], c.SSHKeys[i] = c.SSHKeys[i], c.SSHKeys[len(c.SSHKeys)-1]
			c.SSHKeys = c.SSHKeys[:len(c.SSHKeys)-1]
			return &SimpleResponse{Result: "success"}, nil
		}
	}

	return &SimpleResponse{Result: "failed"}, nil
}
