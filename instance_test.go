package civogo

import (
	"testing"
)

func TestListInstances(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
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

func TestFindInstance(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/instances": `{"page": 1, "per_page": 20, "pages": 2, "items":[{"id": "12345", "hostname": "foo.example.com"}, {"id":"67890", "hostname": "bar.zip.com"}]}`,
	})
	defer server.Close()

	got, _ := client.FindInstance("45")
	if got.ID != "12345" {
		t.Errorf("Expected %s, got %s", "12345", got.ID)
	}

	got, _ = client.FindInstance("89")
	if got.ID != "67890" {
		t.Errorf("Expected %s, got %s", "67890", got.ID)
	}

	got, _ = client.FindInstance("foo")
	if got.ID != "12345" {
		t.Errorf("Expected %s, got %s", "12345", got.ID)
	}

	got, _ = client.FindInstance("bar")
	if got.ID != "67890" {
		t.Errorf("Expected %s, got %s", "67890", got.ID)
	}

	_, err := client.FindInstance("com")
	if err.Error() != "MultipleMatchesError: unable to find com because there were multiple matches" {
		t.Errorf("Expected %s, got %s", "unable to find com because there were multiple matches", err.Error())
	}

	_, err = client.FindInstance("missing")
	if err.Error() != "ZeroMatchesError: unable to find missing, zero matches" {
		t.Errorf("Expected %s, got %s", "unable to find missing, zero matches", err.Error())
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

	if got.Items[0].ID != "12345" {
		t.Errorf("Expected %s, got %s", "12345", got.Items[0].ID)
	}
	if got.Page != 1 {
		t.Errorf("Expected %d, got %d", 1, got.Page)
	}
	if got.Pages != 2 {
		t.Errorf("Expected %d, got %d", 2, got.Pages)
	}
	if got.PerPage != 20 {
		t.Errorf("Expected %d, got %d", 20, got.PerPage)
	}
}

func TestGetInstance(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/instances/12345": `{"id": "12345", "hostname": "foo.example.com", "ipv6":"::1234:5678:9abc:def0"}`,
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
	if got.IPv6 != "::1234:5678:9abc:def0" {
		t.Errorf("Expected %s, got %s", "::1234:5678:9abc:def0", got.IPv6)
	}
}

func TestNewInstanceConfig(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/networks":    `[{"id": "1", "default": true, "name": "Default Network"}]`,
		"/v2/disk_images": `[{ "ID": "77bea4dd-bfd4-492c-823d-f92eb6dd962d", "Name": "ubuntu-focal", "Version": "20.04", "State": "available", "Distribution": "ubuntu", "Description": "", "Label": "focal" }]`,
		"/v2/sshkeys":     `{"items":[{"id": "4", "name": "RSA Key", "default": true}]}`,
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
	if got.TemplateID != "77bea4dd-bfd4-492c-823d-f92eb6dd962d" {
		t.Errorf("Expected %s, got %s", "77bea4dd-bfd4-492c-823d-f92eb6dd962d", got.TemplateID)
	}
	if got.Count != 1 {
		t.Errorf("Expected %d, got %d", 1, got.Count)
	}
	if got.FirewallID != "" {
		t.Errorf("Expected firewall ID to be blank, but got %s", got.FirewallID)
	}
}

func TestCreateInstance(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/instances": `{
		  "id": "b177ae0e-60fa-11e5-be02-5cf9389be614",
		  "openstack_server_id": "369588f7-de40-4eca-bc8d-4c2dbc1cc7f3",
		  "hostname": "b177ae0e-60fa-11e5-be02-5cf9389be614.clients.civo.com",
		  "reverse_dns": null,
		  "size": "g2.xsmall",
		  "region": "lon1",
		  "network_id": "12345",
		  "private_ip": "10.0.0.4",
		  "public_ip": "31.28.66.181",
		  "pseudo_ip": "172.31.0.230",
		  "template_id": "2",
		  "snapshot_id": null,
		  "initial_user": "civo",
		  "initial_password": "password_here",
		  "ssh_key": "61f1b5c8-2c87-4cc7-b1af-6278f3050a28",
		  "ssh_key_id": "816014bd-55de-446a-aedd-2059ff12bb79",
		  "status": "ACTIVE",
		  "notes": null,
		  "firewall_id": "default",
		  "tags": [
			"web",
			"main",
			"linux"
		  ],
		  "civostatsd_token": "f84d920f-c74b-4b48-a21e-5ff7a671e5f9",
		  "civostatsd_stats": null,
		  "civostatsd_stats_per_minute": [],
		  "civostatsd_stats_per_hour": [],
		  "openstack_image_id": null,
		  "rescue_password": null,
		  "volume_backed": true,
		  "script": "#!/bin/bash\necho 'Hello world'",
		  "created_at": "2015-09-20T19:31:36+00:00"
		}`,
	})
	defer server.Close()

	config := &InstanceConfig{
		Hostname:   "b177ae0e-60fa-11e5-be02-5cf9389be614.clients.civo.com",
		Size:       "g2.xsmall",
		NetworkID:  "12345",
		TemplateID: "2",
		SSHKeyID:   "61f1b5c8-2c87-4cc7-b1af-6278f3050a28",
		TagsList:   "web main linux",
	}

	got, err := client.CreateInstance(config)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	if got.Hostname != "b177ae0e-60fa-11e5-be02-5cf9389be614.clients.civo.com" {
		t.Errorf("Expected %s, got %s", "b177ae0e-60fa-11e5-be02-5cf9389be614.clients.civo.com", got.NetworkID)
	}
	if got.NetworkID != "12345" {
		t.Errorf("Expected %s, got %s", "12345", got.NetworkID)
	}
	if got.TemplateID != "2" {
		t.Errorf("Expected %s, got %s", "2", got.TemplateID)
	}
	if got.SSHKey != "61f1b5c8-2c87-4cc7-b1af-6278f3050a28" {
		t.Errorf("Expected %s, got %s", "61f1b5c8-2c87-4cc7-b1af-6278f3050a28", got.SSHKey)
	}
	if got.SSHKeyID != "816014bd-55de-446a-aedd-2059ff12bb79" {
		t.Errorf("Expected %s, got %s", "61f1b5c8-2c87-4cc7-b1af-6278f3050a28", got.SSHKey)
	}
}

func TestSetInstanceTags(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/instances/12345/tags": `{
			"result": "success"
		}`,
	})
	defer server.Close()

	got, err := client.SetInstanceTags(&Instance{ID: "12345"}, "prod lamp")
	EnsureSuccessfulSimpleResponse(t, got, err)
}

func TestUpdateInstance(t *testing.T) {
	client, server, _ := NewAdvancedClientForTesting([]ConfigAdvanceClientForTesting{
		{
			Method: "PUT",
			Value: []ValueAdvanceClientForTesting{
				{
					RequestBody:  `{"hostname":"dummy.example.com","notes":"my notes","reverse_dns":"dummy-reverse.example.com", "subnets": ["test-subnet-1", "test-subnet-2"]}`,
					URL:          "/v2/instances/12345",
					ResponseBody: `{"result": "success"}`,
				},
			},
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
	client, server, _ := NewAdvancedClientForTesting([]ConfigAdvanceClientForTesting{
		{
			Method: "DELETE",
			Value: []ValueAdvanceClientForTesting{
				{
					RequestBody:  "",
					URL:          "/v2/instances/12345",
					ResponseBody: `{"result": "success"}`,
				},
			},
		},
	})
	defer server.Close()

	got, err := client.DeleteInstance("12345")
	EnsureSuccessfulSimpleResponse(t, got, err)
}

func TestRebootInstance(t *testing.T) {
	client, server, _ := NewAdvancedClientForTesting([]ConfigAdvanceClientForTesting{
		{
			Method: "POST",
			Value: []ValueAdvanceClientForTesting{
				{
					RequestBody:  "",
					URL:          "/v2/instances/12345/hard_reboots",
					ResponseBody: `{"result": "success"}`,
				},
			},
		},
	})
	defer server.Close()

	got, err := client.RebootInstance("12345")
	EnsureSuccessfulSimpleResponse(t, got, err)
}

func TestHardRebootInstance(t *testing.T) {
	client, server, _ := NewAdvancedClientForTesting([]ConfigAdvanceClientForTesting{
		{
			Method: "POST",
			Value: []ValueAdvanceClientForTesting{
				{
					RequestBody:  "",
					URL:          "/v2/instances/12345/hard_reboots",
					ResponseBody: `{"result": "success"}`,
				},
			},
		},
	})
	defer server.Close()

	got, err := client.HardRebootInstance("12345")
	EnsureSuccessfulSimpleResponse(t, got, err)
}

func TestSoftRebootInstance(t *testing.T) {
	client, server, _ := NewAdvancedClientForTesting([]ConfigAdvanceClientForTesting{
		{
			Method: "POST",
			Value: []ValueAdvanceClientForTesting{
				{
					RequestBody:  "",
					URL:          "/v2/instances/12345/soft_reboots",
					ResponseBody: `{"result": "success"}`,
				},
			},
		},
	})
	defer server.Close()

	got, err := client.SoftRebootInstance("12345")
	EnsureSuccessfulSimpleResponse(t, got, err)
}

func TestStopInstance(t *testing.T) {
	client, server, _ := NewAdvancedClientForTesting([]ConfigAdvanceClientForTesting{
		{
			Method: "PUT",
			Value: []ValueAdvanceClientForTesting{
				{
					RequestBody:  "",
					URL:          "/v2/instances/12345/stop",
					ResponseBody: `{"result": "success"}`,
				},
			},
		},
	})
	defer server.Close()

	got, err := client.StopInstance("12345")
	EnsureSuccessfulSimpleResponse(t, got, err)
}

func TestStartInstance(t *testing.T) {
	client, server, _ := NewAdvancedClientForTesting([]ConfigAdvanceClientForTesting{
		{
			Method: "PUT",
			Value: []ValueAdvanceClientForTesting{
				{
					RequestBody:  "",
					URL:          "/v2/instances/12345/start",
					ResponseBody: `{"result": "success"}`,
				},
			},
		},
	})
	defer server.Close()

	got, err := client.StartInstance("12345")
	EnsureSuccessfulSimpleResponse(t, got, err)
}

func TestUpgradeInstance(t *testing.T) {
	client, server, _ := NewAdvancedClientForTesting([]ConfigAdvanceClientForTesting{
		{
			Method: "PUT",
			Value: []ValueAdvanceClientForTesting{
				{
					RequestBody:  `{"size":"g99.huge"}`,
					URL:          "/v2/instances/12345/resize",
					ResponseBody: `{"result": "success"}`,
				},
			},
		},
	})
	defer server.Close()

	got, err := client.UpgradeInstance("12345", "g99.huge")
	EnsureSuccessfulSimpleResponse(t, got, err)
}

func TestMovePublicIPToInstance(t *testing.T) {
	client, server, _ := NewAdvancedClientForTesting([]ConfigAdvanceClientForTesting{
		{
			Method: "PUT",
			Value: []ValueAdvanceClientForTesting{
				{
					RequestBody:  `""`,
					URL:          "/v2/instances/12345/ip/1.2.3.4",
					ResponseBody: `{"result": "success"}`,
				},
			},
		},
	})
	defer server.Close()

	got, err := client.MovePublicIPToInstance("12345", "1.2.3.4")
	EnsureSuccessfulSimpleResponse(t, got, err)
}

func TestGetInstanceConsoleURL(t *testing.T) {
	client, server, _ := NewAdvancedClientForTesting([]ConfigAdvanceClientForTesting{
		{
			Method: "GET",
			Value: []ValueAdvanceClientForTesting{
				{
					RequestBody:  `""`,
					URL:          "/v2/instances/12345/console",
					ResponseBody: `{"url": "https://console.example.com/12345"}`,
				},
			},
		},
	})
	defer server.Close()

	got, _ := client.GetInstanceConsoleURL("12345")

	if got != "https://console.example.com/12345" {
		t.Errorf("Expected %s, got %s", "https://console.example.com/12345", got)
	}
}

func TestSetInstanceFirewall(t *testing.T) {
	client, server, _ := NewAdvancedClientForTesting([]ConfigAdvanceClientForTesting{
		{
			Method: "PUT",
			Value: []ValueAdvanceClientForTesting{
				{
					RequestBody:  `{"firewall_id":"67890"}`,
					URL:          "/v2/instances/12345/firewall",
					ResponseBody: `{"result": "success"}`,
				},
			},
		},
	})
	defer server.Close()

	got, err := client.SetInstanceFirewall("12345", "67890")
	EnsureSuccessfulSimpleResponse(t, got, err)
}
