package civogo

import (
	"reflect"
	"testing"

	corev1 "k8s.io/api/core/v1"
)

func TestCreateKubernetesClusterPool(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/kubernetes/clusters/e733ea47-cc80-443b-b3f2-cccfe7a61ef5/pools": `{"result": "success"}`,
	})
	defer server.Close()

	newPool := &KubernetesClusterPoolUpdateConfig{
		ID:     "8a849cc5-bd51-45ce-814a-c378b09dcb06",
		Count:  &[]int{3}[0],
		Size:   "g4s.kube.small",
		Labels: map[string]string{},
		Taints: []corev1.Taint{},
		Region: "LON1",
	}

	got, err := client.CreateKubernetesClusterPool("e733ea47-cc80-443b-b3f2-cccfe7a61ef5", newPool)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &SimpleResponse{Result: "success"}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}
func TestGetKubernetesPool(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/kubernetes/clusters/e733ea47-cc80-443b-b3f2-cccfe7a61ef5/pools/fad8638d-efac-41e0-8787-23d37d845685": `{
			"id": "fad8638d-efac-41e0-8787-23d37d845685",
			"size": "g4s.kube.small",
			"count": 1,
			"instance_names": [
					"k3s-second-cluster-5359-217631-node-pool-b786-3ez7r"
			],
			"instances": [
					{
							"id": "15e0a7dd-744e-4558-8113-2192e4eca040",
							"name": "k3s-second-cluster-5359-217631-node-pool-b786-3ez7r",
							"hostname": "k3s-second-cluster-5359-217631-node-pool-b786-3ez7r",
							"account_id": "eaef1dd6-1cec-4d9c-8480-96452bd94dea",
							"size": "g4s.kube.small",
							"firewall_id": "57416f92-1f81-4111-b87d-77ba77724e5b",
							"source_type": "diskimage",
							"source_id": "1.26.4-k3s1",
							"network_id": "01d7d363-93c5-4203-9f8f-39ea9e7d0098",
							"initial_user": "root",
							"initial_password": "IXnb@e#FGPB95@H6iHJL",
							"ssh_key": "-",
							"tags": [
									"k3s"
							],
							"script": "",
							"status": "ACTIVE",
							"civostatsd_token": "d14047b3-27ae-45da-ae11-f0ff906eb4d8",
							"public_ip": "74.220.20.152",
							"private_ip": "192.168.1.16",
							"namespace_id": "cust-default-eaef1dd6-7505",
							"notes": "",
							"reverse_dns": "",
							"cpu_cores": 1,
							"ram_mb": 2048,
							"disk_gb": 40,
							"created_at": "2023-09-15T10:48:13Z"
					}
			]
	}`,
	})
	defer server.Close()

	got, err := client.GetKubernetesClusterPool("e733ea47-cc80-443b-b3f2-cccfe7a61ef5", "fad8638d-efac-41e0-8787-23d37d845685")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := KubernetesPool{
		ID:               "fad8638d-efac-41e0-8787-23d37d845685",
		Count:            1,
		Size:             "g4s.kube.small",
		InstanceNames:    []string{"k3s-second-cluster-5359-217631-node-pool-b786-3ez7r"},
		Instances:        []KubernetesInstance{},
		Labels:           map[string]string{},
		Taints:           []corev1.Taint{},
		PublicIPNodePool: false,
	}

	if got.ID != expected.ID {
		t.Errorf("Expected %+v, got %+v", expected.ID, got.ID)
		return
	}

	if got.Count != expected.Count {
		t.Errorf("Expected %+v, got %+v", expected.Count, got.Count)
		return
	}

	if got.Size != expected.Size {
		t.Errorf("Expected %+v, got %+v", expected.Size, got.Size)
		return
	}

	if len(got.InstanceNames) != len(expected.InstanceNames) {
		t.Errorf("Expected %+v, got %+v", len(expected.InstanceNames), len(got.InstanceNames))
		return
	}
}

func TestDeleteKubernetesPool(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/kubernetes/clusters/e733ea47-cc80-443b-b3f2-cccfe7a61ef5/pools/fad8638d-efac-41e0-8787-23d37d845685": `{"result": "success"}`,
	})
	defer server.Close()

	got, err := client.DeleteKubernetesClusterPool("e733ea47-cc80-443b-b3f2-cccfe7a61ef5", "fad8638d-efac-41e0-8787-23d37d845685")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &SimpleResponse{Result: "success"}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestUpdateKubernetesPool(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/kubernetes/clusters/e733ea47-cc80-443b-b3f2-cccfe7a61ef5/pools/fad8638d-efac-41e0-8787-23d37d845685": `{
			"id": "fad8638d-efac-41e0-8787-23d37d845685",
			"size": "g4s.kube.small",
			"count": 3,
			"instance_names": [
					"k3s-second-cluster-5359-217631-node-pool-b786-3ez7r",
					"k3s-second-cluster-5359-217631-node-pool-b786-3ez7r",
					"k3s-second-cluster-5359-217631-node-pool-b786-3ez7r"
			],
			"labels": {
					"label1": "label-result"
			},
			"taints": [
				{
					"key": "app",
					"value": "frontend",
					"effect": "NoSchedule"
				}
			],
			"instances": [
					{
							"id": "15e0a7dd-744e-4558-8113-2192e4eca040",
							"name": "k3s-second-cluster-5359-217631-node-pool-b786-3ez7r",
							"hostname": "k3s-second-cluster-5359-217631-node-pool-b786-3ez7r",
							"account_id": "eaef1dd6-1cec-4d9c-8480-96452bd94dea",
							"size": "g4s.kube.small",
							"firewall_id": "57416f92-1f81-4111-b87d-77ba77724e5b",
							"source_type": "diskimage",
							"source_id": "1.26.4-k3s1",
							"network_id": "01d7d363-93c5-4203-9f8f-39ea9e7d0098",
							"initial_user": "root",
							"initial_password": "IXnb@e#FGPB95@H6iHJL",
							"ssh_key": "-",
							"tags": [
									"k3s"
							],
							"script": "",
							"status": "ACTIVE",
							"civostatsd_token": "d14047b3-27ae-45da-ae11-f0ff906eb4d8",
							"public_ip": "74.220.20.152",
							"private_ip": "192.168.1.16",
							"namespace_id": "cust-default-eaef1dd6-7505",
							"notes": "",
							"reverse_dns": "",
							"cpu_cores": 1,
							"ram_mb": 2048,
							"disk_gb": 40,
							"created_at": "2023-09-15T10:48:13Z"
					}
			]
	}`,
	})
	defer server.Close()

	updatePool := &KubernetesClusterPoolUpdateConfig{
		Count: &[]int{3}[0],
		Labels: map[string]string{
			"label1": "label",
		},
		Taints: []corev1.Taint{
			{
				Key:    "app",
				Value:  "frontend",
				Effect: "NoSchedule",
			},
		},
	}

	got, err := client.UpdateKubernetesClusterPool("e733ea47-cc80-443b-b3f2-cccfe7a61ef5", "fad8638d-efac-41e0-8787-23d37d845685", updatePool)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := KubernetesPool{
		ID:            "fad8638d-efac-41e0-8787-23d37d845685",
		Count:         3,
		Size:          "g4s.kube.small",
		InstanceNames: []string{"k3s-second-cluster-5359-217631-node-pool-b786-3ez7r", "k3s-second-cluster-5359-217631-node-pool-b786-3ez7r", "k3s-second-cluster-5359-217631-node-pool-b786-3ez7r"},
		Instances:     []KubernetesInstance{},
		Labels: map[string]string{
			"label1": "label-result",
		},
		Taints: []corev1.Taint{
			{
				Key:    "app",
				Value:  "frontend",
				Effect: "NoSchedule",
			},
		},
		PublicIPNodePool: false,
	}

	if got.ID != expected.ID {
		t.Errorf("Expected %+v, got %+v", expected.ID, got.ID)
		return
	}

	if got.Count != expected.Count {
		t.Errorf("Expected %+v, got %+v", expected.Count, got.Count)
		return
	}

	if got.Size != expected.Size {
		t.Errorf("Expected %+v, got %+v", expected.Size, got.Size)
		return
	}

	if len(got.InstanceNames) != len(expected.InstanceNames) {
		t.Errorf("Expected %+v, got %+v", len(expected.InstanceNames), len(got.InstanceNames))
		return
	}

	if !reflect.DeepEqual(got.Labels, expected.Labels) {
		t.Errorf("Expected %+v, got %+v", expected.Labels, got.Labels)
	}

	if !reflect.DeepEqual(got.Taints, expected.Taints) {
		t.Errorf("Expected %+v, got %+v", expected.Taints, got.Taints)
	}
}

func TestDeleteKubernetesPoolInstance(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/kubernetes/clusters/e733ea47-cc80-443b-b3f2-cccfe7a61ef5/pools/fad8638d-efac-41e0-8787-23d37d845685/instances/15e0a7dd-744e-4558-8113-2192e4eca040": `{"result": "success"}`,
	})
	defer server.Close()

	got, err := client.DeleteKubernetesClusterPoolInstance("e733ea47-cc80-443b-b3f2-cccfe7a61ef5", "fad8638d-efac-41e0-8787-23d37d845685", "15e0a7dd-744e-4558-8113-2192e4eca040")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &SimpleResponse{Result: "success"}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}
