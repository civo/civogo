package civogo

import (
	"reflect"
	"testing"
)

func TestListIPs(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/ips": `{
			"page": 1,
			"per_page": 20,
			"pages": 1,
			"items": [
			  {
				"id": "7bb2c574-7b34-4de4-9111-4ac2b5653efa",
				"name": "test-ip",
				"ip": "127.0.0.1"
			  }
			]
		  }`,
	})
	defer server.Close()

	got, err := client.ListIPs()
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &PaginatedIPs{
		Page:    1,
		PerPage: 20,
		Pages:   1,
		Items: []IP{
			{
				ID:   "7bb2c574-7b34-4de4-9111-4ac2b5653efa",
				Name: "test-ip",
				IP:   "127.0.0.1",
			},
		},
	}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestFindIP(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/ips": `{
			"page": 1,
			"per_page": 20,
			"pages": 1,
			"items": [
			  {
				"id": "7bb2c574-7b34-4de4-9111-4ac2b5653efa",
				"name": "test-ip",
				"ip": "127.0.0.1"
			  }
			]
		  }`,
	})
	defer server.Close()

	got, _ := client.FindIP("7bb2c574")
	if got.ID != "7bb2c574-7b34-4de4-9111-4ac2b5653efa" {
		t.Errorf("Expected %s, got %s", "16549d23-1957-4dea-a3d5-6282d19a4c00", got.ID)
	}

	got, _ = client.FindIP("127.0.0.1")
	if got.ID != "7bb2c574-7b34-4de4-9111-4ac2b5653efa" {
		t.Errorf("Expected %s, got %s", "b421b075-ff18-48b4-a092-604bca968f49", got.ID)
	}

	_, err := client.FindIP("missing")
	if err.Error() != "ZeroMatchesError: unable to find missing, zero matches" {
		t.Errorf("Expected %s, got %s", "unable to find missing, zero matches", err.Error())
	}
}

func TestCreateIP(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/ips": `{
			"id": "56dca3ae-ea3f-480f-9b25-abf90b439729",
			"name": "test-ip"		
		}`,
	})
	defer server.Close()

	cfg := &CreateIPRequest{
		Name: "test-ip",
	}
	got, err := client.NewIP(cfg)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &IP{
		ID:   "56dca3ae-ea3f-480f-9b25-abf90b439729",
		Name: "test-ip",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestUpdateIP(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/ips/a1bd123c-b7e2-4d4f-9fda-7940c7e06b38": `{
			"id": "a1bd123c-b7e2-4d4f-9fda-7940c7e06b38",
			"name": "test-ip-updated"
		}`,
	})
	defer server.Close()

	cfg := &UpdateIPRequest{
		Name: "test-ip-updated",
	}
	got, err := client.UpdateIP("a1bd123c-b7e2-4d4f-9fda-7940c7e06b38", cfg)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &IP{
		ID:   "a1bd123c-b7e2-4d4f-9fda-7940c7e06b38",
		Name: "test-ip-updated",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestDeleteIP(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/ips/12345": `{"result": "success"}`,
	})
	defer server.Close()

	got, err := client.DeleteIP("12345")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &SimpleResponse{Result: "success"}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestAssignIP(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/ips/12345/actions": `{"result": "success"}`,
	})
	defer server.Close()

	got, err := client.AssignIP("12345", "234567", "instance", "TEST")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &SimpleResponse{Result: "success"}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestUnAssignIP(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/ips/12345/actions": `{"result": "success"}`,
	})
	defer server.Close()

	got, err := client.UnassignIP("12345", "TEST")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &SimpleResponse{Result: "success"}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}
