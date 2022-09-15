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

func TestListVolumesForCluster(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{

		"/v2/kubernetes/clusters": `{"page":1,"per_page":20,"pages":1,"items":[
			{ "id": "69a23478-a89e-41d2-97b1-6f4c341cee70", "name": "your-first-cluster-name", "version": "2", "status": "ACTIVE", "ready": true, "num_target_nodes": 1, "target_nodes_size": "g2.xsmall", "built_at": "2019-09-23T13:04:23.000+01:00", "kubeconfig": "YAML_VERSION_OF_KUBECONFIG_HERE\n", "kubernetes_version": "0.8.1", "api_endpoint": "https://your.cluster.ip.address:6443", "dns_entry": "69a23478-a89e-41d2-97b1-6f4c341cee70.k8s.civo.com", "tags": [], "created_at": "2019-09-23T13:02:59.000+01:00", "firewall_id": "42118911-44c2-4cab-ad77-bcae062815b3", "instances": [{ "hostname": "kube-master-HEXDIGITS", "size": "g2.xsmall", "region": "lon1", "created_at": "2019-09-23T13:03:00.000+01:00", "status": "ACTIVE", "firewall_id": "5f0ba9ed-5ca7-4e14-9a09-449a84196d64", "public_ip": "your.cluster.ip.address", "tags": ["civo-kubernetes:installed", "civo-kubernetes:master"] }], "installed_applications": [{ "application": "Traefik", "title": null, "version": "(default)", "dependencies": null, "maintainer": "@Rancher_Labs", "description": "A reverse proxy/load-balancer that's easy, dynamic, automatic, fast, full-featured, open source, production proven and provides metrics.", "post_install": "Some documentation here\n", "installed": true, "url": "https://traefik.io", "category": "architecture", "updated_at": "2019-09-23T13:02:59.000+01:00", "image_url": "https://api.civo.com/k3s-marketplace/traefik.png", "plan": null, "configuration": {} }] }
		]}`,
		"/v2/volumes": `[{ "id": "12345", "name": "my-volume", "size_gb": 25, "bootable": false, "cluster_id": "69a23478-a89e-41d2-97b1-6f4c341cee70"}]`,
	})
	defer server.Close()
	got, err := client.ListVolumesForCluster("69a23478-a89e-41d2-97b1-6f4c341cee70")

	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	expected := []Volume{{ID: "12345", Name: "my-volume", SizeGigabytes: 25, Bootable: false, ClusterID: "69a23478-a89e-41d2-97b1-6f4c341cee70"}}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestListDanglingVolumes(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/volumes": `[{
			"id": "12345",
			"name": "my-volume",
			"cluster_id": "69a23478-a89e-41d2-97b1-6f4c341cee70",
			"size_gb": 25,
			"bootable": false
		  },
		  {
		  	"id": "34567",
			"name": "my-volume-two",
			"size_gb": 25,
			"bootable": false
		  }]`,
	})
	defer server.Close()
	got, err := client.ListDanglingVolumes()

	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	// Should only return the dangling volume (Vol with a cluster id but cluster doesn't exist)
	expected := []Volume{{ID: "12345", Name: "my-volume", ClusterID: "69a23478-a89e-41d2-97b1-6f4c341cee70", SizeGigabytes: 25, Bootable: false}}
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

	got, _ := client.FindVolume("34")
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

	_, err := client.FindVolume("volume")
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
		"/v2/volumes": `{
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

func TestResizeVolume(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/volumes/12346/resize": `{"result": "success"}`,
	})
	defer server.Close()
	got, err := client.ResizeVolume("12346", 20)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &SimpleResponse{Result: "success"}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}
