package civogo

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Domain represents a domain registered within Civo's infrastructure
type Domain struct {
	// The ID of the domain
	ID string `json:"id"`

	// The ID of the account
	AccountID string `json:"account_id"`

	// The Name of the domain
	Name string `json:"name"`

	// The Result of the operation
	Result string `json:"result"`
}

// RecordType represents the allowed record types: a, cname, mx or txt
type RecordType string

// Record represents a DNS record registered within Civo's infrastructure
type Record struct {
	ID       string     `json:"id"`
	DomainID string     `json:"domain_id"`
	Name     string     `json:"name"`
	Value    string     `json:"value"`
	Type     RecordType `json:"type"`
	Priority int        `json:"priority"`
	TTL      int        `json:"ttl"`
}

// RecordConfig describes the parameters for a new DNS record
// none of the fields are mandatory and will be automatically
// set with default values
type RecordConfig struct {
	DomainID string     `json:"-"`
	Type     RecordType `json:"type"`
	Name     string     `json:"name"`
	Value    string     `json:"value"`
	Priority int        `json:"priority"`
	TTL      int        `json:"ttl"`
}

const (
	// RecordTypeA represents an A record
	RecordTypeA = "a"

	// RecordTypeCName represents an CNAME record
	RecordTypeCName = "cname"

	// RecordTypeMX represents an MX record
	RecordTypeMX = "mx"

	// RecordTypeTXT represents an TXT record
	RecordTypeTXT = "txt"
)

// ListDomains returns all Domains owned by the calling API account
func (c *Client) ListDomains() ([]Domain, error) {
	url := "/v2/dns"

	resp, err := c.SendGetRequest(url)
	if err != nil {
		return nil, err
	}

	var ds []Domain
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&ds); err != nil {
		return nil, err

	}

	return ds, nil
}

// GetDomain returns the Domain that matches the name
func (c *Client) GetDomain(name string) (*Domain, error) {
	ds, err := c.ListDomains()
	if err != nil {
		return nil, err
	}

	for _, d := range ds {
		if d.Name == name {
			return &d, nil
		}
	}

	return nil, fmt.Errorf("domain not found")
}

// DeleteDomain deletes the Domain that matches the name
func (c *Client) DeleteDomain(d *Domain) (*SimpleResponse, error) {
	url := fmt.Sprintf("/v2/dns/%s", d.ID)
	resp, err := c.SendDeleteRequest(url)
	if err != nil {
		return nil, err
	}

	return c.DecodeSimpleResponse(resp)
}

// NewRecord creates a new DNS record
func (c *Client) NewRecord(r *RecordConfig) (*Record, error) {
	if len(r.DomainID) == 0 {
		return nil, fmt.Errorf("r.DomainID is empty")
	}

	url := fmt.Sprintf("/v2/dns/%s/records", r.DomainID)
	body, err := c.SendPostRequest(url, r)
	if err != nil {
		return nil, err
	}

	var record *Record
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(record); err != nil {
		return nil, err
	}

	return record, nil
}

// DeleteRecord deletes the DNS record
func (c *Client) DeleteRecord(r *Record) (*SimpleResponse, error) {
	if len(r.ID) == 0 {
		return nil, fmt.Errorf("r.ID is empty")
	}

	if len(r.DomainID) == 0 {
		return nil, fmt.Errorf("r.DomainID is empty")
	}

	url := fmt.Sprintf("/v2/dns/%s/records/%s", r.DomainID, r.ID)
	resp, err := c.SendDeleteRequest(url)
	if err != nil {
		return nil, err
	}

	return c.DecodeSimpleResponse(resp)
}
