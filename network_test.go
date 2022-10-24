package civogo

import (
	"reflect"
	"testing"
)

func TestGetDefaultNetwork(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/networks": `[{"id": "12345", "default": true, "name": "Default Network", "status": "Active"}]`,
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
	if got.Status != "Active" {
		t.Errorf("Expected %s, got %s", "Active", got.Status)
	}
}

func TestGetNetwork(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/networks/12345": `{"id": "12345", "name": "test-network"}`,
	})
	defer server.Close()

	got, err := client.GetNetwork("12345")
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
			"label": "development",
			"status": "Deleting"
		  }]`,
	})
	defer server.Close()
	got, err := client.ListNetworks()

	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	expected := []Network{{ID: "12345", Name: "my-net", Default: false, CIDR: "0.0.0.0/0", Label: "development", Status: "Deleting"}}
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

	got, _ := client.FindNetwork("34")
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

	got, _ = client.FindNetwork("production")
	if got.ID != "67890" {
		t.Errorf("Expected %s, got %s", "67890", got.ID)
	}

	_, err := client.FindNetwork("net")
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

func TestGetSubnet(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/networks/12345/subnets/6789": `{"network_id": "12345", "subnetID": "6789", "name": "test-subnet"}`,
	})
	defer server.Close()

	got, err := client.GetSubnet("12345", "6789")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	if got.Name != "test-subnet" {
		t.Errorf("Expected %s, got %s", "test-subnet", got.Name)
	}
}

func TestFindSubnet(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/networks/12345/subnets": `[
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

	got, _ := client.FindSubnet("6789", "12345")
	if got.ID != "6789" {
		t.Errorf("Expected %s, got %s", "6789", got.ID)
	}

	got, _ = client.FindSubnet("test-subnet-2", "12345")
	if got.ID != "67890" {
		t.Errorf("Expected %s, got %s", "67890", got.ID)
	}
}

func TestCreateSubnet(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/networks/12345/subnets": `{"id": "76cc107f-fbef-4e2b-b97f-f5d34f4075d3","network_id": "12345","name": "test-subnet","status": "success"}`,
	})
	defer server.Close()

	subnet := SubnetConfig{
		Name: "test-subnet",
	}

	got, err := client.CreateSubnet("12345", subnet)
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

func TestListSubnets(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/networks/12345/subnets": `[{
			"id": "6789",
			"name": "test-subnet",
			"network_id": "12345",
			"label": "test-subnet"
		  }]`,
	})
	defer server.Close()
	got, err := client.ListSubnets("12345")

	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	expected := []Subnet{{ID: "6789", Name: "test-subnet", NetworkID: "12345"}}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestDeleteSubnet(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/networks/12345/subnets/6789": `{"result": "success"}`,
	})
	defer server.Close()

	got, err := client.DeleteSubnet("12345", "6789")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &SimpleResponse{Result: "success"}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}
