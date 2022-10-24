package civogo

import (
	"reflect"
	"testing"
)

func TestListDatabases(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/databases": `{"page": 1, "per_page": 20, "pages": 2, "items":[{"id": "12345", "name": "test-db"}]}`,
	})
	defer server.Close()

	got, err := client.ListDatabases()
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &PaginatedDatabases{
		Page:    1,
		PerPage: 20,
		Pages:   2,
		Items: []Database{
			{
				ID:   "12345",
				Name: "test-db",
			},
		},
	}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestFindDatabase(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/databases": `{
			"page": 1,
			"per_page": 20,
			"pages": 1,
			"items": [
			  {
				"id": "12345",
				"name": "test-db"
			  }
			]
		  }`,
	})
	defer server.Close()

	got, _ := client.FindDatabase("test-db")
	if got.ID != "12345" {
		t.Errorf("Expected %s, got %s", "12345", got.ID)
	}
}

func TestNewDatabase(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/databases": `{
			"id": "12345",
			"name": "test-db",
			"size": "g3.db.xsmall",
			"status" : "Ready"
		}`,
	})
	defer server.Close()

	cfg := &CreateDatabaseRequest{
		Name: "test-db",
		Size: "g3.db.xsmall",
	}
	got, err := client.NewDatabase(cfg)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &Database{
		ID:     "12345",
		Name:   "test-db",
		Size:   "g3.db.xsmall",
		Status: "Ready",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestDeleteDatabase(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/databases/12345": `{"result": "success"}`,
	})
	defer server.Close()

	got, err := client.DeleteDatabase("12345")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &SimpleResponse{Result: "success"}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}
