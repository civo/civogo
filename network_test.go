package civogo

import (
	"reflect"
	"testing"
)

func TestGetDefaultNetwork(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/networks": `[{"id": "12345", "default": true, "name": "Default Network"}]`,
	})
	defer server.Close()

	got, err := client.GetDefaultNetwork()
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	if got.ID != "12345" {
		t.Errorf("Expected %s, got %s", "12345", got.ID)
	}
}

func TestNewNetwork(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/networks": `{
			"id": "76cc107f-fbef-4e2b-b97f-f5d34f4075d3",
			"label": "private-net",
			"result": "success"
		}`,
	})
	defer server.Close()

	cfg := &NetworkConfig{Label: "private-net"}
	got, err := client.NewNetwork(cfg)
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

func TestListNetworks(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/networks": `[{
			"id": "12345",
			"name": "my-net",
			"region": "lon1",
			"default": false,
			"cidr": "0.0.0.0/0",
			"label": "development"
		  }]`,
	})
	defer server.Close()
	got, err := client.ListNetworks()

	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	expected := []Network{{ID: "12345", Name: "my-net", Region: "lon1", Default: false, CIDR: "0.0.0.0/0", Label: "development"}}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestRenameNetwork(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/networks/76cc107f-fbef-4e2b-b97f-f5d34f4075d3": `{
			"id": "76cc107f-fbef-4e2b-b97f-f5d34f4075d3",
			"label": "private-net",
			"result": "success"
		}`,
	})
	defer server.Close()

	cfg := &NetworkConfig{Label: "private-net"}
	got, err := client.RenameNetwork(cfg, "76cc107f-fbef-4e2b-b97f-f5d34f4075d3")
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

func TestDeleteNetwork(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/networks/12345": `{"result": "success"}`,
	})
	defer server.Close()

	got, err := client.DeleteNetwork("12345")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &SimpleResponse{Result: "success"}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}
