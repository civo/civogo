package civogo

import (
	"testing"
)

func TestListInstances(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/instances": `{"page": 1, "per_page": 20, "pages": 2, "items":[{"id": "12345", "hostname": "foo.example.com"}]}`,
	})
	defer server.Close()

	got, err := client.ListInstances()
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	if got.Items[0].ID != "12345" {
		t.Errorf("Expected %s, got %s", "12345", got.Items[0].ID)
	}
	if got.Page != 1 {
		t.Errorf("Expected %d, got %d", 1, got.Page)
	}
	if got.Pages != 2 {
		t.Errorf("Expected %d, got %d", 2, got.Pages)
	}
	if got.PerPage != 20 {
		t.Errorf("Expected %d, got %d", 20, got.PerPage)
	}
}
