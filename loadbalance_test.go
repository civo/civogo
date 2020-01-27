package civogo

import (
	"reflect"
	"testing"
)

func TestListLoadBalance(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/loadbalancers": `[
		  {
			"id": "542e9eca-539d-45e6-b629-2f905d0b5f93",
			"hostname": "www.example.com",
			"protocol": "https",
			"port": "443",
			"max_request_size": 20,
			"tls_certificate": "...base64-encoded...",
			"tls_key": "...base64-encoded...",
			"policy": "random",
			"health_check_path": "/",
			"fail_timeout": 30,
			"max_conns": 10,
			"ignore_invalid_backend_tls": true,
			"backends": [
			  {
				"instance_id": "82ef8d8e-688c-4fc3-a31c-41746f27b074",
				"protocol": "http",
				"port": 3000
			  }
			]
		  }
		]`,
	})
	defer server.Close()
	got, err := client.ListLoadBalance()

	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	expected := []LoadBalance{{
		ID:                      "542e9eca-539d-45e6-b629-2f905d0b5f93",
		Hostname:                "www.example.com",
		Protocol:                "https",
		Port:                    "443",
		MaxRequestSize:          20,
		TlsCertificate:          "...base64-encoded...",
		TlsKey:                  "...base64-encoded...",
		Policy:                  "random",
		HealthCheckPath:         "/",
		FailTimeout:             30,
		MaxConns:                10,
		IgnoreInvalidBackendTls: true,
		Backends:                []LoadBalanceBackend{{InstanceID: "82ef8d8e-688c-4fc3-a31c-41746f27b074", Protocol: "http", Port: 3000}},
	}}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestNewLoadbalance(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/loadbalancers": `{
			"id": "542e9eca-539d-45e6-b629-2f905d0b5f93",
			"hostname": "www.example.com",
			"protocol": "https",
			"port": "443",
			"max_request_size": 20,
			"tls_certificate": "...base64-encoded...",
			"tls_key": "...base64-encoded...",
			"policy": "random",
			"health_check_path": "/",
			"fail_timeout": 30,
			"max_conns": 10,
			"ignore_invalid_backend_tls": true,
			"backends": [
			  {
				"instance_id": "82ef8d8e-688c-4fc3-a31c-41746f27b074",
				"protocol": "http",
				"port": 3000
			  }
			]
		  }`,
	})
	defer server.Close()

	cfg := &LoadBalanceConfig{
		Hostname:                "www.example.com",
		Protocol:                "https",
		Port:                    "443",
		MaxRequestSize:          20,
		TlsCertificate:          "...base64-encoded...",
		TlsKey:                  "...base64-encoded...",
		Policy:                  "random",
		HealthCheckPath:         "/",
		FailTimeout:             30,
		MaxConns:                10,
		IgnoreInvalidBackendTls: true,
		Backends: []LoadBalanceBackendConfig{{
			InstanceID: "82ef8d8e-688c-4fc3-a31c-41746f27b074",
			Protocol:   "http",
			Port:       3000,
		},
		},
	}
	got, err := client.NewLoadBalance(cfg)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &LoadBalance{
		ID:                      "542e9eca-539d-45e6-b629-2f905d0b5f93",
		Hostname:                "www.example.com",
		Protocol:                "https",
		Port:                    "443",
		MaxRequestSize:          20,
		TlsCertificate:          "...base64-encoded...",
		TlsKey:                  "...base64-encoded...",
		Policy:                  "random",
		HealthCheckPath:         "/",
		FailTimeout:             30,
		MaxConns:                10,
		IgnoreInvalidBackendTls: true,
		Backends:                []LoadBalanceBackend{{InstanceID: "82ef8d8e-688c-4fc3-a31c-41746f27b074", Protocol: "http", Port: 3000}},
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestUpdateLoadbalance(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/loadbalancers/542e9eca-539d-45e6-b629-2f905d0b5f93": `{
			"id": "542e9eca-539d-45e6-b629-2f905d0b5f93",
			"hostname": "www.example.com",
			"protocol": "https",
			"port": "443",
			"max_request_size": 20,
			"tls_certificate": "...base64-encoded...",
			"tls_key": "...base64-encoded...",
			"policy": "random",
			"health_check_path": "/",
			"fail_timeout": 30,
			"max_conns": 10,
			"ignore_invalid_backend_tls": true,
			"backends": [
			  {
				"instance_id": "82ef8d8e-688c-4fc3-a31c-41746f27b074",
				"protocol": "http",
				"port": 3000
			  },
			  {
				"instance_id": "85der56e-688c-4fc3-a31c-41746f27b074",
				"protocol": "http",
				"port": 3001
			  }
			]
		  }`,
	})
	defer server.Close()

	cfg := &LoadBalanceConfig{
		Backends: []LoadBalanceBackendConfig{{
			InstanceID: "85der56e-688c-4fc3-a31c-41746f27b074",
			Protocol:   "http",
			Port:       3001,
		},
		},
	}
	got, err := client.UpdateLoadBalance("542e9eca-539d-45e6-b629-2f905d0b5f93", cfg)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &LoadBalance{
		ID:                      "542e9eca-539d-45e6-b629-2f905d0b5f93",
		Hostname:                "www.example.com",
		Protocol:                "https",
		Port:                    "443",
		MaxRequestSize:          20,
		TlsCertificate:          "...base64-encoded...",
		TlsKey:                  "...base64-encoded...",
		Policy:                  "random",
		HealthCheckPath:         "/",
		FailTimeout:             30,
		MaxConns:                10,
		IgnoreInvalidBackendTls: true,
		Backends: []LoadBalanceBackend{
			{InstanceID: "82ef8d8e-688c-4fc3-a31c-41746f27b074", Protocol: "http", Port: 3000},
			{InstanceID: "85der56e-688c-4fc3-a31c-41746f27b074", Protocol: "http", Port: 3001},
		},
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestDeleteLoadBalance(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/loadbalancers/12345": `{"result": "success"}`,
	})
	defer server.Close()
	got, err := client.DeleteLoadBalance("12345")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &SimpleResponse{Result: "success"}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}
