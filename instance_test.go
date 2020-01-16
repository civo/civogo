package civogo

import (
	"fmt"
	"testing"
)

func TestListInstances(t *testing.T) {
	client, server, err := NewClientForTesting(map[string]string{
		"/v2/instances": `{"page": 1, "per_page": 20, "pages": 2, "items":[{"id": "12345", "hostname": "foo.example.com"}]}`,
	})
	defer server.Close()

	items, err := client.ListAllInstances()
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	if items[0].ID != "12345" {
		t.Errorf("Expected %s, got %s", "12345", items[0].ID)
	}
}

func TestListInstancesWithPage(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/instances?page=2&per_page=10": `{"page": 1, "per_page": 20, "pages": 2, "items":[{"id": "12345", "hostname": "foo.example.com"}]}`,
	})
	defer server.Close()

	got, err := client.ListInstances(2, 10)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	fmt.Println(got)
	if got.Items[0].ID != "12345" {
		t.Errorf("Expected %s, got %s", "12345", got.Items[0].ID)
	}
	if got.Page != 2 {
		t.Errorf("Expected %d, got %d", 1, got.Page)
	}
	if got.Pages != 2 {
		t.Errorf("Expected %d, got %d", 2, got.Pages)
	}
	if got.PerPage != 10 {
		t.Errorf("Expected %d, got %d", 20, got.PerPage)
	}
}

func TestGetInstance(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/instances/12345": `{"id": "12345", "hostname": "foo.example.com"}`,
	})
	defer server.Close()

	got, err := client.GetInstance("12345")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	if got.ID != "12345" {
		t.Errorf("Expected %s, got %s", "12345", got.ID)
	}
	if got.Hostname != "foo.example.com" {
		t.Errorf("Expected %s, got %s", "foo.example.com", got.Hostname)
	}
}

func TestNewInstanceConfig(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/networks":  `[{"id": "1", "default": true, "name": "Default Network"}]`,
		"/v2/templates": `[{"id": "2", "code": "centos-7"},{"id": "3", "code": "ubuntu-18.04"}]`,
		"/v2/sshkeys":   `{"items":[{"id": "4", "name": "RSA Key"}]}`,
	})
	defer server.Close()

	got, err := client.NewInstanceConfig()
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	if got.Hostname == "" {
		t.Errorf("Expected hostname not to be blank, but it was")
	}
	if got.NetworkID != "1" {
		t.Errorf("Expected %s, got %s", "1", got.NetworkID)
	}
	if got.TemplateID != "3" {
		t.Errorf("Expected %s, got %s", "3", got.TemplateID)
	}
	if got.SSHKeyID != "4" {
		t.Errorf("Expected %s, got %s", "3", got.TemplateID)
	}
	if got.Count != 1 {
		t.Errorf("Expected %d, got %d", 1, got.Count)
	}
}

func TestCreateInstance(t *testing.T) {
	client, server, _ := NewAdvancedClientForTesting(map[string]map[string]string{
		"/v2/networks": map[string]string{
			"requestBody":  "",
			"method":       "GET",
			"responseBody": `[{"id": "1", "default": true, "name": "Default Network"}]`,
		},
		"/v2/templates": map[string]string{
			"requestBody":  "",
			"method":       "GET",
			"responseBody": `[{"id": "2", "code": "centos-7"},{"id": "3", "code": "ubuntu-18.04"}]`,
		},
		"/v2/sshkeys": map[string]string{
			"requestBody":  "",
			"method":       "GET",
			"responseBody": `{"items":[{"id": "4", "name": "RSA Key"}]}`,
		},
		"/v2/instances": map[string]string{
			"requestBody":  "count=1&hostname=foo.example.com&initial_user=civo&network_id=1&public_ip_required=true&region=lon1&reverse_dns=&script=&size=g2.xsmall&snapshot_id=&ssh_key_id=4&tags=&template_id=3",
			"method":       "POST",
			"responseBody": `{"id": "12345", "hostname": "foo.example.com", "network_id": "1", "ssh_key": "4", "template_id": "3"}`,
		},
	})
	defer server.Close()

	config, err := client.NewInstanceConfig()
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	config.Hostname = "foo.example.com"

	got, err := client.CreateInstance(config)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	if got.Hostname != "foo.example.com" {
		t.Errorf("Expected %s, got %s", "1", got.NetworkID)
	}
	if got.NetworkID != "1" {
		t.Errorf("Expected %s, got %s", "1", got.NetworkID)
	}
	if got.TemplateID != "3" {
		t.Errorf("Expected %s, got %s", "3", got.TemplateID)
	}
	if got.SSHKey != "4" {
		t.Errorf("Expected %s, got %s", "3", got.TemplateID)
	}
}

func TestSetInstanceTags(t *testing.T) {
	client, server, _ := NewAdvancedClientForTesting(map[string]map[string]string{
		"/v2/instances/12345/tags": map[string]string{
			"requestBody":  `tags=prod+lamp`,
			"method":       "PUT",
			"responseBody": `{"result": "success"}`,
		},
	})
	defer server.Close()

	got, err := client.SetInstanceTags("12345", "prod lamp")
	EnsureSuccessfulSimpleResponse(t, got, err)
}

func TestUpdateInstance(t *testing.T) {
	client, server, _ := NewAdvancedClientForTesting(map[string]map[string]string{
		"/v2/instances/12345": map[string]string{
			"requestBody":  `hostname=dummy.example.com&notes=my+notes&reverse_dns=dummy-reverse.example.com`,
			"method":       "PUT",
			"responseBody": `{"result": "success"}`,
		},
	})
	defer server.Close()

	i := Instance{
		ID:         "12345",
		Hostname:   "dummy.example.com",
		ReverseDNS: "dummy-reverse.example.com",
		Notes:      "my notes",
	}
	got, err := client.UpdateInstance(&i)
	EnsureSuccessfulSimpleResponse(t, got, err)
}

func TestDeleteInstance(t *testing.T) {
	client, server, _ := NewAdvancedClientForTesting(map[string]map[string]string{
		"/v2/instances/12345": map[string]string{
			"requestBody":  ``,
			"method":       "DELETE",
			"responseBody": `{"result": "success"}`,
		},
	})
	defer server.Close()

	got, err := client.DeleteInstance("12345")
	EnsureSuccessfulSimpleResponse(t, got, err)
}

func TestRebootInstance(t *testing.T) {
	client, server, _ := NewAdvancedClientForTesting(map[string]map[string]string{
		"/v2/instances/12345/hard_reboots": map[string]string{
			"requestBody":  ``,
			"method":       "POST",
			"responseBody": `{"result": "success"}`,
		},
	})
	defer server.Close()

	got, err := client.RebootInstance("12345")
	EnsureSuccessfulSimpleResponse(t, got, err)
}

func TestHardRebootInstance(t *testing.T) {
	client, server, _ := NewAdvancedClientForTesting(map[string]map[string]string{
		"/v2/instances/12345/hard_reboots": map[string]string{
			"requestBody":  ``,
			"method":       "POST",
			"responseBody": `{"result": "success"}`,
		},
	})
	defer server.Close()

	got, err := client.HardRebootInstance("12345")
	EnsureSuccessfulSimpleResponse(t, got, err)
}

func TestSoftRebootInstance(t *testing.T) {
	client, server, _ := NewAdvancedClientForTesting(map[string]map[string]string{
		"/v2/instances/12345/soft_reboots": map[string]string{
			"requestBody":  ``,
			"method":       "POST",
			"responseBody": `{"result": "success"}`,
		},
	})
	defer server.Close()

	got, err := client.SoftRebootInstance("12345")
	EnsureSuccessfulSimpleResponse(t, got, err)
}

func TestStopInstance(t *testing.T) {
	client, server, _ := NewAdvancedClientForTesting(map[string]map[string]string{
		"/v2/instances/12345/stop": map[string]string{
			"requestBody":  ``,
			"method":       "PUT",
			"responseBody": `{"result": "success"}`,
		},
	})
	defer server.Close()

	got, err := client.StopInstance("12345")
	EnsureSuccessfulSimpleResponse(t, got, err)
}

func TestStartInstance(t *testing.T) {
	client, server, _ := NewAdvancedClientForTesting(map[string]map[string]string{
		"/v2/instances/12345/start": map[string]string{
			"requestBody":  ``,
			"method":       "PUT",
			"responseBody": `{"result": "success"}`,
		},
	})
	defer server.Close()

	got, err := client.StartInstance("12345")
	EnsureSuccessfulSimpleResponse(t, got, err)
}

func TestUpgradeInstance(t *testing.T) {
	client, server, _ := NewAdvancedClientForTesting(map[string]map[string]string{
		"/v2/instances/12345/resize": map[string]string{
			"requestBody":  `size=g99.huge`,
			"method":       "PUT",
			"responseBody": `{"result": "success"}`,
		},
	})
	defer server.Close()

	got, err := client.UpgradeInstance("12345", "g99.huge")
	EnsureSuccessfulSimpleResponse(t, got, err)
}

func TestMovePublicIPToInstance(t *testing.T) {
	client, server, _ := NewAdvancedClientForTesting(map[string]map[string]string{
		"/v2/instances/12345/ip/1.2.3.4": map[string]string{
			"requestBody":  ``,
			"method":       "PUT",
			"responseBody": `{"result": "success"}`,
		},
	})
	defer server.Close()

	got, err := client.MovePublicIPToInstance("12345", "1.2.3.4")
	EnsureSuccessfulSimpleResponse(t, got, err)
}

func TestSetInstanceFirewall(t *testing.T) {
	client, server, _ := NewAdvancedClientForTesting(map[string]map[string]string{
		"/v2/instances/12345/firewall": map[string]string{
			"requestBody":  `firewall_id=67890`,
			"method":       "PUT",
			"responseBody": `{"result": "success"}`,
		},
	})
	defer server.Close()

	got, err := client.SetInstanceFirewall("12345", "67890")
	EnsureSuccessfulSimpleResponse(t, got, err)
}
