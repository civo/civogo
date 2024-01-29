package civogo

import (
	"fmt"
	"strings"
)

// CreateWebhook implemented in a fake way for automated tests
func (c *FakeClient) CreateWebhook(r *WebhookConfig) (*Webhook, error) {
	webhook := Webhook{
		ID:     c.generateID(),
		Events: r.Events,
		Secret: r.Secret,
		URL:    r.URL,
	}
	c.Webhooks = append(c.Webhooks, webhook)

	return &webhook, nil
}

// ListWebhooks implemented in a fake way for automated tests
func (c *FakeClient) ListWebhooks() ([]Webhook, error) {
	return c.Webhooks, nil
}

// FindWebhook implemented in a fake way for automated tests
func (c *FakeClient) FindWebhook(search string) (*Webhook, error) {
	for _, webhook := range c.Webhooks {
		if strings.Contains(webhook.Secret, search) || strings.Contains(webhook.URL, search) {
			return &webhook, nil
		}
	}

	err := fmt.Errorf("unable to find %s, zero matches", search)
	return nil, ZeroMatchesError.wrap(err)
}

// UpdateWebhook implemented in a fake way for automated tests
func (c *FakeClient) UpdateWebhook(id string, r *WebhookConfig) (*Webhook, error) {
	for i, webhook := range c.Webhooks {
		if webhook.ID == id {
			c.Webhooks[i].Events = r.Events
			c.Webhooks[i].Secret = r.Secret
			c.Webhooks[i].URL = r.URL

			return &webhook, nil
		}
	}

	err := fmt.Errorf("unable to find %s, zero matches", id)
	return nil, ZeroMatchesError.wrap(err)
}

// DeleteWebhook implemented in a fake way for automated tests
func (c *FakeClient) DeleteWebhook(id string) (*SimpleResponse, error) {
	for i, webhook := range c.Webhooks {
		if webhook.ID == id {
			c.Webhooks[len(c.Webhooks)-1], c.Webhooks[i] = c.Webhooks[i], c.Webhooks[len(c.Webhooks)-1]
			c.Webhooks = c.Webhooks[:len(c.Webhooks)-1]
			return &SimpleResponse{Result: "success"}, nil
		}
	}

	return &SimpleResponse{Result: "failed"}, nil
}
