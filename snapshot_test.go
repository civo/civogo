package civogo

import (
	"reflect"
	"testing"
)

func TestCreateSnapshot(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/snapshots/my-backup": `{
		  "id": "0ca69adc-ff39-4fc1-8f08-d91434e86fac",
		  "instance_id": "44aab548-61ca-11e5-860e-5cf9389be614",
		  "hostname": "server1.prod.example.com",
		  "template_id": "0b213794-d795-4483-8982-9f249c0262b9",
		  "openstack_snapshot_id": null,
		  "region": "lon1",
		  "name": "my-backup",
		  "safe": 1,
		  "size_gb": 0,
		  "state": "new",
		  "cron_timing": null,
		  "requested_at": null,
		  "completed_at": null
		}`,
	})
	defer server.Close()

	cfg := &SnapshotsConfig{
		InstanceID: "44aab548-61ca-11e5-860e-5cf9389be614",
		Safe:       true,
		Cron:       "",
	}
	got, err := client.CreateSnapshot("my-backup", cfg)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &Snapshot{
		ID:         "0ca69adc-ff39-4fc1-8f08-d91434e86fac",
		InstanceID: "44aab548-61ca-11e5-860e-5cf9389be614",
		Hostname:   "server1.prod.example.com",
		Template:   "0b213794-d795-4483-8982-9f249c0262b9",
		Region:     "lon1",
		Name:       "my-backup",
		Safe:       1,
		SizeGB:     0,
		State:      "new",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}
