package civogo

import (
	"testing"
)

func TestCreateInstanceSnapshot(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/instances/12345/snapshots": `{
			"id": "snapshot-123",
			"name": "test-snapshot",
			"description": "Test snapshot",
			"included_volumes": ["vol-1", "vol-2"],
			"status": {
				"state": "pending",
				"message": "Creating snapshot",
				"volumes": [
					{
						"id": "vol-1",
						"state": "pending"
					},
					{
						"id": "vol-2",
						"state": "pending"
					}
				]
			},
			"created_at": "2024-03-20T10:00:00Z"
		}`,
	})
	defer server.Close()

	config := &CreateInstanceSnapshotConfig{
		Name:           "test-snapshot",
		Description:    "Test snapshot",
		IncludeVolumes: true,
	}

	got, err := client.CreateInstanceSnapshot("12345", config)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	if got.ID != "snapshot-123" {
		t.Errorf("Expected %s, got %s", "snapshot-123", got.ID)
	}
	if got.Name != "test-snapshot" {
		t.Errorf("Expected %s, got %s", "test-snapshot", got.Name)
	}
	if got.Description != "Test snapshot" {
		t.Errorf("Expected %s, got %s", "Test snapshot", got.Description)
	}
	if len(got.IncludedVolumes) != 2 {
		t.Errorf("Expected 2 volumes, got %d", len(got.IncludedVolumes))
	}
	if got.Status.State != "pending" {
		t.Errorf("Expected %s, got %s", "pending", got.Status.State)
	}
}

func TestGetInstanceSnapshot(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/instances/12345/snapshots/snapshot-123": `{
			"id": "snapshot-123",
			"name": "test-snapshot",
			"description": "Test snapshot",
			"included_volumes": ["vol-1", "vol-2"],
			"status": {
				"state": "completed",
				"message": "Snapshot completed successfully",
				"volumes": [
					{
						"id": "vol-1",
						"state": "completed"
					},
					{
						"id": "vol-2",
						"state": "completed"
					}
				]
			},
			"created_at": "2024-03-20T10:00:00Z"
		}`,
	})
	defer server.Close()

	got, err := client.GetInstanceSnapshot("12345", "snapshot-123")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	if got.ID != "snapshot-123" {
		t.Errorf("Expected %s, got %s", "snapshot-123", got.ID)
	}
	if got.Name != "test-snapshot" {
		t.Errorf("Expected %s, got %s", "test-snapshot", got.Name)
	}
	if got.Description != "Test snapshot" {
		t.Errorf("Expected %s, got %s", "Test snapshot", got.Description)
	}
	if len(got.IncludedVolumes) != 2 {
		t.Errorf("Expected 2 volumes, got %d", len(got.IncludedVolumes))
	}
	if got.Status.State != "completed" {
		t.Errorf("Expected %s, got %s", "completed", got.Status.State)
	}
}

func TestListInstanceSnapshots(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/instances/12345/snapshots": `[
			{
				"id": "snapshot-123",
				"name": "test-snapshot-1",
				"description": "Test snapshot 1",
				"included_volumes": ["vol-1"],
				"status": {
					"state": "completed",
					"message": "Snapshot completed successfully",
					"volumes": [
						{
							"id": "vol-1",
							"state": "completed"
						}
					]
				},
				"created_at": "2024-03-20T10:00:00Z"
			},
			{
				"id": "snapshot-124",
				"name": "test-snapshot-2",
				"description": "Test snapshot 2",
				"included_volumes": ["vol-2"],
				"status": {
					"state": "completed",
					"message": "Snapshot completed successfully",
					"volumes": [
						{
							"id": "vol-2",
							"state": "completed"
						}
					]
				},
				"created_at": "2024-03-20T11:00:00Z"
			}
		]`,
	})
	defer server.Close()

	got, err := client.ListInstanceSnapshots("12345")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	if len(got) != 2 {
		t.Errorf("Expected 2 snapshots, got %d", len(got))
	}

	if got[0].ID != "snapshot-123" {
		t.Errorf("Expected %s, got %s", "snapshot-123", got[0].ID)
	}
	if got[0].Name != "test-snapshot-1" {
		t.Errorf("Expected %s, got %s", "test-snapshot-1", got[0].Name)
	}

	if got[1].ID != "snapshot-124" {
		t.Errorf("Expected %s, got %s", "snapshot-124", got[1].ID)
	}
	if got[1].Name != "test-snapshot-2" {
		t.Errorf("Expected %s, got %s", "test-snapshot-2", got[1].Name)
	}
}

func TestUpdateInstanceSnapshot(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/instances/12345/snapshots/snapshot-123": `{
			"id": "snapshot-123",
			"name": "updated-snapshot",
			"description": "Updated snapshot description",
			"included_volumes": ["vol-1", "vol-2"],
			"status": {
				"state": "completed",
				"message": "Snapshot completed successfully",
				"volumes": [
					{
						"id": "vol-1",
						"state": "completed"
					},
					{
						"id": "vol-2",
						"state": "completed"
					}
				]
			},
			"created_at": "2024-03-20T10:00:00Z"
		}`,
	})
	defer server.Close()

	config := &UpdateInstanceSnapshotConfig{
		Name:        "updated-snapshot",
		Description: "Updated snapshot description",
	}

	got, err := client.UpdateInstanceSnapshot("12345", "snapshot-123", config)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	if got.ID != "snapshot-123" {
		t.Errorf("Expected %s, got %s", "snapshot-123", got.ID)
	}
	if got.Name != "updated-snapshot" {
		t.Errorf("Expected %s, got %s", "updated-snapshot", got.Name)
	}
	if got.Description != "Updated snapshot description" {
		t.Errorf("Expected %s, got %s", "Updated snapshot description", got.Description)
	}
	if got.Status.State != "completed" {
		t.Errorf("Expected %s, got %s", "completed", got.Status.State)
	}
}

func TestDeleteInstanceSnapshot(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/instances/12345/snapshots/snapshot-123": `{
			"result": "success"
		}`,
	})
	defer server.Close()

	got, err := client.DeleteInstanceSnapshot("12345", "snapshot-123")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	if got.Result != "success" {
		t.Errorf("Expected %s, got %s", "success", got.Result)
	}
}
