package civogo

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestListActions(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/actions": `{"page":1,"per_page":5,"pages":1,"items":[{"id":267531707,"created_at":"2022-10-10T16:30:11Z","updated_at":"2022-10-10T16:30:11Z","account_id":"eaef1dd6-1cec-4d9c-8480-96452bd94dea","user_id":"","type":"cluster-create","details":"Created a new cluster called cluster-kubectl","related_id":"259e96c5-ecd4-43c6-be35-d21dcf05b650","related_type":"cluster","debug":false},{"id":267531696,"created_at":"2022-10-10T16:29:50Z","updated_at":"2022-10-10T16:29:50Z","account_id":"eaef1dd6-1cec-4d9c-8480-96452bd94dea","user_id":"6b6e9801-8bcd-4049-84d7-028c8e748f58","type":"cluster-delete","details":"Deleted cluster : cluster-kubectl-c387-e36f8c","related_id":"7b107c00-28b1-49b8-b609-c850aeb2d72e","related_type":"cluster","debug":false},{"id":267531012,"created_at":"2022-10-10T12:08:09Z","updated_at":"2022-10-10T12:08:09Z","account_id":"eaef1dd6-1cec-4d9c-8480-96452bd94dea","user_id":"","type":"cluster-create","details":"Created a new cluster called cluster-kubectl","related_id":"7b107c00-28b1-49b8-b609-c850aeb2d72e","related_type":"cluster","debug":false},{"id":267527196,"created_at":"2022-10-09T12:40:37Z","updated_at":"2022-10-09T12:40:37Z","account_id":"eaef1dd6-1cec-4d9c-8480-96452bd94dea","user_id":"6b6e9801-8bcd-4049-84d7-028c8e748f58","type":"network-delete","details":"Deleted a network called cust-test-kubectl-eaef1dd6-5b25-58204a62","related_id":"0d7be2cc-515d-44c2-8b84-1dbd686806eb","related_type":"network","debug":false},{"id":267527195,"created_at":"2022-10-09T12:40:22Z","updated_at":"2022-10-09T12:40:22Z","account_id":"eaef1dd6-1cec-4d9c-8480-96452bd94dea","user_id":"6b6e9801-8bcd-4049-84d7-028c8e748f58","type":"volume-delete","details":"Deleted volume : test-server-665b-21da40","related_id":"0b5313a7-9b59-43c0-b749-0a0195407a62","related_type":"volume","debug":false}]}`,
	})
	defer server.Close()

	actionListRequest := &ActionListRequest{}
	allActions, err := client.ListActions(actionListRequest)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	if len(allActions.Items) != 5 {
		t.Errorf("Expected %d, got %d", 5, len(allActions.Items))
	}
}

func TestListActionsWithFilter(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/actions": `{"page":1,"per_page":5,"pages":1,"items":[{"id":267531707,"created_at":"2022-10-10T16:30:11Z","updated_at":"2022-10-10T16:30:11Z","account_id":"eaef1dd6-1cec-4d9c-8480-96452bd94dea","user_id":"","type":"cluster-create","details":"Created a new cluster called cluster-kubectl","related_id":"259e96c5-ecd4-43c6-be35-d21dcf05b650","related_type":"cluster","debug":false}]}`,
	})
	defer server.Close()

	actionListRequest := &ActionListRequest{
		RelatedID: "259e96c5-ecd4-43c6-be35-d21dcf05b650",
	}
	allActions, err := client.ListActions(actionListRequest)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	if len(allActions.Items) != 1 {
		t.Errorf("Expected %d, got %d", 1, len(allActions.Items))
	}
}

// TestListAllActions verifies the basic API contract of the new
// auto-iterating wrapper: with a server that reports a single page, the
// returned []Action contains every item.
func TestListAllActions(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/actions": `{"page":1,"per_page":100,"pages":1,"items":[{"id":1,"type":"instance-create","related_type":"instance"},{"id":2,"type":"instance-delete","related_type":"instance"}]}`,
	})
	defer server.Close()

	got, err := client.ListAllActions(nil)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	if len(got) != 2 {
		t.Errorf("Expected 2 actions, got %d", len(got))
	}
	if got[0].ID != 1 || got[1].ID != 2 {
		t.Errorf("Expected actions with IDs [1, 2], got [%d, %d]", got[0].ID, got[1].ID)
	}
}

// TestListAllActions_AcrossPages verifies that ListAllActions iterates a
// multi-page response and returns every item in order. Uses the multi-page
// mock server helper from pagination_test.go (same package).
func TestListAllActions_AcrossPages(t *testing.T) {
	const total = paginationPerPage*2 + 7 // 3 pages
	client, server := newMultiPageServer(t, "/v2/actions", total, paginationPerPage, func(i int) string {
		return fmt.Sprintf(`{"id":%d,"type":"instance-create","related_type":"instance"}`, i)
	})
	defer server.Close()

	got, err := client.ListAllActions(nil)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	if len(got) != total {
		t.Errorf("Expected %d actions, got %d", total, len(got))
	}
	for i, action := range got {
		if action.ID != i {
			t.Errorf("Expected action[%d].ID == %d, got %d", i, i, action.ID)
		}
	}
}

// TestListAllActions_PreservesFilters verifies that filter fields on the
// caller-supplied ActionListRequest are forwarded to every page request,
// while any Page/PerPage values on the filters are overridden by the
// iterator (so a caller can't accidentally short-circuit pagination by
// passing a small per_page).
func TestListAllActions_PreservesFilters(t *testing.T) {
	const total = paginationPerPage + 3 // 2 pages
	var seenResourceType, seenPerPage string

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if rt := req.URL.Query().Get("resource_type"); rt != "" {
			seenResourceType = rt
		}
		if pp := req.URL.Query().Get("per_page"); pp != "" {
			seenPerPage = pp
		}
		page := req.URL.Query().Get("page")
		if page == "" {
			page = "1"
		}
		start := 0
		if page == "2" {
			start = paginationPerPage
		}
		var items []string
		for i := start; i < start+paginationPerPage && i < total; i++ {
			items = append(items, fmt.Sprintf(`{"id":%d,"type":"instance-create"}`, i))
		}
		fmt.Fprintf(rw, `{"page":%s,"per_page":%d,"pages":2,"items":[%s]}`,
			page, paginationPerPage, strings.Join(items, ","))
	}))
	defer server.Close()
	client, err := NewClientForTestingWithServer(server)
	if err != nil {
		t.Fatalf("NewClientForTestingWithServer: %v", err)
	}

	// Caller passes a filter AND bogus Page/PerPage values. Filter should
	// be forwarded; iterator should override Page/PerPage.
	filters := &ActionListRequest{
		ResourceType: "instance",
		Page:         999,
		PerPage:      1,
	}
	got, err := client.ListAllActions(filters)
	if err != nil {
		t.Fatalf("ListAllActions: %v", err)
	}
	if len(got) != total {
		t.Errorf("Expected %d actions, got %d", total, len(got))
	}
	if seenResourceType != "instance" {
		t.Errorf("Expected resource_type=instance to be forwarded, got %q", seenResourceType)
	}
	if seenPerPage != fmt.Sprintf("%d", paginationPerPage) {
		t.Errorf("Expected per_page=%d (iterator override), got %q", paginationPerPage, seenPerPage)
	}
}
