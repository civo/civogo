package civogo

import (
	"testing"
)

func TestGetDefaultNetwork(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/networks": `[{"id": "12345", "default": true, "name": "Default Network"}]`,
	})
	defer server.Close()

	got, err := client.GetDefaultNetwork()
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	if got.ID != "12345" {
		t.Errorf("Expected %s, got %s", "12345", got.ID)
	}
}
