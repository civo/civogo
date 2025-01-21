package civogo

import (
	"reflect"
	"testing"
)

func TestListVolumeSnapshots(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/snapshots?region=TEST&resource_type=volume": `[{
			"name": "test-snapshot",
			"snapshot_id": "12345",
			"snapshot_description": "snapshot for test",
			"volume_id": "12345",
			"source_volume_name": "test-volume",
			"instance_id": "ins1234",
			"restore_size": 20,
			"state": "Ready",
			"creation_time": "2020-01-01T00:00:00Z"
		}]`,
	})
	defer server.Close()

	got, err := client.ListVolumeSnapshots()

	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := []VolumeSnapshot{
		{
			Name:                "test-snapshot",
			SnapshotID:          "12345",
			SnapshotDescription: "snapshot for test",
			VolumeID:            "12345",
			SourceVolumeName:    "test-volume",
			InstanceID:          "ins1234",
			RestoreSize:         20,
			State:               "Ready",
			CreationTime:        "2020-01-01T00:00:00Z",
		},
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestGetVolumeSnapshot(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/snapshots/snapshot-uuid?region=TEST&resource_type=volume": `{
			"name": "test-snapshot",
			"snapshot_id": "snapshot-uuid",
			"snapshot_description": "snapshot for testing",
			"volume_id": "12345",
			"source_volume_name": "test-volume",
			"instance_id": "ins1234",
			"restore_size": 20,
			"state": "Ready",
			"creation_time": "2020-01-01T00:00:00Z"
		}`,
	})
	defer server.Close()
	got, err := client.GetVolumeSnapshot("snapshot-uuid")

	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &VolumeSnapshot{
		Name:                "test-snapshot",
		SnapshotID:          "snapshot-uuid",
		SnapshotDescription: "snapshot for testing",
		VolumeID:            "12345",
		SourceVolumeName:    "test-volume",
		InstanceID:          "ins1234",
		RestoreSize:         20,
		State:               "Ready",
		CreationTime:        "2020-01-01T00:00:00Z",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestDeleteVolumeSnapshot(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/snapshots/12346": `{"result": "success"}`,
	})
	defer server.Close()
	got, err := client.DeleteVolumeSnapshot("12346")

	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &SimpleResponse{Result: "success"}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}
