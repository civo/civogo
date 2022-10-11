package civogo

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"strconv"
	"sync"
	"time"

	"github.com/civo/civogo/utils"
	"github.com/google/uuid"
)

const (
	libraryVersion = "1.0.0"
	defaultBaseURL = "https://api.civo.com/"
	userAgent      = "civogo/" + libraryVersion
	mediaType      = "application/json"

	headerRateLimit     = "X-Ratelimit-Limit"
	headerRateRemaining = "X-Ratelimit-Remaining"
	headerRateReset     = "X-Ratelimit-Reset"
	headerRequestID     = "X-Request-ID"
)

// Interface is the interface for the Civo API
type ClientInterface interface {
	SSHKeyGetter
	NetworkGetter
	DNSGetter
	InstancesGetter
}

// Client manages communication with the Civo API.
type Client struct {
	// HTTP client used to communicate with the Civo API.
	client *http.Client

	// Base URL for API requests.
	BaseURL *url.URL

	// User agent for client
	userAgent string

	// APIkey to be used for authentication
	APIKey string

	// Region to be used for authentication
	Region string

	// Services used for communicating with the API
	Rate    Rate
	ratemtx sync.Mutex

	// Services used for communicating with the API
	/*
		Account AccountService
		Application ApplicationService
		Charge ChargeService
		DiskImage DiskImageService

		Firewall FirewallService
		Instance InstanceService
		IP IPService
		Kubernetes KubernetesService
		LoadBalancer LoadBalancerService
		ObjectStorage ObjectStorageService
		Qouta QuotaService
		Region RegionService
		Team TeamService
		User UserService
		Volume VolumeService
		Webhook WebhookService
	*/

	// Optional function called after every successful request made to the Civo APIs
	onRequestCompleted RequestCompletionCallback

	// Optional extra HTTP headers to set on every request to the API.
	headers map[string]string
}

// Data is to manage the data returned by the API
type Data struct {
	Meta Metadata
	Data interface{}
}

// Metadata is a Civo API response. This wraps the standard http.Response
type Metadata struct {
	Rate
	Header        http.Header
	Status        string
	StatusCode    int
	ContentLength int64
}

// An ErrorResponse reports the error caused by an API request
type ErrorResponse struct {
	// HTTP response that caused this error
	Response   *http.Response
	RequestID  string
	Code       string `json:"code,omitempty"`
	Statuscode int
	Status     string
	Reason     string `json:"reason,omitempty"`
}

// SimpleResponse is a structure that returns success and/or any error
type SimpleResponse struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Result string `json:"result"`
}

// Rate contains the rate limit for the current client.
type Rate struct {
	// The number of request per hour the client is currently limited to.
	Limit int `json:"limit"`

	// The number of remaining requests the client can make this hour.
	Remaining int `json:"remaining"`

	// The time at which the current rate limit will reset.
	Reset utils.Timestamp `json:"reset"`
}

// RequestCompletionCallback is the type of the function called after every request to the Civo API.
type RequestCompletionCallback func(*http.Request, *http.Response)

// ListOptions struct used for pagination
type ListOptions struct {
	// For paginated result sets, page of results to retrieve.
	Page int `url:"page,omitempty"`

	// For paginated result sets, the number of results to include per page.
	PerPage int `url:"per_page,omitempty"`
}

// NewClient returns a new Civo API client.
func NewClient(apiKey, region string) (*Client, error) {
	return NewClientWithURL(apiKey, region, defaultBaseURL)
}

// NewClientWithURL initializes a Client with a specific API URL
func NewClientWithURL(apiKey, region, civoAPIURL string) (*Client, error) {
	if apiKey == "" {
		err := errors.New("no API Key supplied, this is required")
		return nil, err
	}

	parsedURL, err := url.Parse(civoAPIURL)
	if err != nil {
		return nil, err
	}

	var httpTransport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
	}

	httpClient := http.DefaultClient
	httpClient.Transport = httpTransport

	client := &Client{
		client:    httpClient,
		BaseURL:   parsedURL,
		userAgent: userAgent,
		APIKey:    apiKey,
		Region:    region,
	}

	return client, nil
}

func (c *Client) SSHKey() SSHKeyService {
	return newSSHKey(c)
}

func (c *Client) Network() NetworkService {
	return newNetwork(c)
}

func (c *Client) DNS() DNSService {
	return newDNS(c)
}

func (c *Client) Instances(network string) InstancesService {
	return newInstances(c, network)
}

// NewClientWithOptions is a fcuntion to create a new client wiyth the given options
func NewClientWithOptions(apiKey, region, civoAPIURL string, options ...ClientOptions) (*Client, error) {
	client, err := NewClientWithURL(apiKey, region, civoAPIURL)
	if err != nil {
		return nil, err
	}

	for _, option := range options {
		if err := option(client); err != nil {
			return nil, err
		}
	}
	return client, nil
}

// ClientOptions are options for New.
type ClientOptions func(*Client) error

// SetUserAgent is a client option for setting the user agent.
func SetUserAgent(userAgent string) ClientOptions {
	return func(c *Client) error {
		c.userAgent = fmt.Sprintf("%s %s", userAgent, c.userAgent)
		return nil
	}
}

// SetRequestHeaders sets optional HTTP headers on the client that are
// sent on each HTTP request.
func SetRequestHeaders(headers map[string]string) ClientOptions {
	return func(c *Client) error {
		for k, v := range headers {
			c.headers[k] = v
		}
		return nil
	}
}

// SetBaseURL is a client option for set the URL
func SetBaseURL(baseURL string) ClientOptions {
	return func(c *Client) error {
		u, err := url.Parse(baseURL)
		if err != nil {
			return err
		}

		c.BaseURL = u
		return nil
	}
}

// NewRequest creates an API request. A relative URL can be provided in urlStr, which will be resolved to the
// BaseURL of the Client. Relative URLs should always be specified without a preceding slash.
func (c *Client) NewRequest(ctx context.Context, method, urlStr string, body interface{}) (*http.Request, error) {
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var req *http.Request
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		req, err = http.NewRequest(method, u.String(), nil)
		if err != nil {
			return nil, err
		}

	default:
		buf := new(bytes.Buffer)
		if body != nil {
			err = json.NewEncoder(buf).Encode(body)
			if err != nil {
				return nil, err
			}
		}

		req, err = http.NewRequest(method, u.String(), buf)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", mediaType)
	}

	for k, v := range c.headers {
		req.Header.Add(k, v)
	}

	if req.Method == http.MethodGet || req.Method == http.MethodDelete {
		// add the region param
		param := req.URL.Query()
		param.Add("region", c.Region)
		req.URL.RawQuery = param.Encode()
	}

	req.Header.Set("Accept", mediaType)
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", c.APIKey))
	req.Header.Set("x-Request-ID", uuid.NewString())

	return req, nil
}

// GetRate returns the rate limit for the current client.
func (c *Client) GetRate() Rate {
	c.ratemtx.Lock()
	defer c.ratemtx.Unlock()
	return c.Rate
}

// newResponse creates a new Response for the provided http.Response
func newResponse(r *http.Response) *Metadata {
	response := Metadata{Header: r.Header, StatusCode: r.StatusCode, ContentLength: r.ContentLength, Status: r.Status}
	response.populateRate()

	return &response
}

// populateRate parses the rate related headers and populates the response Rate.
func (r *Metadata) populateRate() {
	if limit := r.Header.Get(headerRateLimit); limit != "" {
		r.Rate.Limit, _ = strconv.Atoi(limit)
	}
	if remaining := r.Header.Get(headerRateRemaining); remaining != "" {
		r.Rate.Remaining, _ = strconv.Atoi(remaining)
	}
	if reset := r.Header.Get(headerRateReset); reset != "" {
		if v, _ := strconv.ParseInt(reset, 10, 64); v != 0 {
			r.Rate.Reset = utils.Timestamp{Time: time.Unix(v, 0)}
		}
	}
}

// Do sends an API request and returns the API response. The API response is JSON decoded and stored in the value
// pointed to by v, or returned as an error if an API error has occurred. If v implements the io.Writer interface,
// the raw response will be written to v, without attempting to decode it.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*Metadata, error) {
	resp, err := DoRequestWithClient(ctx, c.client, req)
	if err != nil {
		return nil, err
	}
	if c.onRequestCompleted != nil {
		c.onRequestCompleted(req, resp)
	}

	defer func() {
		// Ensure the response body is fully read and closed
		// before we reconnect, so that we reuse the same TCPConnection.
		// Close the previous response's body. But read at least some of
		// the body so if it's small the underlying TCP connection will be
		// re-used. No need to check for errors: if it fails, the Transport
		// won't reuse it anyway.
		const maxBodySlurpSize = 2 << 10
		if resp.ContentLength == -1 || resp.ContentLength <= maxBodySlurpSize {
			io.CopyN(io.Discard, resp.Body, maxBodySlurpSize)
		}

		if rerr := resp.Body.Close(); err == nil {
			err = rerr
		}
	}()

	response := newResponse(resp)
	c.ratemtx.Lock()
	c.Rate = response.Rate
	c.ratemtx.Unlock()

	err = CheckResponse(resp)
	if err != nil {
		return response, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
			if err != nil {
				return nil, err
			}
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
			if err != nil {
				return nil, err
			}
		}
	}

	return response, err
}

// DoRequest submits an HTTP request.
func DoRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	return DoRequestWithClient(ctx, http.DefaultClient, req)
}

// DoRequestWithClient submits an HTTP request using the specified client.
func DoRequestWithClient(
	ctx context.Context,
	client *http.Client,
	req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)
	return client.Do(req)
}

func (r *ErrorResponse) Error() string {
	if r.RequestID != "" {
		return fmt.Sprintf("%v %v: %d (request %q) %v",
			r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.RequestID, r.Reason)
	}
	return fmt.Sprintf("%v %v: %d %v",
		r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.Reason)
}

// CheckResponse checks the API response for errors, and returns them if present. A response is considered an
// error if it has a status code outside the 200 range. API error responses are expected to have either no response
// body, or a JSON response body that maps to ErrorResponse. Any other response body will be silently ignored.
// If the API error response does not include the request ID in its body, the one from its header will be used.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; c >= 200 && c <= 299 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: r, Statuscode: r.StatusCode, Status: r.Status}
	data, err := io.ReadAll(r.Body)
	if err == nil && len(data) > 0 {
		err := json.Unmarshal(data, errorResponse)
		if err != nil {
			errorResponse.Statuscode = r.StatusCode
			errorResponse.Reason = string(data)
			errorResponse.Status = r.Status
			errorResponse.RequestID = r.Header.Get(headerRequestID)
		}
	}

	return errorResponse
}

func (r Rate) String() string {
	return utils.Stringify(r)
}

// String is a helper routine that allocates a new string value
// to store v and returns a pointer to it.
func String(v string) *string {
	p := new(string)
	*p = v
	return p
}

// Int is a helper routine that allocates a new int32 value
// to store v and returns a pointer to it, but unlike Int32
// its argument value is an int.
func Int(v int) *int {
	p := new(int)
	*p = v
	return p
}

// Bool is a helper routine that allocates a new bool value
// to store v and returns a pointer to it.
func Bool(v bool) *bool {
	p := new(bool)
	*p = v
	return p
}

// StreamToString converts a reader to a string
func StreamToString(stream io.Reader) string {
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(stream)
	return buf.String()
}
