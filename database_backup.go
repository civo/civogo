package civogo

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// DatabaseBackup represents a backup
type DatabaseBackup struct {
	Name         string   `json:"name"`
	DatabaseName string   `json:"database_name"`
	DatabaseID   string   `json:"database_id"`
	Software     string   `json:"software"`
	Schedule     string   `json:"schedule"`
	Count        int32    `json:"count"`
	Backups      []string `json:"backups"`
}

// DatabaseBackupCreateRequest represents a backup create request
type DatabaseBackupCreateRequest struct {
	Name     string `json:"name"`
	Schedule string `json:"schedule"`
	Count    int32  `json:"count"`
	Region   string `json:"region"`
}

// DatabaseBackupUpdateRequest represents a backup update request
type DatabaseBackupUpdateRequest struct {
	Name     string `json:"name"`
	Schedule string `json:"schedule"`
	Count    int32  `json:"count"`
	Region   string `json:"region"`
}

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
