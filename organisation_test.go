package civogo

import (
	"testing"
)

func TestGetOrganisation(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/organisation": `{"id":"12345", "name":"Some Org Ltd"}`,
	})
	defer server.Close()

	got, err := client.GetOrganisation()
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	if got.Name != "Some Org Ltd" {
		t.Errorf("Expected %s, got %s", "Some Org Ltd", got.Name)
	}
	if got.ID != "12345" {
		t.Errorf("Expected %s, got %s", "12345", got.ID)
	}
}

func TestCreateOrganisation(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/organisation": `{"id":"12345", "name":"Some Org Ltd"}`,
	})
	defer server.Close()

	got, err := client.CreateOrganisation("Some Org Ltd")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	if got.Name != "Some Org Ltd" {
		t.Errorf("Expected %s, got %s", "Some Org Ltd", got.Name)
	}
	if got.ID != "12345" {
		t.Errorf("Expected %s, got %s", "12345", got.ID)
	}
}

func TestRenameOrganisation(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/organisation": `{"id":"12345", "name":"Some Org Ltd"}`,
	})
	defer server.Close()

	got, err := client.RenameOrganisation("Some Org Ltd")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	if got.Name != "Some Org Ltd" {
		t.Errorf("Expected %s, got %s", "Some Org Ltd", got.Name)
	}
	if got.ID != "12345" {
		t.Errorf("Expected %s, got %s", "12345", got.ID)
	}
}

func TestAddAccountToOrganisation(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/organisation/accounts": `[{"id":"abcde", "token":"fghij"}]`,
	})
	defer server.Close()

	got, err := client.AddAccountToOrganisation("abcde", "fghij")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	if got[0].ID != "abcde" {
		t.Errorf("Expected %s, got %s", "abcde", got[0].ID)
	}
}

func TestListAccountsInOrganisation(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/organisation/accounts": `[{"id":"abcde"}]`,
	})
	defer server.Close()

	got, err := client.ListAccountsInOrganisation()
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	if got[0].ID != "abcde" {
		t.Errorf("Expected %s, got %s", "abcde", got[0].ID)
	}
}
