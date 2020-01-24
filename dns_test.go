package civogo

import (
	"reflect"
	"testing"
)

func TestListDomains(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/dns": `[{"id": "12345", "account_id": "1", "name": "example.com"}, {"id": "12346", "account_id": "1", "name": "example.net"}]`,
	})
	defer server.Close()
	got, err := client.ListDomains()

	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	expected := []Domain{{ID: "12345", AccountID: "1", Name: "example.com"}, {ID: "12346", AccountID: "1", Name: "example.net"}}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestGetDomain(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/dns": `[{"id": "12345", "account_id": "1", "name": "example.com"}, {"id": "12346", "account_id": "1", "name": "example.net"}]`,
	})
	defer server.Close()
	got, err := client.GetDomain("example.net")

	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	expected := &Domain{ID: "12346", AccountID: "1", Name: "example.net"}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}

	got, err = client.GetDomain("example.io")
	if err != ErrDomainNotFound {
		t.Errorf("Expected %+v, got %+v", ErrDomainNotFound, got)
	}
}

func TestDeleteDomain(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/dns/12346": `{"result": "success"}`,
	})
	defer server.Close()
	got, err := client.DeleteDomain(&Domain{ID: "12346"})
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
			"type": "mx",
			"priority": 10,
			"ttl": 600
		}`,
	})
	defer server.Close()

	cfg := &RecordConfig{DomainID: "12346", Name: "mail", Type: RecordTypeMX, Value: "10.0.0.1", Priority: 10}
	got, err := client.NewRecord(cfg)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &Record{
		ID:       "76cc107f-fbef-4e2b-b97f-f5d34f4075d3",
		DomainID: "12346",
		Name:     "mail",
		Value:    "10.0.0.1",
		Type:     "mx",
		Priority: 10,
		TTL:      600,
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

	r := &Record{
		ID:       "76cc107f-fbef-4e2b-b97f-f5d34f4075d3",
		DomainID: "12346",
		Name:     "mail",
		Value:    "10.0.0.1",
		Type:     "mx",
		Priority: 10,
		TTL:      600,
	}

	got, err := client.DeleteRecord(r)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &SimpleResponse{Result: ResultSuccess}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestListRecords(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/dns/1111/records": `[{"id": "12345", "domain_id":"1111", "account_id": "1", "name": "www", "type": "cname", "value": "10.0.0.0", "ttl": 600}, {"id": "12346", "account_id": "1", "domain_id":"1111", "name": "mail", "type": "mx", "value": "10.0.0.1", "ttl": 600, "priority": 10}]`,
	})
	defer server.Close()
	got, err := client.ListRecords("1111")

	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	expected := []Record{
		{ID: "12345", AccountID: "1", DomainID: "1111", Name: "www", Value: "10.0.0.0", Type: RecordTypeCName, TTL: 600},
		{ID: "12346", AccountID: "1", DomainID: "1111", Name: "mail", Value: "10.0.0.1", Type: RecordTypeMX, TTL: 600, Priority: 10},
	}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestGetRecord(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/dns/1111/records": `[{"id": "12345", "domain_id":"1111", "account_id": "1", "name": "www", "type": "cname", "value": "10.0.0.0", "ttl": 600}, {"id": "12346", "account_id": "1", "domain_id":"1111", "name": "mail", "type": "mx", "value": "10.0.0.1", "ttl": 600, "priority": 10}]`,
	})

	defer server.Close()
	got, err := client.GetRecord("1111", "mail")

	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	expected := &Record{ID: "12346", AccountID: "1", DomainID: "1111", Name: "mail", Value: "10.0.0.1", Type: RecordTypeMX, TTL: 600, Priority: 10}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}

	got, err = client.GetRecord("1111", "hello")
	if err != ErrRecordNotFound {
		t.Errorf("Expected %+v, got %+v", ErrDomainNotFound, got)
		return
	}
}
