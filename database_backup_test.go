package civogo

import (
	"reflect"
	"testing"
)

func TestCreateDatabaseBackup(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/databases/12345/backups": `{
			"id": "backup123",
			"name": "initial-backup",
			"software": "MySQL",
			"status": "In Progress",
			"schedule": "manual",
			"database_name": "test-db",
			"database_id": "12345",
			"is_scheduled": false
		}`,
	})
	defer server.Close()

	v := &DatabaseBackupCreateRequest{
		Name:     "initial-backup",
		Schedule: "manual",
		Count:    1,
		Type:     "full",
		Region:   "us-east",
	}
	got, err := client.CreateDatabaseBackup("12345", v)
	if err != nil {
		t.Errorf("CreateDatabaseBackup returned an error: %s", err)
		return
	}

	expected := &DatabaseBackup{
		ID:           "backup123",
		Name:         "initial-backup",
		Software:     "MySQL",
		Status:       "In Progress",
		Schedule:     "manual",
		DatabaseName: "test-db",
		DatabaseID:   "12345",
		IsScheduled:  false,
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestUpdateDatabaseBackup(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/databases/12345/backups": `{
			"id": "backup123",
			"name": "updated-backup",
			"schedule": "weekly",
			"count": 2,
			"region": "us-west"
		}`,
	})
	defer server.Close()

	v := &DatabaseBackupUpdateRequest{
		Name:     "updated-backup",
		Schedule: "weekly",
		Count:    2,
		Region:   "us-west",
	}
	got, err := client.UpdateDatabaseBackup("12345", v)
	if err != nil {
		t.Errorf("UpdateDatabaseBackup returned an error: %s", err)
		return
	}

	expected := &DatabaseBackup{
		ID:       "backup123",
		Name:     "updated-backup",
		Schedule: "weekly",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestGetDatabaseBackup(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/databases/12345/backups/backup123": `{
			"id": "backup123",
			"name": "nightly-backup",
			"software": "MySQL",
			"status": "Completed",
			"schedule": "daily",
			"database_name": "test-db",
			"database_id": "12345",
			"backup": "url-to-backup",
			"is_scheduled": true
		}`,
	})
	defer server.Close()

	got, err := client.GetDatabaseBackup("12345", "backup123")
	if err != nil {
		t.Errorf("GetDatabaseBackup returned an error: %s", err)
		return
	}

	expected := &DatabaseBackup{
		ID:           "backup123",
		Name:         "nightly-backup",
		Software:     "MySQL",
		Status:       "Completed",
		Schedule:     "daily",
		DatabaseName: "test-db",
		DatabaseID:   "12345",
		Backup:       "url-to-backup",
		IsScheduled:  true,
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestListDatabaseBackup(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/databases/12345/backups": `{
			"page": 1,
			"per_page": 2,
			"pages": 1,
			"items": [
				{
					"id": "backup123",
					"name": "initial-backup"
				},
				{
					"id": "backup124",
					"name": "second-backup"
				}
			]
		}`,
	})
	defer server.Close()

	got, err := client.ListDatabaseBackup("12345")
	if err != nil {
		t.Errorf("ListDatabaseBackup returned an error: %s", err)
		return
	}

	expected := &PaginatedDatabaseBackup{
		Page:    1,
		PerPage: 2,
		Pages:   1,
		Items: []DatabaseBackup{
			{ID: "backup123", Name: "initial-backup"},
			{ID: "backup124", Name: "second-backup"},
		},
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}
