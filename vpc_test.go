package civogo

import (
	"reflect"
	"testing"
)

// =============================================================================
// VPC Networks Tests
// =============================================================================

func TestGetDefaultVPCNetwork(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/vpc/networks": `[{"id": "12345", "default": true, "name": "Default Network", "status": "Active"}]`,
	})
	defer server.Close()

	got, err := client.GetDefaultVPCNetwork()
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	if got.ID != "12345" {
		t.Errorf("Expected %s, got %s", "12345", got.ID)
	}
	if got.Status != "Active" {
		t.Errorf("Expected %s, got %s", "Active", got.Status)
	}
}

func TestGetVPCNetwork(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/vpc/networks/12345": `{"id": "12345", "name": "test-network"}`,
	})
	defer server.Close()

	got, err := client.GetVPCNetwork("12345")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	if got.ID != "12345" {
		t.Errorf("Expected %s, got %s", "12345", got.ID)
	}
	if got.Name != "test-network" {
		t.Errorf("Expected %s, got %s", "test-network", got.Name)
	}
}

func TestNewVPCNetwork(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/vpc/networks": `{
			"id": "76cc107f-fbef-4e2b-b97f-f5d34f4075d3",
			"label": "private-net",
			"result": "success"
		}`,
	})
	defer server.Close()

	got, err := client.NewVPCNetwork("private-net")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &NetworkResult{
		ID:     "76cc107f-fbef-4e2b-b97f-f5d34f4075d3",
		Label:  "private-net",
		Result: "success",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestCreateVPCNetwork(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/vpc/networks": `{
			"id": "41e4b4f5-5be0-4ac1-8c62-7e58f14f9155",
			"result": "success",
			"label": "private-net"
		}`,
	})
	defer server.Close()

	configs := NetworkConfig{
		Label: "private-net",
	}
	got, err := client.CreateVPCNetwork(configs)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &NetworkResult{
		ID:     "41e4b4f5-5be0-4ac1-8c62-7e58f14f9155",
		Label:  "private-net",
		Result: "success",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestListVPCNetworks(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/vpc/networks": `[{
			"id": "12345",
			"name": "my-net",
			"default": false,
			"cidr": "0.0.0.0/0",
			"label": "development",
			"status": "Deleting"
		  }]`,
	})
	defer server.Close()
	got, err := client.ListVPCNetworks()

	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	expected := []Network{{ID: "12345", Name: "my-net", Default: false, CIDR: "0.0.0.0/0", Label: "development", Status: "Deleting"}}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestFindVPCNetwork(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/vpc/networks": `[
			{
				"id": "12345",
				"name": "my-net",
				"region": "lon1",
				"default": false,
				"cidr": "0.0.0.0/0",
				"label": "development"
			},
			{
				"id": "67890",
				"name": "other-net",
				"region": "lon1",
				"default": false,
				"cidr": "0.0.0.0/0",
				"label": "production"
			}
			]`,
	})
	defer server.Close()

	got, _ := client.FindVPCNetwork("34")
	if got.ID != "12345" {
		t.Errorf("Expected %s, got %s", "12345", got.ID)
	}

	got, _ = client.FindVPCNetwork("89")
	if got.ID != "67890" {
		t.Errorf("Expected %s, got %s", "67890", got.ID)
	}

	got, _ = client.FindVPCNetwork("my")
	if got.ID != "12345" {
		t.Errorf("Expected %s, got %s", "12345", got.ID)
	}

	got, _ = client.FindVPCNetwork("other")
	if got.ID != "67890" {
		t.Errorf("Expected %s, got %s", "67890", got.ID)
	}

	got, _ = client.FindVPCNetwork("production")
	if got.ID != "67890" {
		t.Errorf("Expected %s, got %s", "67890", got.ID)
	}

	_, err := client.FindVPCNetwork("net")
	if err.Error() != "MultipleMatchesError: unable to find net because there were multiple matches" {
		t.Errorf("Expected %s, got %s", "unable to find net because there were multiple matches", err.Error())
	}

	_, err = client.FindVPCNetwork("missing")
	if err.Error() != "ZeroMatchesError: unable to find missing, zero matches" {
		t.Errorf("Expected %s, got %s", "ZeroMatchesError: unable to find missing, zero matches", err.Error())
	}
}

func TestRenameVPCNetwork(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/vpc/networks/76cc107f-fbef-4e2b-b97f-f5d34f4075d3": `{
			"id": "76cc107f-fbef-4e2b-b97f-f5d34f4075d3",
			"label": "private-net",
			"result": "success"
		}`,
	})
	defer server.Close()

	got, err := client.RenameVPCNetwork("private-net", "76cc107f-fbef-4e2b-b97f-f5d34f4075d3")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &NetworkResult{
		ID:     "76cc107f-fbef-4e2b-b97f-f5d34f4075d3",
		Label:  "private-net",
		Result: "success",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestUpdateVPCNetwork(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/vpc/networks/76cc107f-fbef-4e2b-b97f-f5d34f4075d3": `{
			"id": "76cc107f-fbef-4e2b-b97f-f5d34f4075d3",
			"label": "updated-net",
			"result": "success"
		}`,
	})
	defer server.Close()

	configs := NetworkConfig{
		Label: "updated-net",
	}
	got, err := client.UpdateVPCNetwork("76cc107f-fbef-4e2b-b97f-f5d34f4075d3", configs)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &NetworkResult{
		ID:     "76cc107f-fbef-4e2b-b97f-f5d34f4075d3",
		Label:  "updated-net",
		Result: "success",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestDeleteVPCNetwork(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/vpc/networks/12345": `{"result": "success"}`,
	})
	defer server.Close()

	got, err := client.DeleteVPCNetwork("12345")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &SimpleResponse{Result: "success"}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

// =============================================================================
// VPC Subnets Tests
// =============================================================================

func TestGetVPCSubnet(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/vpc/networks/12345/subnets/6789": `{"network_id": "12345", "id": "6789", "name": "test-subnet"}`,
	})
	defer server.Close()

	got, err := client.GetVPCSubnet("12345", "6789")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	if got.Name != "test-subnet" {
		t.Errorf("Expected %s, got %s", "test-subnet", got.Name)
	}
}

func TestFindVPCSubnet(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/vpc/networks/12345/subnets": `[
			{
				"id": "6789",
				"name": "test-subnet",
				"network_id": "12345"
			},
			{
				"id": "67890",
				"name": "test-subnet-2",
				"network_id": "12345"
			}
			]`,
	})
	defer server.Close()

	got, _ := client.FindVPCSubnet("6789", "12345")
	if got.ID != "6789" {
		t.Errorf("Expected %s, got %s", "6789", got.ID)
	}

	got, _ = client.FindVPCSubnet("test-subnet-2", "12345")
	if got.ID != "67890" {
		t.Errorf("Expected %s, got %s", "67890", got.ID)
	}
}

func TestCreateVPCSubnet(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/vpc/networks/12345/subnets": `{"id": "76cc107f-fbef-4e2b-b97f-f5d34f4075d3","network_id": "12345","name": "test-subnet","status": "success"}`,
	})
	defer server.Close()

	subnet := SubnetConfig{
		Name: "test-subnet",
	}

	got, err := client.CreateVPCSubnet("12345", subnet)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &Subnet{
		ID:        "76cc107f-fbef-4e2b-b97f-f5d34f4075d3",
		Name:      "test-subnet",
		NetworkID: "12345",
		Status:    "success",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestListVPCSubnets(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/vpc/networks/12345/subnets": `[{
			"id": "6789",
			"name": "test-subnet",
			"network_id": "12345"
		  }]`,
	})
	defer server.Close()
	got, err := client.ListVPCSubnets("12345")

	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	expected := []Subnet{{ID: "6789", Name: "test-subnet", NetworkID: "12345"}}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestDeleteVPCSubnet(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/vpc/networks/12345/subnets/6789": `{"result": "success"}`,
	})
	defer server.Close()

	got, err := client.DeleteVPCSubnet("12345", "6789")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &SimpleResponse{Result: "success"}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

// =============================================================================
// VPC Firewalls Tests
// =============================================================================

func TestListVPCFirewalls(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/vpc/firewalls": `[{"id": "12345", "name": "instance-123456", "rules_count": 3, "instance_count": 10, "cluster_count": 2, "loadbalancer_count": 1}, {"id": "67789", "name": "instance-7890", "rules_count": 1, "instance_count": 2, "cluster_count": 1, "loadbalancer_count": 0}]`,
	})
	defer server.Close()
	got, err := client.ListVPCFirewalls()

	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	expected := []Firewall{{ID: "12345", Name: "instance-123456", RulesCount: 3, InstanceCount: 10, ClusterCount: 2, LoadBalancerCount: 1}, {ID: "67789", Name: "instance-7890", RulesCount: 1, InstanceCount: 2, ClusterCount: 1, LoadBalancerCount: 0}}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestFindVPCFirewall(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/vpc/firewalls": `[{"id": "12345", "name": "web-instance", "rules_count": 3, "instance_count": 10, "cluster_count": 2, "loadbalancer_count": 1, "region": "lon1"}, {"id": "67789", "name": "web-node", "rules_count": 1, "instances_count": 2, "region": "lon1"}]`,
	})
	defer server.Close()

	got, _ := client.FindVPCFirewall("45")
	if got.ID != "12345" {
		t.Errorf("Expected %s, got %s", "12345", got.ID)
	}

	got, _ = client.FindVPCFirewall("89")
	if got.ID != "67789" {
		t.Errorf("Expected %s, got %s", "67789", got.ID)
	}

	got, _ = client.FindVPCFirewall("inst")
	if got.ID != "12345" {
		t.Errorf("Expected %s, got %s", "12345", got.ID)
	}

	got, _ = client.FindVPCFirewall("nod")
	if got.ID != "67789" {
		t.Errorf("Expected %s, got %s", "67789", got.ID)
	}

	_, err := client.FindVPCFirewall("web")
	if err.Error() != "MultipleMatchesError: unable to find web because there were multiple matches" {
		t.Errorf("Expected %s, got %s", "unable to find web because there were multiple matches", err.Error())
	}

	_, err = client.FindVPCFirewall("missing")
	if err.Error() != "ZeroMatchesError: unable to find missing, zero matches" {
		t.Errorf("Expected %s, got %s", "unable to find missing, zero matches", err.Error())
	}
}

func TestNewVPCFirewall(t *testing.T) {
	client, server, _ := NewAdvancedClientForTesting([]ConfigAdvanceClientForTesting{
		{
			Method: "POST",
			Value: []ValueAdvanceClientForTesting{
				{
					RequestBody:  `{"name":"fw-mail","region":"LON1","network_id":"1234-5698-9874-98","create_rules":true}`,
					URL:          "/v2/vpc/firewalls",
					ResponseBody: `{"id": "76cc107f-fbef-4e2b-b97f-f5d34f4075d3","name": "fw-mail","result": "success"}`,
				},
			},
		},
	})
	defer server.Close()

	CreateRules := true
	firewallConfig := &FirewallConfig{
		Name:        "fw-mail",
		NetworkID:   "1234-5698-9874-98",
		Region:      "LON1",
		CreateRules: &CreateRules,
	}
	got, err := client.NewVPCFirewall(firewallConfig)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &FirewallResult{
		ID:     "76cc107f-fbef-4e2b-b97f-f5d34f4075d3",
		Name:   "fw-mail",
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

func TestRenameVPCFirewall(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/vpc/firewalls/12346": `{"result": "success"}`,
	})
	defer server.Close()
	rename := &FirewallConfig{
		Name: "new_name",
	}
	got, err := client.RenameVPCFirewall("12346", rename)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &SimpleResponse{Result: "success"}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestDeleteVPCFirewall(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/vpc/firewalls/12346": `{"result": "success"}`,
	})
	defer server.Close()
	got, err := client.DeleteVPCFirewall("12346")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &SimpleResponse{Result: "success"}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

// =============================================================================
// VPC Firewall Rules Tests
// =============================================================================

func TestNewVPCFirewallRule(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/vpc/firewalls/78901/rules": `{
		  "id": "123456",
		  "firewall_id": "78901",
		  "protocol": "tcp",
		  "start_port": "443",
		  "end_port": "443",
		  "cidr": [
			"0.0.0.0/0"
		  ],
		  "direction": "ingress",
		  "action": "allow",
		  "label": null
		}`,
	})
	defer server.Close()

	cfg := &FirewallRuleConfig{FirewallID: "78901", Protocol: "tcp", StartPort: "443", EndPort: "443", Cidr: []string{"0.0.0.0/0"}, Direction: "inbound", Label: "", Action: "allow"}
	got, err := client.NewVPCFirewallRule(cfg)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &FirewallRule{
		ID:         "123456",
		FirewallID: "78901",
		Protocol:   "tcp",
		StartPort:  "443",
		EndPort:    "443",
		Cidr:       []string{"0.0.0.0/0"},
		Direction:  "ingress",
		Action:     "allow",
		Label:      "",
	}

	if expected.ID != got.ID {
		t.Errorf("Expected %s, got %s", expected.ID, got.ID)
	}

	if expected.FirewallID != got.FirewallID {
		t.Errorf("Expected %s, got %s", expected.FirewallID, got.FirewallID)
	}

	if expected.Protocol != got.Protocol {
		t.Errorf("Expected %s, got %s", expected.Protocol, got.Protocol)
	}

	if !reflect.DeepEqual(expected.Cidr, got.Cidr) {
		t.Errorf("Expected %q, got %q", expected.Cidr, got.Cidr)
	}
}

func TestFindVPCFirewallRule(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/vpc/firewalls/22/rules": `[{
			"id": "21",
			"firewall_id": "22",
			"protocol": "tcp",
			"start_port": "443",
			"end_port": "443",
			"cidr": [
			  "0.0.0.0/0"
			],
			"direction": "ingress",
			"action": "allow",
			"label": "My Rule"
		  },{
			"id": "22",
			"firewall_id": "22",
			"protocol": "tcp",
			"start_port": "80",
			"end_port": "80",
			"cidr": [
			  "0.0.0.0/0"
			],
			"direction": "ingress",
			"action": "allow",
			"label": "My Rule"
		  }]`,
	})
	defer server.Close()

	got, _ := client.FindVPCFirewallRule("22", "21")
	if got.ID != "21" {
		t.Errorf("Expected %s, got %s", "21", got.ID)
	}

	got, _ = client.FindVPCFirewallRule("22", "22")
	if got.ID != "22" {
		t.Errorf("Expected %s, got %s", "22", got.ID)
	}

	_, err := client.FindVPCFirewallRule("22", "2")
	if err.Error() != "MultipleMatchesError: unable to find 2 because there were multiple matches" {
		t.Errorf("Expected %s, got %s", "unable to find 2 because there were multiple matches", err.Error())
	}

	_, err = client.FindVPCFirewallRule("22", "missing")
	if err.Error() != "ZeroMatchesError: unable to find missing, zero matches" {
		t.Errorf("Expected %s, got %s", "unable to find missing, zero matches", err.Error())
	}
}

func TestListVPCFirewallRules(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/vpc/firewalls/22/rules": `[{
			"id": "1",
			"firewall_id": "22",
			"protocol": "tcp",
			"start_port": "443",
			"end_port": "443",
			"cidr": [
			  "0.0.0.0/0"
			],
			"direction": "ingress",
			"action": "allow",
			"label": "My Rule"
		  },{
			"id": "2",
			"firewall_id": "22",
			"protocol": "tcp",
			"start_port": "80",
			"end_port": "80",
			"cidr": [
			  "0.0.0.0/0"
			],
			"direction": "ingress",
			"action": "allow",
			"label": "My Rule"
		  }]`,
	})
	defer server.Close()
	got, err := client.ListVPCFirewallRules("22")

	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	expected := []FirewallRule{{ID: "1", FirewallID: "22", Protocol: "tcp", StartPort: "443", EndPort: "443", Cidr: []string{"0.0.0.0/0"}, Direction: "ingress", Label: "My Rule", Action: "allow"}, {ID: "2", FirewallID: "22", Protocol: "tcp", StartPort: "80", EndPort: "80", Cidr: []string{"0.0.0.0/0"}, Direction: "ingress", Label: "My Rule", Action: "allow"}}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestDeleteVPCFirewallRule(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/vpc/firewalls/12346/rules/12345": `{"result": "success"}`,
	})
	defer server.Close()
	got, err := client.DeleteVPCFirewallRule("12346", "12345")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &SimpleResponse{Result: "success"}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

// =============================================================================
// VPC DNS Domains Tests
// =============================================================================

func TestListVPCDNSDomains(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/vpc/dns": `[{"id": "12345", "account_id": "1", "name": "example.com"}]`,
	})
	defer server.Close()

	got, err := client.ListVPCDNSDomains()
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	if len(got) != 1 {
		t.Errorf("Expected 1 domain, got %d", len(got))
		return
	}

	if got[0].Name != "example.com" {
		t.Errorf("Expected %s, got %s", "example.com", got[0].Name)
	}
}

func TestFindVPCDNSDomain(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/vpc/dns": `[{"id": "12345", "account_id": "1", "name": "example.com"}, {"id": "67890", "account_id": "1", "name": "test.com"}]`,
	})
	defer server.Close()

	got, _ := client.FindVPCDNSDomain("example")
	if got.ID != "12345" {
		t.Errorf("Expected %s, got %s", "12345", got.ID)
	}

	_, err := client.FindVPCDNSDomain(".com")
	if err.Error() != "MultipleMatchesError: unable to find .com because there were multiple matches" {
		t.Errorf("Expected multiple matches error, got %s", err.Error())
	}

	_, err = client.FindVPCDNSDomain("missing")
	if err.Error() != "ZeroMatchesError: unable to find missing, zero matches" {
		t.Errorf("Expected zero matches error, got %s", err.Error())
	}
}

func TestCreateVPCDNSDomain(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/vpc/dns": `{"id": "12345", "account_id": "1", "name": "example.com"}`,
	})
	defer server.Close()

	got, err := client.CreateVPCDNSDomain("example.com")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	if got.Name != "example.com" {
		t.Errorf("Expected %s, got %s", "example.com", got.Name)
	}
}

func TestGetVPCDNSDomain(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/vpc/dns": `[{"id": "12345", "account_id": "1", "name": "example.com"}]`,
	})
	defer server.Close()

	got, err := client.GetVPCDNSDomain("example.com")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	if got.ID != "12345" {
		t.Errorf("Expected %s, got %s", "12345", got.ID)
	}
}

func TestUpdateVPCDNSDomain(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/vpc/dns/12345": `{"id": "12345", "account_id": "1", "name": "newdomain.com"}`,
	})
	defer server.Close()

	domain := &DNSDomain{ID: "12345", AccountID: "1", Name: "example.com"}
	got, err := client.UpdateVPCDNSDomain(domain, "newdomain.com")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	if got.Name != "newdomain.com" {
		t.Errorf("Expected %s, got %s", "newdomain.com", got.Name)
	}
}

func TestDeleteVPCDNSDomain(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/vpc/dns/12345": `{"result": "success"}`,
	})
	defer server.Close()

	domain := &DNSDomain{ID: "12345", AccountID: "1", Name: "example.com"}
	got, err := client.DeleteVPCDNSDomain(domain)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &SimpleResponse{Result: "success"}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

// =============================================================================
// VPC DNS Records Tests
// =============================================================================

func TestCreateVPCDNSRecord(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/vpc/dns/12345/records": `{"id": "record123", "domain_id": "12345", "name": "www", "value": "1.2.3.4", "type": "A"}`,
	})
	defer server.Close()

	cfg := &DNSRecordConfig{Type: DNSRecordTypeA, Name: "www", Value: "1.2.3.4"}
	got, err := client.CreateVPCDNSRecord("12345", cfg)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	if got.Name != "www" {
		t.Errorf("Expected %s, got %s", "www", got.Name)
	}
}

func TestListVPCDNSRecords(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/vpc/dns/12345/records": `[{"id": "record123", "domain_id": "12345", "name": "www", "value": "1.2.3.4", "type": "A"}]`,
	})
	defer server.Close()

	got, err := client.ListVPCDNSRecords("12345")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	if len(got) != 1 {
		t.Errorf("Expected 1 record, got %d", len(got))
	}
}

func TestGetVPCDNSRecord(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/vpc/dns/12345/records": `[{"id": "record123", "domain_id": "12345", "name": "www", "value": "1.2.3.4", "type": "A"}]`,
	})
	defer server.Close()

	got, err := client.GetVPCDNSRecord("12345", "record123")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	if got.ID != "record123" {
		t.Errorf("Expected %s, got %s", "record123", got.ID)
	}
}

func TestUpdateVPCDNSRecord(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/vpc/dns/12345/records/record123": `{"id": "record123", "domain_id": "12345", "name": "www", "value": "5.6.7.8", "type": "A"}`,
	})
	defer server.Close()

	record := &DNSRecord{ID: "record123", DNSDomainID: "12345", Name: "www"}
	cfg := &DNSRecordConfig{Type: DNSRecordTypeA, Name: "www", Value: "5.6.7.8"}
	got, err := client.UpdateVPCDNSRecord(record, cfg)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	if got.Value != "5.6.7.8" {
		t.Errorf("Expected %s, got %s", "5.6.7.8", got.Value)
	}
}

func TestDeleteVPCDNSRecord(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/vpc/dns/12345/records/record123": `{"result": "success"}`,
	})
	defer server.Close()

	record := &DNSRecord{ID: "record123", DNSDomainID: "12345", Name: "www"}
	got, err := client.DeleteVPCDNSRecord(record)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &SimpleResponse{Result: "success"}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

// =============================================================================
// VPC Load Balancers Tests
// =============================================================================

func TestListVPCLoadBalancers(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/vpc/loadbalancers": `[{"id": "12345", "name": "my-lb", "algorithm": "round_robin", "public_ip": "1.2.3.4", "state": "available"}]`,
	})
	defer server.Close()

	got, err := client.ListVPCLoadBalancers()
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	if len(got) != 1 {
		t.Errorf("Expected 1 load balancer, got %d", len(got))
		return
	}

	if got[0].Name != "my-lb" {
		t.Errorf("Expected %s, got %s", "my-lb", got[0].Name)
	}
}

func TestGetVPCLoadBalancer(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/vpc/loadbalancers/12345": `{"id": "12345", "name": "my-lb", "algorithm": "round_robin", "public_ip": "1.2.3.4", "state": "available"}`,
	})
	defer server.Close()

	got, err := client.GetVPCLoadBalancer("12345")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	if got.ID != "12345" {
		t.Errorf("Expected %s, got %s", "12345", got.ID)
	}
}

func TestFindVPCLoadBalancer(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/vpc/loadbalancers": `[{"id": "12345", "name": "my-lb", "algorithm": "round_robin", "public_ip": "1.2.3.4", "state": "available"}, {"id": "67890", "name": "other-lb", "algorithm": "round_robin", "public_ip": "5.6.7.8", "state": "available"}]`,
	})
	defer server.Close()

	got, _ := client.FindVPCLoadBalancer("my")
	if got.ID != "12345" {
		t.Errorf("Expected %s, got %s", "12345", got.ID)
	}

	_, err := client.FindVPCLoadBalancer("lb")
	if err.Error() != "MultipleMatchesError: unable to find lb because there were multiple matches" {
		t.Errorf("Expected multiple matches error, got %s", err.Error())
	}

	_, err = client.FindVPCLoadBalancer("missing")
	if err.Error() != "ZeroMatchesError: unable to find missing, zero matches" {
		t.Errorf("Expected zero matches error, got %s", err.Error())
	}
}

func TestCreateVPCLoadBalancer(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/vpc/loadbalancers": `{"id": "12345", "name": "my-lb", "algorithm": "round_robin", "public_ip": "1.2.3.4", "state": "available"}`,
	})
	defer server.Close()

	cfg := &LoadBalancerConfig{Name: "my-lb", Region: "LON1"}
	got, err := client.CreateVPCLoadBalancer(cfg)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	if got.Name != "my-lb" {
		t.Errorf("Expected %s, got %s", "my-lb", got.Name)
	}
}

func TestUpdateVPCLoadBalancer(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/vpc/loadbalancers/12345": `{"id": "12345", "name": "updated-lb", "algorithm": "round_robin", "public_ip": "1.2.3.4", "state": "available"}`,
	})
	defer server.Close()

	cfg := &LoadBalancerUpdateConfig{Name: "updated-lb", Region: "LON1"}
	got, err := client.UpdateVPCLoadBalancer("12345", cfg)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	if got.Name != "updated-lb" {
		t.Errorf("Expected %s, got %s", "updated-lb", got.Name)
	}
}

func TestDeleteVPCLoadBalancer(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/vpc/loadbalancers/12345": `{"result": "success"}`,
	})
	defer server.Close()

	got, err := client.DeleteVPCLoadBalancer("12345")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &SimpleResponse{Result: "success"}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

// =============================================================================
// VPC Reserved IPs Tests
// =============================================================================

func TestListVPCIPs(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/vpc/ips": `{"page": 1, "per_page": 20, "pages": 1, "items": [{"id": "12345", "name": "my-ip", "ip": "1.2.3.4"}]}`,
	})
	defer server.Close()

	got, err := client.ListVPCIPs()
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	if len(got.Items) != 1 {
		t.Errorf("Expected 1 IP, got %d", len(got.Items))
		return
	}

	if got.Items[0].IP != "1.2.3.4" {
		t.Errorf("Expected %s, got %s", "1.2.3.4", got.Items[0].IP)
	}
}

func TestGetVPCIP(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/vpc/ips/12345": `{"id": "12345", "name": "my-ip", "ip": "1.2.3.4"}`,
	})
	defer server.Close()

	got, err := client.GetVPCIP("12345")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	if got.ID != "12345" {
		t.Errorf("Expected %s, got %s", "12345", got.ID)
	}
}

func TestFindVPCIP(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/vpc/ips": `{"page": 1, "per_page": 20, "pages": 1, "items": [{"id": "12345", "name": "my-ip", "ip": "1.2.3.4"}, {"id": "67890", "name": "other-ip", "ip": "5.6.7.8"}]}`,
	})
	defer server.Close()

	got, _ := client.FindVPCIP("1.2.3.4")
	if got.ID != "12345" {
		t.Errorf("Expected %s, got %s", "12345", got.ID)
	}

	_, err := client.FindVPCIP("ip")
	if err.Error() != "MultipleMatchesError: unable to find ip because there were multiple matches" {
		t.Errorf("Expected multiple matches error, got %s", err.Error())
	}

	_, err = client.FindVPCIP("missing")
	if err.Error() != "ZeroMatchesError: unable to find missing, zero matches" {
		t.Errorf("Expected zero matches error, got %s", err.Error())
	}
}

func TestNewVPCIP(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/vpc/ips": `{"id": "12345", "name": "my-ip", "ip": "1.2.3.4"}`,
	})
	defer server.Close()

	cfg := &CreateIPRequest{Name: "my-ip", Region: "LON1"}
	got, err := client.NewVPCIP(cfg)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	if got.Name != "my-ip" {
		t.Errorf("Expected %s, got %s", "my-ip", got.Name)
	}
}

func TestUpdateVPCIP(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/vpc/ips/12345": `{"id": "12345", "name": "updated-ip", "ip": "1.2.3.4"}`,
	})
	defer server.Close()

	cfg := &UpdateIPRequest{Name: "updated-ip", Region: "LON1"}
	got, err := client.UpdateVPCIP("12345", cfg)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	if got.Name != "updated-ip" {
		t.Errorf("Expected %s, got %s", "updated-ip", got.Name)
	}
}

func TestDeleteVPCIP(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/vpc/ips/12345": `{"result": "success"}`,
	})
	defer server.Close()

	got, err := client.DeleteVPCIP("12345")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &SimpleResponse{Result: "success"}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}
