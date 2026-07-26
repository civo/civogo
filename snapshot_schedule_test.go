package civogo

import (
	"testing"
)

func TestCreateSnapshotSchedule(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/resourcesnapshotschedules": `{
			"id": "schedule-123",
			"name": "daily-schedule",
			"description": "Daily snapshot schedule for instances",
			"cron_expression": "0 0 * * *",
			"paused": false,
			"retention": {
				"period": "48h",
				"max_snapshots": 7
			},
			"instances": [
				{
					"id": "instance-123",
					"size": "g3.small",
					"included_volumes": ["volume-1", "volume-2"]
				}
			],
			"status": {
				"state": "active",
				"last_snapshot": {
					"id": "",
					"name": "",
					"state": ""
				}
			},
			"created_at": "2025-04-01T09:11:14Z"
		}`})
	defer server.Close()

	createReq := &CreateSnapshotScheduleRequest{
		Name:           "daily-schedule",
		Description:    "Daily snapshot schedule for instances",
		CronExpression: "0 0 * * *",
		Retention: SnapshotRetention{
			Period:       "48h",
			MaxSnapshots: 7,
		},
		Instances: []CreateSnapshotInstance{
			{
				InstanceID:     "instance-123",
				IncludeVolumes: true,
			},
		},
	}

	got, err := client.CreateSnapshotSchedule(createReq)

	if err != nil {
		t.Errorf("CreateSnapshotSchedule returned an error: %s", err)
	}

	if got.Name != "daily-schedule" {
		t.Errorf("Expected Name daily-schedule, got %s", got.Name)
	}

	if got.CronExpression != "0 0 * * *" {
		t.Errorf("Expected CronExpression 0 0 * * *, got %s", got.CronExpression)
	}
}

func TestListSnapshotSchedules(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/resourcesnapshotschedules": `[
			{
				"id": "schedule-123",
				"name": "daily-schedule",
				"description": "Daily snapshot schedule",
				"cron_expression": "0 0 * * *",
				"paused": false,
				"retention": {
					"period": "48h",
					"max_snapshots": 7
				},
				"instances": [],
				"status": {
					"state": "active",
					"last_snapshot": {
						"id": "",
						"name": "",
						"state": ""
					}
				},
				"created_at": "2025-04-01T09:11:14Z"
			}
		]`})
	defer server.Close()

	got, err := client.ListSnapshotSchedules()

	if err != nil {
		t.Errorf("ListSnapshotSchedules returned an error: %s", err)
	}

	if len(got) != 1 {
		t.Errorf("Expected 1 snapshot schedule, got %d", len(got))
	}

	if got[0].Name != "daily-schedule" {
		t.Errorf("Expected Name daily-schedule, got %s", got[0].Name)
	}
}

func TestGetSnapshotSchedule(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/resourcesnapshotschedules/schedule-123": `{
			"id": "schedule-123",
			"name": "daily-schedule",
			"description": "Daily snapshot schedule",
			"cron_expression": "0 0 * * *",
			"paused": false,
			"retention": {
				"period": "48h",
				"max_snapshots": 7
			},
			"instances": [],
			"status": {
				"state": "active",
				"last_snapshot": {
					"id": "",
					"name": "",
					"state": ""
				}
			},
			"created_at": "2025-04-01T09:11:14Z"
		}`})
	defer server.Close()

	got, err := client.GetSnapshotSchedule("schedule-123")

	if err != nil {
		t.Errorf("GetSnapshotSchedule returned an error: %s", err)
	}

	if got.ID != "schedule-123" {
		t.Errorf("Expected ID schedule-123, got %s", got.ID)
	}
}

func TestDeleteSnapshotSchedule(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/resourcesnapshotschedules/schedule-123": `{"result": "success"}`})
	defer server.Close()

	got, err := client.DeleteSnapshotSchedule("schedule-123")

	if err != nil {
		t.Errorf("DeleteSnapshotSchedule returned an error: %s", err)
	}

	if got.Result != "success" {
		t.Errorf("Expected result success, got %s", got.Result)
	}
}
