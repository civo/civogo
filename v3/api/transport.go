package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
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

// DefaultRegion returns the client-level default region, or "".
func (t *Transport) DefaultRegion() string { return t.defaultRegion }

// Do performs one API call: it encodes body as JSON (when non-nil),
// sends method to path with query, and decodes the JSON response into
// out (when non-nil). Any non-2xx response is returned as *Error.
func (t *Transport) Do(ctx context.Context, method, path string, query url.Values, body, out any) error {
	u := *t.baseURL
	u.Path = strings.TrimSuffix(u.Path, "/") + path
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
