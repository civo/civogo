package civogo

import (
	"testing"
)

func TestGetDefaultSSHKey(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/sshkeys": `{"items":[{"id": "12345", "name": "RSA Key", "default": true}]}`,
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

func TestListSSHKeys(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/sshkeys": `{"page": 1, "per_page": 25, "items":[{"id": "12345", "name": "RSA Key", "default": true}]}`,
	})
	defer server.Close()

	got, err := client.ListSSHKeys(1, 25)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	if got.Items[0].ID != "12345" {
		t.Errorf("Expected %s, got %s", "12345", got.Items[0].ID)
	}
}

func TestFindSSHKey(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/sshkeys": `{"page": 1, "per_page": 25, "items":[{"id": "12345", "name": "RSA Key", "default": true},{"id": "67890", "name": "DSS Key", "default": true}]}`,
	})
	defer server.Close()

	got, err := client.FindSSHKey("34")
	if got.ID != "12345" {
		t.Errorf("Expected %s, got %s", "12345", got.ID)
	}

	got, _ = client.FindSSHKey("89")
	if got.ID != "67890" {
		t.Errorf("Expected %s, got %s", "67890", got.ID)
	}

	got, _ = client.FindSSHKey("RSA")
	if got.ID != "12345" {
		t.Errorf("Expected %s, got %s", "12345", got.ID)
	}

	got, _ = client.FindSSHKey("DSS")
	if got.ID != "67890" {
		t.Errorf("Expected %s, got %s", "67890", got.ID)
	}

	_, err = client.FindSSHKey("Key")
	if err.Error() != "unable to find Key because there were multiple matches" {
		t.Errorf("Expected %s, got %s", "unable to find Key because there were multiple matches", err.Error())
	}

	_, err = client.FindSSHKey("missing")
	if err.Error() != "unable to find missing, zero matches" {
		t.Errorf("Expected %s, got %s", "unable to find missing, zero matches", err.Error())
	}
}
