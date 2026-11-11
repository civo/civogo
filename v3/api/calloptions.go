package api

// CallParams scope a call: the regions it targets and, rarely, the account
// it acts on. Build them with CallOptions.
type CallParams struct {
	Regions   []string
	AccountID string
}

// CallOption customises the scoping of a single API call.
type CallOption func(*CallParams)

// WithRegion targets the call at region instead of the client's default
// region. List calls accept it several times and query every region given;
// any other call fails if given more than one.
func WithRegion(region string) CallOption {
	return func(p *CallParams) { p.Regions = append(p.Regions, region) }
}

// WithAccountID scopes the call to a specific account. Only needed when
// authenticating with a JWT: an API key already identifies its account.
func WithAccountID(accountID string) CallOption {
	return func(p *CallParams) { p.AccountID = accountID }
}
