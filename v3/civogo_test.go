package civogo_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	civogo "github.com/civo/civogo/v3"
	"github.com/civo/civogo/v3/api"
	"github.com/civo/civogo/v3/compute"
)

func testClient(t *testing.T, handler http.HandlerFunc, opts ...civogo.Option) *civogo.Client {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	client, err := civogo.New("test-token", append([]civogo.Option{civogo.WithBaseURL(server.URL)}, opts...)...)
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func TestListInstancesMultiRegion(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v3alpha1/compute/instances" {
			t.Errorf("path = %s", r.URL.Path)
		}
		regions := r.URL.Query()["regions"]
		if len(regions) != 2 || regions[0] != "lon1" || regions[1] != "fra1" {
			t.Errorf("regions = %v", regions)
		}
		if got := r.URL.Query().Get("pagination.per_page"); got != "50" {
			t.Errorf("per_page = %q", got)
		}
		fmt.Fprint(w, `{
			"instances": [{"id": "i-1", "hostname": "a", "region": "lon1"}, {"id": "i-2", "hostname": "b", "region": "fra1"}],
			"page_info": {"page": 1, "per_page": 50, "items_count": 2, "total_items_count": 2, "total_pages": 1}
		}`)
	})

	resp, err := client.Compute.Instances.List(context.Background(), compute.ListInstancesRequest{
		Regions: []string{"lon1", "fra1"},
		Page:    api.PageParams{PerPage: 50},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Instances) != 2 || resp.PageInfo.TotalItemsCount != 2 {
		t.Errorf("resp = %+v", resp)
	}
	if resp.Instances[1].Region != "fra1" {
		t.Errorf("instance = %+v", resp.Instances[1])
	}
}

func TestListInstancesUsesDefaultRegion(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if regions := r.URL.Query()["regions"]; len(regions) != 1 || regions[0] != "nyc1" {
			t.Errorf("regions = %v", regions)
		}
		fmt.Fprint(w, `{"instances": [], "page_info": {"total_pages": 0}}`)
	}, civogo.WithDefaultRegion("nyc1"))

	if _, err := client.Compute.Instances.List(context.Background(), compute.ListInstancesRequest{}); err != nil {
		t.Fatal(err)
	}
}

func TestAllPaginatesLazily(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("pagination.page")
		items := map[string]string{
			"1": `[{"id": "i-1"}, {"id": "i-2"}]`,
			"2": `[{"id": "i-3"}]`,
		}[page]
		fmt.Fprintf(w, `{"instances": %s, "page_info": {"page": %s, "total_pages": 2}}`, items, page)
	}, civogo.WithDefaultRegion("lon1"))

	var ids []string
	for inst, err := range client.Compute.Instances.All(context.Background(), compute.ListInstancesRequest{}) {
		if err != nil {
			t.Fatal(err)
		}
		ids = append(ids, inst.ID)
	}
	if len(ids) != 3 || ids[2] != "i-3" {
		t.Errorf("ids = %v", ids)
	}
}

func TestGetInstanceUnwraps(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v3alpha1/compute/instances/i-123" {
			t.Errorf("path = %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("region"); got != "fra1" {
			t.Errorf("region = %q", got)
		}
		fmt.Fprint(w, `{"instance": {"id": "i-123", "hostname": "web-1", "status": "INSTANCE_STATUS_ACTIVE"}}`)
	})

	inst, err := client.Compute.Instances.Get(context.Background(), "i-123", api.WithRegion("fra1"))
	if err != nil {
		t.Fatal(err)
	}
	if inst.ID != "i-123" || inst.Status != compute.InstanceStatusActive {
		t.Errorf("instance = %+v", inst)
	}
}

func TestCreateInstanceDefaultsRegionWithoutMutatingRequest(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body["region"] != "lon1" || body["hostname"] != "web-1" {
			t.Errorf("body = %v", body)
		}
		fmt.Fprint(w, `{"instances": [{"id": "i-new", "region": "lon1"}]}`)
	}, civogo.WithDefaultRegion("lon1"))

	req := compute.CreateInstanceRequest{Hostname: "web-1"}
	created, err := client.Compute.Instances.Create(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if len(created) != 1 || created[0].ID != "i-new" {
		t.Errorf("created = %+v", created)
	}
	if req.Region != "" {
		t.Errorf("caller's request was mutated: %+v", req)
	}
}

func TestStopSendsRegionInBody(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v3alpha1/compute/instances/i-1/stop" {
			t.Errorf("path = %s", r.URL.Path)
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body["region"] != "fra1" {
			t.Errorf("body = %v", body)
		}
		fmt.Fprint(w, `{}`)
	}, civogo.WithDefaultRegion("lon1"))

	if err := client.Compute.Instances.Stop(context.Background(), "i-1", api.WithRegion("fra1")); err != nil {
		t.Fatal(err)
	}
}

func TestNotFoundErrorFlow(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"error":{"code":404,"status":"NOT_FOUND","message":"instance not found","details":[{"@type":"type.googleapis.com/google.rpc.ErrorInfo","reason":"INSTANCE_NOT_FOUND","domain":"compute.civo.com"}]}}`)
	})

	_, err := client.Compute.Instances.Get(context.Background(), "i-nope")
	if err == nil {
		t.Fatal("want error")
	}
	if !civogo.IsNotFound(err) {
		t.Errorf("IsNotFound = false for %v", err)
	}
	var apiErr *civogo.Error
	if !errors.As(err, &apiErr) || apiErr.HTTPStatus != 404 || apiErr.Reason != "INSTANCE_NOT_FOUND" {
		t.Errorf("error = %#v", err)
	}
}

func TestGetEscapesID(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.EscapedPath(); got != "/v3alpha1/compute/instances/a%20b%2Fc" {
			t.Errorf("path = %s", got)
		}
		fmt.Fprint(w, `{"instance": {"id": "x"}}`)
	})
	if _, err := client.Compute.Instances.Get(context.Background(), "a b/c"); err != nil {
		t.Fatal(err)
	}
}

func TestGetRejectsMissingWrapper(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{}`)
	})
	if inst, err := client.Compute.Instances.Get(context.Background(), "i-1"); err == nil {
		t.Errorf("want error, got %+v", inst)
	}
}

func TestWithUserAgent(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("User-Agent"); got != "myapp/1 civogo/v3" {
			t.Errorf("User-Agent = %q", got)
		}
		fmt.Fprint(w, `{"instance": {"id": "x"}}`)
	}, civogo.WithUserAgent("myapp/1"))
	if _, err := client.Compute.Instances.Get(context.Background(), "x"); err != nil {
		t.Fatal(err)
	}
}
