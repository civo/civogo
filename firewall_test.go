package civogo

import (
	"reflect"
	"testing"
)

func TestListFirewalls(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/firewalls": `[{"id": "12345", "name": "instance-123456", "rules_count": 3, "instances_count": 10, "region": "lon1"}, {"id": "67789", "name": "instance-7890", "rules_count": 1, "instances_count": 2, "region": "lon1"}]`,
	})
	defer server.Close()
	got, err := client.ListFirewalls()

	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	expected := []Firewall{{ID: "12345", Name: "instance-123456", RulesCount: 3, InstancesCount: 10, Region: "lon1"}, {ID: "67789", Name: "instance-7890", RulesCount: 1, InstancesCount: 2, Region: "lon1"}}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestFindFirewall(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/firewalls": `[{"id": "12345", "name": "web-instance", "rules_count": 3, "instances_count": 10, "region": "lon1"}, {"id": "67789", "name": "web-node", "rules_count": 1, "instances_count": 2, "region": "lon1"}]`,
	})
	defer server.Close()

	got, _ := client.FindFirewall("45")
	if got.ID != "12345" {
		t.Errorf("Expected %s, got %s", "12345", got.ID)
	}

	got, _ = client.FindFirewall("89")
	if got.ID != "67789" {
		t.Errorf("Expected %s, got %s", "67789", got.ID)
	}

	got, _ = client.FindFirewall("inst")
	if got.ID != "12345" {
		t.Errorf("Expected %s, got %s", "12345", got.ID)
	}

	got, _ = client.FindFirewall("nod")
	if got.ID != "67789" {
		t.Errorf("Expected %s, got %s", "67789", got.ID)
	}

	_, err := client.FindFirewall("web")
	if err.Error() != "unable to find web because there were multiple matches" {
		t.Errorf("Expected %s, got %s", "unable to find web because there were multiple matches", err.Error())
	}

	_, err = client.FindFirewall("missing")
	if err.Error() != "unable to find missing, zero matches" {
		t.Errorf("Expected %s, got %s", "unable to find missing, zero matches", err.Error())
	}
}

func TestNewFirewall(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/firewalls": `{
			"id": "76cc107f-fbef-4e2b-b97f-f5d34f4075d3",
			"name": "fw-mail",
			"result": "success"
		}`,
	})
	defer server.Close()

	got, err := client.NewFirewall("fw-mail")
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

func TestRenameFirewall(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/firewalls/12346": `{"result": "success"}`,
	})
	defer server.Close()
	got, err := client.RenameFirewall("12346", "new_name")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &SimpleResponse{Result: "success"}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestDeleteFirewall(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/firewalls/12346": `{"result": "success"}`,
	})
	defer server.Close()
	got, err := client.DeleteFirewall("12346")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &SimpleResponse{Result: "success"}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestNewFirewallRule(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/firewalls/78901/rules": `{
		  "id": "123456",
		  "firewall_id": "78901",
		  "openstack_security_group_rule_id": null,
		  "protocol": "tcp",
		  "start_port": "443",
		  "end_port": "443",
		  "cidr": [
			"0.0.0.0/0"
		  ],
		  "direction": "ingress",
		  "label": null
		}`,
	})
	defer server.Close()

	cfg := &FirewallRuleConfig{FirewallID: "78901", Protocol: "tcp", StartPort: "443", EndPort: "443", Cidr: []string{"0.0.0.0/0"}, Direction: "inbound", Label: ""}
	got, err := client.NewFirewallRule(cfg)
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

	if expected.StartPort != got.StartPort {
		t.Errorf("Expected %s, got %s", expected.StartPort, got.StartPort)
	}

	if expected.EndPort != got.EndPort {
		t.Errorf("Expected %s, got %s", expected.EndPort, got.EndPort)
	}

	if !reflect.DeepEqual(expected.Cidr, got.Cidr) {
		t.Errorf("Expected %q, got %q", expected.Cidr, got.Cidr)
	}

	if expected.Direction != got.Direction {
		t.Errorf("Expected %s, got %s", expected.Direction, got.Direction)
	}

	if expected.Label != got.Label {
		t.Errorf("Expected %s, got %s", expected.Label, got.Label)
	}
}

func TestFindFirewallRule(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/firewalls/22/rules": `[{
			"id": "21",
			"firewall_id": "22",
			"openstack_security_group_rule_id": null,
			"protocol": "tcp",
			"start_port": "443",
			"end_port": "443",
			"cidr": [
			  "0.0.0.0/0"
			],
			"direction": "ingress",
			"label": "My Rule"
		  },{
			"id": "22",
			"firewall_id": "22",
			"openstack_security_group_rule_id": null,
			"protocol": "tcp",
			"start_port": "80",
			"end_port": "80",
			"cidr": [
			  "0.0.0.0/0"
			],
			"direction": "ingress",
			"label": "My Rule"
		  }]`,
	})
	defer server.Close()

	got, _ := client.FindFirewallRule("22", "21")
	if got.ID != "21" {
		t.Errorf("Expected %s, got %s", "1", got.ID)
	}

	got, _ = client.FindFirewallRule("22", "22")
	if got.ID != "22" {
		t.Errorf("Expected %s, got %s", "2", got.ID)
	}

	_, err := client.FindFirewallRule("22", "2")
	if err.Error() != "unable to find 2 because there were multiple matches" {
		t.Errorf("Expected %s, got %s", "unable to find 2 because there were multiple matches", err.Error())
	}

	_, err = client.FindFirewallRule("22", "missing")
	if err.Error() != "unable to find missing, zero matches" {
		t.Errorf("Expected %s, got %s", "unable to find missing, zero matches", err.Error())
	}
}

func TestListFirewallRules(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/firewalls/22/rules": `[{
			"id": "1",
			"firewall_id": "22",
			"openstack_security_group_rule_id": null,
			"protocol": "tcp",
			"start_port": "443",
			"end_port": "443",
			"cidr": [
			  "0.0.0.0/0"
			],
			"direction": "ingress",
			"label": "My Rule"
		  },{
			"id": "2",
			"firewall_id": "22",
			"openstack_security_group_rule_id": null,
			"protocol": "tcp",
			"start_port": "80",
			"end_port": "80",
			"cidr": [
			  "0.0.0.0/0"
			],
			"direction": "ingress",
			"label": "My Rule"
		  }]`,
	})
	defer server.Close()
	got, err := client.ListFirewallRules("22")

	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	expected := []FirewallRule{{ID: "1", FirewallID: "22", Protocol: "tcp", StartPort: "443", EndPort: "443", Cidr: []string{"0.0.0.0/0"}, Direction: "ingress", Label: "My Rule"}, {ID: "2", FirewallID: "22", Protocol: "tcp", StartPort: "80", EndPort: "80", Cidr: []string{"0.0.0.0/0"}, Direction: "ingress", Label: "My Rule"}}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestDeleteFirewallRule(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/firewalls/12346/rules/12345": `{"result": "success"}`,
	})
	defer server.Close()
	got, err := client.DeleteFirewallRule("12346", "12345")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &SimpleResponse{Result: "success"}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}
