package civogo

import (
	"reflect"
	"testing"
)

func TestListLoadBalancers(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/loadbalancers": `[
		  {
		     "id": "7bb2c574-7b34-4de4-9111-4ac2b5653efa",
		     "name": "test-lb",
		     "algorithm": "round_robin",
				 "external_traffic_policy": "Cluster",
		     "backends": [
		       {
		         "ip": "192.168.1.3",
		         "protocol": "TCP",
		         "source_port": 80,
		         "target_port": 31579
		       },
		       {
		         "ip": "192.168.1.4",
		         "protocol": "TCP",
		         "source_port": 80,
		         "target_port": 31579
		       }
		     ],
		     "public_ip": "10.90.100.37",
		     "private_ip": "192.168.1.2",
		     "firewall_id": "d8e8bae3-5a76-435e-b96d-24955deb8792",
		     "state": "available"
		  }
		]`,
	})
	defer server.Close()

	got, err := client.ListLoadBalancers()
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := []LoadBalancer{
		{
			ID:                    "7bb2c574-7b34-4de4-9111-4ac2b5653efa",
			Name:                  "test-lb",
			Algorithm:             "round_robin",
			ExternalTrafficPolicy: "Cluster",
			Backends: []LoadBalancerBackend{
				{
					IP:         "192.168.1.3",
					Protocol:   "TCP",
					SourcePort: 80,
					TargetPort: 31579,
				},
				{
					IP:         "192.168.1.4",
					Protocol:   "TCP",
					SourcePort: 80,
					TargetPort: 31579,
				},
			},
			PublicIP:   "10.90.100.37",
			PrivateIP:  "192.168.1.2",
			FirewallID: "d8e8bae3-5a76-435e-b96d-24955deb8792",
			State:      "available",
		},
	}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestFindLoadBalancer(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/loadbalancers": `[
			{	"id": "16549d23-1957-4dea-a3d5-6282d19a4c00",	"name": "test1-lb",	"algorithm": "round_robin",	"external_traffic_policy": "Cluster",	"backends": [ {	"ip": "192.168.1.3", "protocol": "TCP",	"source_port": 80, "target_port": 31579	}, { "ip": "192.168.1.4", "protocol": "TCP", "source_port": 80,	"target_port": 31579}	],	"public_ip": "10.90.100.37",	"private_ip": "192.168.1.2",	"firewall_id": "73e6baab-a319-4ead-8286-19b796e08143",	"state": "available"},
			{	"id": "b421b075-ff18-48b4-a092-604bca968f49",	"name": "test2-lb",	"algorithm": "round_robin",	"external_traffic_policy": "Cluster",	"backends": [ {	"ip": "192.168.1.3", "protocol": "TCP",	"source_port": 80, "target_port": 31579	}, { "ip": "192.168.1.4", "protocol": "TCP", "source_port": 80,	"target_port": 31579}	],	"public_ip": "10.90.100.37",	"private_ip": "192.168.1.2",	"firewall_id": "85fc8018-71d7-4346-ae07-d964864b1569",	"state": "available"}
		]`,
	})
	defer server.Close()

	got, _ := client.FindLoadBalancer("16549d23")
	if got.ID != "16549d23-1957-4dea-a3d5-6282d19a4c00" {
		t.Errorf("Expected %s, got %s", "16549d23-1957-4dea-a3d5-6282d19a4c00", got.ID)
	}

	got, _ = client.FindLoadBalancer("a092")
	if got.ID != "b421b075-ff18-48b4-a092-604bca968f49" {
		t.Errorf("Expected %s, got %s", "b421b075-ff18-48b4-a092-604bca968f49", got.ID)
	}

	got, _ = client.FindLoadBalancer("test1")
	if got.ID != "16549d23-1957-4dea-a3d5-6282d19a4c00" {
		t.Errorf("Expected %s, got %s", "16549d23-1957-4dea-a3d5-6282d19a4c00", got.ID)
	}

	got, _ = client.FindLoadBalancer("test2")
	if got.ID != "b421b075-ff18-48b4-a092-604bca968f49" {
		t.Errorf("Expected %s, got %s", "b421b075-ff18-48b4-a092-604bca968f49", got.ID)
	}

	_, err := client.FindLoadBalancer("missing")
	if err.Error() != "ZeroMatchesError: unable to find missing, zero matches" {
		t.Errorf("Expected %s, got %s", "unable to find missing, zero matches", err.Error())
	}
}

func TestCreateLoadBalancer(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/loadbalancers": `{
			"id": "56dca3ae-ea3f-480f-9b25-abf90b439729",
			"name": "test-lb",
			"network_id": "b064d568-5869-427c-827a-77d48cde6a2e",
			"algorithm": "round_robin",
			"backends": [
				{
					"ip": "192.168.1.3",
					"protocol": "TCP",
					"source_port": 80,
					"target_port": 31579
				},
				{
					"ip": "192.168.1.4",
					"protocol": "TCP",
					"source_port": 80,
					"target_port": 31579
				}
			],
			"firewall_id": "9717fb32-dc0b-49e9-8265-5c84863e2164"
		}`,
	})
	defer server.Close()

	cfg := &LoadBalancerConfig{}
	got, err := client.CreateLoadBalancer(cfg)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &LoadBalancer{
		ID:        "56dca3ae-ea3f-480f-9b25-abf90b439729",
		Name:      "test-lb",
		Algorithm: "round_robin",
		Backends: []LoadBalancerBackend{
			{
				IP:         "192.168.1.3",
				Protocol:   "TCP",
				SourcePort: 80,
				TargetPort: 31579,
			},
			{
				IP:         "192.168.1.4",
				Protocol:   "TCP",
				SourcePort: 80,
				TargetPort: 31579,
			},
		},
		FirewallID: "9717fb32-dc0b-49e9-8265-5c84863e2164",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestUpdateLoadBalancer(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/loadbalancers/a1bd123c-b7e2-4d4f-9fda-7940c7e06b38": `{
			"id": "a1bd123c-b7e2-4d4f-9fda-7940c7e06b38",
			"name": "test-lb-updated",
			"network_id": "b064d568-5869-427c-827a-77d48cde6a2e",
			"algorithm": "round_robin",
			"external_traffic_policy": "Cluster",
			"backends": [
				{
					"ip": "192.168.1.3",
					"protocol": "TCP",
					"source_port": 80,
					"target_port": 31579
				},
				{
					"ip": "192.168.1.4",
					"protocol": "TCP",
					"source_port": 80,
					"target_port": 31579
				}
			],
			"firewall_id": "9717fb32-dc0b-49e9-8265-5c84863e2164"
		}`,
	})
	defer server.Close()

	cfg := &LoadBalancerUpdateConfig{
		Name: "test-lb-updated",
	}
	got, err := client.UpdateLoadBalancer("a1bd123c-b7e2-4d4f-9fda-7940c7e06b38", cfg)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &LoadBalancer{
		ID:                    "a1bd123c-b7e2-4d4f-9fda-7940c7e06b38",
		Name:                  "test-lb-updated",
		Algorithm:             "round_robin",
		ExternalTrafficPolicy: "Cluster",
		Backends: []LoadBalancerBackend{
			{
				IP:         "192.168.1.3",
				Protocol:   "TCP",
				SourcePort: 80,
				TargetPort: 31579,
			},
			{
				IP:         "192.168.1.4",
				Protocol:   "TCP",
				SourcePort: 80,
				TargetPort: 31579,
			},
		},
		FirewallID: "9717fb32-dc0b-49e9-8265-5c84863e2164",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestGetLoadBalancer(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/loadbalancers/a1bd123c-b7e2-4d4f-9fda-7940c7e06b38": `{
			"id": "a1bd123c-b7e2-4d4f-9fda-7940c7e06b38",
			"name": "test-lb-updated",
			"network_id": "b064d568-5869-427c-827a-77d48cde6a2e",
			"algorithm": "round_robin",
			"external_traffic_policy": "Cluster",
			"backends": [
				{
					"ip": "192.168.1.3",
					"protocol": "TCP",
					"source_port": 80,
					"target_port": 31579
				},
				{
					"ip": "192.168.1.4",
					"protocol": "TCP",
					"source_port": 80,
					"target_port": 31579
				}
			],
			"firewall_id": "9717fb32-dc0b-49e9-8265-5c84863e2164"
		}`,
	})
	defer server.Close()

	got, err := client.GetLoadBalancer("a1bd123c-b7e2-4d4f-9fda-7940c7e06b38")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &LoadBalancer{
		ID:                    "a1bd123c-b7e2-4d4f-9fda-7940c7e06b38",
		Name:                  "test-lb-updated",
		Algorithm:             "round_robin",
		ExternalTrafficPolicy: "Cluster",
		Backends: []LoadBalancerBackend{
			{
				IP:         "192.168.1.3",
				Protocol:   "TCP",
				SourcePort: 80,
				TargetPort: 31579,
			},
			{
				IP:         "192.168.1.4",
				Protocol:   "TCP",
				SourcePort: 80,
				TargetPort: 31579,
			},
		},
		FirewallID: "9717fb32-dc0b-49e9-8265-5c84863e2164",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}
func TestDeleteLoadBalancer(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/loadbalancers/12345": `{"result": "success"}`,
	})
	defer server.Close()

	got, err := client.DeleteLoadBalancer("12345")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &SimpleResponse{Result: "success"}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}
