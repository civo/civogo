package civogo

import (
	"reflect"
	"testing"
)

func TestNewSSHKey(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/sshkeys": `{
		  "result": "success",
		  "id": "730c960f-a51f-44e5-9c21-bd135d015d12"
		}`,
	})
	defer server.Close()
	got, err := client.NewSSHKey("test", "730c960f-a51f-44e5-9c21-bd135d015d12")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &SimpleResponse{Result: "success", ID: "730c960f-a51f-44e5-9c21-bd135d015d12"}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestListSSHKeys(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/sshkeys": `[
			{"id": "12345", "name": "RSA Key", "fingerprint": "SHA256:SS4+2d7Zl1Pt5Bc9af9NubyA0yI+fdopOUlEhUoEna0" },
			{"id": "33567", "name": "RSA Key", "fingerprint": "SHA256:SS4+87asdf795Bc9af9NubyA0yI+fdopOUlEhUoEna0" }]`,
	})
	defer server.Close()

	got, err := client.ListSSHKeys()
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	if got[0].ID != "12345" {
		t.Errorf("Expected %s, got %s", "12345", got[0].ID)
	}
}

func TestFindSSHKey(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/sshkeys": `[
			{"id": "12345", "name": "RSA Key", "fingerprint": "SHA256:SS4+2d7Zl1Pt5Bc9af9NubyA0yI+fdopOUlEhUoEna0" },
			{"id": "233567", "name": "Test", "fingerprint": "SHA256:SS4+87asdf795Bc9af9NubyA0yI+fdopOUlEhUoEna0" }]`,
	})
	defer server.Close()

	got, err := client.FindSSHKey("34")
	if got.ID != "12345" {
		t.Errorf("Expected %s, got %s", "12345", got.ID)
	}

	got, _ = client.FindSSHKey("RSA")
	if got.ID != "12345" {
		t.Errorf("Expected %s, got %s", "12345", got.ID)
	}

	_, err = client.FindSSHKey("23")
	if err.Error() != "unable to find 23 because there were multiple matches" {
		t.Errorf("Expected %s, got %s", "unable to find 23 because there were multiple matches", err.Error())
	}

	_, err = client.FindSSHKey("missing")
	if err.Error() != "unable to find missing, zero matches" {
		t.Errorf("Expected %s, got %s", "unable to find missing, zero matches", err.Error())
	}
}
