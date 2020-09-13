package civogo

import (
	"reflect"
	"testing"
	"time"
)

func TestDNSListDomains(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/dns": `[{"id": "12345", "account_id": "1", "name": "example.com"}, {"id": "12346", "account_id": "1", "name": "example.net"}]`,
	})
	defer server.Close()
	got, err := client.ListDNSDomains()

	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	expected := []DNSDomain{{ID: "12345", AccountID: "1", Name: "example.com"}, {ID: "12346", AccountID: "1", Name: "example.net"}}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestFindDNSDomain(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/dns": `[{"id": "12345", "account_id": "1", "name": "example.com"}, {"id": "12346", "account_id": "1", "name": "example.net"}]`,
	})
	defer server.Close()

	got, _ := client.FindDNSDomain("45")
	if got.ID != "12345" {
		t.Errorf("Expected %s, got %s", "12345", got.ID)
	}

	got, _ = client.FindDNSDomain("46")
	if got.ID != "12346" {
		t.Errorf("Expected %s, got %s", "12346", got.ID)
	}

	got, _ = client.FindDNSDomain("com")
	if got.ID != "12345" {
		t.Errorf("Expected %s, got %s", "12345", got.ID)
	}

	got, _ = client.FindDNSDomain("net")
	if got.ID != "12346" {
		t.Errorf("Expected %s, got %s", "12346", got.ID)
	}

	_, err := client.FindDNSDomain("example")
	if err.Error() != "MultipleMatchesError: unable to find example because there were multiple matches" {
		t.Errorf("Expected %s, got %s", "MultipleMatchesError: unable to find example because there were multiple matches", err.Error())
	}

	_, err = client.FindDNSDomain("missing")
	if err.Error() != "ZeroMatchesError: unable to find missing, zero matches" {
		t.Errorf("Expected %s, got %s", "ZeroMatchesError: unable to find missing, zero matches", err.Error())
	}
}

func TestCreateDomain(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/dns": `{"id": "12345", "account_id": "1", "name": "example.com"}`,
	})
	defer server.Close()
	got, err := client.CreateDNSDomain("example.com")

	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	expected := &DNSDomain{ID: "12345", AccountID: "1", Name: "example.com"}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestGetDNSDomain(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/dns": `[{"id": "12345", "account_id": "1", "name": "example.com"}, {"id": "12346", "account_id": "1", "name": "example.net"}]`,
	})
	defer server.Close()
	got, err := client.GetDNSDomain("example.net")

	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	expected := &DNSDomain{ID: "12346", AccountID: "1", Name: "example.net"}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}

	got, err = client.GetDNSDomain("example.io")
	if err != ErrDNSDomainNotFound {
		t.Errorf("Expected %+v, got %+v", ErrDNSDomainNotFound, got)
	}
}

func TestUpdateDomain(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/dns/12345": `{"id": "12345", "account_id": "1", "name": "example.com"}`,
	})
	defer server.Close()
	d := &DNSDomain{ID: "12345", AccountID: "1", Name: "example.com"}
	got, err := client.UpdateDNSDomain(d, "example.net")

	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &DNSDomain{ID: "12345", AccountID: "1", Name: "example.com"}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestDeleteDomain(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/dns/12346": `{"result": "success"}`,
	})
	defer server.Close()
	got, err := client.DeleteDNSDomain(&DNSDomain{ID: "12346"})
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &SimpleResponse{Result: "success"}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestNewRecord(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/dns/12346/records": `{
			"id": "76cc107f-fbef-4e2b-b97f-f5d34f4075d3",
			"created_at": "2019-04-11T12:47:56.000+01:00",
			"updated_at": "2019-04-11T12:47:56.000+01:00",
			"account_id": null,
			"domain_id": "12346",
			"name": "mail",
			"value": "10.0.0.1",
			"type": "MX",
			"priority": 10,
			"ttl": 600
		}`,
	})
	defer server.Close()

	cfg := &DNSRecordConfig{Name: "mail", Type: DNSRecordTypeMX, Value: "10.0.0.1", Priority: 10}
	got, err := client.CreateDNSRecord("12346", cfg)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &DNSRecord{
		ID:          "76cc107f-fbef-4e2b-b97f-f5d34f4075d3",
		DNSDomainID: "12346",
		Name:        "mail",
		Value:       "10.0.0.1",
		Type:        "MX",
		Priority:    10,
		TTL:         600,
	}

	if expected.ID != got.ID {
		t.Errorf("Expected %s, got %s", expected.ID, got.ID)
	}

	if expected.Name != got.Name {
		t.Errorf("Expected %s, got %s", expected.Name, got.Name)
	}

	if expected.Value != got.Value {
		t.Errorf("Expected %s, got %s", expected.Value, got.Value)
	}

	if expected.Type != got.Type {
		t.Errorf("Expected %s, got %s", expected.Type, got.Type)
	}

	if expected.Priority != got.Priority {
		t.Errorf("Expected %d, got %d", expected.Priority, got.Priority)
	}

	if expected.TTL != got.TTL {
		t.Errorf("Expected %d, got %d", expected.TTL, got.TTL)
	}
}

func TestDeleteRecord(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/dns/12346/records/76cc107f-fbef-4e2b-b97f-f5d34f4075d3": `{"result": "success"}`,
	})
	defer server.Close()

	r := &DNSRecord{
		ID:          "76cc107f-fbef-4e2b-b97f-f5d34f4075d3",
		DNSDomainID: "12346",
		Name:        "mail",
		Value:       "10.0.0.1",
		Type:        "mx",
		Priority:    10,
		TTL:         600,
	}

	got, err := client.DeleteDNSRecord(r)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &SimpleResponse{Result: ResultSuccess}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestListDNSRecords(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/dns/1111/records": `[{"id": "12345", "domain_id":"1111", "account_id": "1", "name": "www", "type": "CNAME", "value": "10.0.0.0", "ttl": 600}, {"id": "12346", "account_id": "1", "domain_id":"1111", "name": "mail", "type": "MX", "value": "10.0.0.1", "ttl": 600, "priority": 10}]`,
	})
	defer server.Close()
	got, err := client.ListDNSRecords("1111")

	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	expected := []DNSRecord{
		{ID: "12345", AccountID: "1", DNSDomainID: "1111", Name: "www", Value: "10.0.0.0", Type: DNSRecordTypeCName, TTL: 600},
		{ID: "12346", AccountID: "1", DNSDomainID: "1111", Name: "mail", Value: "10.0.0.1", Type: DNSRecordTypeMX, TTL: 600, Priority: 10},
	}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestUpdateDNSRecord(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/dns/edc5dacf-a2ad-4757-41ee-c12f06259c70/records/76cc107f-fbef-4e2b-b97f-f5d34f4075d3": `{
		  "id": "76cc107f-fbef-4e2b-b97f-f5d34f4075d3",
		  "created_at": "2019-04-11T12:47:56.000+01:00",
		  "updated_at": "2019-04-11T12:47:56.000+01:00",
		  "account_id": null,
		  "domain_id": "edc5dacf-a2ad-4757-41ee-c12f06259c70",
		  "name": "email",
		  "value": "10.0.0.1",
		  "type": "MX",
		  "priority": 10,
		  "ttl": 600
		}`,
	})
	defer server.Close()
	rc := &DNSRecordConfig{Name: "email"}
	r := &DNSRecord{ID: "76cc107f-fbef-4e2b-b97f-f5d34f4075d3", AccountID: "1", Name: "www", DNSDomainID: "edc5dacf-a2ad-4757-41ee-c12f06259c70"}
	got, err := client.UpdateDNSRecord(r, rc)

	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	createdAt, _ := time.Parse(time.RFC3339, "2019-04-11T12:47:56.000+01:00")
	updateAt, _ := time.Parse(time.RFC3339, "2019-04-11T12:47:56.000+01:00")

	expected := &DNSRecord{
		ID:          "76cc107f-fbef-4e2b-b97f-f5d34f4075d3",
		DNSDomainID: "edc5dacf-a2ad-4757-41ee-c12f06259c70",
		Name:        "email",
		Value:       "10.0.0.1",
		Type:        "MX",
		Priority:    10,
		TTL:         600,
		CreatedAt:   createdAt,
		UpdatedAt:   updateAt,
	}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestGetRecord(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/dns/1111/records": `[{"id": "12345", "domain_id":"1111", "account_id": "1", "name": "www", "type": "CNAME", "value": "10.0.0.0", "ttl": 600}, {"id": "12346", "account_id": "1", "domain_id":"1111", "name": "mail", "type": "MX", "value": "10.0.0.1", "ttl": 600, "priority": 10}]`,
	})

	defer server.Close()
	got, err := client.GetDNSRecord("1111", "12346")

	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	expected := &DNSRecord{ID: "12346", AccountID: "1", DNSDomainID: "1111", Name: "mail", Value: "10.0.0.1", Type: DNSRecordTypeMX, TTL: 600, Priority: 10}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}

	got, err = client.GetDNSRecord("1111", "hello")
	if err != ErrDNSRecordNotFound {
		t.Errorf("Expected %+v, got %+v", ErrDNSDomainNotFound, got)
		return
	}
}
