package civogo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

type Snapshot struct {
	ID          string    `json:"id"`
	InstanceID  string    `json:"instance_id"`
	Hostname    string    `json:"hostname"`
	Template    string    `json:"template_id"`
	Region      string    `json:"region"`
	Name        string    `json:"name"`
	Safe        int       `json:"safe"`
	SizeGB      int       `json:"size_gb"`
	State       string    `json:"state"`
	Cron        string    `json:"cron_timing,omitempty"`
	RequestedAt time.Time `json:"requested_at,omitempty"`
	CompletedAt time.Time `json:"completed_at,omitempty"`
}

type SnapshotsConfig struct {
	InstanceID string `form:"instance_id"`
	Safe       bool   `from:"safe"`
	Cron       string `from:"cron_timing"`
}

// CreateSnapshot create a new or update an old snapshot
func (c *Client) CreateSnapshot(name string, r *SnapshotsConfig) (*Snapshot, error) {
	url := fmt.Sprintf("/v2/snapshots/%s", name)
	body, err := c.SendPutRequest(url, r)
	if err != nil {
		return nil, err
	}

	var n = &Snapshot{}
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(n); err != nil {
		return nil, err
	}

	return n, nil
}
