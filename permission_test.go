package civogo

import (
	"testing"
)

func TestListPermissions(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/permissions": `[{"name":"instance.create","description":"Create Compute instances"},{"name":"billing.update","description":"Update billing details"}]`,
	})
	defer server.Close()

	got, err := client.ListPermissions()
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	if got[0].Name != "instance.create" {
		t.Errorf("Expected %s, got %s", "instance.create", got[0].Name)
	}
	if got[0].Description != "Create Compute instances" {
		t.Errorf("Expected %s, got %s", "Create Compute instances", got[0].Description)
	}
}
