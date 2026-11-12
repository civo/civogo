// Package transport is the request path shared by the generated API
// packages: auth, encoding, error decoding and pagination. It is internal
// so it can change without breaking users.
package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"iter"
	"net/http"
	"net/url"
	"strings"

	"github.com/civo/civogo/v3/api"
)

// Transport is the single request path for all v3 API calls. It is
// immutable after construction and safe for concurrent use.
type Transport struct {
	baseURL       *url.URL
	token         string
	userAgent     string
	defaultRegion string
	httpClient    *http.Client
}

// Config is populated by the root package's options.
type Config struct {
	BaseURL       string
	Token         string
	UserAgent     string
	DefaultRegion string
	HTTPClient    *http.Client
}

// NewTransport validates cfg and builds a Transport.
func NewTransport(cfg Config) (*Transport, error) {
	if cfg.Token == "" {
		return nil, fmt.Errorf("civo: no API token supplied")
	}
	u, err := url.Parse(cfg.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("civo: invalid base URL %q: %w", cfg.BaseURL, err)
	}
	if u.Scheme == "" || u.Host == "" {
		return nil, fmt.Errorf("civo: invalid base URL %q: want scheme and host", cfg.BaseURL)
	}
	hc := cfg.HTTPClient
	if hc == nil {
		hc = &http.Client{}
	}
	return &Transport{
		baseURL:       u,
		token:         cfg.Token,
		userAgent:     cfg.UserAgent,
		defaultRegion: cfg.DefaultRegion,
		httpClient:    hc,
	}, nil
}

// Do performs one API call: it encodes body as JSON (when non-nil),
// sends method to path (already escaped) with query, and decodes the JSON response into
// out (when non-nil). Any non-2xx response is returned as *Error.
func (t *Transport) Do(ctx context.Context, method, path string, query url.Values, body, out any) error {
	u := *t.baseURL
	u.RawPath = strings.TrimSuffix(t.baseURL.EscapedPath(), "/") + path
	p, err := url.PathUnescape(u.RawPath)
	if err != nil {
		return fmt.Errorf("civo: invalid path %q: %w", path, err)
	}
	u.Path = p
	if len(query) > 0 {
		u.RawQuery = query.Encode()
	}

	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("civo: encoding request body: %w", err)
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), reqBody)
	if err != nil {
		return fmt.Errorf("civo: building request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+t.token)
	if t.userAgent != "" {
		req.Header.Set("User-Agent", t.userAgent)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("civo: reading response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return decodeError(resp.StatusCode, respBody)
	}

	if out != nil {
		if err := json.Unmarshal(respBody, out); err != nil {
			return fmt.Errorf("civo: decoding response: %w", err)
		}
	}
	return nil
}

// errorEnvelope is the AIP-193 error body the v3 gateway writes.
type errorEnvelope struct {
	Error *struct {
		Code    int               `json:"code"` // HTTP status, not the gRPC code
		Status  string            `json:"status"`
		Message string            `json:"message"`
		Details []json.RawMessage `json:"details"`
	} `json:"error"`
}

// detail is the union of the google.rpc detail types, discriminated by @type.
type detail struct {
	Type            string               `json:"@type"`
	Reason          string               `json:"reason"`
	Domain          string               `json:"domain"`
	Metadata        map[string]string    `json:"metadata"`
	FieldViolations []api.FieldViolation `json:"fieldViolations"`
	Violations      []api.QuotaViolation `json:"violations"`
}

const typePrefix = "type.googleapis.com/google.rpc."

func decodeError(httpStatus int, body []byte) *api.Error {
	var env errorEnvelope
	if err := json.Unmarshal(body, &env); err != nil || env.Error == nil {
		// Not a structured envelope (e.g. a proxy error page):
		// keep the raw body as the message so nothing is lost.
		return &api.Error{HTTPStatus: httpStatus, Message: strings.TrimSpace(string(body))}
	}

	e := &api.Error{
		HTTPStatus: httpStatus,
		Status:     env.Error.Status,
		Message:    env.Error.Message,
		Details:    env.Error.Details,
	}
	for _, raw := range env.Error.Details {
		var d detail
		if json.Unmarshal(raw, &d) != nil {
			continue
		}
		switch d.Type {
		case typePrefix + "ErrorInfo":
			e.Reason, e.Domain, e.Metadata = d.Reason, d.Domain, d.Metadata
			e.RequestID = d.Metadata["request_id"]
		case typePrefix + "BadRequest":
			e.FieldViolations = append(e.FieldViolations, d.FieldViolations...)
		case typePrefix + "QuotaFailure":
			e.QuotaViolations = append(e.QuotaViolations, d.Violations...)
		}
	}
	return e
}

// Scope is the region and account a single-region call targets. It is
// sent as query parameters, or merged into the JSON body.
type Scope struct {
	Region    string `json:"region,omitempty"`
	AccountID string `json:"account_id,omitempty"`
}

// Query renders the scope as URL query values.
func (s Scope) Query() url.Values {
	q := url.Values{}
	if s.Region != "" {
		q.Set("region", s.Region)
	}
	if s.AccountID != "" {
		q.Set("account_id", s.AccountID)
	}
	return q
}

// Resolve applies opts for a single-region call, falling back to the
// client's default region. More than one region is an error.
func Resolve(t *Transport, opts []api.CallOption) (Scope, error) {
	p := t.apply(opts)
	if len(p.Regions) > 1 {
		return Scope{}, fmt.Errorf("civo: this call takes one region, got %d (%s)", len(p.Regions), strings.Join(p.Regions, ", "))
	}
	s := Scope{AccountID: p.AccountID}
	if len(p.Regions) == 1 {
		s.Region = p.Regions[0]
	}
	return s, nil
}

// ListQuery applies opts for a list call: every region given, or the
// client's default region.
func ListQuery(t *Transport, opts []api.CallOption) url.Values {
	p := t.apply(opts)
	q := url.Values{}
	if len(p.Regions) > 0 {
		q["regions"] = p.Regions
	}
	if p.AccountID != "" {
		q.Set("account_id", p.AccountID)
	}
	return q
}

func (t *Transport) apply(opts []api.CallOption) api.CallParams {
	var p api.CallParams
	for _, opt := range opts {
		opt(&p)
	}
	if len(p.Regions) == 0 && t.defaultRegion != "" {
		p.Regions = []string{t.defaultRegion}
	}
	return p
}

const allPerPage = 100

// FetchPage fetches one page of a list endpoint, returning the items on
// that page plus the server's pagination metadata.
type FetchPage[T any] func(p api.PageParams) (items []T, info *api.PageInfo, err error)

// All lazily iterates every item of a paginated list endpoint, fetching
// pages as needed. It stops at the first error, yielded with a zero item.
func All[T any](fetch FetchPage[T]) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		for page := 1; ; page++ {
			items, info, err := fetch(api.PageParams{Page: page, PerPage: allPerPage})
			if err != nil {
				var zero T
				yield(zero, err)
				return
			}
			for _, item := range items {
				if !yield(item, nil) {
					return
				}
			}
			if info == nil || page >= info.TotalPages {
				return
			}
		}
	}
}
