package civogo

import (
	"testing"
)

func TestGetTemplateByCode(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/templates": `[{"id": "1", "code": "centos-7"},{"id": "2", "code": "ubuntu-18.04"}]`,
	})
	defer server.Close()

	got, err := client.GetTemplateByCode("ubuntu-18.04")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	if got.ID != "2" {
		t.Errorf("Expected %s, got %s", "12345", got.ID)
	}
}
