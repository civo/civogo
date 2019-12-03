package civogo

import (
	"bytes"
	"encoding/json"
	"errors"
)

// Template represents a Template for launching instances from
type Template struct {
	ID               string `json:"id"`
	Code             string `json:"code"`
	Name             string `json:"name"`
	AccountID        string `json:"account_id"`
	ImageID          string `json:"image_id"`
	VolumeID         string `json:"volume_id"`
	ShortDescription string `json:"short_description"`
	Description      string `json:"description"`
	DefaultUsername  string `json:"default_username"`
	CloudConfig      string `json:"cloud_config"`
}

// GetTemplateByCode finds the Template for an account with the specified code
func (c *Client) GetTemplateByCode(code string) (*Template, error) {
	resp, err := c.SendGetRequest("/v2/templates")
	if err != nil {
		return nil, err
	}

	templates := make([]Template, 0)
	err = json.NewDecoder(bytes.NewReader(resp)).Decode(&templates)
	for _, template := range templates {
		if template.Code == code {
			return &template, nil
		}
	}

	return nil, errors.New("Template not found")
}
