package civogo

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestListInstances(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/instances": `{"page": 1, "per_page": 20, "pages": 2, "items":[{"id": "12345", "hostname": "foo.example.com"}]}`,
	})
	defer server.Close()

	got, err := client.ListInstances()
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	if got.Items[0].ID != "12345" {
		t.Errorf("Expected %s, got %s", "12345", got.Items[0].ID)
	}
	if got.Page != 1 {
		t.Errorf("Expected %d, got %d", 1, got.Page)
	}
	if got.Pages != 2 {
		t.Errorf("Expected %d, got %d", 2, got.Pages)
	}
	if got.PerPage != 20 {
		t.Errorf("Expected %d, got %d", 20, got.PerPage)
	}
}

func TestGetInstance(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/instances/12345": `{"id": "12345", "hostname": "foo.example.com"}`,
	})
	defer server.Close()

	got, err := client.GetInstance("12345")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	if got.ID != "12345" {
		t.Errorf("Expected %s, got %s", "12345", got.ID)
	}
	if got.Hostname != "foo.example.com" {
		t.Errorf("Expected %s, got %s", "foo.example.com", got.Hostname)
	}
}

func TestNewInstanceConfig(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/networks":  `[{"id": "1", "default": true, "name": "Default Network"}]`,
		"/v2/templates": `[{"id": "2", "code": "centos-7"},{"id": "3", "code": "ubuntu-18.04"}]`,
		"/v2/sshkeys":   `{"items":[{"id": "4", "name": "RSA Key"}]}`,
	})
	defer server.Close()

	got, err := client.NewInstanceConfig()
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	if got.Hostname == "" {
		t.Errorf("Expected hostname not to be blank, but it was")
	}
	if got.NetworkID != "1" {
		t.Errorf("Expected %s, got %s", "1", got.NetworkID)
	}
	if got.TemplateID != "3" {
		t.Errorf("Expected %s, got %s", "3", got.TemplateID)
	}
	if got.SSHKeyID != "4" {
		t.Errorf("Expected %s, got %s", "3", got.TemplateID)
	}
	if got.Count != 1 {
		t.Errorf("Expected %d, got %d", 1, got.Count)
	}
}

func TestCreateInstance(t *testing.T) {
	responses := map[string]string{
		"/v2/networks":  `[{"id": "1", "default": true, "name": "Default Network"}]`,
		"/v2/templates": `[{"id": "2", "code": "centos-7"},{"id": "3", "code": "ubuntu-18.04"}]`,
		"/v2/sshkeys":   `{"items":[{"id": "4", "name": "RSA Key"}]}`,
	}
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		for url, response := range responses {
			if strings.Contains(req.URL.String(), url) {
				rw.Write([]byte(response))
			}
		}

		// Grab the request body and set a new body, which will simulate the same data we read:
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			t.Errorf("Error reading body: %v", err)
			return
		}
		req.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		if strings.Contains(req.URL.String(), "/v2/instances") &&
			req.Method == "POST" &&
			string(body) == "count=1&hostname=foo.example.com&initial_user=civo&network_id=1&public_ip_required=true&region=lon1&reverse_dns=&script=&size=g2.xsmall&snapshot_id=&ssh_key_id=4&tags=&template_id=3" {
			rw.Write([]byte(`{"id": "12345", "hostname": "foo.example.com", "network_id": "1", "ssh_key": "4", "template_id": "3"}`))
		}

	}))
	defer server.Close()

	client, err := NewClientForTestingWithServer(server)
	if err != nil {
		t.Errorf("Creating a client returned an error: %s", err)
		return
	}

	config, err := client.NewInstanceConfig()
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	config.Hostname = "foo.example.com"

	got, err := client.CreateInstance(config)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	if got.Hostname != "foo.example.com" {
		t.Errorf("Expected %s, got %s", "1", got.NetworkID)
	}
	if got.NetworkID != "1" {
		t.Errorf("Expected %s, got %s", "1", got.NetworkID)
	}
	if got.TemplateID != "3" {
		t.Errorf("Expected %s, got %s", "3", got.TemplateID)
	}
	if got.SSHKey != "4" {
		t.Errorf("Expected %s, got %s", "3", got.TemplateID)
	}
}
