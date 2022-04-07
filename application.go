package civogo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

// Application is the struct for the Application model
type Application struct {
	Name        string        `json:"name" validate:"required" schema:"name"`
	AccountID   string        `json:"account_id"`
	ID          string        `json:"id"`
	NetworkID   string        `json:"network_id" validate:"required" schema:"network_id"`
	Description string        `json:"description"`
	ProcessInfo []ProcessInfo `json:"process_info,omitempty"`
	Domains     []string      `json:"domains,omitempty"`
	SSHKeyIDs   []string      `json:"ssh_key_ids,omitempty"`
	Config      []EnvVar      `json:"config,omitempty"`
	// Status can be one of:
	// - "building":  Implies app is building
	// - "ready": Implies app is ready
	Status string `json:"status"`
}

type PaginatedApplications struct {
	Page    int           `json:"page"`
	PerPage int           `json:"per_page"`
	Pages   int           `json:"pages"`
	Items   []Application `json:"items"`
}

// EnvVar holds key-value pairs for an application
type EnvVar struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// ProcessInfo contains the information about the process obtained from Procfile
type ProcessInfo struct {
	ProcessType  string `json:"process_type"`
	ProcessCount int    `json:"process_count"`
}

// ListApplications returns all applications in that specific region
func (c *Client) ListApplications() (*PaginatedApplications, error) {
	resp, err := c.SendGetRequest("/v2/applications")
	if err != nil {
		return nil, decodeError(err)
	}

	application := &PaginatedApplications{}
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&application); err != nil {
		return nil, decodeError(err)
	}

	return application, nil
}

// GetApplication returns an application by ID
func (c *Client) GetApplication(id string) (*Application, error) {
	resp, err := c.SendGetRequest(fmt.Sprintf("/v2/applications/%s", id))
	if err != nil {
		return nil, decodeError(err)
	}

	application := &Application{}
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&application); err != nil {
		return nil, decodeError(err)
	}

	return application, nil
}

// FindApplication finds an application by either part of the ID or part of the name
func (c *Client) FindApplication(search string) (*Application, error) {
	apps, err := c.ListApplications()
	if err != nil {
		return nil, decodeError(err)
	}

	exactMatch := false
	partialMatchesCount := 0
	result := Application{}

	for _, value := range apps.Items {
		if value.Name == search || value.ID == search {
			exactMatch = true
			result = value
		} else if strings.Contains(value.Name, search) || strings.Contains(value.ID, search) {
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

// CreateApplication creates a new application
func (c *Client) CreateApplication(name string) (*Application, error) {
	body, err := c.SendPostRequest("/v2/applications", name)
	if err != nil {
		return nil, decodeError(err)
	}

	application := &Application{}
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(application); err != nil {
		return nil, err
	}

	return application, nil
}

// DeleteApplication deletes an application
func (c *Client) DeleteApplication(id string) (*SimpleResponse, error) {
	resp, err := c.SendDeleteRequest(fmt.Sprintf("/v2/applications/%s", id))
	if err != nil {
		return nil, decodeError(err)
	}

	return c.DecodeSimpleResponse(resp)
}
