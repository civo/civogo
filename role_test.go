package civogo

import (
	"testing"
)

func TestListRoles(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/roles": `[{"id":"12345","name":"admin","permissions":"*.*","built_in":true}]`,
	})
	defer server.Close()

	got, err := client.ListRoles()
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	if got[0].ID != "12345" {
		t.Errorf("Expected %s, got %s", "12345", got[0].ID)
	}
	if got[0].Name != "admin" {
		t.Errorf("Expected %s, got %s", "admin", got[0].Name)
	}
	if got[0].Permissions != "*.*" {
		t.Errorf("Expected %s, got %s", "*.*", got[0].Permissions)
	}
	if got[0].BuiltIn != true {
		t.Errorf("Expected %t, got %t", true, got[0].BuiltIn)
	}
}

func TestCreateRole(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/roles": `{"id":"12345","name":"Org Admin","permissions":"organisation.*","built_in":false}`,
	})
	defer server.Close()

	got, err := client.CreateRole("Org Admin", "organisation.*")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	if got.ID != "12345" {
		t.Errorf("Expected %s, got %s", "12345", got.ID)
	}
	if got.Name != "Org Admin" {
		t.Errorf("Expected %s, got %s", "admin", got.Name)
	}
	if got.Permissions != "organisation.*" {
		t.Errorf("Expected %s, got %s", "*.*", got.Permissions)
	}
	if got.BuiltIn != false {
		t.Errorf("Expected %t, got %t", false, got.BuiltIn)
	}
}

func TestDeleteRole(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/roles/12345": `{"result":"success"}`,
	})
	defer server.Close()

	got, err := client.DeleteRole("12345")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	if got.Result != "success" {
		t.Errorf("Expected %s, got %s", "success", got.Result)
	}
}
