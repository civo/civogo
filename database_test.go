package civogo

import (
	"reflect"
	"testing"
)

func TestListDatabases(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/databases": `{"page": 1, "per_page": 20, "pages": 2, "items":[{"id": "12345", "name": "test-db"}]}`,
	})
	defer server.Close()

	got, err := client.ListDatabases()
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &PaginatedDatabases{
		Page:    1,
		PerPage: 20,
		Pages:   2,
		Items: []Database{
			{
				ID:   "12345",
				Name: "test-db",
			},
		},
	}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestFindDatabase(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/databases": `{
			"page": 1,
			"per_page": 20,
			"pages": 1,
			"items": [
			  {
				"id": "12345",
				"name": "test-db"
			  }
			]
		  }`,
	})
	defer server.Close()

	got, _ := client.FindDatabase("test-db")
	if got.ID != "12345" {
		t.Errorf("Expected %s, got %s", "12345", got.ID)
	}
}

func TestNewDatabase(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/databases": `{
			"id": "12345",
			"name": "test-db",
			"size": "g3.db.xsmall",
			"software": "MySQL",
			"status" : "Ready"
		}`,
	})
	defer server.Close()

	cfg := &CreateDatabaseRequest{
		Name:     "test-db",
		Size:     "g3.db.xsmall",
		Software: "MySQL",
	}
	got, err := client.NewDatabase(cfg)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &Database{
		ID:       "12345",
		Name:     "test-db",
		Size:     "g3.db.xsmall",
		Software: "MySQL",
		Status:   "Ready",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestDeleteDatabase(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/databases/12345": `{"result": "success"}`,
	})
	defer server.Close()

	got, err := client.DeleteDatabase("12345")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &SimpleResponse{Result: "success"}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestRestoreDatabase(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/databases/12345/restore": `{"result": "success"}`,
	})
	defer server.Close()

	// Define the input parameters for RestoreDatabase
	restoreRequest := &RestoreDatabaseRequest{
		Name:      "test-db",
		Software:  "MySQL",
		NetworkID: "network-id",
		Backup:    "backup-file-path",
		Region:    "some-region",
	}

	// Call RestoreDatabase method
	got, err := client.RestoreDatabase("12345", restoreRequest)
	if err != nil {
		t.Errorf("RestoreDatabase returned an error: %s", err)
		return
	}

	expected := &SimpleResponse{Result: "success"}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestGetDatabase(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/databases/12345": `{
			"id": "12345",
			"name": "test-db",
			"nodes": 1,
			"size": "g3.db.small",
			"software": "MySQL",
			"software_version": "8.0",
			"public_ipv4": "127.0.0.1",
			"network_id": "network-id",
			"firewall_id": "firewall-id",
			"port": 3306,
			"username": "user",
			"password": "password",
			"database_user_info": [
				{"username": "user", "password": "password", "port": 3306}
			],
			"status": "Ready"
		}`,
	})
	defer server.Close()

	// Call GetDatabase method
	got, err := client.GetDatabase("12345")
	if err != nil {
		t.Errorf("GetDatabase returned an error: %s", err)
		return
	}

	expected := &Database{
		ID:              "12345",
		Name:            "test-db",
		Nodes:           1,
		Size:            "g3.db.small",
		Software:        "MySQL",
		SoftwareVersion: "8.0",
		PublicIPv4:      "127.0.0.1",
		NetworkID:       "network-id",
		FirewallID:      "firewall-id",
		Port:            3306,
		Username:        "user",
		Password:        "password",
		DatabaseUserInfo: []DatabaseUserInfo{
			{Username: "user", Password: "password", Port: 3306},
		},
		Status: "Ready",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestUpdateDatabase(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/databases/12345": `{
			"id": "12345",
			"name": "updated-test-db",
			"nodes": 2,
			"size": "g3.db.medium",
			"software": "PostgreSQL",
			"software_version": "13",
			"public_ipv4": "127.0.0.1",
			"network_id": "network-id",
			"firewall_id": "firewall-id",
			"port": 5432,
			"username": "user",
			"password": "password",
			"database_user_info": [
				{"username": "user", "password": "password", "port": 5432}
			],
			"status": "Ready"
		}`,
	})
	defer server.Close()

	nodeCount := 2

	// Define the input parameters for UpdateDatabase
	updateRequest := &UpdateDatabaseRequest{
		Name:       "updated-test-db",
		Nodes:      &nodeCount, // Int is a helper function to create an integer pointer
		FirewallID: "firewall-id",
		Region:     "region",
	}

	// Call UpdateDatabase method
	got, err := client.UpdateDatabase("12345", updateRequest)
	if err != nil {
		t.Errorf("UpdateDatabase returned an error: %s", err)
		return
	}

	expected := &Database{
		ID:              "12345",
		Name:            "updated-test-db",
		Nodes:           2,
		Size:            "g3.db.medium",
		Software:        "PostgreSQL",
		SoftwareVersion: "13",
		PublicIPv4:      "127.0.0.1",
		NetworkID:       "network-id",
		FirewallID:      "firewall-id",
		Port:            5432,
		Username:        "user",
		Password:        "password",
		DatabaseUserInfo: []DatabaseUserInfo{
			{Username: "user", Password: "password", Port: 5432},
		},
		Status: "Ready",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}
