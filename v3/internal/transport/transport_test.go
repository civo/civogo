package transport

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/civo/civogo/v3/api"
)

func testTransport(t *testing.T, handler http.HandlerFunc) *Transport {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	tr, err := NewTransport(Config{BaseURL: server.URL, Token: "test-token", UserAgent: "civogo-test", DefaultRegion: "lon1"})
	if err != nil {
		t.Fatal(err)
	}
	return tr
}

func TestDoSendsAuthAndDecodes(t *testing.T) {
	tr := testTransport(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Errorf("Authorization = %q", got)
		}
		if got := r.Header.Get("User-Agent"); got != "civogo-test" {
			t.Errorf("User-Agent = %q", got)
		}
		if r.Header.Get("Content-Type") != "" {
			t.Errorf("Content-Type set on bodyless request")
		}
		fmt.Fprint(w, `{"name":"x"}`)
	})

	var out struct {
		Name string `json:"name"`
	}
	if err := tr.Do(context.Background(), "GET", "/v3/things", nil, nil, &out); err != nil {
		t.Fatal(err)
	}
	if out.Name != "x" {
		t.Errorf("decoded %+v", out)
	}
}

func TestDoDecodesErrorEnvelope(t *testing.T) {
	tr := testTransport(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"error":{"code":404,"status":"NOT_FOUND","message":"instance not found","details":[`+
			`{"@type":"type.googleapis.com/google.rpc.ErrorInfo","reason":"INSTANCE_NOT_FOUND","domain":"compute.civo.com","metadata":{"request_id":"req-1"}}]}}`)
	})

	err := tr.Do(context.Background(), "GET", "/v3/things/nope", nil, nil, nil)
	var apiErr *api.Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("want *Error, got %T: %v", err, err)
	}
	if apiErr.HTTPStatus != 404 || apiErr.Status != api.StatusNotFound || apiErr.Message != "instance not found" ||
		apiErr.Reason != "INSTANCE_NOT_FOUND" || apiErr.Domain != "compute.civo.com" || apiErr.RequestID != "req-1" ||
		len(apiErr.Details) != 1 {
		t.Errorf("unexpected error: %+v", apiErr)
	}
	if !api.IsNotFound(err) || !api.HasReason(err, "INSTANCE_NOT_FOUND") {
		t.Error("IsNotFound/HasReason = false")
	}
	if api.IsRateLimited(err) || api.HasReason(err, "SIZE_NOT_FOUND") {
		t.Error("IsRateLimited/HasReason matched wrongly")
	}
}

func TestDoDecodesViolations(t *testing.T) {
	tr := testTransport(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"error":{"code":400,"status":"INVALID_ARGUMENT","message":"request validation failed","details":[`+
			`{"@type":"type.googleapis.com/google.rpc.ErrorInfo","reason":"REQUEST_VALIDATION_FAILED","domain":"compute.civo.com"},`+
			`{"@type":"type.googleapis.com/google.rpc.BadRequest","fieldViolations":[{"field":"account_id","description":"must be a valid UUID"}]},`+
			`{"@type":"type.googleapis.com/google.rpc.QuotaFailure","violations":[{"subject":"instances","description":"limit 10"}]},`+
			`{"@type":"type.googleapis.com/google.rpc.Unmodelled","x":1}]}}`)
	})

	err := tr.Do(context.Background(), "POST", "/v3/things", nil, nil, nil)
	var apiErr *api.Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("want *Error, got %T: %v", err, err)
	}
	if !api.HasReason(err, "REQUEST_VALIDATION_FAILED") || !api.IsInvalidArgument(err) {
		t.Errorf("reason/status not decoded: %+v", apiErr)
	}
	if len(apiErr.FieldViolations) != 1 || apiErr.FieldViolations[0] != (api.FieldViolation{Field: "account_id", Description: "must be a valid UUID"}) {
		t.Errorf("FieldViolations = %+v", apiErr.FieldViolations)
	}
	if len(apiErr.QuotaViolations) != 1 || apiErr.QuotaViolations[0].Subject != "instances" {
		t.Errorf("QuotaViolations = %+v", apiErr.QuotaViolations)
	}
	if len(apiErr.Details) != 4 {
		t.Errorf("Details kept %d of 4", len(apiErr.Details))
	}
	if want := "account_id: must be a valid UUID"; !strings.Contains(err.Error(), want) {
		t.Errorf("Error() = %q, want it to contain %q", err.Error(), want)
	}
}

func TestDoKeepsUnstructuredErrorBodies(t *testing.T) {
	tr := testTransport(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprint(w, "<html>bad gateway</html>")
	})

	err := tr.Do(context.Background(), "GET", "/", nil, nil, nil)
	var apiErr *api.Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("want *Error, got %T", err)
	}
	if apiErr.HTTPStatus != 502 || apiErr.Message != "<html>bad gateway</html>" {
		t.Errorf("unexpected error: %+v", apiErr)
	}
}

func TestIsNotFoundIgnoresBare404(t *testing.T) {
	tr := testTransport(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "404 page not found")
	})

	err := tr.Do(context.Background(), "GET", "/", nil, nil, nil)
	if err == nil || api.IsNotFound(err) {
		t.Errorf("IsNotFound = true for a bare 404: %v", err)
	}
}

func TestDoHonoursContextCancellation(t *testing.T) {
	tr := testTransport(t, func(w http.ResponseWriter, r *http.Request) {})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := tr.Do(ctx, "GET", "/", nil, nil, nil); !errors.Is(err, context.Canceled) {
		t.Errorf("want context.Canceled, got %v", err)
	}
}

func TestNewTransportRequiresToken(t *testing.T) {
	if _, err := NewTransport(Config{BaseURL: "https://api.civo.com"}); err == nil {
		t.Error("want error for missing token")
	}
}

func TestAllIteratesEveryPage(t *testing.T) {
	fetched := 0
	fetch := func(p api.PageParams) ([]int, *api.PageInfo, error) {
		fetched++
		items := []int{p.Page*10 + 1, p.Page*10 + 2}
		return items, &api.PageInfo{Page: p.Page, TotalPages: 3}, nil
	}

	var got []int
	for v, err := range All(fetch) {
		if err != nil {
			t.Fatal(err)
		}
		got = append(got, v)
	}
	if len(got) != 6 || fetched != 3 {
		t.Errorf("got %v after %d fetches", got, fetched)
	}
}

func TestAllStopsOnBreak(t *testing.T) {
	fetched := 0
	fetch := func(p api.PageParams) ([]int, *api.PageInfo, error) {
		fetched++
		return []int{1, 2}, &api.PageInfo{Page: p.Page, TotalPages: 100}, nil
	}
	for _, err := range All(fetch) {
		if err != nil {
			t.Fatal(err)
		}
		break
	}
	if fetched != 1 {
		t.Errorf("fetched %d pages after break", fetched)
	}
}

func TestAllYieldsFetchError(t *testing.T) {
	boom := errors.New("boom")
	fetch := func(p api.PageParams) ([]int, *api.PageInfo, error) {
		if p.Page == 2 {
			return nil, nil, boom
		}
		return []int{1}, &api.PageInfo{Page: p.Page, TotalPages: 5}, nil
	}
	var errs []error
	var items int
	for _, err := range All(fetch) {
		if err != nil {
			errs = append(errs, err)
			continue
		}
		items++
	}
	if items != 1 || len(errs) != 1 || !errors.Is(errs[0], boom) {
		t.Errorf("items=%d errs=%v", items, errs)
	}
}

func TestResolve(t *testing.T) {
	tr := testTransport(t, func(w http.ResponseWriter, r *http.Request) {})

	sc, err := Resolve(tr, nil)
	if err != nil || sc != (Scope{Region: "lon1"}) {
		t.Errorf("default region not applied: %+v, %v", sc, err)
	}
	sc, err = Resolve(tr, []api.CallOption{api.WithRegion("fra1"), api.WithAccountID("acc-1")})
	if err != nil || sc != (Scope{Region: "fra1", AccountID: "acc-1"}) {
		t.Errorf("options not applied: %+v, %v", sc, err)
	}
	if q := sc.Query(); q.Get("region") != "fra1" || q.Get("account_id") != "acc-1" {
		t.Errorf("query = %v", q)
	}
	if _, err := Resolve(tr, []api.CallOption{api.WithRegion("lon1"), api.WithRegion("fra1")}); err == nil {
		t.Error("want error for two regions on a single-region call")
	}
	sc, err = Resolve(tr, []api.CallOption{api.WithRegion(""), api.WithRegion("fra1")})
	if err != nil || sc.Region != "fra1" {
		t.Errorf("empty region not ignored: %+v, %v", sc, err)
	}
	if sc, _ := Resolve(tr, []api.CallOption{api.WithRegion("")}); sc.Region != "lon1" {
		t.Errorf("empty region dropped the default: %+v", sc)
	}
}

func TestListQuery(t *testing.T) {
	tr := testTransport(t, func(w http.ResponseWriter, r *http.Request) {})

	if q := ListQuery(tr, nil); q.Encode() != "regions=lon1" {
		t.Errorf("default query = %s", q.Encode())
	}
	q := ListQuery(tr, []api.CallOption{api.WithRegion("lon1"), api.WithRegion("fra1"), api.WithAccountID("acc-1")})
	if q.Encode() != "account_id=acc-1&regions=lon1&regions=fra1" {
		t.Errorf("query = %s", q.Encode())
	}
}

func TestAllStopsOnEmptyResult(t *testing.T) {
	for _, info := range []*api.PageInfo{{TotalPages: 0}, nil} {
		fetched := 0
		fetch := func(p api.PageParams) ([]int, *api.PageInfo, error) {
			fetched++
			if fetched > 3 {
				t.Fatal("still fetching")
			}
			return nil, info, nil
		}
		for _, err := range All(fetch) {
			t.Fatalf("yielded on an empty result: %v", err)
		}
		if fetched != 1 {
			t.Errorf("info %+v: fetched %d pages", info, fetched)
		}
	}
}

func TestErrorClassification(t *testing.T) {
	envelope := func(status, reason string) string {
		d := ""
		if reason != "" {
			d = `{"@type":"type.googleapis.com/google.rpc.ErrorInfo","reason":"` + reason + `"}`
		}
		return `{"error":{"status":"` + status + `","message":"m","details":[` + d + `]}}`
	}
	for _, tc := range []struct {
		name   string
		code   int
		body   string
		is     func(error) bool
		expect bool
	}{
		{"not found", 404, envelope("NOT_FOUND", "INSTANCE_NOT_FOUND"), api.IsNotFound, true},
		{"routing 404 without reason", 404, envelope("NOT_FOUND", ""), api.IsNotFound, false},
		{"already exists", 409, envelope("ALREADY_EXISTS", "X_EXISTS"), api.IsAlreadyExists, true},
		{"already exists is not not-found", 409, envelope("ALREADY_EXISTS", "X_EXISTS"), api.IsNotFound, false},
		{"bare 409", 409, "conflict", api.IsAlreadyExists, true},
		{"bare 401", 401, "", api.IsAuthFailure, true},
		{"bare 403", 403, "", api.IsAuthFailure, true},
		{"permission denied", 403, envelope("PERMISSION_DENIED", "R"), api.IsAuthFailure, true},
		{"bare 429", 429, "slow down", api.IsRateLimited, true},
		{"bare 400", 400, "bad", api.IsInvalidArgument, true},
		{"bare 500", 500, "boom", api.IsInvalidArgument, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := decodeError(tc.code, []byte(tc.body))
			if got := tc.is(err); got != tc.expect {
				t.Errorf("got %v for %v", got, err)
			}
		})
	}
}

func TestErrorMessageFallsBackToStatusText(t *testing.T) {
	if got := decodeError(502, nil).Error(); got != "civo: Bad Gateway (http 502)" {
		t.Errorf("Error() = %q", got)
	}
}

func TestDoBuildsURLAndHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.EscapedPath(); got != "/prefix/v3/things/a%20b%2Fc" {
			t.Errorf("path = %s", got)
		}
		if r.Header.Get("Content-Type") != "application/json" || r.Header.Get("Accept") != "application/json" {
			t.Errorf("headers = %v", r.Header)
		}
	}))
	t.Cleanup(server.Close)
	tr, err := NewTransport(Config{BaseURL: server.URL + "/prefix/", Token: "t"})
	if err != nil {
		t.Fatal(err)
	}
	if err := tr.Do(context.Background(), "POST", "/v3/things/a%20b%2Fc", nil, map[string]string{}, nil); err != nil {
		t.Fatal(err)
	}
}

func TestNewTransportRejectsBaseURLWithoutScheme(t *testing.T) {
	if _, err := NewTransport(Config{BaseURL: "api.civo.com", Token: "t"}); err == nil {
		t.Error("want error for a base URL without scheme")
	}
}
