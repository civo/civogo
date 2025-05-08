package civogo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

// InstanceSnapshot represents a snapshot of an instance
type InstanceSnapshot struct {
	ID              string                 `json:"id"`
	Name            string                 `json:"name"`
	Description     string                 `json:"description,omitempty"`
	IncludedVolumes []string               `json:"included_volumes"`
	Status          InstanceSnapshotStatus `json:"status"`
	CreatedAt       time.Time              `json:"created_at"`
}

// InstanceSnapshotStatus represents the status of an instance snapshot
type InstanceSnapshotStatus struct {
	State   string                   `json:"state"`
	Message string                   `json:"message"`
	Volumes []InstanceSnapshotVolume `json:"volumes,omitempty"`
}

// InstanceSnapshotVolume represents the status of a volume in a snapshot
type InstanceSnapshotVolume struct {
	ID    string `json:"id"`
	State string `json:"state"`
}

// CreateInstanceSnapshotConfig represents the configuration for creating a new instance snapshot
type CreateInstanceSnapshotConfig struct {
	Name           string `json:"name"`
	Description    string `json:"description,omitempty"`
	IncludeVolumes bool   `json:"include_volumes"`
}

// UpdateInstanceSnapshotConfig represents the configuration for updating an instance snapshot
type UpdateInstanceSnapshotConfig struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// CreateInstanceSnapshot creates a new snapshot of an instance
func (c *Client) CreateInstanceSnapshot(instanceID string, config *CreateInstanceSnapshotConfig) (*InstanceSnapshot, error) {
	body, err := c.SendPostRequest(fmt.Sprintf("/v2/instances/%s/snapshots", instanceID), config)
	if err != nil {
		return nil, decodeError(err)
	}

	var snapshot InstanceSnapshot
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&snapshot); err != nil {
		return nil, err
	}

	return &snapshot, nil
}

// GetInstanceSnapshot retrieves snapshot of a specific instance
func (c *Client) GetInstanceSnapshot(instanceID, snapshotID string) (*InstanceSnapshot, error) {
	resp, err := c.SendGetRequest(fmt.Sprintf("/v2/instances/%s/snapshots/%s", instanceID, snapshotID))
	if err != nil {
		return nil, decodeError(err)
	}

	var snapshot InstanceSnapshot
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&snapshot); err != nil {
		return nil, err
	}

	return &snapshot, nil
}

// ListInstanceSnapshots retrieves all snapshots for a specific instance
func (c *Client) ListInstanceSnapshots(instanceID string) ([]InstanceSnapshot, error) {
	resp, err := c.SendGetRequest(fmt.Sprintf("/v2/instances/%s/snapshots", instanceID))
	if err != nil {
		return nil, decodeError(err)
	}

	var snapshots []InstanceSnapshot
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&snapshots); err != nil {
		return nil, err
	}

	return snapshots, nil
}

// UpdateInstanceSnapshot updates an existing instance snapshot
func (c *Client) UpdateInstanceSnapshot(instanceID, snapshotID string, config *UpdateInstanceSnapshotConfig) (*InstanceSnapshot, error) {
	body, err := c.SendPutRequest(fmt.Sprintf("/v2/instances/%s/snapshots/%s", instanceID, snapshotID), config)
	if err != nil {
		return nil, decodeError(err)
	}

	var snapshot InstanceSnapshot
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&snapshot); err != nil {
		return nil, err
	}

	return &snapshot, nil
}

// DeleteInstanceSnapshot deletes an instance snapshot by its ID
func (c *Client) DeleteInstanceSnapshot(instanceID, snapshotID string) (*SimpleResponse, error) {
	resp, err := c.SendDeleteRequest(fmt.Sprintf("/v2/instances/%s/snapshots/%s", instanceID, snapshotID))
	if err != nil {
		return nil, decodeError(err)
	}

	var result SimpleResponse
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
