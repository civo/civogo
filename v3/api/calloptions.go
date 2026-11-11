package api

import "net/url"

// CallParams scope a call whose only inputs are the target resource, a
// region and an optional account. They are sent as query parameters or as
// the JSON body, depending on the operation.
type CallParams struct {
	Region    string `json:"region,omitempty"`
	AccountID string `json:"account_id,omitempty"`
}

// CallOption customises the scoping of a single API call.
type CallOption func(*CallParams)

// WithRegion targets the call at a specific region, overriding the
// client-level default.
func WithRegion(region string) CallOption {
	return func(p *CallParams) { p.Region = region }
}

// WithAccountID scopes the call to a specific account.
func WithAccountID(accountID string) CallOption {
	return func(p *CallParams) { p.AccountID = accountID }
}

// ResolveCallParams applies opts and falls back to the transport's
// default region when none was chosen.
func ResolveCallParams(t *Transport, opts ...CallOption) CallParams {
	var p CallParams
	for _, opt := range opts {
		opt(&p)
	}
	if p.Region == "" {
		p.Region = t.DefaultRegion()
	}
	return p
}

// Query renders the parameters as URL query values.
func (p CallParams) Query() url.Values {
	q := url.Values{}
	if p.Region != "" {
		q.Set("region", p.Region)
	}
	if p.AccountID != "" {
		q.Set("account_id", p.AccountID)
	}
	return q
}
