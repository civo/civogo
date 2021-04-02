package civogo

import (
	"reflect"
	"testing"
)

func TestCreateWebhook(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/webhooks": `{
		  "id": "b8de2e4e-72f4-4911-83ee-f4fc030fc4a2",
		  "events": ["*"],
		  "url": "https://api.example.com/webhook",
		  "secret": "DfeFUON8gorc5Zj0hk4GT1M9QImnRL6J",
		  "disabled": false,
		  "failures": 0,
		  "last_failure_reason": ""
		}`,
	})
	defer server.Close()

	cfg := &WebhookConfig{
		Events: []string{"*"},
		URL:    "https://api.example.com/webhook",
		Secret: "DfeFUON8gorc5Zj0hk4GT1M9QImnRL6J",
	}
	got, err := client.CreateWebhook(cfg)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &Webhook{
		ID:                "b8de2e4e-72f4-4911-83ee-f4fc030fc4a2",
		Events:            []string{"*"},
		URL:               "https://api.example.com/webhook",
		Secret:            "DfeFUON8gorc5Zj0hk4GT1M9QImnRL6J",
		Disabled:          false,
		Failures:          0,
		LasrFailureReason: "",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestListWebhooks(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/webhooks": `[{
		  "id": "b8de2e4e-72f4-4911-83ee-f4fc030fc4a2",
		  "events": ["*"],
		  "url": "https://api.example.com/webhook",
		  "secret": "DfeFUON8gorc5Zj0hk4GT1M9QImnRL6J",
		  "disabled": false,
		  "failures": 0,
		  "last_failure_reason": ""
		}]`,
	})
	defer server.Close()

	got, err := client.ListWebhooks()
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := []Webhook{{
		ID:                "b8de2e4e-72f4-4911-83ee-f4fc030fc4a2",
		Events:            []string{"*"},
		URL:               "https://api.example.com/webhook",
		Secret:            "DfeFUON8gorc5Zj0hk4GT1M9QImnRL6J",
		Disabled:          false,
		Failures:          0,
		LasrFailureReason: "",
	}}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestFindWebhook(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/webhooks": `[
			{
				"id": "12345",
				"events": ["*"],
				"url": "https://my.example.com/webhook",
				"secret": "DfeFUON8gorc5Zj0hk4GT1M9QImnRL6J",
				"disabled": false,
				"failures": 0,
				"last_failure_reason": ""
			},
			{
				"id": "67890",
				"events": ["*"],
				"url": "https://other.example.com/webhook",
				"secret": "DfeFUON8gorc5Zj0hk4GT1M9QImnRL6J",
				"disabled": false,
				"failures": 0,
				"last_failure_reason": ""
			}
		]`,
	})
	defer server.Close()

	got, _ := client.FindWebhook("34")
	if got.ID != "12345" {
		t.Errorf("Expected %s, got %s", "12345", got.ID)
	}

	got, _ = client.FindWebhook("89")
	if got.ID != "67890" {
		t.Errorf("Expected %s, got %s", "67890", got.ID)
	}

	got, _ = client.FindWebhook("my")
	if got.ID != "12345" {
		t.Errorf("Expected %s, got %s", "12345", got.ID)
	}

	got, _ = client.FindWebhook("other")
	if got.ID != "67890" {
		t.Errorf("Expected %s, got %s", "67890", got.ID)
	}

	_, err := client.FindWebhook("example")
	if err.Error() != "MultipleMatchesError: unable to find example because there were multiple matches" {
		t.Errorf("Expected %s, got %s", "unable to find example because there were multiple matches", err.Error())
	}

	_, err = client.FindWebhook("missing")
	if err.Error() != "ZeroMatchesError: unable to find missing, zero matches" {
		t.Errorf("Expected %s, got %s", "unable to find missing, zero matches", err.Error())
	}
}

func TestUpdateWebhook(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/webhooks/b8de2e4e-72f4-4911-83ee-f4fc030fc4a2": `{
		  "id": "b8de2e4e-72f4-4911-83ee-f4fc030fc4a2",
		  "events": ["instance.created", "instance.active"],
		  "url": "https://api.example.com/webhook",
		  "secret": "DfeFUON8gorc5Zj0hk4GT1M9QImnRL6J",
		  "disabled": false,
		  "failures": 0,
		  "last_failure_reason": ""
		}`,
	})
	defer server.Close()
	cfg := &WebhookConfig{
		Events: []string{"instance.created", "instance.active"},
	}
	got, err := client.UpdateWebhook("b8de2e4e-72f4-4911-83ee-f4fc030fc4a2", cfg)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &Webhook{
		ID:                "b8de2e4e-72f4-4911-83ee-f4fc030fc4a2",
		Events:            []string{"instance.created", "instance.active"},
		URL:               "https://api.example.com/webhook",
		Secret:            "DfeFUON8gorc5Zj0hk4GT1M9QImnRL6J",
		Disabled:          false,
		Failures:          0,
		LasrFailureReason: "",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestDeleteWebhook(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/webhooks/b8de2e4e-72f4-4911-83ee-f4fc030fc4a2": `{"result": "success"}`,
	})
	defer server.Close()

	got, err := client.DeleteWebhook("b8de2e4e-72f4-4911-83ee-f4fc030fc4a2")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &SimpleResponse{Result: "success"}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}
