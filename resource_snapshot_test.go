package civogo

import (
	"reflect"
	"testing"
	"time"
)

func TestListResourceSnapshots(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/resourcesnapshots": `[
			{
				"id": "12345",
				"name": "test-snapshot",
				"description": "Test snapshot",
				"resource_type": "instance",
				"created_at": "2023-01-01T12:00:00Z",
				"instance": {
					"id": "inst-12345",
					"name": "test-instance",
					"description": "Test instance",
					"status": {"state": "available"},
					"created_at": "2023-01-01T12:00:00Z"
				}
			}
		]`,
	})
	defer server.Close()

	got, err := client.ListResourceSnapshots()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(got) != 1 {
		t.Errorf("Expected 1 snapshot, got %d", len(got))
	}

	if got[0].ID != "12345" {
		t.Errorf("Expected ID 12345, got %s", got[0].ID)
	}
}

func TestGetResourceSnapshot(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/resourcesnapshots/12345": `{
			"id": "12345",
			"name": "test-snapshot",
			"description": "Test snapshot",
			"resource_type": "instance",
			"created_at": "2023-01-01T12:00:00Z"
		}`,
	})
	defer server.Close()

	got, err := client.GetResourceSnapshot("12345")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	expectedTime, _ := time.Parse(time.RFC3339, "2023-01-01T12:00:00Z")
	expected := &ResourceSnapshot{
		ID:           "12345",
		Name:         "test-snapshot",
		Description:  "Test snapshot",
		ResourceType: "instance",
		CreatedAt:    expectedTime,
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestUpdateResourceSnapshot(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/resourcesnapshots/12345": `{
			"id": "12345",
			"name": "updated-snapshot",
			"description": "Updated description",
			"resource_type": "instance",
			"created_at": "2023-01-01T12:00:00Z"
		}`,
	})
	defer server.Close()

	req := &UpdateResourceSnapshotRequest{
		Name:        "updated-snapshot",
		Description: "Updated description",
	}

	got, err := client.UpdateResourceSnapshot("12345", req)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if got.Name != "updated-snapshot" {
		t.Errorf("Expected name 'updated-snapshot', got %s", got.Name)
	}
}

func TestDeleteResourceSnapshot(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/resourcesnapshots/12345": `{"result": "success"}`,
	})
	defer server.Close()

	got, err := client.DeleteResourceSnapshot("12345")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if got.Result != "success" {
		t.Errorf("Expected result 'success', got %s", got.Result)
	}
}

func TestRestoreResourceSnapshot(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/resourcesnapshots/12345/restore": `{
			"id": "12345",
			"name": "restored-snapshot",
			"description": "Restored snapshot",
			"resource_type": "instance",
			"created_at": "2023-01-01T12:00:00Z"
		}`,
	})
	defer server.Close()

	req := &RestoreResourceSnapshotRequest{
		Instance: &RestoreInstanceSnapshotRequest{
			Description:    "Restored snapshot",
			Hostname:       "restored-instance",
			IncludeVolumes: true,
		},
	}

	got, err := client.RestoreResourceSnapshot("12345", req)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if got.Name != "restored-snapshot" {
		t.Errorf("Expected name 'restored-snapshot', got %s", got.Name)
	}
}
