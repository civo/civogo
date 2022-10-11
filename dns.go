package civogo

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const (
	dnsBasePath = "/v2/dns"

	// DNSRecordTypeA represents an A record
	DNSRecordTypeA = "A"

	// DNSRecordTypeCName represents an CNAME record
	DNSRecordTypeCName = "CNAME"

	// DNSRecordTypeMX represents an MX record
	DNSRecordTypeMX = "MX"

	// DNSRecordTypeSRV represents an SRV record
	DNSRecordTypeSRV = "SRV"

	// DNSRecordTypeTXT represents an TXT record
	DNSRecordTypeTXT = "TXT"
)

// NetworkService is an interface for interfacing with the network
type DNSService interface {
	// Domain
	List(ctx context.Context) ([]DNSDomain, *Metadata, error)
	GetByID(ctx context.Context, domainID string) (*DNSDomain, *Metadata, error)
	Find(ctx context.Context, value string) (*DNSDomain, *Metadata, error)
	Create(ctx context.Context, createRequest *DNSDomainCreateRequest) (*DNSDomain, *Metadata, error)
	Update(ctx context.Context, domainID string, updateRequest *DNSDomainUpdateRequest) (*DNSDomain, *Metadata, error)
	Delete(ctx context.Context, domainID string) (*SimpleResponse, *Metadata, error)

	// Record
	Records(ctx context.Context, domainID string) ([]DNSRecord, *Metadata, error)
	// RecordsByType(ctx context.Context, domainID string, recordType string) ([]DNSRecord, *Metadata, error)
	// RecordsByName(ctx context.Context, domainID string, name string) ([]DNSRecord, *Metadata, error)
	// Record(ctx context.Context, domainID string, recordID string) (*DNSRecord, *Metadata, error)
	// CreateRecord(ctx context.Context, domainID string, createRequest *DNSRecordCreateRequest) (*DNSRecord, *Metadata, error)
	// UpdateRecord(ctx context.Context, domainID string, recordID string, updateRequest *DNSRecordUpdateRequest) (*DNSRecord, *Metadata, error)
	// DeleteRecord(ctx context.Context, domainID string, recordID string) (*SimpleResponse, *Metadata, error)

}

// DNSServiceOp Service used for communicating with the API
type DNSServiceOp struct {
	client *Client
}


type DNSGetter interface {
	SSHKey() SSHKeyService
}

// newSSHKey returns a SSHKey
func newDNS(c *Client) *DNSServiceOp {
	return &DNSServiceOp{
		client:  c,
	}
}

// DNSRecordType represents the type of DNS record
type DNSRecordType string

// DNSDomain represents a DNS domain
type DNSDomain struct {
	ID        string `json:"id"`
	AccountID string `json:"account_id"`
	Name      string `json:"name"`
}

// DNSDomainCreateRequest represents a request to create a DNS domain
type DNSDomainCreateRequest struct {
	Name string `json:"name"`
}

// DNSDomainUpdateRequest represents a request to update a DNS domain
type DNSDomainUpdateRequest struct {
	Name string `json:"name"`
}

// DNSRecord represents a DNS record registered within Civo's infrastructure
type DNSRecord struct {
	ID          string        `json:"id"`
	AccountID   string        `json:"account_id,omitempty"`
	DNSDomainID string        `json:"domain_id,omitempty"`
	Name        string        `json:"name,omitempty"`
	Value       string        `json:"value,omitempty"`
	Type        DNSRecordType `json:"type,omitempty"`
	Priority    int           `json:"priority,omitempty"`
	TTL         int           `json:"ttl,omitempty"`
	CreatedAt   time.Time     `json:"created_at,omitempty"`
	UpdatedAt   time.Time     `json:"updated_at,omitempty"`
}

// DNSRecordCreateRequest describes the parameters for a new DNS record
// none of the fields are mandatory and will be automatically
// set with default values
type DNSRecordCreateRequest struct {
	Type     DNSRecordType `json:"type"`
	Name     string        `json:"name"`
	Value    string        `json:"value"`
	Priority int           `json:"priority"`
	TTL      int           `json:"ttl"`
}

type DNSRecordUpdateRequest struct {
	Type     DNSRecordType `json:"type"`
	Name     string        `json:"name"`
	Value    string        `json:"value"`
	Priority int           `json:"priority"`
	TTL      int           `json:"ttl"`
}

// List list all dns domain for an account
func (c *DNSServiceOp) List(ctx context.Context) ([]DNSDomain, *Metadata, error) {
	path := dnsBasePath
	req, err := c.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new([]DNSDomain)
	resp, err := c.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return *root, resp, err
}

// GetByID get a dns domain by id
func (c *DNSServiceOp) GetByID(ctx context.Context, domainID string) (*DNSDomain, *Metadata, error) {
	path := fmt.Sprintf("%s/%s", dnsBasePath, domainID)
	req, err := c.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(DNSDomain)
	resp, err := c.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Find find a dns domain by name
func (c *DNSServiceOp) Find(ctx context.Context, value string) (*DNSDomain, *Metadata, error) {
	if value == "" {
		return nil, nil, errors.New("the search term cannot be empty")
	}

	allDNSDomain, meta, err := c.List(ctx)
	if err != nil {
		return nil, nil, err
	}

	for _, domain := range allDNSDomain {
		if domain.ID == value || domain.Name == value {
			return &domain, meta, nil
		}
	}

	return nil, nil, errors.New("no domain found for that value %s" + value)
}

// Create create a new dns domain
func (c *DNSServiceOp) Create(ctx context.Context, createRequest *DNSDomainCreateRequest) (*DNSDomain, *Metadata, error) {
	if createRequest == nil {
		return nil, nil, errors.New("createRequest is required")
	}

	path := dnsBasePath
	req, err := c.client.NewRequest(ctx, http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(DNSDomain)
	resp, err := c.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Update update a dns domain
func (c *DNSServiceOp) Update(ctx context.Context, domainID string, updateRequest *DNSDomainUpdateRequest) (*DNSDomain, *Metadata, error) {
	if domainID == "" {
		return nil, nil, errors.New("domainID is required")
	}

	if updateRequest == nil {
		return nil, nil, errors.New("updateRequest is required")
	}

	path := fmt.Sprintf("%s/%s", dnsBasePath, domainID)
	req, err := c.client.NewRequest(ctx, http.MethodPut, path, updateRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(DNSDomain)
	resp, err := c.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Delete delete a dns domain
func (c *DNSServiceOp) Delete(ctx context.Context, domainID string) (*SimpleResponse, *Metadata, error) {
	if domainID == "" {
		return nil, nil, errors.New("domainID is required")
	}

	path := fmt.Sprintf("%s/%s", dnsBasePath, domainID)
	req, err := c.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(SimpleResponse)
	resp, err := c.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// ListRecords list all dns records for a domain
func (c *DNSServiceOp) Records(ctx context.Context, domainID string) ([]DNSRecord, *Metadata, error) {
	if domainID == "" {
		return nil, nil, errors.New("domainID is required")
	}

	path := fmt.Sprintf("%s/%s/records", dnsBasePath, domainID)
	req, err := c.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new([]DNSRecord)
	resp, err := c.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return *root, resp, err
}
