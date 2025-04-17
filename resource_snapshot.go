package civogo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

// ResourceSnapshot represents a snapshot of any resource type
type ResourceSnapshot struct {
	ID           string                `json:"id"`
	Name         string                `json:"name"`
	Description  string                `json:"description"`
	ResourceType string                `json:"resource_type"`
	CreatedAt    time.Time             `json:"created_at"`
	Instance     *InstanceSnapshotInfo `json:"instance,omitempty"`
}

// InstanceSnapshotInfo represents instance-specific snapshot details
type InstanceSnapshotInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      struct {
		State string `json:"state"`
	} `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// UpdateResourceSnapshotRequest represents the request to update a resource snapshot
type UpdateResourceSnapshotRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// RestoreInstanceSnapshotRequest represents the request to restore an instance snapshot
type RestoreInstanceSnapshotRequest struct {
	Description       string `json:"description,omitempty"`
	Hostname          string `json:"hostname,omitempty"`
	PrivateIPv4       string `json:"private_ipv4,omitempty"`
	IncludeVolumes    bool   `json:"include_volumes,omitempty"`
	OverwriteExisting bool   `json:"overwrite_existing,omitempty"`
}

// RestoreResourceSnapshotRequest represents the request to restore a resource snapshot
type RestoreResourceSnapshotRequest struct {
	Instance *RestoreInstanceSnapshotRequest `json:"instance,omitempty"`
}

// ListResourceSnapshots returns all resource snapshots
func (c *Client) ListResourceSnapshots() ([]ResourceSnapshot, error) {
	resp, err := c.SendGetRequest("/v2/resourcesnapshots")
	if err != nil {
		return nil, decodeError(err)
	}

	var snapshots []ResourceSnapshot
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&snapshots); err != nil {
		return nil, err
	}

	return snapshots, nil
}

// GetResourceSnapshot retrieves a specific resource snapshot by ID
func (c *Client) GetResourceSnapshot(id string) (*ResourceSnapshot, error) {
	resp, err := c.SendGetRequest(fmt.Sprintf("/v2/resourcesnapshots/%s", id))
	if err != nil {
		return nil, decodeError(err)
	}

	var snapshot ResourceSnapshot
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&snapshot); err != nil {
		return nil, err
	}

	return &snapshot, nil
}

// UpdateResourceSnapshot updates a resource snapshot
func (c *Client) UpdateResourceSnapshot(id string, req *UpdateResourceSnapshotRequest) (*ResourceSnapshot, error) {
	body, err := c.SendPutRequest(fmt.Sprintf("/v2/resourcesnapshots/%s", id), req)
	if err != nil {
		return nil, decodeError(err)
	}

	var snapshot ResourceSnapshot
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&snapshot); err != nil {
		return nil, err
	}

	return &snapshot, nil
}

// DeleteResourceSnapshot deletes a resource snapshot
func (c *Client) DeleteResourceSnapshot(id string) (*SimpleResponse, error) {
	resp, err := c.SendDeleteRequest(fmt.Sprintf("/v2/resourcesnapshots/%s", id))
	if err != nil {
		return nil, decodeError(err)
	}

	return c.DecodeSimpleResponse(resp)
}

// RestoreResourceSnapshot restores a resource from a snapshot
func (c *Client) RestoreResourceSnapshot(id string, req *RestoreResourceSnapshotRequest) (*ResourceSnapshot, error) {
	body, err := c.SendPostRequest(fmt.Sprintf("/v2/resourcesnapshots/%s/restore", id), req)
	if err != nil {
		return nil, decodeError(err)
	}

	var snapshot ResourceSnapshot
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&snapshot); err != nil {
		return nil, err
	}

	return &snapshot, nil
}
