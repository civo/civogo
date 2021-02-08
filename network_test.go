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

	got, err := client.NewNetwork("private-net")
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
	expected := []Network{{ID: "12345", Name: "my-net", Default: false, CIDR: "0.0.0.0/0", Label: "development"}}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestFindNetwork(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/networks": `[
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

	got, err := client.FindNetwork("34")
	if got.ID != "12345" {
		t.Errorf("Expected %s, got %s", "12345", got.ID)
	}

	got, _ = client.FindNetwork("89")
	if got.ID != "67890" {
		t.Errorf("Expected %s, got %s", "67890", got.ID)
	}

	got, _ = client.FindNetwork("my")
	if got.ID != "12345" {
		t.Errorf("Expected %s, got %s", "12345", got.ID)
	}

	got, _ = client.FindNetwork("other")
	if got.ID != "67890" {
		t.Errorf("Expected %s, got %s", "67890", got.ID)
	}

	got, err = client.FindNetwork("production")
	if got.ID != "67890" {
		t.Errorf("Expected %s, got %s", "67890", got.ID)
	}

	_, err = client.FindNetwork("net")
	if err.Error() != "MultipleMatchesError: unable to find net because there were multiple matches" {
		t.Errorf("Expected %s, got %s", "unable to find net because there were multiple matches", err.Error())
	}

	_, err = client.FindNetwork("missing")
	if err.Error() != "ZeroMatchesError: unable to find missing, zero matches" {
		t.Errorf("Expected %s, got %s", "ZeroMatchesError: unable to find missing, zero matches", err.Error())
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

	got, err := client.RenameNetwork("private-net", "76cc107f-fbef-4e2b-b97f-f5d34f4075d3")
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
