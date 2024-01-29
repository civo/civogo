package civogo

import (
	"fmt"
	"strings"
)

// ListDNSDomains implemented in a fake way for automated tests
func (c *FakeClient) ListDNSDomains() ([]DNSDomain, error) {
	return c.Domains, nil
}

// FindDNSDomain implemented in a fake way for automated tests
func (c *FakeClient) FindDNSDomain(search string) (*DNSDomain, error) {
	for _, domain := range c.Domains {
		if strings.Contains(domain.Name, search) {
			return &domain, nil
		}
	}

	err := fmt.Errorf("unable to find %s, zero matches", search)
	return nil, ZeroMatchesError.wrap(err)
}

// CreateDNSDomain implemented in a fake way for automated tests
func (c *FakeClient) CreateDNSDomain(name string) (*DNSDomain, error) {
	domain := DNSDomain{
		ID:   c.generateID(),
		Name: name,
	}
	c.Domains = append(c.Domains, domain)
	return &domain, nil
}

// GetDNSDomain implemented in a fake way for automated tests
func (c *FakeClient) GetDNSDomain(name string) (*DNSDomain, error) {
	for _, domain := range c.Domains {
		if domain.Name == name {
			return &domain, nil
		}
	}

	return nil, ErrDNSDomainNotFound
}

// UpdateDNSDomain implemented in a fake way for automated tests
func (c *FakeClient) UpdateDNSDomain(d *DNSDomain, name string) (*DNSDomain, error) {
	for i, domain := range c.Domains {
		if domain.Name == d.Name {
			c.Domains[i] = *d
			return d, nil
		}
	}

	return nil, ErrDNSDomainNotFound
}

// DeleteDNSDomain implemented in a fake way for automated tests
func (c *FakeClient) DeleteDNSDomain(d *DNSDomain) (*SimpleResponse, error) {
	for i, domain := range c.Domains {
		if domain.Name == d.Name {
			c.Domains[len(c.Domains)-1], c.Domains[i] = c.Domains[i], c.Domains[len(c.Domains)-1]
			c.Domains = c.Domains[:len(c.Domains)-1]
			return &SimpleResponse{Result: "success"}, nil
		}
	}

	return nil, ErrDNSDomainNotFound
}

// CreateDNSRecord implemented in a fake way for automated tests
func (c *FakeClient) CreateDNSRecord(domainID string, r *DNSRecordConfig) (*DNSRecord, error) {
	record := DNSRecord{
		ID:          c.generateID(),
		DNSDomainID: domainID,
		Name:        r.Name,
		Value:       r.Value,
		Type:        r.Type,
	}

	c.DomainRecords = append(c.DomainRecords, record)
	return &record, nil
}

// ListDNSRecords implemented in a fake way for automated tests
func (c *FakeClient) ListDNSRecords(dnsDomainID string) ([]DNSRecord, error) {
	return c.DomainRecords, nil
}

// GetDNSRecord implemented in a fake way for automated tests
func (c *FakeClient) GetDNSRecord(domainID, domainRecordID string) (*DNSRecord, error) {
	for _, record := range c.DomainRecords {
		if record.ID == domainRecordID && record.DNSDomainID == domainID {
			return &record, nil
		}
	}

	return nil, ErrDNSRecordNotFound
}

// UpdateDNSRecord implemented in a fake way for automated tests
func (c *FakeClient) UpdateDNSRecord(r *DNSRecord, rc *DNSRecordConfig) (*DNSRecord, error) {
	for i, record := range c.DomainRecords {
		if record.ID == r.ID {
			record := DNSRecord{
				ID:          c.generateID(),
				DNSDomainID: record.DNSDomainID,
				Name:        rc.Name,
				Value:       rc.Value,
				Type:        rc.Type,
			}

			c.DomainRecords[i] = record
			return &record, nil
		}
	}

	return nil, ErrDNSRecordNotFound
}

// DeleteDNSRecord implemented in a fake way for automated tests
func (c *FakeClient) DeleteDNSRecord(r *DNSRecord) (*SimpleResponse, error) {
	for i, record := range c.DomainRecords {
		if record.ID == r.ID {
			c.DomainRecords[len(c.DomainRecords)-1], c.DomainRecords[i] = c.DomainRecords[i], c.DomainRecords[len(c.DomainRecords)-1]
			c.DomainRecords = c.DomainRecords[:len(c.DomainRecords)-1]
			return &SimpleResponse{Result: "success"}, nil
		}
	}

	return nil, ErrDNSRecordNotFound
}
