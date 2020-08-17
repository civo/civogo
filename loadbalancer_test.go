package civogo

import (
	"reflect"
	"testing"
)

func TestListLoadBalancers(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/loadbalancers": `[
		  {
			"id": "542e9eca-539d-45e6-b629-2f905d0b5f93",
			"hostname": "www.example.com",
			"protocol": "https",
			"port": 443,
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
	got, err := client.ListLoadBalancers()

	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	expected := []LoadBalancer{{
		ID:                      "542e9eca-539d-45e6-b629-2f905d0b5f93",
		Hostname:                "www.example.com",
		Protocol:                "https",
		Port:                    443,
		MaxRequestSize:          20,
		TLSCertificate:          "...base64-encoded...",
		TLSKey:                  "...base64-encoded...",
		Policy:                  "random",
		HealthCheckPath:         "/",
		FailTimeout:             30,
		MaxConns:                10,
		IgnoreInvalidBackendTLS: true,
		Backends:                []LoadBalancerBackend{{InstanceID: "82ef8d8e-688c-4fc3-a31c-41746f27b074", Protocol: "http", Port: 3000}},
	}}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestFindLoadBalancer(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/loadbalancers": `[
		  { "id": "542e9eca-539d-45e6-b629-2f905d0b5f93", "hostname": "www.example.com", "protocol": "https", "port": 443, "max_request_size": 20, "tls_certificate": "...base64-encoded...", "tls_key": "...base64-encoded...", "policy": "random", "health_check_path": "/", "fail_timeout": 30, "max_conns": 10, "ignore_invalid_backend_tls": true, "backends": [ { "instance_id": "82ef8d8e-688c-4fc3-a31c-41746f27b074", "protocol": "http", "port": 3000 } ] },
		  { "id": "c33051ae-f337-45de-a3a5-004d822deff5", "hostname": "other.example.com", "protocol": "https", "port": 443, "max_request_size": 20, "tls_certificate": "...base64-encoded...", "tls_key": "...base64-encoded...", "policy": "random", "health_check_path": "/", "fail_timeout": 30, "max_conns": 10, "ignore_invalid_backend_tls": true, "backends": [ { "instance_id": "82ef8d8e-688c-4fc3-a31c-41746f27b074", "protocol": "http", "port": 3000 } ] }
		]`,
	})
	defer server.Close()

	got, err := client.FindLoadBalancer("542e9eca")
	if got.ID != "542e9eca-539d-45e6-b629-2f905d0b5f93" {
		t.Errorf("Expected %s, got %s", "542e9eca-539d-45e6-b629-2f905d0b5f93", got.ID)
	}

	got, _ = client.FindLoadBalancer("f337")
	if got.ID != "c33051ae-f337-45de-a3a5-004d822deff5" {
		t.Errorf("Expected %s, got %s", "c33051ae-f337-45de-a3a5-004d822deff5", got.ID)
	}

	got, _ = client.FindLoadBalancer("www")
	if got.ID != "542e9eca-539d-45e6-b629-2f905d0b5f93" {
		t.Errorf("Expected %s, got %s", "542e9eca-539d-45e6-b629-2f905d0b5f93", got.ID)
	}

	got, _ = client.FindLoadBalancer("other")
	if got.ID != "c33051ae-f337-45de-a3a5-004d822deff5" {
		t.Errorf("Expected %s, got %s", "c33051ae-f337-45de-a3a5-004d822deff5", got.ID)
	}

	_, err = client.FindLoadBalancer("example")
	if err.Error() != "MultipleMatchesError: unable to find example because there were multiple matches" {
		t.Errorf("Expected %s, got %s", "unable to find example because there were multiple matches", err.Error())
	}

	_, err = client.FindLoadBalancer("missing")
	if err.Error() != "ZeroMatchesError: unable to find missing, zero matches" {
		t.Errorf("Expected %s, got %s", "unable to find missing, zero matches", err.Error())
	}
}

func TestCreateLoadBalancer(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/loadbalancers": `{
			"id": "542e9eca-539d-45e6-b629-2f905d0b5f93",
			"hostname": "www.example.com",
			"protocol": "https",
			"port": 443,
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

	cfg := &LoadBalancerConfig{
		Hostname:                "www.example.com",
		Protocol:                "https",
		Port:                    443,
		MaxRequestSize:          20,
		TLSCertificate:          "...base64-encoded...",
		TLSKey:                  "...base64-encoded...",
		Policy:                  "random",
		HealthCheckPath:         "/",
		FailTimeout:             30,
		MaxConns:                10,
		IgnoreInvalidBackendTLS: true,
		Backends: []LoadBalancerBackendConfig{{
			InstanceID: "82ef8d8e-688c-4fc3-a31c-41746f27b074",
			Protocol:   "http",
			Port:       3000,
		},
		},
	}
	got, err := client.CreateLoadBalancer(cfg)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &LoadBalancer{
		ID:                      "542e9eca-539d-45e6-b629-2f905d0b5f93",
		Hostname:                "www.example.com",
		Protocol:                "https",
		Port:                    443,
		MaxRequestSize:          20,
		TLSCertificate:          "...base64-encoded...",
		TLSKey:                  "...base64-encoded...",
		Policy:                  "random",
		HealthCheckPath:         "/",
		FailTimeout:             30,
		MaxConns:                10,
		IgnoreInvalidBackendTLS: true,
		Backends:                []LoadBalancerBackend{{InstanceID: "82ef8d8e-688c-4fc3-a31c-41746f27b074", Protocol: "http", Port: 3000}},
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestUpdateLoadBalancer(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/loadbalancers/542e9eca-539d-45e6-b629-2f905d0b5f93": `{
			"id": "542e9eca-539d-45e6-b629-2f905d0b5f93",
			"hostname": "www.example.com",
			"protocol": "https",
			"port": 443,
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

	cfg := &LoadBalancerConfig{
		Backends: []LoadBalancerBackendConfig{{
			InstanceID: "85der56e-688c-4fc3-a31c-41746f27b074",
			Protocol:   "http",
			Port:       3001,
		},
		},
	}
	got, err := client.UpdateLoadBalancer("542e9eca-539d-45e6-b629-2f905d0b5f93", cfg)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &LoadBalancer{
		ID:                      "542e9eca-539d-45e6-b629-2f905d0b5f93",
		Hostname:                "www.example.com",
		Protocol:                "https",
		Port:                    443,
		MaxRequestSize:          20,
		TLSCertificate:          "...base64-encoded...",
		TLSKey:                  "...base64-encoded...",
		Policy:                  "random",
		HealthCheckPath:         "/",
		FailTimeout:             30,
		MaxConns:                10,
		IgnoreInvalidBackendTLS: true,
		Backends: []LoadBalancerBackend{
			{InstanceID: "82ef8d8e-688c-4fc3-a31c-41746f27b074", Protocol: "http", Port: 3000},
			{InstanceID: "85der56e-688c-4fc3-a31c-41746f27b074", Protocol: "http", Port: 3001},
		},
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestDeleteLoadBalancer(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/loadbalancers/12345": `{"result": "success"}`,
	})
	defer server.Close()
	got, err := client.DeleteLoadBalancer("12345")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &SimpleResponse{Result: "success"}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}
