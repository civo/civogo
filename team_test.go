package civogo

import (
	"testing"
)

func TestListTeams(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/teams": `[{"id":"12345","name":"admin"}]`,
	})
	defer server.Close()

	got, err := client.ListTeams()
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
}

func TestCreateTeam(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/teams": `{"name":"Org Admin"}`,
	})
	defer server.Close()

	got, err := client.CreateTeam("Org Admin")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	if got.Name != "Org Admin" {
		t.Errorf("Expected %s, got %s", "admin", got.Name)
	}
}

func TestRenameTeam(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/teams/12345": `{"id":"12345","name":"New Admin"}`,
	})
	defer server.Close()

	got, err := client.RenameTeam("12345", "New Admin")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	if got.ID != "12345" {
		t.Errorf("Expected %s, got %s", "12345", got.ID)
	}
	if got.Name != "New Admin" {
		t.Errorf("Expected %s, got %s", "admin", got.Name)
	}
}

func TestDeleteTeam(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/teams/12345": `{"result":"success"}`,
	})
	defer server.Close()

	got, err := client.DeleteTeam("12345")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	if got.Result != "success" {
		t.Errorf("Expected %s, got %s", "success", got.Result)
	}
}

func TestListTeamMembers(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/teams/12345/members": `[{"id":"abcde","user_id":"bcdef"}]`,
	})
	defer server.Close()

	got, err := client.ListTeamMembers("12345")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	if got[0].ID != "abcde" {
		t.Errorf("Expected %s, got %s", "12345", got[0].ID)
	}
	if got[0].UserID != "bcdef" {
		t.Errorf("Expected %s, got %s", "bcdef", got[0].UserID)
	}
}

func TestAddTeamMember(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/teams/12345/members": `[{"id":"abcde","user_id":"bcdef"}]`,
	})
	defer server.Close()

	got, err := client.AddTeamMember("12345", "abcde", "*.*", "")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	if got[0].ID != "abcde" {
		t.Errorf("Expected %s, got %s", "12345", got[0].ID)
	}
	if got[0].UserID != "bcdef" {
		t.Errorf("Expected %s, got %s", "bcdef", got[0].UserID)
	}
}

func TestUpdateTeamMember(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/teams/12345/members/abcde": `{"id":"12345","permissions":"*.*"}`,
	})
	defer server.Close()

	got, err := client.UpdateTeamMember("12345", "abcde", "*.*", "")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	if got.ID != "12345" {
		t.Errorf("Expected %s, got %s", "12345", got.ID)
	}
	if got.Permissions != "*.*" {
		t.Errorf("Expected %s, got %s", "*.*", got.Permissions)
	}
}

func TestRemoveTeamMember(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/teams/12345/members/abcde": `{"result":"success"}`,
	})
	defer server.Close()

	got, err := client.RemoveTeamMember("12345", "abcde")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	if got.Result != "success" {
		t.Errorf("Expected %s, got %s", "success", got.Result)
	}
}
