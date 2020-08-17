package civogo

import (
	"reflect"
	"testing"
)

func TestListVolumes(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/volumes": `[{
			"id": "12345",
			"name": "my-volume",
			"instance_id": "null",
			"mountpoint": "null",
			"openstack_id": "null",
			"size_gb": 25,
			"bootable": false
		  }]`,
	})
	defer server.Close()
	got, err := client.ListVolumes()

	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	expected := []Volume{{ID: "12345", InstanceID: "null", Name: "my-volume", MountPoint: "null", SizeGigabytes: 25, Bootable: false}}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestFindVolume(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/volumes": `[
			{
				"id": "12345",
				"name": "my-volume",
				"instance_id": "null",
				"mountpoint": "null",
				"openstack_id": "null",
				"size_gb": 25,
				"bootable": false
			},
			{
				"id": "67890",
				"name": "other-volume",
				"instance_id": "null",
				"mountpoint": "null",
				"openstack_id": "null",
				"size_gb": 25,
				"bootable": false
			}
		]`,
	})
	defer server.Close()

	got, err := client.FindVolume("34")
	if got.ID != "12345" {
		t.Errorf("Expected %s, got %s", "12345", got.ID)
	}

	got, _ = client.FindVolume("89")
	if got.ID != "67890" {
		t.Errorf("Expected %s, got %s", "67890", got.ID)
	}

	got, _ = client.FindVolume("my")
	if got.ID != "12345" {
		t.Errorf("Expected %s, got %s", "12345", got.ID)
	}

	got, _ = client.FindVolume("other")
	if got.ID != "67890" {
		t.Errorf("Expected %s, got %s", "67890", got.ID)
	}

	_, err = client.FindVolume("volume")
	if err.Error() != "MultipleMatchesError: unable to find volume because there were multiple matches" {
		t.Errorf("Expected %s, got %s", "unable to find volume because there were multiple matches", err.Error())
	}

	_, err = client.FindVolume("missing")
	if err.Error() != "ZeroMatchesError: unable to find missing, zero matches" {
		t.Errorf("Expected %s, got %s", "unable to find missing, zero matches", err.Error())
	}
}

func TestNewVolume(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/volumes/": `{
			"id": "76cc107f-fbef-4e2b-b97f-f5d34f4075d3",
			"name": "my-volume",
			"result": "success"
		}`,
	})
	defer server.Close()

	cfg := &VolumeConfig{Name: "my-volume", SizeGigabytes: 25, Bootable: false}
	got, err := client.NewVolume(cfg)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &VolumeResult{
		ID:     "76cc107f-fbef-4e2b-b97f-f5d34f4075d3",
		Name:   "my-volume",
		Result: "success",
	}

	if expected.ID != got.ID {
		t.Errorf("Expected %s, got %s", expected.ID, got.ID)
	}

	if expected.Name != got.Name {
		t.Errorf("Expected %s, got %s", expected.Name, got.Name)
	}

	if expected.Result != got.Result {
		t.Errorf("Expected %s, got %s", expected.Result, got.Result)
	}
}

func TestResizeVolumes(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/volumes/12346/resize": `{"result": "success"}`,
	})
	defer server.Close()
	got, err := client.ResizeVolume("12346", 25)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &SimpleResponse{Result: "success"}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestAttachVolumes(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/volumes/12346/attach": `{"result": "success"}`,
	})
	defer server.Close()
	got, err := client.AttachVolume("12346", "123456")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &SimpleResponse{Result: "success"}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestDetachVolumes(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/volumes/12346/detach": `{"result": "success"}`,
	})
	defer server.Close()
	got, err := client.DetachVolume("12346")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &SimpleResponse{Result: "success"}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestDeleteVolumes(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/volumes/12346": `{"result": "success"}`,
	})
	defer server.Close()
	got, err := client.DeleteVolume("12346")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &SimpleResponse{Result: "success"}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}
