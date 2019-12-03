package civogo

import (
	"testing"
)

func TestGetDefaultSSHKey(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/sshkeys": `{"items":[{"id": "12345", "name": "RSA Key"}]}`,
	})
	defer server.Close()

	got, err := client.GetDefaultSSHKey()
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	if got.ID != "12345" {
		t.Errorf("Expected %s, got %s", "12345", got.ID)
	}
}
