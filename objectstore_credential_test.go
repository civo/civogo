package civogo

import (
	"reflect"
	"testing"
)

func TestListObjectStoreCredentials(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/objectstore/credentials": `{"page": 1, "per_page": 20, "pages": 2, "items":[{"id": "12345", "name": "test-objectstore-cred"}]}`,
	})
	defer server.Close()

	got, err := client.ListObjectStoreCredentials()
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &PaginatedObjectStoreCredentials{
		Page:    1,
		PerPage: 20,
		Pages:   2,
		Items: []ObjectStoreCredential{
			{
				ID:   "12345",
				Name: "test-objectstore-cred",
			},
		},
	}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestFindObjectStoreCredential(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/objectstore/credentials": `{
			"page": 1,
			"per_page": 20,
			"pages": 1,
			"items": [
			  {
				"id": "12345",
				"name": "test-objectstore-cred"
			  }
			]
		  }`,
	})
	defer server.Close()

	got, _ := client.FindObjectStoreCredential("test-objectstore-cred")
	if got.ID != "12345" {
		t.Errorf("Expected %s, got %s", "12345", got.ID)
	}
}

func TestNewObjectStoreCredential(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/objectstore/credentials": `{
			"id": "12345",
			"name": "test-objectstore-cred",
			"max_size_gb": 500
		}`,
	})
	defer server.Close()

	cfg := &CreateObjectStoreCredentialRequest{
		Name:      "test-objectstore-cred",
		MaxSizeGB: intPtr(500),
	}
	got, err := client.NewObjectStoreCredential(cfg)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &ObjectStoreCredential{
		ID:        "12345",
		Name:      "test-objectstore-cred",
		MaxSizeGB: 500,
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestUpdateObjecStoreCredential(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/objectstore/credentials/12345": `{
			"id": "12345",
			"name": "test-objectstore-cred",
			"max_size_gb": 1000
		}`,
	})
	defer server.Close()

	cfg := &UpdateObjectStoreCredentialRequest{
		MaxSizeGB: intPtr(1000),
	}
	got, err := client.UpdateObjectStoreCredential("12345", cfg)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &ObjectStoreCredential{
		ID:        "12345",
		Name:      "test-objectstore-cred",
		MaxSizeGB: 1000,
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestDeleteObjectStoreCredential(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/objectstore/credentials/12345": `{"result": "success"}`,
	})
	defer server.Close()

	got, err := client.DeleteObjectStoreCredential("12345")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &SimpleResponse{Result: "success"}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func intPtr(i int) *int {
	return &i
}
