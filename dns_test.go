package civogo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestDNSDomain_List(t *testing.T) {
	initServer()
	defer downServer()

	mux.HandleFunc("/v2/dns", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)

		allDNSDomain := []DNSDomain{
			{ID: "a5bd357a-8afd-4f60-b055-ece013646f55", Name: "example.com"},
			{ID: "2ecf2e0a-629c-4d16-9cb9-aa81059c4bad", Name: "example.net"},
		}
		value := toJSON(t, allDNSDomain)
		fmt.Fprint(w, value)
	})

	keys, meta, err := client.DNS().List(ctx)
	if err != nil {
		t.Errorf("DNS.List returned error: %v", err)
	}

	expectedKeys := []DNSDomain{
		{ID: "a5bd357a-8afd-4f60-b055-ece013646f55", Name: "example.com"},
		{ID: "2ecf2e0a-629c-4d16-9cb9-aa81059c4bad", Name: "example.net"},
	}
	if !reflect.DeepEqual(keys, expectedKeys) {
		t.Errorf("DNS.List returned keys %+v, expected %+v", keys, expectedKeys)
	}

	// compare status code
	if meta.StatusCode != http.StatusOK {
		t.Errorf("DNS.List returned status code %d, expected %d", meta.StatusCode, http.StatusOK)
	}
}

func TestDNSDomain_GetByID(t *testing.T) {
	initServer()
	defer downServer()

	mux.HandleFunc("/v2/dns/a5bd357a-8afd-4f60-b055-ece013646f55", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)

		oneDNSDomain := DNSDomain{ID: "a5bd357a-8afd-4f60-b055-ece013646f55", Name: "example.com"}
		value := toJSON(t, oneDNSDomain)
		fmt.Fprint(w, value)
	})

	keys, meta, err := client.DNS().GetByID(ctx, "a5bd357a-8afd-4f60-b055-ece013646f55")
	if err != nil {
		t.Errorf("DNS.GetByID returned error: %v", err)
	}

	expectedDomain := &DNSDomain{ID: "a5bd357a-8afd-4f60-b055-ece013646f55", Name: "example.com"}
	if !reflect.DeepEqual(keys, expectedDomain) {
		t.Errorf("DNS.GetByID returned keys %+v, expected %+v", keys, expectedDomain)
	}

	// compare status code
	if meta.StatusCode != http.StatusOK {
		t.Errorf("DNS.GetByID returned status code %d, expected %d", meta.StatusCode, http.StatusOK)
	}
}

func TestDNSDomain_Find(t *testing.T) {
	initServer()
	defer downServer()

	mux.HandleFunc("/v2/dns", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)

		allDNSDomain := []DNSDomain{
			{ID: "a5bd357a-8afd-4f60-b055-ece013646f55", Name: "example.com"},
			{ID: "2ecf2e0a-629c-4d16-9cb9-aa81059c4bad", Name: "example.net"},
		}
		value := toJSON(t, allDNSDomain)
		fmt.Fprint(w, value)
	})

	domain, meta, err := client.DNS().Find(ctx, "example.com")
	if err != nil {
		t.Errorf("DNS.Find returned error: %v", err)
	}

	expectedDomain := &DNSDomain{ID: "a5bd357a-8afd-4f60-b055-ece013646f55", Name: "example.com"}
	if !reflect.DeepEqual(domain, expectedDomain) {
		t.Errorf("DNS.Find returned keys %+v, expected %+v", domain, expectedDomain)
	}

	// compare status code
	if meta.StatusCode != http.StatusOK {
		t.Errorf("DNS.Find returned status code %d, expected %d", meta.StatusCode, http.StatusOK)
	}
}

func TestDNSDomain_Create(t *testing.T) {
	initServer()
	defer downServer()

	mux.HandleFunc("/v2/dns", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)

		respondDomain := DNSDomain{ID: "a5bd357a-8afd-4f60-b055-ece013646f55", Name: "example-1.com"}
		value := toJSON(t, respondDomain)
		fmt.Fprint(w, value)
	})

	newDomain := &DNSDomainCreateRequest{
		Name: "example-1.com",
	}

	result, meta, err := client.DNS().Create(ctx, newDomain)
	if err != nil {
		t.Errorf("DNS.Create returned error: %v", err)
	}

	expectedKeys := &DNSDomain{ID: "a5bd357a-8afd-4f60-b055-ece013646f55", Name: "example-1.com"}
	if !reflect.DeepEqual(result, expectedKeys) {
		t.Errorf("DNS.Create returned keys %+v, expected %+v", result, expectedKeys)
	}

	// compare status code
	if meta.StatusCode != http.StatusOK {
		t.Errorf("DNS.Create returned status code %d, expected %d", meta.StatusCode, http.StatusOK)
	}
}

func TestDNSDomain_Update(t *testing.T) {
	initServer()
	defer downServer()

	updateDomain := &DNSDomainUpdateRequest{
		Name: "example-1.com",
	}

	mux.HandleFunc("/v2/dns/78f64e5c-abd3-4f4d-85c8-ac63b50caa55", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)

		expectedUpdateDomain := map[string]interface{}{
			"name":  "example-1.com",
		}

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, expectedUpdateDomain) {
			t.Errorf("Request body = %#v, expected %#v", v, updateDomain)
		}

		respondDomain := DNSDomain{ID: "78f64e5c-abd3-4f4d-85c8-ac63b50caa55", Name: "example-1.com"}
		value := toJSON(t, respondDomain)
		fmt.Fprint(w, value)
	})

	result, meta, err := client.DNS().Update(ctx, "78f64e5c-abd3-4f4d-85c8-ac63b50caa55", updateDomain)
	if err != nil {
		t.Errorf("DNS.Update returned error: %v", err)
	}

	expectedNetwork := &DNSDomain{ID: "78f64e5c-abd3-4f4d-85c8-ac63b50caa55", Name: "example-1.com"}
	if !reflect.DeepEqual(result, expectedNetwork) {
		t.Errorf("DNS.Update returned keys %+v, expected %+v", result, expectedNetwork)
	}

	// compare status code
	if meta.StatusCode != http.StatusOK {
		t.Errorf("DNS.Update returned status code %d, expected %d", meta.StatusCode, http.StatusOK)
	}
}

func TestDNSDomain_Delete(t *testing.T) {
	initServer()
	defer downServer()

	mux.HandleFunc("/v2/dns/78f64e5c-abd3-4f4d-85c8-ac63b50caa55", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)

		respondSimple := SimpleResponse{
			ID:     "78f64e5c-abd3-4f4d-85c8-ac63b50caa55",
			Result: "success",
		}
		value := toJSON(t, respondSimple)
		fmt.Fprint(w, value)
	})

	result, meta, err := client.DNS().Delete(ctx, "78f64e5c-abd3-4f4d-85c8-ac63b50caa55")
	if err != nil {
		t.Errorf("DNS.Delete returned error: %v", err)
	}

	expectedNetwork := &SimpleResponse{ID: "78f64e5c-abd3-4f4d-85c8-ac63b50caa55", Result: "success"}
	if !reflect.DeepEqual(result, expectedNetwork) {
		t.Errorf("DNS.Delete returned keys %+v, expected %+v", result, expectedNetwork)
	}

	// compare status code
	if meta.StatusCode != http.StatusOK {
		t.Errorf("DNS.Delete returned status code %d, expected %d", meta.StatusCode, http.StatusOK)
	}
}

func TestDNSRecord_Records(t *testing.T) {
	initServer()
	defer downServer()

	mux.HandleFunc("/v2/dns/78f64e5c-abd3-4f4d-85c8-ac63b50caa55/records", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)

		respondRecords := []DNSRecord{
			{
				ID:   "a5bd357a-8afd-4f60-b055-ece013646f55",
				DNSDomainID: "78f64e5c-abd3-4f4d-85c8-ac63b50caa55",
				Type: "A",
				Name: "www.example.com",
				Value: "192.168.1.1",
			},
			{
				ID:   "8efee15f-89e7-4a13-a2df-dd78060a7185",
				DNSDomainID: "78f64e5c-abd3-4f4d-85c8-ac63b50caa55",
				Type: "A",
				Name: "www1.example.com",
				Value: "192.168.1.2",
			},
		}
		value := toJSON(t, respondRecords)
		fmt.Fprint(w, value)
	})

	keys, meta, err := client.DNS().Records(ctx, "78f64e5c-abd3-4f4d-85c8-ac63b50caa55")
	if err != nil {
		t.Errorf("DNS.ListRecords returned error: %v", err)
	}

	expectedKeys := []DNSRecord{
		{
			ID:   "a5bd357a-8afd-4f60-b055-ece013646f55",
			DNSDomainID: "78f64e5c-abd3-4f4d-85c8-ac63b50caa55",
			Type: "A",
			Name: "www.example.com",
			Value: "192.168.1.1",
		},
		{
			ID:   "8efee15f-89e7-4a13-a2df-dd78060a7185",
			DNSDomainID: "78f64e5c-abd3-4f4d-85c8-ac63b50caa55",
			Type: "A",
			Name: "www1.example.com",
			Value: "192.168.1.2",
		},
	}
	if !reflect.DeepEqual(keys, expectedKeys) {
		t.Errorf("DNS.ListRecords returned keys %+v, expected %+v", keys, expectedKeys)
	}

	// compare status code
	if meta.StatusCode != http.StatusOK {
		t.Errorf("DNS.ListRecords returned status code %d, expected %d", meta.StatusCode, http.StatusOK)
	}
}
