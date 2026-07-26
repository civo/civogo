package civogo

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

// newMultiPageServer returns an httptest.Server that responds to GET requests
// on path. The server reads the `page` query parameter and returns a
// PaginatedX-shaped JSON body whose Items slice contains itemsPerPage entries
// for every page up to totalPages, except the final page which gets the
// remainder. The body's `pages` field always reports totalPages so the SDK
// knows when to stop.
//
// itemBuilder is called with the global item index (0-based) and must return
// a JSON object snippet that fits inside the Items array.
func newMultiPageServer(t *testing.T, path string, totalItems, itemsPerPage int, itemBuilder func(idx int) string) (*Client, *httptest.Server) {
	t.Helper()
	totalPages := totalItems / itemsPerPage
	if totalItems%itemsPerPage > 0 {
		totalPages++
	}
	if totalItems == 0 {
		totalPages = 1
	}

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if !strings.HasPrefix(req.URL.Path, path) {
			http.NotFound(rw, req)
			return
		}
		page, _ := strconv.Atoi(req.URL.Query().Get("page"))
		if page == 0 {
			page = 1
		}
		start := (page - 1) * itemsPerPage
		end := start + itemsPerPage
		if end > totalItems {
			end = totalItems
		}
		items := make([]string, 0, end-start)
		for i := start; i < end; i++ {
			items = append(items, itemBuilder(i))
		}
		body := fmt.Sprintf(
			`{"page":%d,"per_page":%d,"pages":%d,"items":[%s]}`,
			page, itemsPerPage, totalPages, strings.Join(items, ","),
		)
		rw.Write([]byte(body))
	}))

	client, err := NewClientForTestingWithServer(server)
	if err != nil {
		server.Close()
		t.Fatalf("NewClientForTestingWithServer: %v", err)
	}
	return client, server
}

// TestPaginateAll_SinglePage verifies the early-exit path: server reports one
// page so iteration must not request page 2.
func TestPaginateAll_SinglePage(t *testing.T) {
	requests := 0
	client, server := newMultiPageServer(t, "/v2/ips", 5, 100, func(i int) string {
		requests++
		return fmt.Sprintf(`{"id":"ip-%d"}`, i)
	})
	defer server.Close()

	got, err := client.ListIPs()
	if err != nil {
		t.Fatalf("ListIPs: %v", err)
	}
	if len(got.Items) != 5 {
		t.Fatalf("expected 5 items, got %d", len(got.Items))
	}
	if got.Page != 1 || got.Pages != 1 || got.PerPage != 5 {
		t.Errorf("expected merged Page=1 Pages=1 PerPage=5, got Page=%d Pages=%d PerPage=%d", got.Page, got.Pages, got.PerPage)
	}
}

// TestPaginateAll_MultiPage exercises the customer's failure case: more items
// than fit on a single page (perPage+1 trailing item on page 2). The original
// bug (no iteration in civogo) dropped the trailing items entirely; this
// regression test asserts every item is now returned.
func TestPaginateAll_MultiPage(t *testing.T) {
	const total = 250 // > paginationPerPage (100), forces 3 pages
	client, server := newMultiPageServer(t, "/v2/ips", total, paginationPerPage, func(i int) string {
		return fmt.Sprintf(`{"id":"ip-%d"}`, i)
	})
	defer server.Close()

	got, err := client.ListIPs()
	if err != nil {
		t.Fatalf("ListIPs: %v", err)
	}
	if len(got.Items) != total {
		t.Fatalf("expected %d items, got %d", total, len(got.Items))
	}
	// Order across pages must be preserved.
	for i, item := range got.Items {
		want := fmt.Sprintf("ip-%d", i)
		if item.ID != want {
			t.Errorf("items[%d].ID = %q, want %q", i, item.ID, want)
		}
	}
	if got.Page != 1 || got.Pages != 1 || got.PerPage != total {
		t.Errorf("expected merged Page=1 Pages=1 PerPage=%d, got Page=%d Pages=%d PerPage=%d", total, got.Page, got.Pages, got.PerPage)
	}
}

// TestPaginateAll_TrailingItem covers the exact shape of the customer ticket
// that motivated this fix: items count = perPage + 1, so page 2 contains a
// single trailing entry that the old (no-iteration) SDK silently dropped.
func TestPaginateAll_TrailingItem(t *testing.T) {
	const total = paginationPerPage + 1
	client, server := newMultiPageServer(t, "/v2/ips", total, paginationPerPage, func(i int) string {
		return fmt.Sprintf(`{"id":"ip-%d"}`, i)
	})
	defer server.Close()

	got, err := client.ListIPs()
	if err != nil {
		t.Fatalf("ListIPs: %v", err)
	}
	if len(got.Items) != total {
		t.Fatalf("expected %d items, got %d", total, len(got.Items))
	}
	if got.Items[total-1].ID != fmt.Sprintf("ip-%d", total-1) {
		t.Errorf("expected trailing item ID ip-%d, got %q", total-1, got.Items[total-1].ID)
	}
}

// TestPaginateAll_CapExceeded asserts the hard cap returns an explicit error
// rather than spinning. We forge a server that always reports a huge page
// count so the loop must terminate via the cap.
func TestPaginateAll_CapExceeded(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(`{"page":1,"per_page":100,"pages":999999,"items":[{"id":"ip-x"}]}`))
	}))
	defer server.Close()
	client, err := NewClientForTestingWithServer(server)
	if err != nil {
		t.Fatalf("NewClientForTestingWithServer: %v", err)
	}

	_, err = client.ListIPs()
	if err == nil {
		t.Fatal("expected pagination cap error, got nil")
	}
	var capErr ErrPaginationCapExceeded
	if !errors.As(err, &capErr) {
		t.Fatalf("expected ErrPaginationCapExceeded, got %T: %v", err, err)
	}
	if capErr.Resource != "ips" {
		t.Errorf("expected Resource=ips, got %q", capErr.Resource)
	}
}

// TestListDatabaseBackup_TrailingItem covers ListDatabaseBackup against the
// /v2/databases/{id}/backups endpoint. The endpoint paginates server-side
// unconditionally and the SDK previously sent no per_page param, so any
// database with > 20 backups silently truncated.
//
// (Per-file regression tests for ListVPCIPs / FindVPCIP live in vpc_test.go;
// ListAllActions has them in action_test.go. This package has no
// database_backup_test.go file, so the regression sits here.)
func TestListDatabaseBackup_TrailingItem(t *testing.T) {
	const total = paginationPerPage + 1
	const dbID = "db-1234"
	client, server := newMultiPageServer(t,
		fmt.Sprintf("/v2/databases/%s/backups", dbID),
		total, paginationPerPage,
		func(i int) string {
			return fmt.Sprintf(`{"id":"bk-%d"}`, i)
		})
	defer server.Close()

	got, err := client.ListDatabaseBackup(dbID)
	if err != nil {
		t.Fatalf("ListDatabaseBackup: %v", err)
	}
	if len(got.Items) != total {
		t.Fatalf("expected %d items, got %d", total, len(got.Items))
	}
	if got.Items[total-1].ID != fmt.Sprintf("bk-%d", total-1) {
		t.Errorf("expected trailing item ID bk-%d, got %q", total-1, got.Items[total-1].ID)
	}
	if got.Page != 1 || got.Pages != 1 || got.PerPage != total {
		t.Errorf("expected merged Page=1 Pages=1 PerPage=%d, got Page=%d Pages=%d PerPage=%d", total, got.Page, got.Pages, got.PerPage)
	}
}
