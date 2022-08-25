package civogo

import (
	"reflect"
	"testing"
)

func TestListObjectStores(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/objectstores": `{"page": 1, "per_page": 20, "pages": 2, "items":[{"id": "12345", "name": "test-objectstore"}]}`,
	})
	defer server.Close()

	got, err := client.ListObjectStores()
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &PaginatedObjectstores{
		Page:    1,
		PerPage: 20,
		Pages:   2,
		Items: []ObjectStore{
			{
				ID:   "12345",
				Name: "test-objectstore",
			},
		},
	}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestFindObjectStore(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/objectstores": `{
			"page": 1,
			"per_page": 20,
			"pages": 1,
			"items": [
			  {
				"id": "12345",
				"name": "test-objectstore"
			  }
			]
		  }`,
	})
	defer server.Close()

	got, _ := client.FindObjectStore("test-objectstore")
	if got.ID != "12345" {
		t.Errorf("Expected %s, got %s", "12345", got.ID)
	}
}

func TestNewObjectStore(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/objectstores": `{
			"id": "12345",
			"name": "test-objectstore",
			"max_size": 500,
			"status" : "active"
		}`,
	})
	defer server.Close()

	cfg := &CreateObjectStoreRequest{
		Name:      "test-objectstore",
		MaxSizeGB: 500,
	}
	got, err := client.NewObjectStore(cfg)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &ObjectStore{
		ID:      "12345",
		Name:    "test-objectstore",
		MaxSize: 500,
		Status:  "active",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestUpdateObjectStore(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/objectstores/12345": `{
			"id": "12345",
			"name": "test-objectstore",
			"max_size": 1000
		}`,
	})
	defer server.Close()

	cfg := &UpdateObjectStoreRequest{
		MaxSizeGB: 1000,
	}
	got, err := client.UpdateObjectStore("12345", cfg)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &ObjectStore{
		ID:      "12345",
		Name:    "test-objectstore",
		MaxSize: 1000,
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestDeleteObjectStore(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/objectstores/12345": `{"result": "success"}`,
	})
	defer server.Close()

	got, err := client.DeleteObjectStore("12345")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &SimpleResponse{Result: "success"}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}
