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
	expected := []Volumes{{ID: "12345", InstanceID: "null", Name: "my-volume", MountPoint: "null", SizeGB: 25, Bootable: false}}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestNewVolumes(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/volumes/": `{
			"id": "76cc107f-fbef-4e2b-b97f-f5d34f4075d3",
			"name": "my-volume",
			"result": "success"
		}`,
	})
	defer server.Close()

	cfg := &VolumesConfig{Name: "my-volume", SizeGB: 25, Bootable: false}
	got, err := client.NewVolumes(cfg)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &VolumesResult{
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
	got, err := client.ResizeVolumes("12346", 25)
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
	got, err := client.AttachVolumes("12346", "123456")
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
	got, err := client.DetachVolumes("12346")
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
	got, err := client.DeleteVolumes("12346")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &SimpleResponse{Result: "success"}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}
