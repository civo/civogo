package civogo

import (
	"bytes"
	"encoding/json"
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

// ListInstanceSizes returns all availble sizes of instances
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
