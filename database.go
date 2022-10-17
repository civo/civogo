package civogo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Database holds the database information
type Database struct {
	ID                   string             `json:"id"`
	Name                 string             `json:"name"`
	Replicas             int                `json:"replicas"`
	NumSnapshotsToRetain int                `json:"num_snapshots_to_retain"`
	Size                 string             `json:"size"`
	Software             string             `json:"software"`
	SoftwareVersion      string             `json:"software_version"`
	PublicIP             string             `json:"public_ip"`
	NetworkID            string             `json:"network_id"`
	FirewallID           string             `json:"firewall_id"`
	Snapshots            []DatabaseSnapshot `json:"snapshots,omitempty"`
	Status               string             `json:"status"`
}

// DatabaseSnapshot represents a database snapshot
type DatabaseSnapshot struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
}

// PaginatedDatabases is the structure for list response from DB endpoint
type PaginatedDatabases struct {
	Page    int        `json:"page"`
	PerPage int        `json:"per_page"`
	Pages   int        `json:"pages"`
	Items   []Database `json:"items"`
}

// CreateDatabaseRequest holds fields required to creates a new database
type CreateDatabaseRequest struct {
	Name                 string `json:"name" validate:"required"`
	Size                 string `json:"size" validate:"required"`
	Software             string `json:"software" validate:"required,oneof=redis postgresql mysql"`
	SoftwareVersion      string `json:"software_version"`
	NetworkID            string `json:"network_id" validate:"required"`
	Replicas             int    `json:"replicas"`
	NumSnapshotsToRetain int    `json:"num_snapshots_to_retain"`
	PublicIPRequired     bool   `json:"public_ip_required"`
	FirewallID           string `json:"firewall_id"`
}

// UpdateDatabaseRequest holds fields required to update a database
type UpdateDatabaseRequest struct {
	ID                   string `json:"id"`
	Name                 string `json:"name"`
	NumSnapshotsToRetain int    `json:"num_snapshots_to_retain"`
	Size                 string `json:"size"`
	FirewallID           string `json:"firewall_id"`
}

// ListDatabases returns a list of all databases
func (c *Client) ListDatabases() (*PaginatedDatabases, error) {
	resp, err := c.SendGetRequest("/v2/databases")
	if err != nil {
		return nil, decodeError(err)
	}

	databases := PaginatedDatabases{}
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&databases); err != nil {
		return nil, err
	}

	return &databases, nil
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

// FindDatabase finds a database by either part of the ID or part of the name
func (c *Client) FindDatabase(search string) (*Database, error) {
	databases, err := c.ListDatabases()
	if err != nil {
		return nil, decodeError(err)
	}

	exactMatch := false
	partialMatchesCount := 0
	result := Database{}

	for _, value := range databases.Items {
		if strings.EqualFold(value.Name, search) || value.ID == search {
			exactMatch = true
			result = value
		} else if strings.Contains(strings.ToUpper(value.Name), strings.ToUpper(search)) || strings.Contains(value.ID, search) {
			if !exactMatch {
				result = value
				partialMatchesCount++
			}
		}
	}

	if exactMatch || partialMatchesCount == 1 {
		return &result, nil
	} else if partialMatchesCount > 1 {
		err := fmt.Errorf("unable to find %s because there were multiple matches", search)
		return nil, MultipleMatchesError.wrap(err)
	} else {
		err := fmt.Errorf("unable to find %s, zero matches", search)
		return nil, ZeroMatchesError.wrap(err)
	}
}
