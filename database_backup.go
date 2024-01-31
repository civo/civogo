package civogo

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// DatabaseBackup represents a backup
type DatabaseBackup struct {
	ID           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Software     string `json:"software,omitempty"`
	Status       string `json:"status,omitempty"`
	Schedule     string `json:"schedule,omitempty"`
	DatabaseName string `json:"database_name,omitempty"`
	DatabaseID   string `json:"database_id,omitempty"`
	Backup       string `json:"backup,omitempty"`
	IsScheduled  bool   `json:"is_scheduled,omitempty"`
}

// PaginatedDatabaseBackup is the structure for list response from DB endpoint
type PaginatedDatabaseBackup struct {
	Page    int              `json:"page"`
	PerPage int              `json:"per_page"`
	Pages   int              `json:"pages"`
	Items   []DatabaseBackup `json:"items"`
}

// DatabaseBackupCreateRequest represents a backup create request
type DatabaseBackupCreateRequest struct {
	Name     string `json:"name"`
	Schedule string `json:"schedule"`
	Count    int32  `json:"count"`
	Type     string `json:"type"`
	Region   string `json:"region"`
}

// DatabaseBackupUpdateRequest represents a backup update request
type DatabaseBackupUpdateRequest struct {
	Name     string `json:"name"`
	Schedule string `json:"schedule"`
	Count    int32  `json:"count"`
	Region   string `json:"region"`
}

// ListDatabaseBackup lists backups for database
func (c *Client) ListDatabaseBackup(did string) (*PaginatedDatabaseBackup, error) {
	resp, err := c.SendGetRequest(fmt.Sprintf("/v2/databases/%s/backups", did))
	if err != nil {
		return nil, decodeError(err)
	}

	back := &PaginatedDatabaseBackup{}
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&back); err != nil {
		return nil, decodeError(err)
	}

	return back, nil
}

// UpdateDatabaseBackup update database backup
func (c *Client) UpdateDatabaseBackup(did string, v *DatabaseBackupUpdateRequest) (*DatabaseBackup, error) {
	body, err := c.SendPutRequest(fmt.Sprintf("/v2/databases/%s/backups", did), v)
	if err != nil {
		return nil, decodeError(err)
	}

	result := &DatabaseBackup{}
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}

// CreateDatabaseBackup create database backup
func (c *Client) CreateDatabaseBackup(did string, v *DatabaseBackupCreateRequest) (*DatabaseBackup, error) {
	body, err := c.SendPostRequest(fmt.Sprintf("/v2/databases/%s/backups", did), v)
	if err != nil {
		return nil, decodeError(err)
	}

	result := &DatabaseBackup{}
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}
