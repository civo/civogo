package civogo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

// Domain represents a domain registered within Civo's infrastructure
type Domain struct {
	// The ID of the domain
	ID string `json:"id"`

	// The ID of the account
	AccountID string `json:"account_id"`

	// The Name of the domain
	Name string `json:"name"`
}

type domainConfig struct {
	Name string `form:"name"`
}

// RecordType represents the allowed record types: a, cname, mx or txt
type RecordType string

// Record represents a DNS record registered within Civo's infrastructure
type Record struct {
	ID        string     `json:"id"`
	AccountID string     `json:"account_id"`
	DomainID  string     `json:"domain_id"`
	Name      string     `json:"name"`
	Value     string     `json:"value"`
	Type      RecordType `json:"type"`
	Priority  int        `json:"priority"`
	TTL       int        `json:"ttl"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// RecordConfig describes the parameters for a new DNS record
// none of the fields are mandatory and will be automatically
// set with default values
type RecordConfig struct {
	DomainID string     `form:"-"`
	Type     RecordType `form:"type"`
	Name     string     `form:"name"`
	Value    string     `form:"value"`
	Priority int        `form:"priority"`
	TTL      int        `form:"ttl"`
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

var (
	// ErrDomainNotFound is returned when the domain is not found
	ErrDomainNotFound = fmt.Errorf("domain not found")

	// ErrRecordNotFound is returned when the record is not found
	ErrRecordNotFound = fmt.Errorf("record not found")
)

// ListDomains returns all Domains owned by the calling API account
func (c *Client) ListDomains() ([]Domain, error) {
	url := "/v2/dns"

	resp, err := c.SendGetRequest(url)
	if err != nil {
		return nil, err
	}

	var ds = make([]Domain, 0)
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&ds); err != nil {
		return nil, err

	}

	return ds, nil
}

// CreateDomain registers a new Domain
func (c *Client) CreateDomain(name string) (*Domain, error) {
	url := "/v2/dns"
	d := &domainConfig{Name: name}
	body, err := c.SendPostRequest(url, d)
	if err != nil {
		return nil, err
	}

	var n = &Domain{}
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(n); err != nil {
		return nil, err
	}

	return n, nil
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

	return nil, ErrDomainNotFound
}

// UpdateDomain updates the provided domain with name
func (c *Client) UpdateDomain(d *Domain, name string) (*Domain, error) {
	url := fmt.Sprintf("/v2/dns/%s", d.ID)
	dc := &domainConfig{Name: name}
	body, err := c.SendPutRequest(url, dc)
	if err != nil {
		return nil, err
	}

	var r = &Domain{}
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(r); err != nil {
		return nil, err
	}

	return r, nil
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

// CreateRecord creates a new DNS record
func (c *Client) CreateRecord(r *RecordConfig) (*Record, error) {
	if len(r.DomainID) == 0 {
		return nil, fmt.Errorf("r.DomainID is empty")
	}

	url := fmt.Sprintf("/v2/dns/%s/records", r.DomainID)
	body, err := c.SendPostRequest(url, r)
	if err != nil {
		return nil, err
	}

	var record = &Record{}
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(record); err != nil {
		return nil, err
	}

	return record, nil
}

// ListRecords returns all the records associated with domainID
func (c *Client) ListRecords(domainID string) ([]Record, error) {
	url := fmt.Sprintf("/v2/dns/%s/records", domainID)
	resp, err := c.SendGetRequest(url)
	if err != nil {
		return nil, err
	}

	var rs = make([]Record, 0)
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&rs); err != nil {
		return nil, err

	}

	return rs, nil
}

// GetRecord returns the Record that matches the name and the domainID
func (c *Client) GetRecord(domainID, name string) (*Record, error) {
	rs, err := c.ListRecords(domainID)
	if err != nil {
		return nil, err
	}

	for _, r := range rs {
		if r.Name == name {
			return &r, nil
		}
	}

	return nil, ErrRecordNotFound
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
