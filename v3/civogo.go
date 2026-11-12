// Package civogo is the official Go SDK for the Civo v3 API.
//
// A Client is immutable after construction and safe for concurrent use:
//
//	client, err := civogo.New(token, civogo.WithDefaultRegion("lon1"))
//	if err != nil { ... }
//
//	inst, err := client.Compute.Instances.Get(ctx, instanceID)
//
// Every call takes WithRegion to override the default region. List
// endpoints accept it several times to query multiple regions, and report
// real server-side pagination; use All for lazy iteration over every page:
//
//	for inst, err := range client.Compute.Instances.All(ctx, compute.ListInstancesRequest{}, civogo.WithRegion("lon1"), civogo.WithRegion("fra1")) {
//		if err != nil { ... }
//		fmt.Println(inst.Hostname)
//	}
package civogo

import (
	"net/http"

	"github.com/civo/civogo/v3/api"
	"github.com/civo/civogo/v3/compute"
	"github.com/civo/civogo/v3/internal/transport"
)

// DefaultBaseURL is the production Civo API endpoint.
const DefaultBaseURL = "https://api.civo.com"

// Error is the single error type returned for any v3 API failure.
type Error = api.Error

// IsNotFound reports whether err is a v3 API "not found" error.
func IsNotFound(err error) bool { return api.IsNotFound(err) }

// IsAlreadyExists reports whether err is a v3 API "already exists" error.
func IsAlreadyExists(err error) bool { return api.IsAlreadyExists(err) }

// IsAuthFailure reports whether err is an authentication or authorization failure.
func IsAuthFailure(err error) bool { return api.IsAuthFailure(err) }

// IsRateLimited reports whether err is a rate-limit / quota-exhausted error.
func IsRateLimited(err error) bool { return api.IsRateLimited(err) }

// IsInvalidArgument reports whether err is a request-validation error.
func IsInvalidArgument(err error) bool { return api.IsInvalidArgument(err) }

// HasReason reports whether err is a v3 API error with the given reason.
func HasReason(err error, reason string) bool { return api.HasReason(err, reason) }

// WithRegion targets a call at region instead of the client's default
// region. List calls accept it several times to query every region given.
func WithRegion(region string) api.CallOption { return api.WithRegion(region) }

// WithAccountID scopes a call to a specific account. Only needed when
// authenticating with a JWT: an API key already identifies its account.
func WithAccountID(accountID string) api.CallOption { return api.WithAccountID(accountID) }

// Ptr returns a pointer to v, for setting optional request fields:
//
//	client.Compute.Instances.Update(ctx, id, compute.UpdateInstanceRequest{PublicIPRequired: civogo.Ptr(false)})
func Ptr[T any](v T) *T { return &v }

// Client is the entry point to the Civo v3 API. Construct it with New.
type Client struct {
	Compute *compute.Client
}

// Option configures a Client at construction time.
type Option func(*transport.Config)

// WithBaseURL points the client at a non-production API endpoint.
func WithBaseURL(u string) Option { return func(c *transport.Config) { c.BaseURL = u } }

// WithDefaultRegion sets the region used when a request does not
// specify any.
func WithDefaultRegion(region string) Option {
	return func(c *transport.Config) { c.DefaultRegion = region }
}

// WithHTTPClient sets the *http.Client used for requests.
func WithHTTPClient(hc *http.Client) Option { return func(c *transport.Config) { c.HTTPClient = hc } }

// WithUserAgent prepends an application-specific User-Agent component.
func WithUserAgent(ua string) Option {
	return func(c *transport.Config) { c.UserAgent = ua + " " + c.UserAgent }
}

// New builds a Client authenticated with token.
func New(token string, opts ...Option) (*Client, error) {
	cfg := transport.Config{
		BaseURL:   DefaultBaseURL,
		Token:     token,
		UserAgent: "civogo/v3",
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	t, err := transport.NewTransport(cfg)
	if err != nil {
		return nil, err
	}
	return &Client{
		Compute: compute.NewClient(t),
	}, nil
}
