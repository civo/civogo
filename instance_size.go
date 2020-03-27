package civogo

import (
	"bytes"
	"encoding/json"
	"errors"
)

// InstanceSize represents an available size for instances to launch
type InstanceSize struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	NiceName      string `json:"nice_name"`
	CPUCores      int    `json:"cpu_cores"`
	RAMMegabytes  int    `json:"ram_mb"`
	DiskGigabytes int    `json:"disk_gb"`
	Description   string `json:"description"`
	Selectable    bool   `json:"selectable"`
}

// ListInstanceSizes returns all available sizes of instances
func (c *Client) ListInstanceSizes() ([]InstanceSize, error) {
	resp, err := c.SendGetRequest("/v2/sizes")
	if err != nil {
		return nil, err
	}

	sizes := make([]InstanceSize, 0)
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&sizes); err != nil {
		return nil, err
	}

	return sizes, nil
}

// GetInstanceSizeByName finds the instance size by the name
func (c *Client) GetInstanceSizeByName(name string) (*InstanceSize, error) {
	resp, err := c.ListInstanceSizes()
	if err != nil {
		return nil, err
	}

	for _, size := range resp {
		if size.Name == name {
			return &size, nil
		}
	}

	return nil, errors.New("instances size not found")
}
