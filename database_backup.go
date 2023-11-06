package civogo

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Scheduled represents scheduled backups
type Scheduled struct {
	Name     string   `json:"name,omitempty"`
	Schedule string   `json:"schedule,omitempty"`
	Count    int32    `json:"count,omitempty"`
	Backups  []string `json:"backups,omitempty"`
}

// Manual represents manual backups
type Manual struct {
	Backup string `json:"backup,omitempty"`
}

// DatabaseBackup represents a backup
type DatabaseBackup struct {
	DatabaseName string     `json:"database_name"`
	DatabaseID   string     `json:"database_id"`
	Software     string     `json:"software"`
	Scheduled    *Scheduled `json:"scheduled,omitempty"`
	Manual       []Manual   `json:"manual,omitempty"`
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
func (c *Client) ListDatabaseBackup(did string) (*DatabaseBackup, error) {
	resp, err := c.SendGetRequest(fmt.Sprintf("/v2/databases/%s/backups", did))
	if err != nil {
		return nil, decodeError(err)
	}

	back := &DatabaseBackup{}
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
