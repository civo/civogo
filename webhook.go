package civogo

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Webhook is a representation of the one o more actions
type Webhook struct {
	ID                string   `json:"id"`
	Events            []string `json:"events"`
	Url               string   `json:"url"`
	Secret            string   `json:"secret"`
	Disabled          bool     `json:"disabled"`
	Failures          int      `json:"failures"`
	LasrFailureReason string   `json:"last_failure_reason"`
}

// WebhookConfig represents the options required for creating a new webhook
type WebhookConfig struct {
	Events []string `form:"events"`
	Url    string   `form:"url"`
	Secret string   `form:"secret"`
}

// CreateSnapshot create a new or update an old snapshot
func (c *Client) CreateWebhook(r *WebhookConfig) (*Webhook, error) {
	body, err := c.SendPostRequest("/v2/webhooks", r)
	if err != nil {
		return nil, err
	}

	var n = &Webhook{}
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(n); err != nil {
		return nil, err
	}

	return n, nil
}

// ListWebhooks returns a list of all webhook within the current account
func (c *Client) ListWebhooks() ([]Webhook, error) {
	resp, err := c.SendGetRequest("/v2/webhooks")
	if err != nil {
		return nil, err
	}

	webhook := make([]Webhook, 0)
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&webhook); err != nil {
		return nil, err
	}

	return webhook, nil
}

// UpdateWebhook update a webhook
func (c *Client) UpdateWebhook(id string, r *WebhookConfig) (*Webhook, error) {
	body, err := c.SendPutRequest(fmt.Sprintf("/v2/webhooks/%s", id), r)
	if err != nil {
		return nil, err
	}

	var n = &Webhook{}
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(n); err != nil {
		return nil, err
	}

	return n, nil
}

// DeleteWebhook deletes a webhook
func (c *Client) DeleteWebhook(id string) (*SimpleResponse, error) {
	resp, err := c.SendDeleteRequest(fmt.Sprintf("/v2/webhooks/%s", id))
	if err != nil {
		return nil, err
	}

	return c.DecodeSimpleResponse(resp)
}
