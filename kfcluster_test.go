package civogo

import (
	"reflect"
	"testing"
)

func TestListKfClusters(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/kfclusters": `{"page": 1, "per_page": 20, "pages": 2, "items":[{"id": "12345", "name": "test-kfcluster"}]}`,
	})
	defer server.Close()

	got, err := client.ListKfClusters()
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &PaginatedKfClusters{
		Page:    1,
		PerPage: 20,
		Pages:   2,
		Items: []KfCluster{
			{
				ID:   "12345",
				Name: "test-kfcluster",
			},
		},
	}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestFindKfCluster(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/kfclusters": `{
			"page": 1,
			"per_page": 20,
			"pages": 1,
			"items": [
			  {
					"id": "12345",
					"name": "test-kfcluster"
			  },
				{
					"id": "12346",
					"name": "demo-kfcluster"
				}
			]
		  }`,
	})
	defer server.Close()

	// Exact Match
	got, _ := client.FindKfCluster("test-kfcluster")
	if got.ID != "12345" {
		t.Errorf("Expected %s, got %s", "12345", got.ID)
	}

	// Multiple Match
	_, err := client.FindKfCluster("kfcluster")
	if err.Error() != "MultipleMatchesError: unable to find kfcluster because there were multiple matches" {
		t.Errorf("Expected %s, got %s", "unable to find kfcluster as there were multiple matches found", err.Error())
	}

	// Zero Match
	_, err = client.FindKfCluster("missing")
	if err.Error() != "ZeroMatchesError: unable to find missing, zero matches" {
		t.Errorf("Expected %s, got %s", "unable to find missing, zero matches", err.Error())
	}
}

func TestCreateKfCluster(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/kfclusters": `{
			"id": "12345",
			"name": "test-kfcluster",
			"size": "g3.kf.small",
			"network_id": "09090"
		}`,
	})
	defer server.Close()

	cfg := CreateKfClusterReq{
		Name:      "test-kfcluster",
		Size:      "g3.kf.small",
		NetworkID: "09090",
	}
	got, err := client.CreateKfCluster(cfg)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &KfCluster{
		ID:        "12345",
		Name:      "test-kfcluster",
		Size:      "g3.kf.small",
		NetworkID: "09090",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestDeleteKfCluster(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/kfclusters/12345": `{"result": "success"}`,
	})
	defer server.Close()

	got, err := client.DeleteKfCluster("12345")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &SimpleResponse{Result: "success"}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}
