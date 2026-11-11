// Package civogo is the official Go SDK for the Civo v3 API.
//
// A Client is immutable after construction and safe for concurrent use:
//
//	client, err := civogo.New(token, civogo.WithDefaultRegion("lon1"))
//	if err != nil { ... }
//
//	inst, err := client.Instances.Get(ctx, instanceID)
//
// List endpoints accept multiple regions and report real server-side
// pagination; use All for lazy iteration over every page:
//
//	for inst, err := range client.Instances.All(ctx, instances.ListRequest{Regions: []string{"lon1", "fra1"}}) {
//		if err != nil { ... }
//		fmt.Println(inst.Hostname)
//	}
package civogo

import (
	"net/http"

	"github.com/civo/civogo/v3/api"
	"github.com/civo/civogo/v3/instances"
)

// DefaultBaseURL is the production Civo API endpoint.
const DefaultBaseURL = "https://api.civo.com"

// Error is the single error type returned for any v3 API failure.
type Error = api.Error

// Error-classification helpers, re-exported from the api package.
var (
	IsNotFound        = api.IsNotFound
	IsAlreadyExists   = api.IsAlreadyExists
	IsAuthFailure     = api.IsAuthFailure
	IsRateLimited     = api.IsRateLimited
	IsInvalidArgument = api.IsInvalidArgument
	HasReason         = api.HasReason
)

// Ptr returns a pointer to v, for setting optional request fields:
//
//	client.Instances.Update(ctx, id, instances.UpdateInstanceRequest{PublicIPRequired: civogo.Ptr(false)})
func Ptr[T any](v T) *T { return &v }

// Client is the entry point to the Civo v3 API. Construct it with New.
type Client struct {
	Instances *instances.Service
}

// Option configures a Client at construction time.
type Option func(*api.Config)

// WithBaseURL points the client at a non-production API endpoint.
func WithBaseURL(u string) Option { return func(c *api.Config) { c.BaseURL = u } }

// WithDefaultRegion sets the region used when a request does not
// specify any.
func WithDefaultRegion(region string) Option {
	return func(c *api.Config) { c.DefaultRegion = region }
}

// WithHTTPClient sets the *http.Client used for requests.
func WithHTTPClient(hc *http.Client) Option { return func(c *api.Config) { c.HTTPClient = hc } }

// WithUserAgent prepends an application-specific User-Agent component.
func WithUserAgent(ua string) Option {
	return func(c *api.Config) { c.UserAgent = ua + " " + c.UserAgent }
}

// New builds a Client authenticated with token.
func New(token string, opts ...Option) (*Client, error) {
	cfg := api.Config{
		BaseURL:   DefaultBaseURL,
		Token:     token,
		UserAgent: "civogo/v3",
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	t, err := api.NewTransport(cfg)
	if err != nil {
		return nil, err
	}
	return &Client{
		Instances: instances.NewService(t),
	}, nil
}
