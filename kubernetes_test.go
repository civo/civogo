package civogo

import (
	"reflect"
	"testing"
)

func TestListKubernetesCluster(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/kubernetes/clusters": `[{
		  "id": "69a23478-a89e-41d2-97b1-6f4c341cee70",
		  "name": "your-cluster-name",
		  "version": "2",
		  "status": "ACTIVE",
		  "ready": true,
		  "num_target_nodes": 1,
		  "target_nodes_size": "g2.xsmall",
		  "built_at": "2019-09-23T13:04:23.000+01:00",
		  "kubeconfig": "YAML_VERSION_OF_KUBECONFIG_HERE\n",
		  "kubernetes_version": "0.8.1",
		  "api_endpoint": "https://your.cluster.ip.address:6443",
		  "dns_entry": "69a23478-a89e-41d2-97b1-6f4c341cee70.k8s.civo.com",
		  "tags": [],
		  "created_at": "2019-09-23T13:02:59.000+01:00",
		  "instances": [{
			"hostname": "kube-master-HEXDIGITS",
			"size": "g2.xsmall",
			"region": "lon1",
			"created_at": "2019-09-23T13:03:00.000+01:00",
			"status": "ACTIVE",
			"firewall_id": "5f0ba9ed-5ca7-4e14-9a09-449a84196d64",
			"public_ip": "your.cluster.ip.address",
			"tags": ["civo-kubernetes:installed", "civo-kubernetes:master"]
		  }],
		  "installed_applications": [{
			"application": "Traefik",
			"title": null,
			"version": "(default)",
			"dependencies": null,
			"maintainer": "@Rancher_Labs",
			"description": "A reverse proxy/load-balancer that's easy, dynamic, automatic, fast, full-featured, open source, production proven and provides metrics.",
			"post_install": "Some documentation here\n",
			"installed": true,
			"url": "https://traefik.io",
			"category": "architecture",
			"updated_at": "2019-09-23T13:02:59.000+01:00",
			"image_url": "https://api.civo.com/k3s-marketplace/traefik.png",
			"plan": null,
			"configuration": {}
		  }]
		}]`,
	})
	defer server.Close()
	got, err := client.ListKubernetesCluster()

	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	expected := []Kubernetes{{
		ID:                "12345",
		Name:              "instance-123456",
		Version:           "2",
		Status:            "ACTIVE",
		Ready:             true,
		NumTargetNode:     1,
		TargetNodeSize:    "g2.xsmall",
		KubeConfig:        "YAML_VERSION_OF_KUBECONFIG_HERE\n",
		KubernetesVersion: "0.8.1",
		ApiEndPoint:       "https://your.cluster.ip.address:6443",
		DnsEntry:          "69a23478-a89e-41d2-97b1-6f4c341cee70.k8s.civo.com",
		Tags:              [],
		Instances: []KubeInstances{{
			Hostname:   "kube-master-HEXDIGITS",
			Size:       "g2.xsmall",
			Region:     "lon1",
			Status:     "ACTIVE",
			FirewallID: "5f0ba9ed-5ca7-4e14-9a09-449a84196d64",
			PublicIP:   "your.cluster.ip.address",
			Tags:       ["civo-kubernetes:installed", "civo-kubernetes:master"],
		}},
		InstalledApplications: []KubeApplications{{
			Application:   "Traefik",
			Title:         null,
			Version:       "(default)",
			Dependencies:  null,
			Maintainer:    "@Rancher_Labs",
			Description:   "A reverse proxy/load-balancer that's easy, dynamic, automatic, fast, full-featured, open source, production proven and provides metrics.",
			PostInstall:   "Some documentation here\n",
			Installed:     true,
			Url:           "https://traefik.io",
			Category:      "architecture",
			ImageUrl:      "https://api.civo.com/k3s-marketplace/traefik.png",
			Plan:          null,
			Configuration: {}
		}}
	}}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

//func TestNewFirewall(t *testing.T) {
//	client, server, _ := NewClientForTesting(map[string]string{
//		"/v2/firewalls/": `{
//			"id": "76cc107f-fbef-4e2b-b97f-f5d34f4075d3",
//			"name": "fw-mail",
//			"result": "success"
//		}`,
//	})
//	defer server.Close()
//
//	cfg := &FirewallConfig{Name: "fw-mail"}
//	got, err := client.NewFirewall(cfg)
//	if err != nil {
//		t.Errorf("Request returned an error: %s", err)
//		return
//	}
//
//	expected := &FirewallResult{
//		ID:     "76cc107f-fbef-4e2b-b97f-f5d34f4075d3",
//		Name:   "fw-mail",
//		Result: "success",
//	}
//
//	if expected.ID != got.ID {
//		t.Errorf("Expected %s, got %s", expected.ID, got.ID)
//	}
//
//	if expected.Name != got.Name {
//		t.Errorf("Expected %s, got %s", expected.Name, got.Name)
//	}
//
//	if expected.Result != got.Result {
//		t.Errorf("Expected %s, got %s", expected.Result, got.Result)
//	}
//}
//
//func TestDeleteFirewall(t *testing.T) {
//	client, server, _ := NewClientForTesting(map[string]string{
//		"/v2/firewalls/12346": `{"result": "success"}`,
//	})
//	defer server.Close()
//	got, err := client.DeleteFirewall("12346")
//	if err != nil {
//		t.Errorf("Request returned an error: %s", err)
//		return
//	}
//
//	expected := &SimpleResponse{Result: "success"}
//	if !reflect.DeepEqual(got, expected) {
//		t.Errorf("Expected %+v, got %+v", expected, got)
//	}
//}
//
//func TestNewFirewallRule(t *testing.T) {
//	client, server, _ := NewClientForTesting(map[string]string{
//		"/v2/firewalls/78901/rules": `{
//		  "id": "123456",
//		  "firewall_id": "78901",
//		  "openstack_security_group_rule_id": null,
//		  "protocol": "tcp",
//		  "start_port": "443",
//		  "end_port": "443",
//		  "cidr": [
//			"0.0.0.0/0"
//		  ],
//		  "direction": "ingress",
//		  "label": null
//		}`,
//	})
//	defer server.Close()
//
//	cfg := &FirewallRuleConfig{FirewallID: "78901", Protocol: "tcp", StartPort: "443", EndPort: "443", Cidr: []string{"0.0.0.0/0"}, Direction: "inbound", Label: ""}
//	got, err := client.NewFirewallRule(cfg)
//	if err != nil {
//		t.Errorf("Request returned an error: %s", err)
//		return
//	}
//
//	expected := &FirewallRule{
//		ID:         "123456",
//		FirewallID: "78901",
//		Protocol:   "tcp",
//		StartPort:  "443",
//		EndPort:    "443",
//		Cidr:       []string{"0.0.0.0/0"},
//		Direction:  "ingress",
//		Label:      "",
//	}
//
//	if expected.ID != got.ID {
//		t.Errorf("Expected %s, got %s", expected.ID, got.ID)
//	}
//
//	if expected.FirewallID != got.FirewallID {
//		t.Errorf("Expected %s, got %s", expected.FirewallID, got.FirewallID)
//	}
//
//	if expected.Protocol != got.Protocol {
//		t.Errorf("Expected %s, got %s", expected.Protocol, got.Protocol)
//	}
//
//	if expected.StartPort != got.StartPort {
//		t.Errorf("Expected %s, got %s", expected.StartPort, got.StartPort)
//	}
//
//	if expected.EndPort != got.EndPort {
//		t.Errorf("Expected %s, got %s", expected.EndPort, got.EndPort)
//	}
//
//	if !reflect.DeepEqual(expected.Cidr, got.Cidr) {
//		t.Errorf("Expected %q, got %q", expected.Cidr, got.Cidr)
//	}
//
//	if expected.Direction != got.Direction {
//		t.Errorf("Expected %s, got %s", expected.Direction, got.Direction)
//	}
//
//	if expected.Label != got.Label {
//		t.Errorf("Expected %s, got %s", expected.Label, got.Label)
//	}
//}
//
//func TestListFirewallRule(t *testing.T) {
//	client, server, _ := NewClientForTesting(map[string]string{
//		"/v2/firewalls/22/rules": `[{
//			"id": "1",
//			"firewall_id": "22",
//			"openstack_security_group_rule_id": null,
//			"protocol": "tcp",
//			"start_port": "443",
//			"end_port": "443",
//			"cidr": [
//			  "0.0.0.0/0"
//			],
//			"direction": "ingress",
//			"label": "My Rule"
//		  },{
//			"id": "2",
//			"firewall_id": "22",
//			"openstack_security_group_rule_id": null,
//			"protocol": "tcp",
//			"start_port": "80",
//			"end_port": "80",
//			"cidr": [
//			  "0.0.0.0/0"
//			],
//			"direction": "ingress",
//			"label": "My Rule"
//		  }]`,
//	})
//	defer server.Close()
//	got, err := client.ListFirewallRule("22")
//
//	if err != nil {
//		t.Errorf("Request returned an error: %s", err)
//		return
//	}
//	expected := []FirewallRule{{ID: "1", FirewallID: "22", Protocol: "tcp", StartPort: "443", EndPort: "443", Cidr: []string{"0.0.0.0/0"}, Direction: "ingress", Label: "My Rule"}, {ID: "2", FirewallID: "22", Protocol: "tcp", StartPort: "80", EndPort: "80", Cidr: []string{"0.0.0.0/0"}, Direction: "ingress", Label: "My Rule"}}
//	if !reflect.DeepEqual(got, expected) {
//		t.Errorf("Expected %+v, got %+v", expected, got)
//	}
//}
//
//func TestDeleteFirewallRule(t *testing.T) {
//	client, server, _ := NewClientForTesting(map[string]string{
//		"/v2/firewalls/12346/rules/12345": `{"result": "success"}`,
//	})
//	defer server.Close()
//	got, err := client.DeleteFirewallRule("12346", "12345")
//	if err != nil {
//		t.Errorf("Request returned an error: %s", err)
//		return
//	}
//
//	expected := &SimpleResponse{Result: "success"}
//	if !reflect.DeepEqual(got, expected) {
//		t.Errorf("Expected %+v, got %+v", expected, got)
//	}
//}
