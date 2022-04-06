package civogo

import (
	"testing"
)

func TestListApplications(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/applications": `[{"id":"12345","name":"test-app"}]`,
	})
	defer server.Close()

	got, err := client.ListApplications()
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	if got[0].ID != "12345" {
		t.Errorf("Expected %s, got %s", "12345", got[0].ID)
	}
	if got[0].Name != "test-app" {
		t.Errorf("Expected %s, got %s", "test-app", got[0].Name)
	}
}

func TestCreateApplication(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/applications": `{"name":"test-app"}`,
	})
	defer server.Close()

	got, err := client.CreateApplication("test-app")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	if got.Name != "test-app" {
		t.Errorf("Expected %s, got %s", "test-app", got.Name)
	}
}

func TestDeleteApplication(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/applications/12345": `{"result":"success"}`,
	})
	defer server.Close()

	got, err := client.DeleteApplication("12345")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	if got.Result != "success" {
		t.Errorf("Expected %s, got %s", "success", got.Result)
	}
}
