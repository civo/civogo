package civogo

import (
	"testing"
)

func TestListRegions(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/regions": `[{"code":"lon1", "name": "London 1", "default": true}]`,
	})
	defer server.Close()

	got, err := client.ListRegions()
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	if got[0].Code != "lon1" {
		t.Errorf("Expected %s, got %s", "lon1", got[0].Code)
	}
	if got[0].Name != "London 1" {
		t.Errorf("Expected %s, got %s", "London 1", got[0].Name)
	}
	if !got[0].Default {
		t.Errorf("Expected first result to be the default")
	}
}
