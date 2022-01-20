package civogo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

// Database holds the database information recieved in GET/LIST operations
type Database struct {
	ID                   string `json:"id"`
	Name                 string `json:"name"`
	Replicas             int    `json:"replicas"`
	NumSnapshotsToRetain int    `json:"num_snapshots_to_retain"`
	Size                 string `json:"size"`
	Software             string `json:"software"`
	SoftwareVersion      string `json:"software_version"`
	PublicIP             string `json:"public_ip"`
	NetworkID            string `json:"network_id"`
	FirewallID           string `json:"firewall_id"`
	Snapshots            []struct {
		ID        string    `json:"id"`
		Timestamp time.Time `json:"timestamp"`
	} `json:"snapshots,omitempty"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// CreateDatabaseRequest holds fields required to creates a new database
type CreateDatabaseRequest struct {
	Name                 string `json:"name" validate:"required" schema:"name"`
	Size                 string `json:"size" validate:"required" schema:"size"`
	Software             string `json:"software" validate:"required,oneof=Redis Postgres MySQL" schema:"software"`
	SoftwareVersion      string `json:"software_version" validate:"required" schema:"software_version"`
	NetworkID            string `json:"network_id" validate:"required" schema:"network_id"`
	Replicas             int    `json:"replicas" schema:"replicas"`
	NumSnapshotsToRetain int    `json:"num_snapshots_to_retain" schema:"num_snapshots_to_retain"`
	PublicIPRequired     bool   `json:"public_ip_required" schema:"public_ip_required"`
	FirewallID           string `json:"firewall_id" schema:"firewall_id"`
}

// UpdateDatabaseRequest holds fields required to update a database
type UpdateDatabaseRequest struct {
	ID                   string `json:"id" schema:"id"` // used by database_test.go
	Name                 string `json:"name" schema:"name"`
	NumSnapshotsToRetain int    `json:"num_snapshots_to_retain" schema:"num_snapshots_to_retain"`
	Size                 string `json:"size" schema:"size"`
	FirewallID           string `json:"firewall_id" schema:"firewall_id"`
}

// ListDatabases returns a list of all databases
func (c *Client) ListDatabases() ([]Database, error) {
	resp, err := c.SendGetRequest("/v2/databases")
	if err != nil {
		return nil, decodeError(err)
	}

	databases := make([]Database, 0)
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&databases); err != nil {
		return nil, err
	}

	return databases, nil
}

// GetDatabase finds a database by the database UUID
func (c *Client) GetDatabase(id string) (*Database, error) {
	resp, err := c.SendGetRequest(fmt.Sprintf("/v2/databases/%s", id))
	if err != nil {
		return nil, decodeError(err)
	}

	db := &Database{}
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(db); err != nil {
		return nil, err
	}

	return db, nil
}

// DeleteDatabase deletes a database
func (c *Client) DeleteDatabase(id string) (*SimpleResponse, error) {
	resp, err := c.SendDeleteRequest(fmt.Sprintf("/v2/databases/%s", id))
	if err != nil {
		return nil, decodeError(err)
	}

	return c.DecodeSimpleResponse(resp)
}

// NewDatabase creates a new database
func (c *Client) NewDatabase(v *CreateDatabaseRequest) (*Database, error) {
	body, err := c.SendPostRequest("/v2/databases", v)
	if err != nil {
		return nil, decodeError(err)
	}

	result := &Database{}
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}

// UpdateDatabase updates a database
func (c *Client) UpdateDatabase(v *UpdateDatabaseRequest) (*Database, error) {
	body, err := c.SendPutRequest(fmt.Sprintf("/v2/databases/%s", v.ID), v)
	if err != nil {
		return nil, decodeError(err)
	}

	result := &Database{}
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}
